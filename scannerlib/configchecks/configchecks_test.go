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

// Package configchecks_tests contains helper functions and tests for the configchecks package.
package configchecks_test

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"time"

	"github.com/google/localtoast/scanapi"
	"github.com/google/localtoast/scannerlib/configchecks"
	apb "github.com/google/localtoast/scannerlib/proto/api_go_proto"
	ipb "github.com/google/localtoast/scannerlib/proto/scan_instructions_go_proto"
)

const (
	testDirPath         = "/dir"
	testFilePath        = "/dir/file"
	pipelineFileToken   = "%%pipeline%%"
	testFileContent     = "File content"
	emptyTestFilePath   = "/dir/empty_file"
	unreadableFilePath  = "/path/to/unreadable/file"
	nonExistentFilePath = "/path/to/non/existent/file"
	fakeQueryNoRows     = "SELECT 1 WHERE FALSE"
	fakeQueryOneRow     = "SELECT 1"
	fakeQueryError      = "INVALID QUERY"
	queryErrorMsg       = "invalid query"
)

var (
	today = int32(time.Now().Sub(time.Unix(0, 0)).Hours() / 24)
)

type fakeAPI struct {
	fileContent  string
	openFileFunc func(ctx context.Context, filePath string) (io.ReadCloser, error)
	supportedDB  ipb.SQLCheck_SQLDatabase
}

type fakeAPIOpt func(r *fakeAPI)

func withFileContent(content string) fakeAPIOpt {
	return func(r *fakeAPI) {
		r.fileContent = content
	}
}

func withOpenFileFunc(f func(ctx context.Context, filePath string) (io.ReadCloser, error)) fakeAPIOpt {
	return func(r *fakeAPI) {
		r.openFileFunc = f
	}
}

func withSupportedDatabase(db ipb.SQLCheck_SQLDatabase) fakeAPIOpt {
	return func(r *fakeAPI) {
		r.supportedDB = db
	}
}

func newFakeAPI(opts ...fakeAPIOpt) *fakeAPI {
	r := &fakeAPI{
		fileContent:  testFileContent,
		openFileFunc: nil,
		supportedDB:  ipb.SQLCheck_DB_MYSQL,
	}
	for _, opt := range opts {
		opt(r)
	}
	return r
}

func (r *fakeAPI) OpenFile(ctx context.Context, filePath string) (io.ReadCloser, error) {
	if r.openFileFunc != nil {
		return r.openFileFunc(ctx, filePath)
	}
	switch filePath {
	case emptyTestFilePath:
		return io.NopCloser(bytes.NewReader([]byte{})), nil
	case unreadableFilePath:
		return nil, errors.New("io error")
	case nonExistentFilePath:
		return nil, os.ErrNotExist
	default:
		return io.NopCloser(bytes.NewReader([]byte(r.fileContent))), nil
	}
}

func (fakeAPI) OpenDir(ctx context.Context, filePath string) (scanapi.DirReader, error) {
	switch filePath {
	case testDirPath:
		return scanapi.SliceToDirReader([]*apb.DirContent{
			{Name: path.Base(emptyTestFilePath), IsDir: false},
			{Name: path.Base(testFilePath), IsDir: false},
		}), nil
	case nonExistentFilePath:
		return nil, os.ErrNotExist
	default:
		return nil, errors.New("not a directory")
	}
}

func (fakeAPI) FilePermissions(ctx context.Context, filePath string) (*apb.PosixPermissions, error) {
	switch filePath {
	case nonExistentFilePath:
		return nil, os.ErrNotExist
	default:
		return &apb.PosixPermissions{
			PermissionNum: 0644,
			Uid:           0,
			User:          "root",
			Gid:           0,
			Group:         "root",
		}, nil
	}
}

func (fakeAPI) SQLQuery(ctx context.Context, query string) (string, error) {
	switch query {
	case fakeQueryNoRows:
		return "", nil
	case fakeQueryOneRow:
		return "testValue", nil
	case fakeQueryError:
		return "", errors.New(queryErrorMsg)
	default:
		return "", fmt.Errorf("the query %q is not supported by fakeAPI", query)
	}
}

func (r *fakeAPI) SupportedDatabase() (ipb.SQLCheck_SQLDatabase, error) {
	return r.supportedDB, nil
}

func singleComplianceResult(m configchecks.ComplianceMap) (result *apb.ComplianceResult, gotSingleton bool) {
	results := []*apb.ComplianceResult{}
	for _, v := range m {
		results = append(results, v)
	}
	if len(results) != 1 {
		return nil, false
	}
	return results[0], true
}
