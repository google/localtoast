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

// Package sqlquerier provides an utility function for running SQL queries.
package sqlquerier

import (
	"context"
	"database/sql"
	"errors"
)

// Query executes a SQL query and returns the number of rows in the result.
func Query(ctx context.Context, db *sql.DB, query string) (int, error) {
	if db == nil {
		return 0, errors.New("no database specified. Please provide one using the --database flag")
	}
	rows, err := db.Query(query)
	if err != nil {
		return 0, err
	}
	defer rows.Close()
	n := 0
	for rows.Next() {
		n++
	}
	if rows.Err() != nil {
		return 0, rows.Err()
	}
	return n, nil
}
