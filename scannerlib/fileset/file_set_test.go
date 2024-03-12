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

package fileset_test

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/testing/protocmp"
	"github.com/google/localtoast/scanapi"
	"github.com/google/localtoast/scannerlib/fileset"
	apb "github.com/google/localtoast/scannerlib/proto/api_go_proto"
	ipb "github.com/google/localtoast/scannerlib/proto/scan_instructions_go_proto"
)

const (
	procEnvironPath    = "/proc/self/environ"
	procEnvironContent = "SHELL=/bin/bash\x00PATH=/root:/root/file1.txt"
)

type fakeDirectoryReader struct{}

func (fakeDirectoryReader) FilePermissions(ctx context.Context, path string) (*apb.PosixPermissions, error) {
	return nil, errors.New("Not implemented")
}

func (fakeDirectoryReader) OpenFile(ctx context.Context, path string) (io.ReadCloser, error) {
	switch path {
	case procEnvironPath:
		return io.NopCloser(bytes.NewReader([]byte(procEnvironContent))), nil
	default:
		return nil, os.ErrNotExist
	}
}

func (fakeDirectoryReader) OpenDir(ctx context.Context, path string) (scanapi.DirReader, error) {
	// Fake dir structure:
	// root---file1.txt
	//  \ \ \-file2.gif
	//   \ \--symlink
	//    \---subdir--file3.txt
	switch path {
	case "/root":
		return scanapi.SliceToDirReader([]*apb.DirContent{
			&apb.DirContent{Name: "file1.txt", IsDir: false},
			&apb.DirContent{Name: "file2.gif", IsDir: false},
			&apb.DirContent{Name: "symlink", IsDir: false, IsSymlink: true},
			&apb.DirContent{Name: "subdir", IsDir: true},
		}), nil
	case "/root/subdir":
		return scanapi.SliceToDirReader([]*apb.DirContent{
			&apb.DirContent{Name: "file3.txt", IsDir: false},
		}), nil
	default:
		return nil, os.ErrNotExist
	}
}

func TestSingleFile(t *testing.T) {
	expectedPath := "/path/to/file"
	fileSet := &ipb.FileSet{
		FilePath: &ipb.FileSet_SingleFile_{SingleFile: &ipb.FileSet_SingleFile{Path: expectedPath}},
	}

	err := fileset.WalkFiles(context.Background(), fileSet, &fakeDirectoryReader{}, time.Time{}, func(walkedPath string, isDir bool, traversingDir bool) error {
		if expectedPath != walkedPath {
			t.Errorf("fileset.WalkFiles(%v) expected to walk on path %s, got %s",
				fileSet, expectedPath, walkedPath)
		}
		return nil
	})
	if err != nil {
		t.Fatalf("fileset.WalkFiles(%v) returned an error: %v", fileSet, err)
	}
}

type traversal struct {
	Path          string
	IsDir         bool
	TraversingDir bool
}

func TestFilesInDir(t *testing.T) {
	testCases := []struct {
		description       string
		fileSet           *ipb.FileSet
		expectedTraversal []*traversal
	}{
		{
			description: "visit all files recursively",
			fileSet: &ipb.FileSet{
				FilePath: &ipb.FileSet_FilesInDir_{FilesInDir: &ipb.FileSet_FilesInDir{
					DirPath:   "/root",
					Recursive: true,
				}},
			},
			expectedTraversal: []*traversal{
				{Path: "/root", IsDir: true, TraversingDir: true},
				{Path: "/root/file1.txt", IsDir: false, TraversingDir: true},
				{Path: "/root/file2.gif", IsDir: false, TraversingDir: true},
				{Path: "/root/symlink", IsDir: false, TraversingDir: true},
				{Path: "/root/subdir", IsDir: true, TraversingDir: true},
				{Path: "/root/subdir/file3.txt", IsDir: false, TraversingDir: true},
			},
		},
		{
			description: "visit only files",
			fileSet: &ipb.FileSet{
				FilePath: &ipb.FileSet_FilesInDir_{FilesInDir: &ipb.FileSet_FilesInDir{
					DirPath:   "/root",
					Recursive: true,
					FilesOnly: true,
				}},
			},
			expectedTraversal: []*traversal{
				{Path: "/root/file1.txt", IsDir: false, TraversingDir: true},
				{Path: "/root/file2.gif", IsDir: false, TraversingDir: true},
				{Path: "/root/symlink", IsDir: false, TraversingDir: true},
				{Path: "/root/subdir/file3.txt", IsDir: false, TraversingDir: true},
			},
		},
		{
			description: "visit only directories",
			fileSet: &ipb.FileSet{
				FilePath: &ipb.FileSet_FilesInDir_{FilesInDir: &ipb.FileSet_FilesInDir{
					DirPath:   "/root",
					Recursive: true,
					DirsOnly:  true,
				}},
			},
			expectedTraversal: []*traversal{
				{Path: "/root", IsDir: true, TraversingDir: true},
				{Path: "/root/subdir", IsDir: true, TraversingDir: true},
			},
		},
		{
			description: "omit symlinks",
			fileSet: &ipb.FileSet{
				FilePath: &ipb.FileSet_FilesInDir_{FilesInDir: &ipb.FileSet_FilesInDir{
					DirPath:      "/root",
					Recursive:    true,
					SkipSymlinks: true,
				}},
			},
			expectedTraversal: []*traversal{
				{Path: "/root", IsDir: true, TraversingDir: true},
				{Path: "/root/file1.txt", IsDir: false, TraversingDir: true},
				{Path: "/root/file2.gif", IsDir: false, TraversingDir: true},
				{Path: "/root/subdir", IsDir: true, TraversingDir: true},
				{Path: "/root/subdir/file3.txt", IsDir: false, TraversingDir: true},
			},
		},
		{
			description: "visit files in the root dir",
			fileSet: &ipb.FileSet{
				FilePath: &ipb.FileSet_FilesInDir_{FilesInDir: &ipb.FileSet_FilesInDir{
					DirPath:   "/root",
					Recursive: false,
				}},
			},
			expectedTraversal: []*traversal{
				{Path: "/root", IsDir: true, TraversingDir: true},
				{Path: "/root/file1.txt", IsDir: false, TraversingDir: true},
				{Path: "/root/file2.gif", IsDir: false, TraversingDir: true},
				{Path: "/root/symlink", IsDir: false, TraversingDir: true},
				{Path: "/root/subdir", IsDir: true, TraversingDir: true},
			},
		},
		{
			description: "filter file names to visit",
			fileSet: &ipb.FileSet{
				FilePath: &ipb.FileSet_FilesInDir_{FilesInDir: &ipb.FileSet_FilesInDir{
					DirPath:       "/root",
					Recursive:     true,
					FilenameRegex: ".*\\.txt",
				}},
			},
			expectedTraversal: []*traversal{
				{Path: "/root/file1.txt", IsDir: false, TraversingDir: true},
				{Path: "/root/subdir/file3.txt", IsDir: false, TraversingDir: true},
			},
		},
		{
			description: "opt out file paths",
			fileSet: &ipb.FileSet{
				FilePath: &ipb.FileSet_FilesInDir_{FilesInDir: &ipb.FileSet_FilesInDir{
					DirPath:           "/root",
					Recursive:         true,
					OptOutPathRegexes: []string{"/root/sub.*"},
				}},
			},
			expectedTraversal: []*traversal{
				{Path: "/root", IsDir: true, TraversingDir: true},
				{Path: "/root/file1.txt", IsDir: false, TraversingDir: true},
				{Path: "/root/file2.gif", IsDir: false, TraversingDir: true},
				{Path: "/root/symlink", IsDir: false, TraversingDir: true},
			},
		},
		{
			description: "opt out root directory",
			fileSet: &ipb.FileSet{
				FilePath: &ipb.FileSet_FilesInDir_{FilesInDir: &ipb.FileSet_FilesInDir{
					DirPath:           "/root",
					Recursive:         true,
					OptOutPathRegexes: []string{"/root"},
				}},
			},
			expectedTraversal: []*traversal{},
		},
		{
			description: "opt out root directory with regex",
			fileSet: &ipb.FileSet{
				FilePath: &ipb.FileSet_FilesInDir_{FilesInDir: &ipb.FileSet_FilesInDir{
					DirPath:           "/root",
					Recursive:         true,
					OptOutPathRegexes: []string{"/roo.*"},
				}},
			},
			expectedTraversal: []*traversal{},
		},
		{
			description: "directory doesn't exist",
			fileSet: &ipb.FileSet{
				FilePath: &ipb.FileSet_FilesInDir_{FilesInDir: &ipb.FileSet_FilesInDir{
					DirPath:   "/non-existent",
					Recursive: true,
				}},
			},
			// The directory is still traversed so that the checks can report
			// non-compliance.
			expectedTraversal: []*traversal{
				{Path: "/non-existent", IsDir: true, TraversingDir: true},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			gotTraversal := []*traversal{}
			err := fileset.WalkFiles(context.Background(), tc.fileSet, &fakeDirectoryReader{}, time.Time{}, func(walkedPath string, isDir bool, traversingDir bool) error {
				gotTraversal = append(gotTraversal, &traversal{walkedPath, isDir, traversingDir})
				return nil
			})
			if err != nil {
				t.Fatalf("fileset.WalkFiles(%v) returned an error: %v", tc.fileSet, err)
			}
			if diff := cmp.Diff(tc.expectedTraversal, gotTraversal); diff != "" {
				t.Errorf("fileset.WalkFiles(%v) made an unexpected traversal diff (-want +got):\n%s",
					tc.fileSet, diff)
			}
		})
	}
}

type infiniteLoopFSReader struct{}

func (infiniteLoopFSReader) OpenFile(ctx context.Context, path string) (io.ReadCloser, error) {
	return nil, errors.New("Not implemented")
}

func (infiniteLoopFSReader) FilePermissions(ctx context.Context, path string) (*apb.PosixPermissions, error) {
	return nil, errors.New("Not implemented")
}

func (infiniteLoopFSReader) OpenDir(ctx context.Context, path string) (scanapi.DirReader, error) {
	return scanapi.SliceToDirReader([]*apb.DirContent{
		&apb.DirContent{Name: "dir", IsDir: true},
	}), nil
}

func TestTraverseFilesystemWithInfiniteLoop(t *testing.T) {
	files := &ipb.FileSet{FilePath: &ipb.FileSet_FilesInDir_{FilesInDir: &ipb.FileSet_FilesInDir{
		DirPath:   "/",
		Recursive: true,
	}}}
	err := fileset.WalkFiles(context.Background(), files, &infiniteLoopFSReader{}, time.Time{}, func(walkedPath string, isDir bool, traversingDir bool) error { return nil })
	if err == nil {
		t.Fatalf("fileset.WalkFiles(%v) didn't return an error", files)
	}
}

type fakeProcessPathReader struct {
	pidToName    map[int]string
	pidToCLIArgs map[int]string
	// If true, tests the race condition when the files become unavailable after
	// they're listed in the /proc directory.
	removeFilesAfterQuery bool
}

func (fakeProcessPathReader) FilePermissions(ctx context.Context, path string) (*apb.PosixPermissions, error) {
	return nil, errors.New("Not implemented")
}

func (r fakeProcessPathReader) OpenFile(ctx context.Context, path string) (io.ReadCloser, error) {
	if r.removeFilesAfterQuery {
		return nil, os.ErrNotExist
	}

	var re *regexp.Regexp
	if strings.HasSuffix(path, "cmdline") {
		re = regexp.MustCompile(`^/proc/([0-9]+)/cmdline`)
	} else {
		re = regexp.MustCompile(`^/proc/([0-9]+)/stat$`)
	}

	gs := re.FindStringSubmatch(path)
	if gs == nil {
		return nil, os.ErrNotExist
	}
	pid, err := strconv.Atoi(gs[1])
	if err != nil {
		return nil, os.ErrNotExist
	}

	if strings.HasSuffix(path, "cmdline") {
		return io.NopCloser(bytes.NewBufferString(r.pidToCLIArgs[pid])), nil
	}

	name, ok := r.pidToName[pid]
	if !ok {
		return nil, os.ErrNotExist
	}
	return io.NopCloser(bytes.NewBufferString(fmt.Sprintf(
		"%d (%s) I 2 0 0 0 -1 69238880 0 0 0 0 0 0 0 0 0 -20 1 0 250 0 0 18446744073709551615 0 0 0 0 0 0 0 2147483647 0 0 0 0 17 3 0 0 0 0 0 0 0 0 0 0 0 0 0",
		pid,
		name,
	))), nil
}

func (r fakeProcessPathReader) OpenDir(ctx context.Context, path string) (scanapi.DirReader, error) {
	if path != "/proc/" {
		return scanapi.SliceToDirReader([]*apb.DirContent{}), nil
	}
	paths := []*apb.DirContent{
		&apb.DirContent{
			Name:  "cpuinfo",
			IsDir: false,
		},
		&apb.DirContent{
			Name:  "net",
			IsDir: true,
		},
	}
	// Make order deterministic
	pids := make([]int, 0, len(r.pidToName))
	for pid := range r.pidToName {
		pids = append(pids, pid)
	}
	sort.Ints(pids)
	for _, pid := range pids {
		paths = append(paths, &apb.DirContent{
			Name:  fmt.Sprintf("%d", pid),
			IsDir: true,
		})
	}
	return scanapi.SliceToDirReader(paths), nil
}

func TestProcessPath(t *testing.T) {
	testCases := []struct {
		name              string
		pidToName         map[int]string
		pidToCLIArgs      map[int]string
		fileSet           *ipb.FileSet
		expectedTraversal []*traversal
	}{
		{
			name: "process name does not exist",
			pidToName: map[int]string{
				1: "foo",
				2: "bar",
			},
			fileSet: &ipb.FileSet{
				FilePath: &ipb.FileSet_ProcessPath_{ProcessPath: &ipb.FileSet_ProcessPath{
					ProcName: "foobar",
				}},
			},
			expectedTraversal: []*traversal{},
		},
		{
			name: "process name exists once",
			pidToName: map[int]string{
				1:    "foo",
				2:    "bar",
				1337: "foobar",
			},
			fileSet: &ipb.FileSet{
				FilePath: &ipb.FileSet_ProcessPath_{ProcessPath: &ipb.FileSet_ProcessPath{
					ProcName: "foobar",
				}},
			},
			expectedTraversal: []*traversal{
				{Path: "/proc/1337", IsDir: true},
			},
		},
		{
			name: "process name exists twice",
			pidToName: map[int]string{
				1:    "foo",
				2:    "bar",
				42:   "foobar",
				1337: "foobar",
			},
			fileSet: &ipb.FileSet{
				FilePath: &ipb.FileSet_ProcessPath_{ProcessPath: &ipb.FileSet_ProcessPath{
					ProcName: "foobar",
				}},
			},
			expectedTraversal: []*traversal{
				{Path: "/proc/42", IsDir: true},
				{Path: "/proc/1337", IsDir: true},
			},
		},
		{
			name: "filename is set",
			pidToName: map[int]string{
				1:    "foo",
				2:    "bar",
				42:   "foobar",
				1337: "foobar",
			},
			fileSet: &ipb.FileSet{
				FilePath: &ipb.FileSet_ProcessPath_{ProcessPath: &ipb.FileSet_ProcessPath{
					ProcName: "foobar",
					FileName: "environ",
				}},
			},
			expectedTraversal: []*traversal{
				{Path: "/proc/42/environ", IsDir: false},
				{Path: "/proc/1337/environ", IsDir: false},
			},
		},
		{
			name: "CLI args are set",
			pidToName: map[int]string{
				1: "foo",
				2: "foo",
			},
			pidToCLIArgs: map[int]string{
				1: "env vars 1",
				2: "env vars 2",
			},
			fileSet: &ipb.FileSet{
				FilePath: &ipb.FileSet_ProcessPath_{ProcessPath: &ipb.FileSet_ProcessPath{
					ProcName:    "foo",
					CliArgRegex: ".* 2",
				}},
			},
			expectedTraversal: []*traversal{
				{Path: "/proc/2", IsDir: true},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := make([]*traversal, 0)
			err := fileset.WalkFiles(
				context.Background(),
				tc.fileSet,
				&fakeProcessPathReader{pidToName: tc.pidToName, pidToCLIArgs: tc.pidToCLIArgs},
				time.Time{},
				func(path string, isDir bool, traversingDir bool) error {
					got = append(got, &traversal{path, isDir, traversingDir})
					return nil
				},
			)
			if err != nil {
				t.Fatalf("fileset.WalkFiles(%v) returned an error: %v", tc.fileSet, err)
			}

			if diff := cmp.Diff(tc.expectedTraversal, got); diff != "" {
				t.Errorf(
					"fileset.WalkFunc(%v) made an unexpected traversal diff (-want +got):\n%s",
					tc.fileSet,
					diff,
				)
			}
		})
	}
}

func TestProcessPathRemovedAfterQuerying(t *testing.T) {
	fileSet := &ipb.FileSet{
		FilePath: &ipb.FileSet_ProcessPath_{ProcessPath: &ipb.FileSet_ProcessPath{
			ProcName: "foo",
		}},
	}
	expectedTraversal := []*traversal{}
	got := make([]*traversal, 0)
	err := fileset.WalkFiles(
		context.Background(),
		fileSet,
		&fakeProcessPathReader{pidToName: map[int]string{1: "foo"}, removeFilesAfterQuery: true},
		time.Time{},
		func(path string, isDir bool, traversingDir bool) error {
			got = append(got, &traversal{path, isDir, traversingDir})
			return nil
		},
	)
	if err != nil {
		t.Fatalf("fileset.WalkFiles(%v) returned an error: %v", fileSet, err)
	}

	if diff := cmp.Diff(expectedTraversal, got); diff != "" {
		t.Errorf(
			"fileset.WalkFunc(%v) made an unexpected traversal diff (-want +got):\n%s",
			fileSet,
			diff,
		)
	}
}

func TestUnixEnvVarPaths(t *testing.T) {
	testCases := []struct {
		name              string
		fileSet           *ipb.FileSet
		expectedTraversal []*traversal
		expectError       bool
	}{
		{
			name: "read a valid var",
			fileSet: &ipb.FileSet{
				FilePath: &ipb.FileSet_UnixEnvVarPaths_{UnixEnvVarPaths: &ipb.FileSet_UnixEnvVarPaths{
					VarName: "PATH",
				}},
			},
			expectedTraversal: []*traversal{
				{Path: "/root", IsDir: true},
				{Path: "/root/file1.txt", IsDir: false},
			},
		},
		{
			name: "read a non-existent var",
			fileSet: &ipb.FileSet{
				FilePath: &ipb.FileSet_UnixEnvVarPaths_{UnixEnvVarPaths: &ipb.FileSet_UnixEnvVarPaths{
					VarName: "NONEXISTENT",
				}},
			},
			expectError: true,
		},
		{
			name: "files only",
			fileSet: &ipb.FileSet{
				FilePath: &ipb.FileSet_UnixEnvVarPaths_{UnixEnvVarPaths: &ipb.FileSet_UnixEnvVarPaths{
					VarName:   "PATH",
					FilesOnly: true,
				}},
			},
			expectedTraversal: []*traversal{
				{Path: "/root/file1.txt", IsDir: false},
			},
		},
		{
			name: "directories only",
			fileSet: &ipb.FileSet{
				FilePath: &ipb.FileSet_UnixEnvVarPaths_{UnixEnvVarPaths: &ipb.FileSet_UnixEnvVarPaths{
					VarName:  "PATH",
					DirsOnly: true,
				}},
			},
			expectedTraversal: []*traversal{
				{Path: "/root", IsDir: true},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var got []*traversal
			err := fileset.WalkFiles(
				context.Background(),
				tc.fileSet,
				&fakeDirectoryReader{},
				time.Time{},
				func(path string, isDir bool, traversingDir bool) error {
					got = append(got, &traversal{path, isDir, traversingDir})
					return nil
				},
			)
			if tc.expectError {
				if err == nil {
					t.Errorf("fileset.WalkFiles(%v) didn't return an error", tc.fileSet)
				}
			} else {
				if err != nil {
					t.Fatalf("fileset.WalkFiles(%v) returned an error: %v", tc.fileSet, err)
				}
				if diff := cmp.Diff(tc.expectedTraversal, got); diff != "" {
					t.Errorf("fileset.WalkFiles(%v) made an unexpected traversal diff (-want +got):\n%s", tc.fileSet, diff)
				}
			}
		})
	}
}

func TestTimeout(t *testing.T) {
	testCases := []struct {
		description string
		fileSet     *ipb.FileSet
	}{
		{
			description: "single file",
			fileSet: &ipb.FileSet{
				FilePath: &ipb.FileSet_SingleFile_{SingleFile: &ipb.FileSet_SingleFile{Path: "/path/to/file"}},
			},
		},
		{
			description: "files in directory",
			fileSet: &ipb.FileSet{
				FilePath: &ipb.FileSet_FilesInDir_{FilesInDir: &ipb.FileSet_FilesInDir{
					DirPath:   "/root",
					Recursive: true,
				}},
			},
		},
		{
			description: "env var paths",
			fileSet: &ipb.FileSet{
				FilePath: &ipb.FileSet_UnixEnvVarPaths_{UnixEnvVarPaths: &ipb.FileSet_UnixEnvVarPaths{
					VarName: "PATH",
				}},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			timeout := time.Now().Add(-1 * time.Second)
			err := fileset.WalkFiles(context.Background(), tc.fileSet, &fakeDirectoryReader{}, timeout, func(walkedPath string, isDir bool, traversingDir bool) error { return nil })
			if err == nil {
				t.Fatalf("fileset.WalkFiles(%v) didn't return an error, expected one", tc.fileSet)
			}
		})
	}
}

func TestApplyReplacementConfig(t *testing.T) {
	replacements := map[string]string{"/old": "/new"}
	testCases := []struct {
		name    string
		config  *apb.ReplacementConfig
		fileSet *ipb.FileSet
		want    *ipb.FileSet
	}{
		{
			name:   "Replace prefix in SingleFile",
			config: &apb.ReplacementConfig{PathPrefixReplacements: replacements},
			fileSet: &ipb.FileSet{
				FilePath: &ipb.FileSet_SingleFile_{SingleFile: &ipb.FileSet_SingleFile{Path: "/old/path"}},
			},
			want: &ipb.FileSet{
				FilePath: &ipb.FileSet_SingleFile_{SingleFile: &ipb.FileSet_SingleFile{Path: "/new/path"}},
			},
		},
		{
			name:   "Replace prefix in FilesInDir",
			config: &apb.ReplacementConfig{PathPrefixReplacements: replacements},
			fileSet: &ipb.FileSet{
				FilePath: &ipb.FileSet_FilesInDir_{FilesInDir: &ipb.FileSet_FilesInDir{DirPath: "/old/path"}},
			},
			want: &ipb.FileSet{
				FilePath: &ipb.FileSet_FilesInDir_{FilesInDir: &ipb.FileSet_FilesInDir{DirPath: "/new/path"}},
			},
		},
		{
			name:   "Replacements with leading slashes",
			config: &apb.ReplacementConfig{PathPrefixReplacements: map[string]string{"/old/": "/new/"}},
			fileSet: &ipb.FileSet{
				FilePath: &ipb.FileSet_SingleFile_{SingleFile: &ipb.FileSet_SingleFile{Path: "/old/path"}},
			},
			want: &ipb.FileSet{
				FilePath: &ipb.FileSet_SingleFile_{SingleFile: &ipb.FileSet_SingleFile{Path: "/new/path"}},
			},
		},
		{
			name:   "Don't replace if prefix not found",
			config: &apb.ReplacementConfig{PathPrefixReplacements: replacements},
			fileSet: &ipb.FileSet{
				FilePath: &ipb.FileSet_SingleFile_{SingleFile: &ipb.FileSet_SingleFile{Path: "/some/path"}},
			},
			want: &ipb.FileSet{
				FilePath: &ipb.FileSet_SingleFile_{SingleFile: &ipb.FileSet_SingleFile{Path: "/some/path"}},
			},
		},
		{
			name:   "Don't replace if not prefix",
			config: &apb.ReplacementConfig{PathPrefixReplacements: replacements},
			fileSet: &ipb.FileSet{
				FilePath: &ipb.FileSet_SingleFile_{SingleFile: &ipb.FileSet_SingleFile{Path: "/some/old/path"}},
			},
			want: &ipb.FileSet{
				FilePath: &ipb.FileSet_SingleFile_{SingleFile: &ipb.FileSet_SingleFile{Path: "/some/old/path"}},
			},
		},
		{
			name:   "Don't replace if only part of the file matches",
			config: &apb.ReplacementConfig{PathPrefixReplacements: replacements},
			fileSet: &ipb.FileSet{
				FilePath: &ipb.FileSet_SingleFile_{SingleFile: &ipb.FileSet_SingleFile{Path: "/old-path/file"}},
			},
			want: &ipb.FileSet{
				FilePath: &ipb.FileSet_SingleFile_{SingleFile: &ipb.FileSet_SingleFile{Path: "/old-path/file"}},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := proto.Clone(tc.fileSet).(*ipb.FileSet)
			fileset.ApplyReplacementConfig(got, tc.config)
			if diff := cmp.Diff(tc.want, got, protocmp.Transform()); diff != "" {
				t.Errorf(
					"fileset.ApplyReplacementConfig(%v, %v) unexpected results, diff (-want +got):\n%s", tc.fileSet, tc.config, diff,
				)
			}
		})
	}
}
