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
	"path"

	"github.com/google/localtoast/localfilereader"
	"github.com/google/localtoast/scannercommon"
	apb "github.com/google/localtoast/scannerlib/proto/api_go_proto"
)

// localScanAPIProvider provides access to the local filesystem.
type localScanAPIProvider struct {
	chrootPath string
}

func (a *localScanAPIProvider) OpenFile(ctx context.Context, filePath string) (io.ReadCloser, error) {
	return localfilereader.OpenFile(ctx, path.Join(a.chrootPath, filePath))
}

func (a *localScanAPIProvider) FilesInDir(ctx context.Context, dirPath string) ([]*apb.DirContent, error) {
	return localfilereader.FilesInDir(ctx, path.Join(a.chrootPath, dirPath))
}

func (a *localScanAPIProvider) FilePermissions(ctx context.Context, filePath string) (*apb.PosixPermissions, error) {
	return localfilereader.FilePermissions(ctx, path.Join(a.chrootPath, filePath))
}

func (localScanAPIProvider) SQLQuery(ctx context.Context, query string) (int, error) {
	// This is intentionally not implemented for the scanner version without SQL.
	return 0, errors.New("not implemented")
}

func main() {
	flags := scannercommon.ParseFlags()
	provider := &localScanAPIProvider{chrootPath: flags.ChrootPath}
	scannercommon.RunScan(flags, provider)
}
