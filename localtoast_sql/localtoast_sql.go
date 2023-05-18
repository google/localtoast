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
	"errors"
	"io"
	"log"
	"os"
	"path"

	"github.com/google/localtoast/cqlquerier"
	"github.com/google/localtoast/localfilereader"
	"github.com/google/localtoast/scanapi"
	"github.com/google/localtoast/scannercommon"
	apb "github.com/google/localtoast/scannerlib/proto/api_go_proto"
	ipb "github.com/google/localtoast/scannerlib/proto/scan_instructions_go_proto"
	"github.com/google/localtoast/sqlquerier"

	// Import Cassandra connector
	"github.com/gocql/gocql"

	// We need this import to call sql.Open with the "mysql" driver.
	_ "github.com/go-sql-driver/mysql"
)

// localScanAPIProvider provides access to the local filesystem and to the
// local SQL database for the scanning library.
type localScanAPIProvider struct {
	chrootPath     string
	sqldb          *sql.DB
	cqldb          *gocql.Session

	dbtype         ipb.SQLCheck_SQLDatabase
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

func (a *localScanAPIProvider) OpenDir(ctx context.Context, path string) (scanapi.DirReader, error) {
	return localfilereader.OpenDir(ctx, a.fullPath(path))
}

func (a *localScanAPIProvider) FilePermissions(ctx context.Context, filePath string) (*apb.PosixPermissions, error) {
	return localfilereader.FilePermissions(ctx, a.fullPath(filePath))
}

func (a *localScanAPIProvider) SupportedDatabase() (ipb.SQLCheck_SQLDatabase, error) {
	if a.dbtype == ipb.SQLCheck_DB_UNSPECIFIED {
		return a.dbtype, errors.New("no database specified.")
	}
	return a.dbtype, nil
}

func (a *localScanAPIProvider) SQLQuery(ctx context.Context, query string) (int, error) {
	dbtype, err := a.SupportedDatabase()
	if err != nil {
		return 0, err
	}
	if dbtype == ipb.SQLCheck_DB_MYSQL {
		return sqlquerier.Query(ctx, a.sqldb, query)
	}
	if dbtype == ipb.SQLCheck_DB_CASSANDRA {
		return cqlquerier.Query(ctx, a.cqldb, query)
	}
	return 0, errors.New("no database specified. Please provide one using --database or --cassandra flag")
}

func main() {
	flags := scannercommon.ParseFlags()

	var sqldb *sql.DB
	var cqldb *gocql.Session
	var err error

	dbtype := ipb.SQLCheck_DB_UNSPECIFIED

	if flags.Database != "" {
		// We assume that the database is MySQL-compatible.
		sqldb, err = sql.Open("mysql", flags.Database)
		if err != nil {
			log.Fatalf("Error connecting to the database: %v\n", err)
		}
		defer sqldb.Close()

		dbtype = ipb.SQLCheck_DB_MYSQL
	} else if flags.Cassandra != "" {
		cluster := gocql.NewCluster(flags.Cassandra)

		// connect to the cluster
		cqldb, err = cluster.CreateSession()
		if err != nil {
			log.Fatalf("Error connecting to Cassandra: %v\n", err)
		}
		defer cqldb.Close()

		dbtype = ipb.SQLCheck_DB_CASSANDRA
	}
	provider := &localScanAPIProvider{
		chrootPath: flags.ChrootPath,
		sqldb:      sqldb,
		cqldb:      cqldb,
		dbtype:     dbtype,
	}
	os.Exit(scannercommon.RunScan(flags, provider))
}
