// Copyright 2021 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package fileset provides a utility for resolving the file paths defined by a FileSet proto.
package fileset

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"regexp"
	"strings"
	"time"

	"google.golang.org/protobuf/encoding/prototext"
	"github.com/google/localtoast/scanapi"
	apb "github.com/google/localtoast/scannerlib/proto/api_go_proto"
	ipb "github.com/google/localtoast/scannerlib/proto/scan_instructions_go_proto"
)

const maxTraversalDepth = 100

// PipelineToken is a wildcard token value to be used in scan instructions which will
// be replaced with the previous check result at runtime.
const PipelineToken = "%%pipeline%%"

var (
	procDirMatcher     = regexp.MustCompile(`^\d+$`)
	procEnvironMatcher = regexp.MustCompile("^(.*)=(.*)$")
)

// FileSetToString returns the string representation of a FileSet.
func FileSetToString(fileSet *ipb.FileSet) string {
	bytes, err := prototext.Marshal(fileSet)
	if err != nil {
		return "unknown file"
	}
	return string(bytes)
}

// ApplyOptOutConfig applies the directory traversal opt-out settings from an
// OptOutConfig to a FileSet.
func ApplyOptOutConfig(fileSet *ipb.FileSet, config *apb.OptOutConfig) {
	if fileSet.GetFilesInDir() == nil {
		return
	}
	fileSet.GetFilesInDir().OptOutPathRegexes =
		append(fileSet.GetFilesInDir().GetOptOutPathRegexes(), config.GetTraversalOptoutRegexes()...)
}

// ApplyReplacementConfig applies the path replacement settings from a
// ReplacementConfig to a FileSet.
func ApplyReplacementConfig(fileSet *ipb.FileSet, config *apb.ReplacementConfig) {
	if config == nil {
		return
	}
	for prefix, replacement := range config.PathPrefixReplacements {
		applyPathPrefixReplacement(fileSet, prefix, replacement)
	}
}

// ApplyPipelineTokenReplacement replaces the File Path if the wildcard is set
func ApplyPipelineTokenReplacement(fileSet *ipb.FileSet, prvRes string) {
	if fileSet.GetSingleFile() != nil {

		if fileSet.GetSingleFile().Path == PipelineToken {
			fileSet.GetSingleFile().Path = prvRes
		}
	}
}

func applyPathPrefixReplacement(fileSet *ipb.FileSet, prefix, replacement string) {
	switch {
	case fileSet.GetSingleFile() != nil:
		fileSet.GetSingleFile().Path = replacePrefix(fileSet.GetSingleFile().GetPath(), prefix, replacement)
	case fileSet.GetFilesInDir() != nil:
		fileSet.GetFilesInDir().DirPath = replacePrefix(fileSet.GetFilesInDir().GetDirPath(), prefix, replacement)
	}
}

func replacePrefix(str, prefix, replacement string) string {
	if !strings.HasSuffix(prefix, "/") {
		prefix += "/"
	}
	if !strings.HasSuffix(replacement, "/") {
		replacement += "/"
	}
	if !strings.HasPrefix(str, prefix) {
		return str
	}
	return replacement + strings.TrimPrefix(str, prefix)
}

// WalkFunc is the type of the function called by WalkFiles to visit each file
// or directory in the FileSet. If the function returns an error, WalkFiles
// stops and returns the same error.
type WalkFunc func(path string, isDir bool, traversingDir bool) error

// WalkFiles calls walkFunc for each file described by the provided FileSet.
func WalkFiles(ctx context.Context, fileSet *ipb.FileSet, fs scanapi.Filesystem, timeout time.Time, walkFunc WalkFunc) error {
	if err := checkTimeout(timeout); err != nil {
		return err
	}
	switch {
	case fileSet.GetSingleFile() != nil:
		return walkFunc(fileSet.GetSingleFile().GetPath(), false, false)
	case fileSet.GetFilesInDir() != nil:
		f := fileSet.GetFilesInDir()
		filenameRegex, optOutPathRegexes, err := createFilterRegexes(f)
		if err != nil {
			return err
		}
		if pathInOptOutList(f.GetDirPath(), optOutPathRegexes) {
			return nil
		}
		// Walk the root directory first (unless it's filtered out).
		matchesRegex := filenameRegex == nil || filenameRegex.MatchString(path.Base(f.GetDirPath()))
		if !f.GetFilesOnly() && matchesRegex {
			if err := walkFunc(f.GetDirPath(), true, true); err != nil {
				return err
			}
		}
		// Walk the sub-directories next.
		return walkFilesInDir(&walkFilesInDirOptions{
			ctx:               ctx,
			dirPath:           f.GetDirPath(),
			depth:             1,
			recursive:         f.GetRecursive(),
			filesOnly:         f.GetFilesOnly(),
			dirsOnly:          f.GetDirsOnly(),
			skipSymlinks:      f.GetSkipSymlinks(),
			filenameRegex:     filenameRegex,
			optOutPathRegexes: optOutPathRegexes,
			timeout:           timeout,
			fs:                fs,
			walkFunc:          walkFunc,
		})
	case fileSet.GetProcessPath() != nil:
		return walkProcessPaths(ctx, fileSet.GetProcessPath().GetProcName(), fileSet.GetProcessPath().GetFileName(), fileSet.GetProcessPath().GetCliArgRegex(), timeout, fs, walkFunc)
	case fileSet.GetUnixEnvVarPaths() != nil:
		return walkVarPaths(ctx, fileSet.GetUnixEnvVarPaths(), timeout, fs, walkFunc)
	default:
		return fmt.Errorf("Unknown FilePath type %v", fileSet.GetFilePath())
	}
}

func createFilterRegexes(f *ipb.FileSet_FilesInDir) (*regexp.Regexp, []*regexp.Regexp, error) {
	var filenameRegex *regexp.Regexp
	if len(f.GetFilenameRegex()) > 0 {
		var err error
		filenameRegex, err = regexp.Compile("^" + f.GetFilenameRegex() + "$")
		if err != nil {
			return nil, nil, err
		}
	}

	optOutPathRegexes := make([]*regexp.Regexp, 0, len(f.GetOptOutPathRegexes()))
	for _, o := range f.GetOptOutPathRegexes() {
		re, err := regexp.Compile("^" + o + "$")
		if err != nil {
			return nil, nil, err
		}
		optOutPathRegexes = append(optOutPathRegexes, re)
	}

	return filenameRegex, optOutPathRegexes, nil
}

type walkFilesInDirOptions struct {
	ctx               context.Context
	dirPath           string
	depth             int
	recursive         bool
	filesOnly         bool
	dirsOnly          bool
	skipSymlinks      bool
	filenameRegex     *regexp.Regexp
	optOutPathRegexes []*regexp.Regexp
	fs                scanapi.Filesystem
	timeout           time.Time
	walkFunc          WalkFunc
}

func walkFilesInDir(opts *walkFilesInDirOptions) error {
	if opts.depth > maxTraversalDepth {
		return fmt.Errorf("exceeded max traversal depth while traversing %s", opts.dirPath)
	}
	dirPath := path.Clean(opts.dirPath)
	d, err := opts.fs.OpenDir(opts.ctx, dirPath)
	if err != nil {
		// If the directory doesn't exist the checks are marked as non-compliant
		// instead of failing.
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}
		return err
	}
	defer d.Close()
	for d.Next() {
		c, err := d.Entry()
		if err != nil {
			return err
		}
		contentPath := path.Join(dirPath, c.GetName())
		if pathInOptOutList(contentPath, opts.optOutPathRegexes) {
			continue
		}
		matchesRegex := opts.filenameRegex == nil || opts.filenameRegex.MatchString(c.GetName())
		skipDirectory := c.GetIsDir() && opts.filesOnly
		skipFile := !c.GetIsDir() && opts.dirsOnly
		skipSymlink := c.GetIsSymlink() && opts.skipSymlinks
		if !skipDirectory && !skipFile && !skipSymlink && matchesRegex {
			if err := opts.walkFunc(contentPath, c.GetIsDir(), true); err != nil {
				return err
			}
			if err := checkTimeout(opts.timeout); err != nil {
				return err
			}
		}
		if opts.recursive && c.GetIsDir() {
			if err := walkFilesInDir(&walkFilesInDirOptions{
				ctx:               opts.ctx,
				dirPath:           contentPath,
				depth:             opts.depth + 1,
				recursive:         opts.recursive,
				filesOnly:         opts.filesOnly,
				dirsOnly:          opts.dirsOnly,
				skipSymlinks:      opts.skipSymlinks,
				filenameRegex:     opts.filenameRegex,
				optOutPathRegexes: opts.optOutPathRegexes,
				timeout:           opts.timeout,
				fs:                opts.fs,
				walkFunc:          opts.walkFunc,
			}); err != nil {
				return err
			}
		}
	}
	return nil
}

func pathInOptOutList(dirPath string, optOutPathRegexes []*regexp.Regexp) bool {
	for _, re := range optOutPathRegexes {
		if re.MatchString(dirPath) {
			return true
		}
	}
	return false
}

// walkProcessPaths iterates over all directories in /proc/ that have a numeric identifier
// and calls the walkFunc on all of them for which the procName is set in the stat file.
//
// Please note means that all those folders in /proc/ are traversed every time this function
// is called. This is fine as long as there are not many checks using the ProcessPath option.
func walkProcessPaths(ctx context.Context, procName string, fileName string, cliArgRegex string, timeout time.Time, fs scanapi.Filesystem, walkFunc WalkFunc) error {
	d, err := fs.OpenDir(ctx, "/proc/")
	if err != nil {
		return fmt.Errorf("unable to enumerate /proc/: %v", err)
	}
	defer d.Close()

	for d.Next() {
		f, err := d.Entry()
		if err != nil {
			return err
		}
		if !f.GetIsDir() || !procDirMatcher.MatchString(f.GetName()) {
			continue
		}
		dirName := path.Join("/proc", f.GetName())
		fh, err := fs.OpenFile(ctx, path.Join(dirName, "stat"))
		if err != nil {
			if errors.Is(err, os.ErrNotExist) {
				// The file got removed since we queried it, ignore.
				continue
			}
			return fmt.Errorf("unable to read file %s/stat: %v", dirName, err)
		}
		defer fh.Close()
		stat, err := ioutil.ReadAll(fh)
		if err != nil {
			// The file got removed since we queried it, ignore.
			continue
		}
		if procName != findProcName(string(stat)) {
			continue
		}

		// If the optional argument "cli_arg_regex" is specified,
		// run the relative regexp on the process cmdline content
		if cliArgRegex != "" {
			// Compile regex
			compiledCliArgRegex, err := regexp.Compile("^" + cliArgRegex + "$")
			if err != nil {
				return fmt.Errorf("unable to compile cli arg regex %q:\n%v", cliArgRegex, err)
			}

			// Open cmdline file and get content
			fh, err := fs.OpenFile(ctx, path.Join(dirName, "cmdline"))
			if err != nil {
				if errors.Is(err, os.ErrNotExist) {
					// The file got removed since we queried it, ignore.
					continue
				}
				return fmt.Errorf("unable to read file %s/cmdline: %v", dirName, err)
			}
			defer fh.Close()
			cmdline, err := ioutil.ReadAll(fh)
			if err != nil {
				// The file got removed since we queried it, ignore.
				continue
			}

			// Skip this process if regex does not match cmdline
			if !compiledCliArgRegex.MatchString(string(cmdline)) {
				continue
			}
		}

		if fileName == "" {
			// No filename was specified, check the directory itself.
			if err := walkFunc(dirName, true, false); err != nil {
				return err
			}
		} else {
			if err := walkFunc(path.Join(dirName, fileName), false, false); err != nil {
				return err
			}
		}

		if err := checkTimeout(timeout); err != nil {
			return err
		}
	}
	return nil
}

// walkVarPaths calls the walkFunc on all paths stored inside a UNIX environment
// variable such as $PATH. The paths are assumed to be separated by ':'s.
func walkVarPaths(ctx context.Context, evp *ipb.FileSet_UnixEnvVarPaths, timeout time.Time, fs scanapi.Filesystem, walkFunc WalkFunc) error {
	envVar, err := readEnvVar(ctx, evp.GetVarName(), fs)
	if err != nil {
		return err
	}

	for _, path := range strings.Split(envVar, ":") {
		isDir := true
		if d, err := fs.OpenDir(ctx, path); err != nil {
			// If the error happens because of other reasons, the FileChecks
			// will catch it later.
			isDir = false
		} else {
			d.Close()
		}

		if isDir && evp.GetFilesOnly() {
			continue
		}
		if !isDir && evp.GetDirsOnly() {
			continue
		}

		if err := walkFunc(path, isDir, false); err != nil {
			return err
		}

		if err := checkTimeout(timeout); err != nil {
			return err
		}
	}

	return nil
}

// readEnvVar reads the value of the specified environment variable by parsing
// the contents of /proc/self/environ.
func readEnvVar(ctx context.Context, varName string, fs scanapi.Filesystem) (string, error) {
	r, err := fs.OpenFile(ctx, "/proc/self/environ")
	if err != nil {
		return "", err
	}
	defer r.Close()

	scanner := bufio.NewScanner(r)
	// Entries are separated by null bytes.
	scanner.Split(func(data []byte, atEOF bool) (advance int, token []byte, err error) {
		if atEOF && len(data) == 0 {
			return 0, nil, nil
		}
		if i := bytes.IndexByte(data, 0); i >= 0 {
			return i + 1, data[0:i], nil
		}
		if atEOF {
			return len(data), data, nil
		}
		return 0, nil, nil
	})

	for scanner.Scan() {
		if scanner.Err() != nil {
			return "", scanner.Err()
		}
		groups := procEnvironMatcher.FindStringSubmatch(scanner.Text())
		if groups != nil && groups[1] == varName {
			return groups[2], nil
		}
	}
	return "", fmt.Errorf("%s not found among environment variables", varName)
}

func findProcName(stat string) string {
	start := strings.Index(stat, "(")
	end := strings.Index(stat, ")")
	if start == -1 || end == -1 {
		return ""
	}
	return stat[start+1 : end]
}

func checkTimeout(timeout time.Time) error {
	if timeout.IsZero() || time.Now().Before(timeout) {
		return nil
	}
	return errors.New("scan timed out")
}
