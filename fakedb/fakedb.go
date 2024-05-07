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

// Package fakedb provides a minimal fake implementation for a database/sql database,
// to be used in tests.
package fakedb

import (
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"
	"io"
)

// Queries supported by the fake database. Each query will return a different number
// of rows, but no actual content.
const (
	QueryNoRows = "SELECT 1 WHERE FALSE"
	QueryOneRow = "SELECT 1"
	QueryError  = "INVALID QUERY"
	ErrorMsg    = "invalid query"
)

func init() {
	sql.Register("fakedb", &fakeDriver{})
}

// FakeDB is a fake implementation of sql.DB.
type FakeDB struct{}

// Open creates a new sql.DB object that supports the queries defined in fakedb.
func Open(*FakeDB) (*sql.DB, error) {
	return sql.Open("fakedb", "fakedsn")
}

// fakeDriver is a fake implementation of sql.Driver.
type fakeDriver struct{}

// Open is a fake implementation for the sql.Driver interface.
func (d *fakeDriver) Open(name string) (driver.Conn, error) {
	return &fakeConn{&FakeDB{}}, nil
}

// fakeConn is a fake implementation of sql.Conn.
type fakeConn struct {
	db *FakeDB
}

// Close is a fake implementation for the sql.Conn interface.
func (*fakeConn) Close() error { return nil }

// Begin is a fake implementation for the sql.Conn interface.
func (*fakeConn) Begin() (driver.Tx, error) { return nil, nil }

// Prepare is a fake implementation for the sql.Conn interface.
func (c *fakeConn) Prepare(query string) (driver.Stmt, error) {
	return &fakeStmt{db: c.db, query: query}, nil
}

// fakeStmt is a fake implementation of driver.Stmt.
type fakeStmt struct {
	db    *FakeDB
	query string
}

// Close is a fake implementation for the driver.Stmt interface.
func (*fakeStmt) Close() error { return nil }

// Exec is a fake implementation for the driver.Stmt interface.
func (*fakeStmt) Exec(args []driver.Value) (driver.Result, error) { return nil, nil }

// NumInput is a fake implementation for the driver.Stmt interface.
func (*fakeStmt) NumInput() int { return 0 }

// Query supports a few selected queries defined in the fakedb package.
func (s *fakeStmt) Query(args []driver.Value) (driver.Rows, error) {
	switch s.query {
	case QueryNoRows:
		return &fakeRows{rowsLeft: 0}, nil
	case QueryOneRow:
		return &fakeRows{rowsLeft: 1, columns: []string{"fakeColumn"}}, nil
	case QueryError:
		return nil, errors.New(ErrorMsg)
	default:
		return nil, fmt.Errorf("the query %q is not supported by fakeStmt", s.query)
	}
}

// fakeRows is a fake implementation of driver.Rows.
// We only care about the number of returned rows, so we don't model their content.
type fakeRows struct {
	rowsLeft int
	columns  []string
}

// Close is a fake implementation for the driver.Rows interface.
func (*fakeRows) Close() error { return nil }

// Columns is a fake implementation for the driver.Rows interface.
func (r *fakeRows) Columns() []string { return r.columns }

// Next is a fake implementation for the driver.Rows interface.
// The function iterates through the number of known rows, but returns no values.
func (r *fakeRows) Next(dst []driver.Value) error {
	if r.rowsLeft > 0 {
		r.rowsLeft--
		dst[0] = "fakeValue"
		return nil
	}
	return io.EOF
}
