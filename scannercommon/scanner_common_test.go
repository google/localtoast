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

package scannercommon_test

import (
	"context"
	"errors"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
	"google.golang.org/protobuf/testing/protocmp"
	"github.com/google/localtoast/cli"
	"github.com/google/localtoast/localfilereader"
	"github.com/google/localtoast/protofilehandler"
	"github.com/google/localtoast/scannercommon"
	apb "github.com/google/localtoast/scannerlib/proto/api_go_proto"
	ipb "github.com/google/localtoast/scannerlib/proto/scan_instructions_go_proto"
	"github.com/google/localtoast/scannerlib/testconfigcreator"
)

var testDirPath = ""

// A ScanApiProvider implementation that only provides filesystem access.
type testAPIProvider struct{}

func (testAPIProvider) OpenFile(ctx context.Context, filePath string) (io.ReadCloser, error) {
	return localfilereader.OpenFile(ctx, path.Join(testDirPath, filePath))
}

func (testAPIProvider) FilesInDir(ctx context.Context, dirPath string) ([]*apb.DirContent, error) {
	return localfilereader.FilesInDir(ctx, path.Join(testDirPath, dirPath))
}

func (testAPIProvider) FilePermissions(ctx context.Context, filePath string) (*apb.PosixPermissions, error) {
	return localfilereader.FilePermissions(ctx, path.Join(testDirPath, filePath))
}

func (testAPIProvider) SQLQuery(ctx context.Context, query string) (int, error) {
	return 0, errors.New("not implemented")
}

func TestRunScan(t *testing.T) {
	testDirPath = os.Getenv("TEST_TMPDIR")
	configPath := filepath.Join(testDirPath, "config.textproto")
	resultPath := filepath.Join(testDirPath, "result.textproto")

	testFile := "test.txt"
	if err := ioutil.WriteFile(filepath.Join(testDirPath, testFile), []byte("text"), 0644); err != nil {
		t.Fatalf("Error while creating file %s: %v", testFile, err)
	}

	flags := &cli.Flags{
		ConfigFile:              configPath,
		ResultFile:              resultPath,
		ShowCompliantBenchmarks: true,
	}
	provider := &testAPIProvider{}

	testCases := []struct {
		description               string
		check                     *ipb.FileCheck
		expectedCompliantCount    int
		expectedNonCompliantCount int
	}{
		{
			description: "compliant scan",
			check: &ipb.FileCheck{
				FilesToCheck: []*ipb.FileSet{testconfigcreator.SingleFileWithPath(testFile)},
				CheckType:    &ipb.FileCheck_Existence{Existence: &ipb.ExistenceCheck{ShouldExist: true}},
			},
			expectedCompliantCount:    1,
			expectedNonCompliantCount: 0,
		},
		{
			description: "non-compliant scan",
			check: &ipb.FileCheck{
				FilesToCheck: []*ipb.FileSet{testconfigcreator.SingleFileWithPath(testFile)},
				CheckType:    &ipb.FileCheck_Existence{Existence: &ipb.ExistenceCheck{ShouldExist: false}},
			},
			expectedCompliantCount:    0,
			expectedNonCompliantCount: 1,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			config := &apb.ScanConfig{
				BenchmarkConfigs: []*apb.BenchmarkConfig{
					testconfigcreator.NewBenchmarkConfig(t, "test", testconfigcreator.NewFileScanInstruction([]*ipb.FileCheck{tc.check})),
				},
			}
			if err := protofilehandler.WriteProtoToFile(configPath, config); err != nil {
				t.Fatalf("Error writing scan config: %v", err)
			}

			scannercommon.RunScan(flags, provider)

			result := &apb.ScanResults{}
			if err := protofilehandler.ReadProtoFromFile(resultPath, result); err != nil {
				t.Fatalf("Error reading scan results: %v", err)
			}

			if tc.expectedCompliantCount != len(result.GetCompliantBenchmarks()) {
				t.Errorf("unexpected compliant benchmark count, want %d got %d", tc.expectedCompliantCount, len(result.GetCompliantBenchmarks()))
			}
			if tc.expectedNonCompliantCount != len(result.GetNonCompliantBenchmarks()) {
				t.Errorf("unexpected non-compliant benchmark count, want %d got %d", tc.expectedNonCompliantCount, len(result.GetNonCompliantBenchmarks()))
			}
		})
	}
}

func TestApplyCLIFlagsToConfig(t *testing.T) {
	testCases := []struct {
		desc   string
		flags  *cli.Flags
		config *apb.ScanConfig
		want   *apb.ScanConfig
	}{
		{
			desc:  "benchmark opt-out",
			flags: &cli.Flags{BenchmarkOptOutIDs: "id1,id3"},
			config: &apb.ScanConfig{
				BenchmarkConfigs: []*apb.BenchmarkConfig{
					&apb.BenchmarkConfig{Id: "id1"},
					&apb.BenchmarkConfig{Id: "id2"},
					&apb.BenchmarkConfig{Id: "id3"},
					&apb.BenchmarkConfig{Id: "id4"},
				},
			},
			want: &apb.ScanConfig{
				BenchmarkConfigs: []*apb.BenchmarkConfig{
					&apb.BenchmarkConfig{Id: "id2"},
					&apb.BenchmarkConfig{Id: "id4"},
				},
			},
		},
		{
			desc:   "file content opt-out",
			flags:  &cli.Flags{ContentOptOutRegexes: "regex1,regex2"},
			config: &apb.ScanConfig{},
			want: &apb.ScanConfig{
				OptOutConfig: &apb.OptOutConfig{
					ContentOptoutRegexes:  []string{"regex1", "regex2"},
					FilenameOptoutRegexes: []string{},
				},
			},
		},
		{
			desc:   "file name opt-out",
			flags:  &cli.Flags{FilenameOptOutRegexes: "regex1,regex2"},
			config: &apb.ScanConfig{},
			want: &apb.ScanConfig{
				OptOutConfig: &apb.OptOutConfig{
					ContentOptoutRegexes:  []string{},
					FilenameOptoutRegexes: []string{"regex1", "regex2"},
				},
			},
		},
		{
			desc:   "file traversal opt-out",
			flags:  &cli.Flags{TraversalOptOutRegexes: "regex1,regex2"},
			config: &apb.ScanConfig{},
			want: &apb.ScanConfig{
				OptOutConfig: &apb.OptOutConfig{
					ContentOptoutRegexes:   []string{},
					TraversalOptoutRegexes: []string{"regex1", "regex2"},
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			scannercommon.ApplyCLIFlagsToConfig(tc.config, tc.flags)

			if diff := cmp.Diff(tc.want, tc.config, protocmp.Transform()); diff != "" {
				t.Errorf("applyCLIFLagsToConfig(%v) returned unexpected diff (-want +got):\n%s",
					tc.flags, diff)
			}
		})
	}
}
