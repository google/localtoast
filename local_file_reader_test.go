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

package localfilereader_test

import (
	"context"
	"fmt"
	"io"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"google.golang.org/protobuf/testing/protocmp"
	apb "github.com/google/localtoast/library/proto/api_go_proto"
	"github.com/google/localtoast/localfilereader"
)

var testDirPath string

const (
	fileName        = "file"
	fileContent     = "FILE_CONTENT"
	filePermission  = 0644
	dirName         = "dir"
	fileSymlinkName = "file_symlink"
	dirSymlinkName  = "dir_symlink"
)

func TestMain(m *testing.M) {
	// Create some temporary files before running the tests.
	testDirPath = os.Getenv("TEST_TMPDIR")
	if err := ioutil.WriteFile(filepath.Join(testDirPath, fileName), []byte(fileContent), filePermission); err != nil {
		panic(fmt.Sprintf("Error while creating file %s: %v", filepath.Join(testDirPath, fileName), err))
	}
	if err := os.Mkdir(filepath.Join(testDirPath, dirName), filePermission); err != nil {
		panic(fmt.Sprintf("Error while creating directory %s: %v", filepath.Join(testDirPath, dirName), err))
	}
	if err := os.Symlink(filepath.Join(testDirPath, fileName), filepath.Join(testDirPath, fileSymlinkName)); err != nil {
		panic(fmt.Sprintf("Error while creating symlink %s: %v", filepath.Join(testDirPath, fileSymlinkName), err))
	}
	if err := os.Symlink(filepath.Join(testDirPath, dirName), filepath.Join(testDirPath, dirSymlinkName)); err != nil {
		panic(fmt.Sprintf("Error while creating symlink %s: %v", filepath.Join(testDirPath, dirSymlinkName), err))
	}

	exitCode := m.Run()
	os.Exit(exitCode)
}

func TestOpenFile(t *testing.T) {
	testFilePath := filepath.Join(testDirPath, fileName)
	reader, err := localfilereader.OpenFile(context.Background(), testFilePath)
	if err != nil {
		t.Fatalf("localfilereader.OpenFile(%s) had unexpected error: %v", testFilePath, err)
	}
	content, err := io.ReadAll(reader)
	if err != nil {
		t.Fatalf("Reading file %s had unexpected error: %v", testFilePath, err)
	}
	if string(content) != fileContent {
		t.Errorf("Got file content %q, expected %q", content, fileContent)
	}
}

func TestOpenPropagatesError(t *testing.T) {
	nonExistentFilePath := filepath.Join(testDirPath, "non-existent-file")
	_, err := localfilereader.OpenFile(context.Background(), nonExistentFilePath)
	if err == nil {
		t.Errorf("localfilereader.OpenFile(%s) didn't return an error", nonExistentFilePath)
	}
}

func TestFilesInDir(t *testing.T) {
	files, err := localfilereader.FilesInDir(context.Background(), testDirPath)
	if err != nil {
		t.Fatalf("localfilereader.FilesInDir(%s) had unexpected error: %v", testDirPath, err)
	}
	expected := []*apb.DirContent{
		&apb.DirContent{Name: fileName, IsDir: false, IsSymlink: false},
		&apb.DirContent{Name: dirName, IsDir: true, IsSymlink: false},
		&apb.DirContent{Name: fileSymlinkName, IsDir: false, IsSymlink: true},
		&apb.DirContent{Name: dirSymlinkName, IsDir: false, IsSymlink: true},
	}
	sortProtosOpt := cmpopts.SortSlices(func(c1, c2 *apb.DirContent) bool {
		return c1.String() < c2.String()
	})
	if diff := cmp.Diff(expected, files, protocmp.Transform(), sortProtosOpt); diff != "" {
		t.Errorf("localfilereader.FilesInDir(%s) returned unexpected diff (-want +got):\n%s", testDirPath, diff)
	}
}

func TestFilesInDirPropagatesError(t *testing.T) {
	nonExistentDirPath := filepath.Join(testDirPath, "non-existent-dir")
	_, err := localfilereader.FilesInDir(context.Background(), nonExistentDirPath)
	if err == nil {
		t.Errorf("localfilereader.FilesInDir(%s) didn't return an error", nonExistentDirPath)
	}
}

func TestFilePermissionsCorrectPermissionNumbers(t *testing.T) {
	testFilePath := filepath.Join(testDirPath, fileName)
	permission, err := localfilereader.FilePermissions(context.Background(), testFilePath)
	if err != nil {
		t.Fatalf("localfilereader.FilePermissions(%s) had unexpected error: %v", testFilePath, err)
	}
	if permission.GetPermissionNum() != filePermission {
		t.Fatalf("localfilereader.FilePermissions(%s) returned %o, expected %o",
			testFilePath, permission.GetPermissionNum(), filePermission)
	}
}

func TestFilePermissionsCorrectSpecialFlags(t *testing.T) {
	testFilePath := filepath.Join(testDirPath, fileName)
	testCases := []struct {
		flagToAdd       fs.FileMode
		expectedSetFlag int32
	}{
		{flagToAdd: fs.ModeSetuid, expectedSetFlag: localfilereader.SetuidFlag},
		{flagToAdd: fs.ModeSetgid, expectedSetFlag: localfilereader.SetgidFlag},
		{flagToAdd: fs.ModeSticky, expectedSetFlag: localfilereader.StickyFlag},
	}
	for _, tc := range testCases {
		os.Chmod(testFilePath, fs.FileMode(filePermission)|tc.flagToAdd)
		permission, err := localfilereader.FilePermissions(context.Background(), testFilePath)
		if err != nil {
			t.Fatalf("localfilereader.FilePermissions(%s) had unexpected error: %v", testFilePath, err)
		}
		if permission.GetPermissionNum()&tc.expectedSetFlag == 0 {
			t.Fatalf("localfilereader.FilePermissions(%s) returned %o, expected %o flag to be set",
				testFilePath, permission.GetPermissionNum(), tc.expectedSetFlag)
		}
	}
}

func TestFilePermissionsPropagatesError(t *testing.T) {
	nonExistentFilePath := filepath.Join(testDirPath, "non-existent-file")
	_, err := localfilereader.FilePermissions(context.Background(), nonExistentFilePath)
	if err == nil {
		t.Errorf("localfilereader.FilesInDir(%s) didn't return an error", nonExistentFilePath)
	}
}
