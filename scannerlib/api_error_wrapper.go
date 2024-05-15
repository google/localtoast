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

package scannerlib

import (
	"context"
	"fmt"
	"io"

	"github.com/google/localtoast/scanapi"
	apb "github.com/google/localtoast/scannerlib/proto/api_go_proto"
	ipb "github.com/google/localtoast/scannerlib/proto/scan_instructions_go_proto"
)

// apiErrorWrapper is a wrapper around the scanner's scan API that makes the error
// messages more verbose.
type apiErrorWrapper struct {
	api scanapi.ScanAPI
}

func (w *apiErrorWrapper) OpenFile(ctx context.Context, path string) (io.ReadCloser, error) {
	rc, err := w.api.OpenFile(ctx, path)
	if err != nil {
		err = fmt.Errorf("api.OpenFile(%q): %w", path, err)
	}
	return rc, err
}

func (w *apiErrorWrapper) OpenDir(ctx context.Context, path string) (scanapi.DirReader, error) {
	d, err := w.api.OpenDir(ctx, path)
	if err != nil {
		err = fmt.Errorf("api.OpenDir(%q): %w", path, err)
	}
	return d, err
}

func (w *apiErrorWrapper) FilePermissions(ctx context.Context, path string) (*apb.PosixPermissions, error) {
	p, err := w.api.FilePermissions(ctx, path)
	if err != nil {
		err = fmt.Errorf("api.FilePermissions(%q): %w", path, err)
	}
	return p, err
}

func (w *apiErrorWrapper) SQLQuery(ctx context.Context, query string) (string, error) {
	res, err := w.api.SQLQuery(ctx, query)
	if err != nil {
		err = fmt.Errorf("api.SQLQuery(%q): %w", query, err)
	}
	return res, err
}

func (w *apiErrorWrapper) SupportedDatabase() (ipb.SQLCheck_SQLDatabase, error) {
	l, err := w.api.SupportedDatabase()
	if err != nil {
		err = fmt.Errorf("api.SupportedDatabase(): %w", err)
	}
	return l, err
}
