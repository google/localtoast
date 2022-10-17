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

// Package localfilereader provides utility functions for reading files and permissions from the local filesystem.
package localfilereader

import (
	"context"
	"io"
	"io/fs"
	"os"
	"os/user"
	"strconv"
	"strings"
	"syscall"

	"github.com/google/localtoast/scanapi"
	apb "github.com/google/localtoast/scannerlib/proto/api_go_proto"
)

// Flags for special UNIX permission bits.
const (
	SetuidFlag = 04000
	SetgidFlag = 02000
	StickyFlag = 01000
)

var (
	userIDLookup = newCachedIDLookup(func(id int) (string, error) {
		usr, err := user.LookupId(strconv.Itoa(id))
		if err != nil {
			return "", err
		}
		return usr.Username, nil
	})
	groupIDLookup = newCachedIDLookup(func(id int) (string, error) {
		grp, err := user.LookupGroupId(strconv.Itoa(id))
		if err != nil {
			return "", err
		}
		return grp.Name, nil
	})
)

type cachedIDLookup struct {
	lookupFunc func(int) (string, error)
	valueCache map[int]string
	errorCache map[int]error
}

func (l *cachedIDLookup) Lookup(id int) (string, error) {
	if val, ok := l.valueCache[id]; ok {
		return val, nil
	}
	if err, ok := l.errorCache[id]; ok {
		return "", err
	}
	val, err := l.lookupFunc(id)
	if err == nil {
		l.valueCache[id] = val
	} else {
		l.errorCache[id] = err
	}
	return val, err
}

func newCachedIDLookup(lookupFunc func(int) (string, error)) cachedIDLookup {
	return cachedIDLookup{
		lookupFunc: lookupFunc,
		valueCache: make(map[int]string),
		errorCache: make(map[int]error),
	}
}

// OpenFile opens the specified file for reading.
func OpenFile(ctx context.Context, path string) (io.ReadCloser, error) {
	return os.Open(path)
}

// FilePermissions returns unix permission-related data for the specified file or directory.
func FilePermissions(ctx context.Context, path string) (*apb.PosixPermissions, error) {
	fi, err := os.Lstat(path)
	if err != nil {
		return nil, err
	}
	sys := fi.Sys()
	uid := int(sys.(*syscall.Stat_t).Uid)
	gid := int(sys.(*syscall.Stat_t).Gid)

	username, err := userIDLookup.Lookup(uid)
	if err != nil {
		// "unknown userid" means the file is unowned (uid not found
		// in /etc/group, possibly because the user got deleted). Leave
		// the username empty to signal this.
		if !strings.Contains(err.Error(), "unknown userid") {
			return nil, err
		}
	}

	groupname, err := groupIDLookup.Lookup(gid)
	if err != nil {
		// "unknown groupid" means the file is ungrouped (gid not found
		// in /etc/group, possibly because the group got deleted). Leave
		// the groupname empty to signal this.
		if !strings.Contains(err.Error(), "unknown groupid") {
			return nil, err
		}
	}

	perms := int32(fi.Mode().Perm())
	// Mode().Perm() only contains the regular permission bits, so add the
	// special flag bits separately.
	if fi.Mode()&fs.ModeSetuid != 0 {
		perms |= SetuidFlag
	}
	if fi.Mode()&fs.ModeSetgid != 0 {
		perms |= SetgidFlag
	}
	if fi.Mode()&fs.ModeSticky != 0 {
		perms |= StickyFlag
	}
	return &apb.PosixPermissions{
		PermissionNum: perms,
		Uid:           int32(uid),
		User:          username,
		Gid:           int32(gid),
		Group:         groupname,
	}, nil
}

// OpenDir opens the specified directory to list its content.
func OpenDir(ctx context.Context, dirPath string) (scanapi.DirReader, error) {
	f, err := os.Open(dirPath)
	if err != nil {
		return nil, err
	}
	return &localDirReader{file: f, currErr: scanapi.ErrEntryBeforeNext}, nil
}

type localDirReader struct {
	file      *os.File
	currEntry *apb.DirContent
	currErr   error
}

func (d *localDirReader) nextEntry() (*apb.DirContent, error) {
	// Read the next entry until the EOF is reached, there was an error,
	// or a valid entry is found (either a dir, a regular file or a symlink).
	for {
		entries, err := d.file.ReadDir(1)
		if err != nil {
			return nil, err
		}
		e := entries[0]
		if e.IsDir() || e.Type().IsRegular() || e.Type()&fs.ModeSymlink == fs.ModeSymlink {
			return &apb.DirContent{
				Name:      e.Name(),
				IsDir:     e.IsDir(),
				IsSymlink: e.Type()&fs.ModeSymlink == fs.ModeSymlink,
			}, nil
		}
	}
}

func (d *localDirReader) Next() bool {
	d.currEntry, d.currErr = d.nextEntry()
	return d.currErr != io.EOF
}

func (d *localDirReader) Entry() (*apb.DirContent, error) {
	return d.currEntry, d.currErr
}

func (d *localDirReader) Close() error {
	return d.file.Close()
}
