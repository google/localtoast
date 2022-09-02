package scanapi_test

import (
	"errors"
	"testing"

	"github.com/google/go-cmp/cmp"
	"google.golang.org/protobuf/testing/protocmp"
	"github.com/google/localtoast/scanapi"
	apb "github.com/google/localtoast/scannerlib/proto/api_go_proto"
)

var (
	testEntries = []*apb.DirContent{{Name: "e1"}, {Name: "e2"}}
)

func TestSliceToDirReader(t *testing.T) {
	d := scanapi.SliceToDirReader(testEntries)
	defer d.Close()

	dirReaderEntries := []*apb.DirContent{}
	for d.Next() {
		e, err := d.Entry()
		if err != nil {
			t.Fatalf("DirReader.Entry() had unexpected error: %v", err)
		}
		dirReaderEntries = append(dirReaderEntries, e)
	}
	if diff := cmp.Diff(testEntries, dirReaderEntries, protocmp.Transform()); diff != "" {
		t.Errorf("DirReader returned unexpected entries (-want +got):\n%s", diff)
	}
}

func TestSliceToDirReaderEmptySlice(t *testing.T) {
	d := scanapi.SliceToDirReader([]*apb.DirContent{})
	defer d.Close()
	if got := d.Next(); got {
		t.Errorf("DirReader.Next() got: %v, want: false", got)
	}
}

func TestSliceToDirReaderErrEntryBeforeNext(t *testing.T) {
	d := scanapi.SliceToDirReader(testEntries)
	defer d.Close()
	if _, err := d.Entry(); err == nil {
		t.Errorf("DirReader.Entry() err got: %v, want: error", err)
	}
}

func TestSliceToDirReaderErrNoMoreEntries(t *testing.T) {
	d := scanapi.SliceToDirReader([]*apb.DirContent{})
	defer d.Close()
	if got := d.Next(); got {
		t.Fatalf("DirReader.Next() got: %v, want: false", got)
	}
	if _, err := d.Entry(); err == nil {
		t.Errorf("DirReader.Entry() err got: %v, want: error", err)
	}
}

func TestSliceToDirReaderErrAfterClose(t *testing.T) {
	d := scanapi.SliceToDirReader(testEntries)
	d.Close()
	if got := d.Next(); got {
		t.Fatalf("DirReader.Next() got: %v, want: false", got)
	}
	if _, err := d.Entry(); err == nil {
		t.Errorf("DirReader.Entry() err got: %v, want: error", err)
	}
}

func TestSliceToDirReaderToSlice(t *testing.T) {
	d := scanapi.SliceToDirReader(testEntries)
	gotEntries, err := scanapi.DirReaderToSlice(d)
	if err != nil {
		t.Fatalf("DirReaderToSlice() had unexpected error: %v", err)
	}
	if diff := cmp.Diff(testEntries, gotEntries, protocmp.Transform()); diff != "" {
		t.Errorf("DirReaderToSlice() returned unexpected entries (-want +got):\n%s", diff)
	}
}

type errorDirReader struct{}

func (errorDirReader) Next() bool {
	return true
}

func (errorDirReader) Entry() (*apb.DirContent, error) {
	return nil, errors.New("error")
}

func (errorDirReader) Close() error {
	return nil
}

func TestDirReaderToSliceError(t *testing.T) {
	if _, err := scanapi.DirReaderToSlice(errorDirReader{}); err == nil {
		t.Errorf("DirReaderToSlice() err got: %v, want: error", err)
	}
}
