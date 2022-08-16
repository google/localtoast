// Copyright 2022 Google LLC
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

// Package scanapi defines the Localtoast scan API
// used to provide access to a local or remote filesystem and database to perform scans.
package scanapi

import (
	"context"
	"io"

	apb "github.com/google/localtoast/scannerlib/proto/api_go_proto"
)

// Filesystem is an interface that gives read access to the filesystem of the machine to scan.
type Filesystem interface {
	// OpenFile opens the specified file for reading.
	// It should return an os.IsNotExist error if the file doesn't exist.
	OpenFile(ctx context.Context, path string) (io.ReadCloser, error)
	// FilePermissions returns unix permission-related data for the specified file or directory.
	FilePermissions(ctx context.Context, path string) (*apb.PosixPermissions, error)
	// FilesInDir lists the contents of the specified directory.
	FilesInDir(ctx context.Context, path string) ([]*apb.DirContent, error)
}

// SQLQuerier is an interface that supports SQL queries to a target SQL database.
type SQLQuerier interface {
	// SQLQuery executes SQL queries to a target SQL database and returns the number of result rows.
	SQLQuery(ctx context.Context, query string) (int, error)
}

// ScanAPI is an interface that gives read access to the filesystem of
// the machine to scan and can execute SQL queries on a single database.
type ScanAPI interface {
	Filesystem
	SQLQuerier
}
