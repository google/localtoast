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

// The localtoast command wraps around the scanner library to create a standalone
// CLI for the config scanner with direct access to the local machine's filesystem.
package main

import (
	"context"
	"errors"
	"io"
	"os"
	"path"
	"runtime/debug"

	"github.com/google/localtoast/localfilereader"
	"github.com/google/localtoast/scanapi"
	"github.com/google/localtoast/scannercommon"
	apb "github.com/google/localtoast/scannerlib/proto/api_go_proto"
	ipb "github.com/google/localtoast/scannerlib/proto/scan_instructions_go_proto"
)

// localScanAPIProvider provides access to the local filesystem.
type localScanAPIProvider struct {
	chrootPath string
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

func (localScanAPIProvider) SQLQuery(ctx context.Context, query string) (string, error) {
	// This is intentionally not implemented for the scanner version without SQL.
	return "", errors.New("not implemented")
}

func (localScanAPIProvider) SupportedDatabase() (ipb.SQLCheck_SQLDatabase, error) {
	// This is intentionally not implemented for the scanner version without SQL.
	return ipb.SQLCheck_DB_UNSPECIFIED, errors.New("not implemented")
}

func main() {
	// Change GCPercent to lower the peak memory usage.
	// Make sure we are not overwriting a custom value. We only want to change the default.
	if os.Getenv("GOGC") == "" {
		debug.SetGCPercent(1)
	}
	flags := scannercommon.ParseFlags()
	provider := &localScanAPIProvider{chrootPath: flags.ChrootPath}
	os.Exit(scannercommon.RunScan(flags, provider))
}
