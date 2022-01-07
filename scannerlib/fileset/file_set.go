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
	"io"
	"io/ioutil"
	"os"
	"path"
	"regexp"
	"runtime/debug"
	"strings"

	"google.golang.org/protobuf/encoding/prototext"
	apb "github.com/google/localtoast/scannerlib/proto/api_go_proto"
	ipb "github.com/google/localtoast/scannerlib/proto/scan_instructions_go_proto"
)

const maxTraversalDepth = 100

var (
	procDirMatcher     = regexp.MustCompile(`^\d+$`)
	procEnvironMatcher = regexp.MustCompile("^(.*)=(.*)$")
	walkCounter        = 0
)

type filesystemReader interface {
	OpenFile(ctx context.Context, path string) (io.ReadCloser, error)
	FilesInDir(ctx context.Context, path string) ([]*apb.DirContent, error)
}

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

// WalkFunc is the type of the function called by WalkFiles to visit each file
// or directory in the FileSet. If the function returns an error, WalkFiles
// stops and returns the same error.
type WalkFunc func(path string, isDir bool) error

// WalkFiles calls walkFunc for each file described by the provided FileSet.
func WalkFiles(ctx context.Context, fileSet *ipb.FileSet, fsReader filesystemReader, walkFunc WalkFunc) error {
	switch {
	case fileSet.GetSingleFile() != nil:
		return walkFunc(fileSet.GetSingleFile().GetPath(), false)
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
			if err := walkFunc(f.GetDirPath(), true); err != nil {
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
			fsReader:          fsReader,
			walkFunc:          walkFunc,
		})
	case fileSet.GetProcessPath() != nil:
		return walkProcessPaths(ctx, fileSet.GetProcessPath().GetProcName(), fileSet.GetProcessPath().GetFileName(), fsReader, walkFunc)
	case fileSet.GetUnixEnvVarPaths() != nil:
		return walkVarPaths(ctx, fileSet.GetUnixEnvVarPaths(), fsReader, walkFunc)
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
	fsReader          filesystemReader
	walkFunc          WalkFunc
}

func walkFilesInDir(opts *walkFilesInDirOptions) error {
	if opts.depth > maxTraversalDepth {
		return fmt.Errorf("exceeded max traversal depth while traversing %s", opts.dirPath)
	}
	dirPath := path.Clean(opts.dirPath)
	dirContents, err := opts.fsReader.FilesInDir(opts.ctx, dirPath)
	if err != nil {
		// If the directory doesn't exist the checks are marked as non-compliant
		// instead of failing.
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}
		return err
	}
	for _, c := range dirContents {
		contentPath := path.Join(dirPath, c.GetName())
		if pathInOptOutList(contentPath, opts.optOutPathRegexes) {
			continue
		}
		matchesRegex := opts.filenameRegex == nil || opts.filenameRegex.MatchString(c.GetName())
		skipDirectory := c.GetIsDir() && opts.filesOnly
		skipFile := !c.GetIsDir() && opts.dirsOnly
		skipSymlink := c.GetIsSymlink() && opts.skipSymlinks
		if !skipDirectory && !skipFile && !skipSymlink && matchesRegex {
			if err := opts.walkFunc(contentPath, c.GetIsDir()); err != nil {
				return err
			}
			walkCounter++
			if walkCounter > 100 {
				walkCounter = 0
				debug.FreeOSMemory()
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
				fsReader:          opts.fsReader,
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
func walkProcessPaths(ctx context.Context, procName string, fileName string, fsReader filesystemReader, walkFunc WalkFunc) error {
	procDir, err := fsReader.FilesInDir(ctx, "/proc/")
	if err != nil {
		return fmt.Errorf("unable to enumerate /proc/: %v", err)
	}

	for _, f := range procDir {
		if !f.GetIsDir() || !procDirMatcher.MatchString(f.GetName()) {
			continue
		}
		dirName := path.Join("/proc", f.GetName())
		fh, err := fsReader.OpenFile(ctx, path.Join(dirName, "stat"))
		if err != nil {
			return fmt.Errorf("unable to read file %s/stat: %v", dirName, err)
		}
		defer fh.Close()
		stat, err := ioutil.ReadAll(fh)
		if err != nil {
			return fmt.Errorf("unable to read from file %s/stat: %v", dirName, err)
		}
		if procName != findProcName(string(stat)) {
			continue
		}
		if fileName == "" {
			// No filename was specified, check the directory itself.
			if err := walkFunc(dirName, true); err != nil {
				return err
			}
		} else {
			if err := walkFunc(path.Join(dirName, fileName), false); err != nil {
				return err
			}
		}
	}
	return nil
}

// walkVarPaths calls the walkFunc on all paths stored inside a UNIX environment
// variable such as $PATH. The paths are assumed to be separated by ':'s.
func walkVarPaths(ctx context.Context, evp *ipb.FileSet_UnixEnvVarPaths, fsReader filesystemReader, walkFunc WalkFunc) error {
	envVar, err := readEnvVar(ctx, evp.GetVarName(), fsReader)
	if err != nil {
		return err
	}

	for _, path := range strings.Split(envVar, ":") {
		isDir := true
		if _, err := fsReader.FilesInDir(ctx, path); err != nil {
			// If the error happens because of other reasons, the FileChecks
			// will catch it later.
			isDir = false
		}

		if isDir && evp.GetFilesOnly() {
			continue
		}
		if !isDir && evp.GetDirsOnly() {
			continue
		}

		if err := walkFunc(path, isDir); err != nil {
			return err
		}
	}

	return nil
}

// readEnvVar reads the value of the specified environment variable by parsing
// the contents of /proc/self/environ.
func readEnvVar(ctx context.Context, varName string, fsReader filesystemReader) (string, error) {
	r, err := fsReader.OpenFile(ctx, "/proc/self/environ")
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
