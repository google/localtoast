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

package scanner_test

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/testing/protocmp"
	gpb "google.golang.org/genproto/googleapis/grafeas/v1"
	apb "github.com/google/localtoast/library/proto/api_go_proto"
	ipb "github.com/google/localtoast/library/proto/scan_instructions_go_proto"
	"github.com/google/localtoast/library/scanner"
	"github.com/google/localtoast/library/testing/testconfigcreator"
)

const (
	testFilePath1          = "/path/to/test/file1"
	testFileContent1       = "File content 1"
	testFilePath2          = "/path/to/test/file2"
	testFileContent2       = "File content 2"
	regexMatchingTestFiles = "/path/to/test.*"
	testQueryNoRows        = "SELECT 1 WHERE FALSE"
	testQueryOneRow        = "SELECT 1"
)

type fakeAPIProvider struct{}

func (fakeAPIProvider) OpenFile(ctx context.Context, path string) (io.ReadCloser, error) {
	if path == testFilePath1 {
		return ioutil.NopCloser(bytes.NewReader([]byte(testFileContent1))), nil
	} else if path == testFilePath2 {
		return ioutil.NopCloser(bytes.NewReader([]byte(testFileContent2))), nil
	}
	return nil, errors.New("File not found")
}
func (fakeAPIProvider) FilesInDir(ctx context.Context, path string) ([]*apb.DirContent, error) {
	return nil, errors.New("not implemented")
}
func (fakeAPIProvider) FilePermissions(ctx context.Context, path string) (*apb.PosixPermissions, error) {
	return nil, errors.New("not implemented")
}
func (fakeAPIProvider) SQLQuery(ctx context.Context, query string) (int, error) {
	switch query {
	case testQueryNoRows:
		return 0, nil
	case testQueryOneRow:
		return 1, nil
	default:
		return 0, fmt.Errorf("the query %q is not supported by fakeAPIProvider", query)
	}
}

func TestScannerVersion(t *testing.T) {
	config := &apb.ScanConfig{}
	result, err := scanner.Scanner{}.Scan(context.Background(), config, fakeAPIProvider{})

	if err != nil {
		t.Fatalf("scanner.Scan() had unexpected error: %v", err)
	}
	if result.GetStatus().GetStatus() != apb.ScanStatus_SUCCEEDED {
		t.Fatalf("scanner.Scan() returned unsuccessful scan status: %v",
			result.GetStatus().GetStatus())
	}

	if result.GetScannerVersion() != scanner.ScannerVersion {
		t.Errorf("scanner.Scan() returned scanner version: %s, expected %s",
			result.GetScannerVersion(), scanner.ScannerVersion)
	}
}

func TestCompliantScan(t *testing.T) {
	compliantCheck := []*ipb.FileCheck{
		&ipb.FileCheck{
			FilesToCheck: []*ipb.FileSet{testconfigcreator.SingleFileWithPath(testFilePath1)},
			CheckType:    &ipb.FileCheck_Content{Content: &ipb.ContentCheck{Content: testFileContent1}},
		},
	}
	config := &apb.ScanConfig{
		BenchmarkConfigs: []*apb.BenchmarkConfig{
			testconfigcreator.NewBenchmarkConfig(
				t, "id", testconfigcreator.NewFileScanInstruction(compliantCheck)),
		},
	}

	result, err := scanner.Scanner{}.Scan(context.Background(), config, fakeAPIProvider{})

	if err != nil {
		t.Fatalf("scanner.Scan(%v) had unexpected error: %v", config, err)
	}
	if result.GetStatus().GetStatus() != apb.ScanStatus_SUCCEEDED {
		t.Fatalf("scanner.Scan(%v) returned unsuccessful scan status: %v",
			config, result.GetStatus().GetStatus())
	}

	if len(result.GetCompliantBenchmarks()) != 1 {
		t.Errorf("scanner.Scan(%v) returned check result: %v, expected 1 compliant check.",
			config, result)
	}
	if len(result.GetNonCompliantBenchmarks()) != 0 {
		t.Errorf("scanner.Scan(%v) returned check result: %v, expected 0 non-compliant configchecks.",
			config, result)
	}
}

func TestNonCompliantScan(t *testing.T) {
	testCases := []struct {
		desc   string
		config *apb.ScanConfig
	}{
		{
			desc: "non compliant files",
			config: &apb.ScanConfig{
				BenchmarkConfigs: []*apb.BenchmarkConfig{
					testconfigcreator.NewBenchmarkConfig(
						t, "id", testconfigcreator.NewFileScanInstruction([]*ipb.FileCheck{
							{
								FilesToCheck: []*ipb.FileSet{testconfigcreator.SingleFileWithPath(testFilePath1)},
								CheckType:    &ipb.FileCheck_Content{Content: &ipb.ContentCheck{Content: "Different content"}},
							},
						})),
				},
			},
		},
		{
			desc: "non compliance reason",
			config: &apb.ScanConfig{
				BenchmarkConfigs: []*apb.BenchmarkConfig{
					testconfigcreator.NewBenchmarkConfig(
						t, "id", testconfigcreator.NewSQLScanInstruction([]*ipb.SQLCheck{
							&ipb.SQLCheck{
								TargetDatabase: ipb.SQLCheck_DB_MYSQL,
								Query:          testQueryOneRow,
								ExpectResults:  false,
							},
						})),
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			result, err := scanner.Scanner{}.Scan(context.Background(), tc.config, fakeAPIProvider{})

			if err != nil {
				t.Fatalf("scanner.Scan(%v) had unexpected error: %v", tc.config, err)
			}
			if result.GetStatus().GetStatus() != apb.ScanStatus_SUCCEEDED {
				t.Fatalf("scanner.Scan(%v) returned unsuccessful scan status: %v",
					tc.config, result.GetStatus().GetStatus())
			}

			if len(result.GetCompliantBenchmarks()) != 0 {
				t.Errorf("scanner.Scan(%v) returned check result: %v, expected 0 compliant configchecks.", tc.config, result)
			}
			if len(result.GetNonCompliantBenchmarks()) != 1 {
				t.Errorf("scanner.Scan(%v) returned check result: %v, expected 1 non-compliant check.",
					tc.config, result)
			}
		})
	}
}

func TestFailingScan(t *testing.T) {
	nonExistentPath := "/non/existent/file"
	failingCheck := []*ipb.FileCheck{
		&ipb.FileCheck{
			FilesToCheck: []*ipb.FileSet{testconfigcreator.SingleFileWithPath(nonExistentPath)},
			CheckType:    &ipb.FileCheck_Content{Content: &ipb.ContentCheck{Content: "Content"}},
		},
	}
	config := &apb.ScanConfig{
		BenchmarkConfigs: []*apb.BenchmarkConfig{
			testconfigcreator.NewBenchmarkConfig(
				t, "id", testconfigcreator.NewFileScanInstruction(failingCheck)),
		},
	}

	result, err := scanner.Scanner{}.Scan(context.Background(), config, fakeAPIProvider{})

	if err != nil {
		t.Fatalf("scanner.Scan(%v) had unexpected error: %v", config, err)
	}
	if result.GetStatus().GetStatus() != apb.ScanStatus_FAILED {
		t.Errorf("scanner.Scan(%v) returned scan status: %v, expected ScanStatus_FAILED",
			config, result.GetStatus().GetStatus())
	}

	expectedFailureReason := fmt.Sprintf(
		"Compliance state of the following benchmarks couldn't be determined: [id]\n"+
			"The following errors were encountered while running the checks:\n"+
			"[content check on single_file:{path:%q}]: api.OpenFile(%q): File not found\n",
		nonExistentPath, nonExistentPath)

	if diff := cmp.Diff(expectedFailureReason, result.GetStatus().GetFailureReason(),
		protocmp.Transform()); diff != "" {
		t.Errorf("scanner.Scan(%v) returned unexpected failure reason, (-want +got):\n%s", config, diff)
	}
}

func TestNonCompliantFileCheckResultsAreAggregated(t *testing.T) {
	testCases := []struct {
		desc                           string
		checkAlternative               *ipb.CheckAlternative
		expectedNonCompliantBenchmarks []*apb.ComplianceResult
	}{
		{
			desc: "File checks",
			checkAlternative: &ipb.CheckAlternative{
				FileChecks: []*ipb.FileCheck{
					&ipb.FileCheck{
						FilesToCheck: []*ipb.FileSet{testconfigcreator.SingleFileWithPath(testFilePath1)},
						CheckType:    &ipb.FileCheck_Content{Content: &ipb.ContentCheck{Content: "Different content 1"}},
					},
					&ipb.FileCheck{
						FilesToCheck: []*ipb.FileSet{testconfigcreator.SingleFileWithPath(testFilePath2)},
						CheckType:    &ipb.FileCheck_Content{Content: &ipb.ContentCheck{Content: "Different content 2"}},
					},
				},
			},
			expectedNonCompliantBenchmarks: []*apb.ComplianceResult{
				&apb.ComplianceResult{
					Id: "id",
					ComplianceOccurrence: &gpb.ComplianceOccurrence{
						NonCompliantFiles: []*gpb.NonCompliantFile{
							&gpb.NonCompliantFile{
								Path:   testFilePath1,
								Reason: fmt.Sprintf("Got content %q, expected \"Different content 1\"", testFileContent1),
							},
							&gpb.NonCompliantFile{
								Path:   testFilePath2,
								Reason: fmt.Sprintf("Got content %q, expected \"Different content 2\"", testFileContent2),
							},
						},
					},
				},
			},
		},
		{
			desc: "SQL checks",
			checkAlternative: &ipb.CheckAlternative{
				SqlChecks: []*ipb.SQLCheck{
					&ipb.SQLCheck{
						TargetDatabase: ipb.SQLCheck_DB_MYSQL,
						Query:          testQueryOneRow,
						ExpectResults:  false,
					},
					&ipb.SQLCheck{
						TargetDatabase: ipb.SQLCheck_DB_MYSQL,
						Query:          testQueryNoRows,
						ExpectResults:  true,
					},
				},
			},
			expectedNonCompliantBenchmarks: []*apb.ComplianceResult{
				&apb.ComplianceResult{
					Id: "id",
					ComplianceOccurrence: &gpb.ComplianceOccurrence{
						NonCompliantFiles:   []*gpb.NonCompliantFile{},
						NonComplianceReason: fmt.Sprintf("Expected no results for query %q, but got 1 rows.\nExpected results for query %q, but got none.", testQueryOneRow, testQueryNoRows),
					},
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			config := &apb.ScanConfig{
				BenchmarkConfigs: []*apb.BenchmarkConfig{
					testconfigcreator.NewBenchmarkConfig(
						t, "id", &ipb.BenchmarkScanInstruction{
							CheckAlternatives: []*ipb.CheckAlternative{tc.checkAlternative},
						}),
				},
			}

			result, err := scanner.Scanner{}.Scan(context.Background(), config, fakeAPIProvider{})

			if err != nil {
				t.Fatalf("scanner.Scan(%v) had unexpected error: %v", config, err)
			}
			if result.GetStatus().GetStatus() != apb.ScanStatus_SUCCEEDED {
				t.Fatalf("scanner.Scan(%v) returned unsuccessful scan status: %v",
					config, result.GetStatus().GetStatus())
			}

			if len(result.GetCompliantBenchmarks()) != 0 {
				t.Errorf("scanner.Scan(%v) returned check result: %v, expected 0 compliant configchecks.",
					config, result)
			}

			sortProtos := cmpopts.SortSlices(func(m1, m2 protocmp.Message) bool { return m1.String() < m2.String() })
			if diff := cmp.Diff(tc.expectedNonCompliantBenchmarks, result.GetNonCompliantBenchmarks(),
				protocmp.Transform(), sortProtos); diff != "" {
				t.Errorf("scanner.Scan(%v) returned unexpected results (-want +got):\n%s", config, diff)
			}
		})
	}
}

func TestBenchmarkIsNonCompliantIfOneCheckIsNonCompliant(t *testing.T) {
	nonCompliantCheck := &ipb.FileCheck{
		FilesToCheck: []*ipb.FileSet{testconfigcreator.SingleFileWithPath(testFilePath1)},
		CheckType:    &ipb.FileCheck_Content{Content: &ipb.ContentCheck{Content: "Different content 1"}},
	}
	compliantCheck := &ipb.FileCheck{
		FilesToCheck: []*ipb.FileSet{testconfigcreator.SingleFileWithPath(testFilePath2)},
		CheckType:    &ipb.FileCheck_Content{Content: &ipb.ContentCheck{Content: testFileContent2}},
	}
	checks := []*ipb.FileCheck{nonCompliantCheck, compliantCheck}
	config := &apb.ScanConfig{
		BenchmarkConfigs: []*apb.BenchmarkConfig{
			testconfigcreator.NewBenchmarkConfig(
				t, "id", testconfigcreator.NewFileScanInstruction(checks)),
		},
	}

	result, err := scanner.Scanner{}.Scan(context.Background(), config, fakeAPIProvider{})

	if err != nil {
		t.Fatalf("scanner.Scan(%v) had unexpected error: %v", config, err)
	}
	if result.GetStatus().GetStatus() != apb.ScanStatus_SUCCEEDED {
		t.Fatalf("scanner.Scan(%v) returned unsuccessful scan status: %v",
			config, result.GetStatus().GetStatus())
	}

	if len(result.GetCompliantBenchmarks()) != 0 {
		t.Errorf("scanner.Scan(%v) returned check result: %v, expected 0 compliant configchecks.",
			config, result)
	}
	if len(result.GetNonCompliantBenchmarks()) != 1 {
		t.Errorf("scanner.Scan(%v) returned check result: %v, expected 1 non-compliant check.",
			config, result)
	}
}

func TestDuplicateBenchmarkIDs(t *testing.T) {
	check := []*ipb.FileCheck{&ipb.FileCheck{
		FilesToCheck: []*ipb.FileSet{testconfigcreator.SingleFileWithPath(testFilePath1)},
		CheckType:    &ipb.FileCheck_Content{Content: &ipb.ContentCheck{Content: testFileContent1}},
	}}
	config := &apb.ScanConfig{
		BenchmarkConfigs: []*apb.BenchmarkConfig{
			testconfigcreator.NewBenchmarkConfig(
				t, "id", testconfigcreator.NewFileScanInstruction(check)),
			testconfigcreator.NewBenchmarkConfig(
				t, "id", testconfigcreator.NewFileScanInstruction(check)),
		},
	}

	if _, err := (scanner.Scanner{}).Scan(context.Background(), config, fakeAPIProvider{}); err == nil {
		t.Fatalf("scanner.Scan(%v) didn't return an error", config)
	}
}

func TestAlternativeWithFileAndDBChecks(t *testing.T) {
	testCases := []struct {
		desc             string
		fileCheck        *ipb.FileCheck
		sqlCheck         *ipb.SQLCheck
		expectCompliance bool
	}{
		{
			desc: "Both checks compliant",
			fileCheck: &ipb.FileCheck{
				FilesToCheck: []*ipb.FileSet{testconfigcreator.SingleFileWithPath(testFilePath1)},
				CheckType:    &ipb.FileCheck_Existence{Existence: &ipb.ExistenceCheck{ShouldExist: true}},
			},
			sqlCheck: &ipb.SQLCheck{
				TargetDatabase: ipb.SQLCheck_DB_MYSQL,
				Query:          "SELECT 1",
				ExpectResults:  true,
			},
			expectCompliance: true,
		},
		{
			desc: "One check non-compliant",
			fileCheck: &ipb.FileCheck{
				FilesToCheck: []*ipb.FileSet{testconfigcreator.SingleFileWithPath(testFilePath1)},
				CheckType:    &ipb.FileCheck_Existence{Existence: &ipb.ExistenceCheck{ShouldExist: true}},
			},
			sqlCheck: &ipb.SQLCheck{
				TargetDatabase: ipb.SQLCheck_DB_MYSQL,
				Query:          "SELECT 1",
				ExpectResults:  false,
			},
			expectCompliance: false,
		},
		{
			desc: "Both checks non-compliant",
			fileCheck: &ipb.FileCheck{
				FilesToCheck: []*ipb.FileSet{testconfigcreator.SingleFileWithPath(testFilePath1)},
				CheckType:    &ipb.FileCheck_Existence{Existence: &ipb.ExistenceCheck{ShouldExist: false}},
			},
			sqlCheck: &ipb.SQLCheck{
				TargetDatabase: ipb.SQLCheck_DB_MYSQL,
				Query:          "SELECT 1",
				ExpectResults:  false,
			},
			expectCompliance: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			config := &apb.ScanConfig{
				BenchmarkConfigs: []*apb.BenchmarkConfig{
					testconfigcreator.NewBenchmarkConfig(
						t, "id", &ipb.BenchmarkScanInstruction{
							CheckAlternatives: []*ipb.CheckAlternative{
								&ipb.CheckAlternative{
									FileChecks: []*ipb.FileCheck{tc.fileCheck},
									SqlChecks:  []*ipb.SQLCheck{tc.sqlCheck},
								},
							},
						}),
				},
			}

			result, err := scanner.Scanner{}.Scan(context.Background(), config, fakeAPIProvider{})
			if err != nil {
				t.Fatalf("scanner.Scan(%v) returned an error: %v", config, err)
			}
			if result.GetStatus().GetStatus() != apb.ScanStatus_SUCCEEDED {
				t.Fatalf("scanner.Scan(%v) returned unsuccessful scan status: %v",
					config, result.GetStatus().GetStatus())
			}

			compliant := len(result.GetNonCompliantBenchmarks()) == 0
			if tc.expectCompliance != compliant {
				t.Errorf("scanner.Scan(%v) expected to return compliance status %t, got %t:\n%v", config, tc.expectCompliance, compliant, result)
			}
		})
	}
}

func TestInstructionsWithUnknownField(t *testing.T) {
	instructions := fmt.Sprintf("check_alternatives:{unknown_field: \"value\" file_checks:{files_to_check:{single_file:{path:%q}} existence:{should_exist: true}}}", testFilePath1)
	config := &apb.ScanConfig{
		BenchmarkConfigs: []*apb.BenchmarkConfig{
			&apb.BenchmarkConfig{
				Id: "test",
				ComplianceNote: &gpb.ComplianceNote{
					ScanInstructions: []byte(instructions),
				},
			},
		},
	}

	result, err := scanner.Scanner{}.Scan(context.Background(), config, fakeAPIProvider{})
	if err != nil {
		t.Fatalf("scanner.Scan(%v) returned an error: %v", config, err)
	}
	if result.GetStatus().GetStatus() != apb.ScanStatus_SUCCEEDED {
		t.Errorf("scanner.Scan(%v) returned unsuccessful scan status: %v", config, result.GetStatus().GetStatus())
	}
}

func TestInstructionsWithBinarySerialization(t *testing.T) {
	compliantCheck := []*ipb.FileCheck{{
		FilesToCheck: []*ipb.FileSet{testconfigcreator.SingleFileWithPath(testFilePath1)},
		CheckType:    &ipb.FileCheck_Content{Content: &ipb.ContentCheck{Content: testFileContent1}},
	}}
	instructions := testconfigcreator.NewFileScanInstruction(compliantCheck)
	serInstructions, err := proto.Marshal(instructions)
	if err != nil {
		t.Fatalf("proto.Marshal(%v) returned error: %v", instructions, err)
	}

	config := &apb.ScanConfig{
		BenchmarkConfigs: []*apb.BenchmarkConfig{
			&apb.BenchmarkConfig{
				Id: "test",
				ComplianceNote: &gpb.ComplianceNote{
					ScanInstructions: serInstructions,
				},
			},
		},
	}

	result, err := scanner.Scanner{}.Scan(context.Background(), config, fakeAPIProvider{})
	if err != nil {
		t.Fatalf("scanner.Scan(%v) returned an error: %v", config, err)
	}
	if result.GetStatus().GetStatus() != apb.ScanStatus_SUCCEEDED {
		t.Errorf("scanner.Scan(%v) returned unsuccessful scan status: %v", config, result.GetStatus().GetStatus())
	}

	if len(result.GetNonCompliantBenchmarks()) > 0 {
		t.Errorf("scanner.Scan(%v) returned non-compliant benchmarks, expected none: %v", config, result.GetNonCompliantBenchmarks())
	}
}

func TestCheckAlternativeAggregation(t *testing.T) {
	compliantChecks := []*ipb.FileCheck{&ipb.FileCheck{
		FilesToCheck: []*ipb.FileSet{testconfigcreator.SingleFileWithPath(testFilePath2)},
		CheckType:    &ipb.FileCheck_Content{Content: &ipb.ContentCheck{Content: testFileContent2}},
	}}
	nonCompliantChecks := []*ipb.FileCheck{&ipb.FileCheck{
		FilesToCheck: []*ipb.FileSet{testconfigcreator.SingleFileWithPath(testFilePath1)},
		CheckType:    &ipb.FileCheck_Content{Content: &ipb.ContentCheck{Content: "Different content"}},
	}}
	nonCompliantChecks2 := []*ipb.FileCheck{&ipb.FileCheck{
		FilesToCheck: []*ipb.FileSet{testconfigcreator.SingleFileWithPath(testFilePath2)},
		CheckType:    &ipb.FileCheck_Content{Content: &ipb.ContentCheck{Content: "Different content"}},
	}}
	failingChecks := []*ipb.FileCheck{&ipb.FileCheck{
		FilesToCheck: []*ipb.FileSet{testconfigcreator.SingleFileWithPath("/non/existent/file")},
		CheckType:    &ipb.FileCheck_Content{Content: &ipb.ContentCheck{Content: "Content"}},
	}}

	testCases := []struct {
		desc                           string
		instructions                   *ipb.BenchmarkScanInstruction
		expectedCompliantBenchmarks    []*apb.ComplianceResult
		expectedNonCompliantBenchmarks []*apb.ComplianceResult
		expectedScanStatus             apb.ScanStatus_ScanStatusEnum
	}{
		{
			desc: "Two compliant alternatives",
			instructions: &ipb.BenchmarkScanInstruction{
				CheckAlternatives: []*ipb.CheckAlternative{
					&ipb.CheckAlternative{
						FileChecks: compliantChecks,
					},
					&ipb.CheckAlternative{
						FileChecks: compliantChecks,
					},
				},
			},
			expectedCompliantBenchmarks: []*apb.ComplianceResult{&apb.ComplianceResult{
				Id: "id",
				ComplianceOccurrence: &gpb.ComplianceOccurrence{
					NonCompliantFiles: []*gpb.NonCompliantFile{},
				},
			}},
			expectedNonCompliantBenchmarks: []*apb.ComplianceResult{},
			expectedScanStatus:             apb.ScanStatus_SUCCEEDED,
		},
		{
			desc: "One compliant and one non-compliant alternative",
			instructions: &ipb.BenchmarkScanInstruction{
				CheckAlternatives: []*ipb.CheckAlternative{
					&ipb.CheckAlternative{
						FileChecks: compliantChecks,
					},
					&ipb.CheckAlternative{
						FileChecks: nonCompliantChecks,
					},
				},
			},
			expectedCompliantBenchmarks: []*apb.ComplianceResult{&apb.ComplianceResult{
				Id: "id",
				ComplianceOccurrence: &gpb.ComplianceOccurrence{
					NonCompliantFiles: []*gpb.NonCompliantFile{},
				},
			}},
			expectedNonCompliantBenchmarks: []*apb.ComplianceResult{},
			expectedScanStatus:             apb.ScanStatus_SUCCEEDED,
		},
		{
			desc: "Two non-compliant alternatives",
			instructions: &ipb.BenchmarkScanInstruction{
				CheckAlternatives: []*ipb.CheckAlternative{
					&ipb.CheckAlternative{
						FileChecks: nonCompliantChecks,
					},
					&ipb.CheckAlternative{
						FileChecks: nonCompliantChecks2,
					},
				},
			},
			expectedCompliantBenchmarks: []*apb.ComplianceResult{},
			expectedNonCompliantBenchmarks: []*apb.ComplianceResult{&apb.ComplianceResult{
				Id: "id",
				ComplianceOccurrence: &gpb.ComplianceOccurrence{
					NonCompliantFiles: []*gpb.NonCompliantFile{
						&gpb.NonCompliantFile{
							Path:   testFilePath1,
							Reason: fmt.Sprintf("Got content %q, expected \"Different content\"", testFileContent1),
						},
						&gpb.NonCompliantFile{
							Path:   testFilePath2,
							Reason: fmt.Sprintf("Got content %q, expected \"Different content\"", testFileContent2),
						},
					},
				},
			}},
			expectedScanStatus: apb.ScanStatus_SUCCEEDED,
		},
		{
			desc: "One failing, one compliant, one non-compliant alternative",
			instructions: &ipb.BenchmarkScanInstruction{
				CheckAlternatives: []*ipb.CheckAlternative{
					&ipb.CheckAlternative{
						FileChecks: compliantChecks,
					},
					&ipb.CheckAlternative{
						FileChecks: nonCompliantChecks,
					},
					&ipb.CheckAlternative{
						FileChecks: failingChecks,
					},
				},
			},
			expectedCompliantBenchmarks:    []*apb.ComplianceResult{},
			expectedNonCompliantBenchmarks: []*apb.ComplianceResult{},
			expectedScanStatus:             apb.ScanStatus_FAILED,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			config := &apb.ScanConfig{
				BenchmarkConfigs: []*apb.BenchmarkConfig{
					testconfigcreator.NewBenchmarkConfig(t, "id", tc.instructions),
				},
			}
			result, err := scanner.Scanner{}.Scan(context.Background(), config, fakeAPIProvider{})

			if err != nil {
				t.Fatalf("scanner.Scan(%v) had unexpected error: %v", config, err)
			}
			if result.GetStatus().GetStatus() != tc.expectedScanStatus {
				t.Fatalf("expected scan status %v, scanner.Scan(%v) returned: %v",
					tc.expectedScanStatus, config, result.GetStatus().GetStatus())
			}

			if diff := cmp.Diff(tc.expectedCompliantBenchmarks, result.GetCompliantBenchmarks(), protocmp.Transform()); diff != "" {
				t.Errorf("scanner.Scan(%v) returned unexpected compliant files (-want +got):\n%s", config, diff)
			}
			if diff := cmp.Diff(tc.expectedNonCompliantBenchmarks, result.GetNonCompliantBenchmarks(), protocmp.Transform()); diff != "" {
				t.Errorf("scanner.Scan(%v) returned unexpected non-compliant files (-want +got):\n%s", config, diff)
			}
		})
	}
}

func TestDuplicateFindingsRemoved(t *testing.T) {
	nonComplianceMsg := "unexpected content"
	displayCmd := "ls /path/to/test"
	checks := []*ipb.FileCheck{
		&ipb.FileCheck{
			FilesToCheck:       []*ipb.FileSet{testconfigcreator.SingleFileWithPath(testFilePath1)},
			CheckType:          &ipb.FileCheck_Content{Content: &ipb.ContentCheck{Content: "Different content"}},
			NonComplianceMsg:   nonComplianceMsg,
			FileDisplayCommand: displayCmd,
		},
		&ipb.FileCheck{
			FilesToCheck:       []*ipb.FileSet{testconfigcreator.SingleFileWithPath(testFilePath2)},
			CheckType:          &ipb.FileCheck_Content{Content: &ipb.ContentCheck{Content: "Different content"}},
			NonComplianceMsg:   nonComplianceMsg,
			FileDisplayCommand: displayCmd,
		},
	}
	config := &apb.ScanConfig{
		BenchmarkConfigs: []*apb.BenchmarkConfig{
			testconfigcreator.NewBenchmarkConfig(t, "id", testconfigcreator.NewFileScanInstruction(checks)),
		},
	}
	result, err := scanner.Scanner{}.Scan(context.Background(), config, fakeAPIProvider{})

	if err != nil {
		t.Fatalf("scanner.Scan(%v) had unexpected error: %v", config, err)
	}

	// Expect only a single non-compliance finding.
	want := []*apb.ComplianceResult{&apb.ComplianceResult{
		Id: "id",
		ComplianceOccurrence: &gpb.ComplianceOccurrence{
			NonCompliantFiles: []*gpb.NonCompliantFile{&gpb.NonCompliantFile{
				DisplayCommand: displayCmd,
				Reason:         nonComplianceMsg,
			}},
		},
	}}

	if diff := cmp.Diff(want, result.GetNonCompliantBenchmarks(), protocmp.Transform()); diff != "" {
		t.Errorf("scanner.Scan(%v) returned unexpected non-compliant files (-want +got):\n%s", config, diff)
	}
}

func benchmarkConfigsWithVersions(t *testing.T, versions [][]string) []*apb.BenchmarkConfig {
	instruction := testconfigcreator.NewFileScanInstruction([]*ipb.FileCheck{
		&ipb.FileCheck{
			FilesToCheck: []*ipb.FileSet{testconfigcreator.SingleFileWithPath(testFilePath1)},
			CheckType:    &ipb.FileCheck_Content{Content: &ipb.ContentCheck{Content: testFileContent1}},
		},
	})
	benchmarkConfigs := make([]*apb.BenchmarkConfig, 0, len(versions))
	for i, v := range versions {
		versions := make([]*gpb.ComplianceVersion, 0, len(v))
		for _, vv := range v {
			versions = append(versions, &gpb.ComplianceVersion{Version: vv})
		}
		bc := testconfigcreator.NewBenchmarkConfig(t, fmt.Sprintf("id%d", i+1), instruction)
		bc.GetComplianceNote().Version = versions
		benchmarkConfigs = append(benchmarkConfigs, bc)
	}
	return benchmarkConfigs
}

func TestOldestBenchmarkVersionInScanResult(t *testing.T) {
	testCases := []struct {
		desc            string
		versions        [][]string
		expectedVersion string
	}{
		{
			desc:            "basic version test",
			versions:        [][]string{{"1.0.1"}, {"1.1.2"}},
			expectedVersion: "1.0.1",
		},
		{
			desc:            "multiple-digit versions",
			versions:        [][]string{{"2.11.14"}, {"2.11.12"}},
			expectedVersion: "2.11.12",
		},
		{
			desc:            "versions substrings of each other",
			versions:        [][]string{{"2.1"}, {"2.1.1"}},
			expectedVersion: "2.1",
		},
		{
			desc: "benchmarks with multiple versions",
			versions: [][]string{
				{"1.0.0", "1.0.1", "1.1.0"},
				{"1.2.0"},
			},
			expectedVersion: "1.1.0",
		},
		{
			desc:            "empty version",
			versions:        [][]string{{""}, {"2.1.1"}},
			expectedVersion: "0.0.0",
		},
		{
			desc:            "invalid format",
			versions:        [][]string{{"invalid"}, {"2.1.1"}},
			expectedVersion: "0.0.0",
		},
		{
			desc:            "non-numeric version fragment",
			versions:        [][]string{{"1.2.invalid"}, {"1.2.3"}},
			expectedVersion: "0.0.0",
		},
	}

	for _, tc := range testCases {
		config := &apb.ScanConfig{
			BenchmarkConfigs: benchmarkConfigsWithVersions(t, tc.versions),
		}

		result, err := scanner.Scanner{}.Scan(context.Background(), config, fakeAPIProvider{})

		if err != nil {
			t.Fatalf("scanner.Scan(%v) had unexpected error: %v", config, err)
		}
		if result.GetStatus().GetStatus() != apb.ScanStatus_SUCCEEDED {
			t.Fatalf("scanner.Scan(%v) returned unsuccessful scan status: %v",
				config, result.GetStatus().GetStatus())
		}

		if result.GetBenchmarkVersion() != tc.expectedVersion {
			t.Errorf("%v: scanner.Scan(%v) returned benchmark version %s, expected %s",
				tc.desc, config, result.GetBenchmarkVersion(), tc.expectedVersion)
		}
	}
}

func TestFilesInOptOutConfigRedacted(t *testing.T) {
	wrongContent := "Wrong content"
	testCases := []struct {
		desc                     string
		optOutConfig             *apb.OptOutConfig
		expectedNonCompliantFile *gpb.NonCompliantFile
	}{
		{
			desc: "non-compliance reason not displayed",
			optOutConfig: &apb.OptOutConfig{
				ContentOptoutRegexes: []string{regexMatchingTestFiles},
			},
			expectedNonCompliantFile: &gpb.NonCompliantFile{
				Path:   testFilePath1,
				Reason: "[redacted due to opt-out config]",
			},
		},
		{
			desc: "path not displayed",
			optOutConfig: &apb.OptOutConfig{
				FilenameOptoutRegexes: []string{regexMatchingTestFiles},
			},
			expectedNonCompliantFile: &gpb.NonCompliantFile{
				Path:   "[redacted due to opt-out config]",
				Reason: fmt.Sprintf("Got content %q, expected %q", testFileContent1, wrongContent),
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			instruction := testconfigcreator.NewFileScanInstruction([]*ipb.FileCheck{
				&ipb.FileCheck{
					FilesToCheck: []*ipb.FileSet{testconfigcreator.SingleFileWithPath(testFilePath1)},
					CheckType:    &ipb.FileCheck_Content{Content: &ipb.ContentCheck{Content: wrongContent}},
				},
			})
			config := &apb.ScanConfig{
				BenchmarkConfigs: []*apb.BenchmarkConfig{
					testconfigcreator.NewBenchmarkConfig(t, "id", instruction),
				},
				OptOutConfig: tc.optOutConfig,
			}

			result, err := scanner.Scanner{}.Scan(context.Background(), config, fakeAPIProvider{})

			if err != nil {
				t.Fatalf("scanner.Scan(%v) had unexpected error: %v", config, err)
			}
			if result.GetStatus().GetStatus() != apb.ScanStatus_SUCCEEDED {
				t.Fatalf("scanner.Scan(%v) returned unsuccessful scan status: %v",
					config, result.GetStatus().GetStatus())
			}

			want := []*apb.ComplianceResult{&apb.ComplianceResult{
				Id: "id",
				ComplianceOccurrence: &gpb.ComplianceOccurrence{
					NonCompliantFiles: []*gpb.NonCompliantFile{tc.expectedNonCompliantFile},
				},
			}}
			if diff := cmp.Diff(want, result.GetNonCompliantBenchmarks(), protocmp.Transform()); diff != "" {
				t.Errorf("scanner.Scan(%v) returned unexpected non-compliant files (-want +got):\n%s", config, diff)
			}
		})
	}
}
