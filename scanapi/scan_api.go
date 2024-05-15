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
	"errors"
	"io"

	apb "github.com/google/localtoast/scannerlib/proto/api_go_proto"
	ipb "github.com/google/localtoast/scannerlib/proto/scan_instructions_go_proto"
)

var (
	// ErrEntryBeforeNext is the error returned if Entry() is called
	// and Next() was never called for a DirReader.
	ErrEntryBeforeNext = errors.New("Entry called before Next")
	// ErrNoMoreEntries is the error returned if Entry() is called
	// after Next() returned false for a DirReader.
	ErrNoMoreEntries = errors.New("Entry called with no more entries")
)

// Filesystem is an interface that gives read access to the filesystem of the machine to scan.
type Filesystem interface {
	// OpenFile opens the specified file for reading.
	// It should return an os.IsNotExist error if the file doesn't exist.
	OpenFile(ctx context.Context, path string) (io.ReadCloser, error)
	// FilePermissions returns unix permission-related data for the specified file or directory.
	FilePermissions(ctx context.Context, path string) (*apb.PosixPermissions, error)
	// OpenDir opens the specified directory to list its content.
	OpenDir(ctx context.Context, path string) (DirReader, error)
}

// SQLQuerier is an interface that supports SQL queries to a target SQL database.
type SQLQuerier interface {
	// SQLQuery executes SQL queries to a target SQL database and returns first result
	// row as a string.
	SQLQuery(ctx context.Context, query string) (string, error)
	// Returns the supported database type
	SupportedDatabase() (ipb.SQLCheck_SQLDatabase, error)
}

// ScanAPI is an interface that gives read access to the filesystem of
// the machine to scan and can execute SQL queries on a single database.
type ScanAPI interface {
	Filesystem
	SQLQuerier
}

// DirReader is an interface to iterate the entries inside a directory.
type DirReader interface {
	// Next reads the next entry in the directory which can then be accessed using Entry.
	// It must be called at least once before calling Entry.
	// Returns false if there are no more entries in the directory.
	Next() bool
	// Entry returns the last entry read using Next or an error if it failed.
	Entry() (*apb.DirContent, error)
	// Close must be called to correctly dispose of the underlying reader.
	Close() error
}

type sliceDirReader struct {
	entries []*apb.DirContent
	idx     int
}

func (s *sliceDirReader) Next() bool {
	s.idx++
	return s.idx < len(s.entries)
}

func (s *sliceDirReader) Entry() (*apb.DirContent, error) {
	if s.idx == -1 {
		return nil, ErrEntryBeforeNext
	}
	if s.idx >= len(s.entries) {
		return nil, ErrNoMoreEntries
	}
	return s.entries[s.idx], nil
}

func (s *sliceDirReader) Close() error {
	s.idx = len(s.entries)
	return nil
}

// SliceToDirReader returns a DirReader given a slice of entries.
func SliceToDirReader(entries []*apb.DirContent) DirReader {
	return &sliceDirReader{entries: entries, idx: -1}
}

// DirReaderToSlice returns a slice of all the entries left in the given DirReader.
// The DirReader is automatically disposed of by calling Close() at the end.
func DirReaderToSlice(d DirReader) ([]*apb.DirContent, error) {
	defer d.Close()
	entries := make([]*apb.DirContent, 0)
	for d.Next() {
		e, err := d.Entry()
		if err != nil {
			return nil, err
		}
		entries = append(entries, e)
	}
	return entries, nil
}
