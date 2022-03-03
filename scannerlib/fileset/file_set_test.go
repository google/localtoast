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
	"io/ioutil"
	"os"
	"regexp"
	"sort"
	"strconv"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/localtoast/scannerlib/fileset"
	apb "github.com/google/localtoast/scannerlib/proto/api_go_proto"
	ipb "github.com/google/localtoast/scannerlib/proto/scan_instructions_go_proto"
)

const (
	procEnvironPath    = "/proc/self/environ"
	procEnvironContent = "SHELL=/bin/bash\x00PATH=/root:/root/file1.txt"
)

type fakeDirectoryReader struct{}

func (fakeDirectoryReader) OpenFile(ctx context.Context, path string) (io.ReadCloser, error) {
	switch path {
	case procEnvironPath:
		return ioutil.NopCloser(bytes.NewReader([]byte(procEnvironContent))), nil
	default:
		return nil, os.ErrNotExist
	}
}

func (fakeDirectoryReader) FilesInDir(ctx context.Context, path string) ([]*apb.DirContent, error) {
	// Fake dir structure:
	// root---file1.txt
	//  \ \ \-file2.gif
	//   \ \--symlink
	//    \---subdir--file3.txt
	switch path {
	case "/root":
		return []*apb.DirContent{
			&apb.DirContent{Name: "file1.txt", IsDir: false},
			&apb.DirContent{Name: "file2.gif", IsDir: false},
			&apb.DirContent{Name: "symlink", IsDir: false, IsSymlink: true},
			&apb.DirContent{Name: "subdir", IsDir: true},
		}, nil
	case "/root/subdir":
		return []*apb.DirContent{
			&apb.DirContent{Name: "file3.txt", IsDir: false},
		}, nil
	default:
		return nil, os.ErrNotExist
	}
}

func TestSingleFile(t *testing.T) {
	expectedPath := "/path/to/file"
	fileSet := &ipb.FileSet{
		FilePath: &ipb.FileSet_SingleFile_{SingleFile: &ipb.FileSet_SingleFile{Path: expectedPath}},
	}

	err := fileset.WalkFiles(context.Background(), fileSet, &fakeDirectoryReader{}, func(walkedPath string, isDir bool) error {
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
	Path  string
	IsDir bool
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
				{Path: "/root", IsDir: true},
				{Path: "/root/file1.txt", IsDir: false},
				{Path: "/root/file2.gif", IsDir: false},
				{Path: "/root/symlink", IsDir: false},
				{Path: "/root/subdir", IsDir: true},
				{Path: "/root/subdir/file3.txt", IsDir: false},
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
				{Path: "/root/file1.txt", IsDir: false},
				{Path: "/root/file2.gif", IsDir: false},
				{Path: "/root/symlink", IsDir: false},
				{Path: "/root/subdir/file3.txt", IsDir: false},
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
				{Path: "/root", IsDir: true},
				{Path: "/root/subdir", IsDir: true},
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
				{Path: "/root", IsDir: true},
				{Path: "/root/file1.txt", IsDir: false},
				{Path: "/root/file2.gif", IsDir: false},
				{Path: "/root/subdir", IsDir: true},
				{Path: "/root/subdir/file3.txt", IsDir: false},
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
				{Path: "/root", IsDir: true},
				{Path: "/root/file1.txt", IsDir: false},
				{Path: "/root/file2.gif", IsDir: false},
				{Path: "/root/symlink", IsDir: false},
				{Path: "/root/subdir", IsDir: true},
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
				{Path: "/root/file1.txt", IsDir: false},
				{Path: "/root/subdir/file3.txt", IsDir: false},
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
				{Path: "/root", IsDir: true},
				{Path: "/root/file1.txt", IsDir: false},
				{Path: "/root/file2.gif", IsDir: false},
				{Path: "/root/symlink", IsDir: false},
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
				{Path: "/non-existent", IsDir: true},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			gotTraversal := []*traversal{}
			err := fileset.WalkFiles(context.Background(), tc.fileSet, &fakeDirectoryReader{}, func(walkedPath string, isDir bool) error {
				gotTraversal = append(gotTraversal, &traversal{Path: walkedPath, IsDir: isDir})
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

func (infiniteLoopFSReader) FilesInDir(ctx context.Context, path string) ([]*apb.DirContent, error) {
	return []*apb.DirContent{
		&apb.DirContent{Name: "dir", IsDir: true},
	}, nil
}

func TestTraverseFilesystemWithInfiniteLoop(t *testing.T) {
	files := &ipb.FileSet{FilePath: &ipb.FileSet_FilesInDir_{FilesInDir: &ipb.FileSet_FilesInDir{
		DirPath:   "/",
		Recursive: true,
	}}}
	err := fileset.WalkFiles(context.Background(), files, &infiniteLoopFSReader{}, func(walkedPath string, isDir bool) error { return nil })
	if err == nil {
		t.Fatalf("fileset.WalkFiles(%v) didn't return an error", files)
	}
}

type fakeProcessPathReader struct {
	pidToName map[int]string
	// If true, tests the race condition when the files become unavailable after
	// they're listed in the /proc directory.
	removeFilesAfterQuery bool
}

func (r fakeProcessPathReader) OpenFile(ctx context.Context, path string) (io.ReadCloser, error) {
	if r.removeFilesAfterQuery {
		return nil, os.ErrNotExist
	}

	re := regexp.MustCompile(`^/proc/([0-9]+)/stat$`)

	gs := re.FindStringSubmatch(path)
	if gs == nil {
		return nil, os.ErrNotExist
	}
	pid, err := strconv.Atoi(gs[1])
	if err != nil {
		return nil, os.ErrNotExist
	}
	name, ok := r.pidToName[pid]
	if !ok {
		return nil, os.ErrNotExist
	}
	return ioutil.NopCloser(bytes.NewBufferString(fmt.Sprintf(
		"%d (%s) I 2 0 0 0 -1 69238880 0 0 0 0 0 0 0 0 0 -20 1 0 250 0 0 18446744073709551615 0 0 0 0 0 0 0 2147483647 0 0 0 0 17 3 0 0 0 0 0 0 0 0 0 0 0 0 0",
		pid,
		name,
	))), nil
}

func (r fakeProcessPathReader) FilesInDir(ctx context.Context, path string) ([]*apb.DirContent, error) {
	if path != "/proc/" {
		return []*apb.DirContent{}, nil
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
	return paths, nil
}

func TestProcessPath(t *testing.T) {
	testCases := []struct {
		name              string
		pidToName         map[int]string
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
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := make([]*traversal, 0)
			err := fileset.WalkFiles(
				context.Background(),
				tc.fileSet,
				&fakeProcessPathReader{pidToName: tc.pidToName},
				func(path string, isDir bool) error {
					got = append(got, &traversal{path, isDir})
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
		func(path string, isDir bool) error {
			got = append(got, &traversal{path, isDir})
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
				func(path string, isDir bool) error {
					got = append(got, &traversal{path, isDir})
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
