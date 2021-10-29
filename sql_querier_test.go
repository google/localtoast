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

package sqlquerier_test

import (
	"context"
	"strings"
	"testing"

	"github.com/google/localtoast/sqlquerier"
	"github.com/google/localtoast/testing/fakedb"
)

func TestSQLCheckWithNoDatabaseFlag(t *testing.T) {
	q := "query"
	_, err := sqlquerier.Query(context.Background(), nil, q)
	if err == nil {
		t.Errorf("sqlquerier.Query(context.Background(), nil, %q) expected to return an error but got none", q)
	}
}

func TestSQLCheck(t *testing.T) {
	testCases := []struct {
		desc        string
		query       string
		want        int
		expectError bool
		errorMsg    string
	}{
		{
			desc:  "SQLCheck one row returned",
			query: fakedb.QueryOneRow,
			want:  1,
		},
		{
			desc:  "SQLCheck no rows returned",
			query: fakedb.QueryNoRows,
			want:  0,
		},
		{
			desc:        "SQLCheck propagates errors",
			query:       fakedb.QueryError,
			want:        0,
			expectError: true,
			errorMsg:    fakedb.ErrorMsg,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			db, err := fakedb.Open(&fakedb.FakeDB{})
			if err != nil {
				t.Errorf("fakedb.Open had an unexpected error: %v", err)
			}
			var got int

			got, err = sqlquerier.Query(context.Background(), db, tc.query)
			if err != nil {
				if !tc.expectError {
					t.Errorf("sqlquerier.Query(ctx, db, %q) had an unexpected error: %v", tc.query, err)
				}
				if !strings.Contains(err.Error(), tc.errorMsg) {
					t.Errorf("sqlquerier.Query(ctx, db, %q) returned the wrong error: want %q, got %v", tc.query, tc.errorMsg, err)
				}
			}
			if got != tc.want {
				t.Errorf("sqlquerier.Query(ctx, db, %q) returned wrong result: want %d rows, got %d rows", tc.query, tc.want, got)
			}
		})
	}
}
