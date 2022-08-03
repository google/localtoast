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

// The cis_scanner command wraps around the scanner library to create a standalone
// CLI for the scanner with direct access to the local machine's filesystem.
package main

import (
	"context"
	"database/sql"
	"io"
	"log"
	"os"
	"path"

	"github.com/google/localtoast/localfilereader"
	"github.com/google/localtoast/scannercommon"
	apb "github.com/google/localtoast/scannerlib/proto/api_go_proto"
	"github.com/google/localtoast/sqlquerier"

	// We need this import to call sql.Open with the "mysql" driver.
	_ "github.com/go-sql-driver/mysql"
)

// localScanAPIProvider provides access to the local filesystem and to the
// local SQL database for the scanning library.
type localScanAPIProvider struct {
	chrootPath string
	db         *sql.DB
}

func (a *localScanAPIProvider) fullPath(entryPath string) string {
	if a.chrootPath == "" {
		return entryPath
	}
	return path.Join(a.chrootPath, entryPath)
}

func (a *localScanAPIProvider) OpenFile(ctx context.Context, filePath string) (io.ReadCloser, error) {
	return localfilereader.OpenFile(ctx, a.fullPath(filePath))
}

func (a *localScanAPIProvider) FilesInDir(ctx context.Context, dirPath string) ([]*apb.DirContent, error) {
	return localfilereader.FilesInDir(ctx, a.fullPath(dirPath))
}

func (a *localScanAPIProvider) FilePermissions(ctx context.Context, filePath string) (*apb.PosixPermissions, error) {
	return localfilereader.FilePermissions(ctx, a.fullPath(filePath))
}

func (a *localScanAPIProvider) SQLQuery(ctx context.Context, query string) (int, error) {
	return sqlquerier.Query(ctx, a.db, query)
}

func main() {
	flags := scannercommon.ParseFlags()

	var db *sql.DB
	if flags.Database != "" {
		var err error
		// We assume that the database is MySQL-compatible.
		db, err = sql.Open("mysql", flags.Database)
		if err != nil {
			log.Fatalf("Error connecting to the database: %v\n", err)
		}
		defer db.Close()
	}
	provider := &localScanAPIProvider{
		chrootPath: flags.ChrootPath,
		db:         db,
	}
	os.Exit(scannercommon.RunScan(flags, provider))
}
