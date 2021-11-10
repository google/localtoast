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

package configchecks_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	"google.golang.org/protobuf/testing/protocmp"
	gpb "google.golang.org/genproto/googleapis/grafeas/v1"
	"github.com/google/localtoast/library/configchecks"
	apb "github.com/google/localtoast/library/proto/api_go_proto"
	ipb "github.com/google/localtoast/library/proto/scan_instructions_go_proto"
	"github.com/google/localtoast/library/testing/testconfigcreator"
)

func TestFileContentEntryCheckComplianceResults(t *testing.T) {
	testCases := []struct {
		description               string
		fileContent               string
		check                     *ipb.ContentEntryCheck
		expectedNonCompliantFiles []*gpb.NonCompliantFile
	}{
		{
			description: "single criterion matches",
			fileContent: "VALUE1=true\n" +
				"VALUE2=false\n" +
				"VALUE1=true",
			check: &ipb.ContentEntryCheck{
				MatchType: ipb.ContentEntryCheck_ALL_MATCH_ANY_ORDER,
				MatchCriteria: []*ipb.MatchCriterion{{
					FilterRegex:   "VALUE1=.*",
					ExpectedRegex: "VALUE1=true",
				}},
			},
			expectedNonCompliantFiles: nil,
		},
		{
			description: "single criterion doesn't match",
			fileContent: "VALUE1=true\n" +
				"VALUE2=false",
			check: &ipb.ContentEntryCheck{
				MatchType: ipb.ContentEntryCheck_ALL_MATCH_ANY_ORDER,
				MatchCriteria: []*ipb.MatchCriterion{{
					FilterRegex:   "VALUE2=.*",
					ExpectedRegex: "VALUE2=true",
				}},
			},
			expectedNonCompliantFiles: []*gpb.NonCompliantFile{
				&gpb.NonCompliantFile{
					Path:   testFilePath,
					Reason: "File contains entry \"VALUE2=false\", expected \"(?s)^VALUE2=true$\"",
				},
			},
		},
		{
			description: "single criterion doesn't always match",
			fileContent: "VALUE1=true\n" +
				"VALUE1=false",
			check: &ipb.ContentEntryCheck{
				MatchType: ipb.ContentEntryCheck_ALL_MATCH_ANY_ORDER,
				MatchCriteria: []*ipb.MatchCriterion{{
					FilterRegex:   "VALUE1=.*",
					ExpectedRegex: "VALUE1=true",
				}},
			},
			expectedNonCompliantFiles: []*gpb.NonCompliantFile{
				&gpb.NonCompliantFile{
					Path:   testFilePath,
					Reason: "File contains entry \"VALUE1=false\", expected \"(?s)^VALUE1=true$\"",
				},
			},
		},
		{
			description: "criterion not found among files",
			fileContent: "VALUE1=true",
			check: &ipb.ContentEntryCheck{
				MatchType: ipb.ContentEntryCheck_ALL_MATCH_ANY_ORDER,
				MatchCriteria: []*ipb.MatchCriterion{{
					FilterRegex:   "VALUE2=.*",
					ExpectedRegex: "VALUE2=true",
				}},
			},
			expectedNonCompliantFiles: []*gpb.NonCompliantFile{
				&gpb.NonCompliantFile{
					Path:   fmt.Sprintf("single_file:{path:%q}", testFilePath),
					Reason: "No entry matching \"(?s)^VALUE2=.*$\" found among files",
				},
			},
		},
		{
			description: "strict order criteria matched in order",
			fileContent: "VALUE1=true\n" +
				"VALUE2=true\n" +
				"VALUE3=true\n",
			check: &ipb.ContentEntryCheck{
				MatchType: ipb.ContentEntryCheck_ALL_MATCH_STRICT_ORDER,
				MatchCriteria: []*ipb.MatchCriterion{{
					FilterRegex:   "VALUE1=.*",
					ExpectedRegex: "VALUE1=true",
				}, {
					FilterRegex:   "VALUE3=.*",
					ExpectedRegex: "VALUE3=true",
				}},
			},
			expectedNonCompliantFiles: nil,
		},
		{
			description: "strict order criteria matched out of order",
			fileContent: "VALUE1=true\n" +
				"VALUE2=true\n" +
				"VALUE3=true\n",
			check: &ipb.ContentEntryCheck{
				MatchType: ipb.ContentEntryCheck_ALL_MATCH_STRICT_ORDER,
				MatchCriteria: []*ipb.MatchCriterion{{
					FilterRegex:   "VALUE3=.*",
					ExpectedRegex: "VALUE3=true",
				}, {
					FilterRegex:   "VALUE1=.*",
					ExpectedRegex: "VALUE1=true",
				}},
			},
			expectedNonCompliantFiles: []*gpb.NonCompliantFile{
				&gpb.NonCompliantFile{
					Path:   testFilePath,
					Reason: "Criteria expected to match in order but file entry \"VALUE1=true\", matched \"(?s)^VALUE1=true$\" before \"(?s)^VALUE3=true$\" was matched",
				},
			},
		},
		{
			description: "strict order criteria matched a second time",
			fileContent: "VALUE1=true\n" +
				"VALUE2=true\n" +
				"VALUE1=true\n",
			check: &ipb.ContentEntryCheck{
				MatchType: ipb.ContentEntryCheck_ALL_MATCH_STRICT_ORDER,
				MatchCriteria: []*ipb.MatchCriterion{{
					FilterRegex:   "VALUE1=.*",
					ExpectedRegex: "VALUE1=true",
				}, {
					FilterRegex:   "VALUE2=.*",
					ExpectedRegex: "VALUE2=true",
				}},
			},
			expectedNonCompliantFiles: []*gpb.NonCompliantFile{
				&gpb.NonCompliantFile{
					Path:   testFilePath,
					Reason: "Criteria expected to match in order but file entry \"VALUE1=true\", matched \"(?s)^VALUE1=true$\" after \"(?s)^VALUE2=true$\" was matched",
				},
			},
		},
		{
			description: "nothing should match, nothing does",
			fileContent: "VALUE1=true\n" +
				"VALUE2=true",
			check: &ipb.ContentEntryCheck{
				MatchType: ipb.ContentEntryCheck_NONE_MATCH,
				MatchCriteria: []*ipb.MatchCriterion{{
					FilterRegex:   "VALUE1=.*",
					ExpectedRegex: "VALUE1=false",
				}},
			},
			expectedNonCompliantFiles: nil,
		},
		{
			description: "nothing should match, something does",
			fileContent: "VALUE1=true\n" +
				"VALUE2=true",
			check: &ipb.ContentEntryCheck{
				MatchType: ipb.ContentEntryCheck_NONE_MATCH,
				MatchCriteria: []*ipb.MatchCriterion{{
					FilterRegex:   "VALUE1=.*",
					ExpectedRegex: "VALUE1=true",
				}},
			},
			expectedNonCompliantFiles: []*gpb.NonCompliantFile{
				&gpb.NonCompliantFile{
					Path:   testFilePath,
					Reason: "File contains entry \"VALUE1=true\", didn't expect any entries matching \"(?s)^VALUE1=true$\"",
				},
			},
		},
		{
			description: "split by other delimiter",
			fileContent: "VALUE1=true\t" +
				"VALUE2=true",
			check: &ipb.ContentEntryCheck{
				Delimiter: []byte{'\t'},
				MatchType: ipb.ContentEntryCheck_ALL_MATCH_ANY_ORDER,
				MatchCriteria: []*ipb.MatchCriterion{{
					FilterRegex:   "VALUE1=.*",
					ExpectedRegex: "VALUE1=true",
				}},
			},
			expectedNonCompliantFiles: nil,
		},
		{
			description: "split by other delimiter with trailing",
			fileContent: "VALUE1=true\t" +
				"VALUE2=true\t",
			check: &ipb.ContentEntryCheck{
				Delimiter: []byte{'\t'},
				MatchType: ipb.ContentEntryCheck_ALL_MATCH_ANY_ORDER,
				MatchCriteria: []*ipb.MatchCriterion{{
					FilterRegex:   "VALUE1=.*",
					ExpectedRegex: "VALUE1=true",
				}},
			},
			expectedNonCompliantFiles: nil,
		},
		{
			description: "match across lines",
			fileContent: "VALUE1=true\nVALUE2=true",
			check: &ipb.ContentEntryCheck{
				Delimiter: []byte{0},
				MatchType: ipb.ContentEntryCheck_ALL_MATCH_ANY_ORDER,
				MatchCriteria: []*ipb.MatchCriterion{{
					FilterRegex:   "VALUE1=.*\nVALUE2=.*",
					ExpectedRegex: "VALUE1=true\nVALUE2=true",
				}},
			},
			expectedNonCompliantFiles: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			scanInstruction := testconfigcreator.NewFileScanInstruction(
				[]*ipb.FileCheck{&ipb.FileCheck{
					FilesToCheck: []*ipb.FileSet{testconfigcreator.SingleFileWithPath(testFilePath)},
					CheckType:    &ipb.FileCheck_ContentEntry{ContentEntry: tc.check},
				}})
			config := testconfigcreator.NewBenchmarkConfig(t, "id", scanInstruction)
			checks, err := configchecks.CreateChecksFromConfig(
				context.Background(),
				&apb.ScanConfig{
					BenchmarkConfigs: []*apb.BenchmarkConfig{config},
				},
				newFakeAPI(withFileContent(tc.fileContent)),
			)
			if err != nil {
				t.Fatalf("CreateChecksFromConfig([%v]) returned an error: %v", config, err)
			}
			if len(checks) != 1 {
				t.Fatalf("Created %d checks, expected only 1", len(checks))
			}

			resultMap, err := checks[0].Exec()
			if err != nil {
				t.Fatalf("checks[0].Exec() returned an error: %v", err)
			}
			result, gotSingleton := singleComplianceResult(resultMap)
			if !gotSingleton {
				t.Fatalf("checks[0].Exec() expected to return 1 result, got %d", len(resultMap))
			}

			want := &apb.ComplianceResult{
				Id: "id",
				ComplianceOccurrence: &gpb.ComplianceOccurrence{
					NonCompliantFiles: tc.expectedNonCompliantFiles,
				},
			}
			if diff := cmp.Diff(want, result, protocmp.Transform()); diff != "" {
				t.Errorf("check.Exec() returned unexpected diff (-want +got):\n%s", diff)
			}
		})
	}
}

func TestFileContentEntryFileDoesntExist(t *testing.T) {
	testCases := []struct {
		name                      string
		check                     *ipb.ContentEntryCheck
		expectedNonCompliantFiles []*gpb.NonCompliantFile
	}{
		{
			name: "all_match_criterion_is_non_compliant",
			check: &ipb.ContentEntryCheck{
				MatchType: ipb.ContentEntryCheck_ALL_MATCH_ANY_ORDER,
				MatchCriteria: []*ipb.MatchCriterion{{
					FilterRegex:   "VALUE1=.*",
					ExpectedRegex: "VALUE1=true",
				}},
			},
			expectedNonCompliantFiles: []*gpb.NonCompliantFile{
				&gpb.NonCompliantFile{
					Path:   nonExistentFilePath,
					Reason: "File doesn't exist",
				},
				&gpb.NonCompliantFile{
					Path:   fmt.Sprintf("single_file:{path:%q}", nonExistentFilePath),
					Reason: "No entry matching \"(?s)^VALUE1=.*$\" found among files",
				},
			},
		},
		{
			name: "none_match_criterion_is_compliant",
			check: &ipb.ContentEntryCheck{
				MatchType: ipb.ContentEntryCheck_NONE_MATCH,
				MatchCriteria: []*ipb.MatchCriterion{{
					FilterRegex:   "VALUE1=.*",
					ExpectedRegex: "VALUE1=true",
				}},
			},
			expectedNonCompliantFiles: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			scanInstruction := testconfigcreator.NewFileScanInstruction(
				[]*ipb.FileCheck{&ipb.FileCheck{
					FilesToCheck: []*ipb.FileSet{testconfigcreator.SingleFileWithPath(nonExistentFilePath)},
					CheckType:    &ipb.FileCheck_ContentEntry{ContentEntry: tc.check},
				}})
			config := testconfigcreator.NewBenchmarkConfig(t, "id", scanInstruction)
			checks, err := configchecks.CreateChecksFromConfig(
				context.Background(),
				&apb.ScanConfig{
					BenchmarkConfigs: []*apb.BenchmarkConfig{config},
				},
				newFakeAPI(),
			)
			if err != nil {
				t.Fatalf("CreateChecksFromConfig([%v]) returned an error: %v", config, err)
			}
			if len(checks) != 1 {
				t.Fatalf("Created %d checks, want only 1", len(checks))
			}

			resultMap, err := checks[0].Exec()
			if err != nil {
				t.Fatalf("checks[0].Exec() returned an error: %v", err)
			}
			result, gotSingleton := singleComplianceResult(resultMap)
			if !gotSingleton {
				t.Fatalf("checks[0].Exec() got %d results, want 1 result", len(resultMap))
			}

			want := &apb.ComplianceResult{
				Id: "id",
				ComplianceOccurrence: &gpb.ComplianceOccurrence{
					NonCompliantFiles: tc.expectedNonCompliantFiles,
				},
			}
			if diff := cmp.Diff(want, result, protocmp.Transform()); diff != "" {
				t.Errorf("check.Exec() returned unexpected diff (-want +got):\n%s", diff)
			}
		})
	}
}

func TestFileContentEntryCheckOnDirectory(t *testing.T) {
	testContentEntry := &ipb.ContentEntryCheck{
		MatchType: ipb.ContentEntryCheck_ALL_MATCH_ANY_ORDER,
		MatchCriteria: []*ipb.MatchCriterion{{
			FilterRegex:   "VALUE1=.*",
			ExpectedRegex: "VALUE1=true",
		}},
	}
	directory := &ipb.FileSet{
		FilePath: &ipb.FileSet_FilesInDir_{FilesInDir: &ipb.FileSet_FilesInDir{DirPath: testDirPath}},
	}
	scanInstruction := testconfigcreator.NewFileScanInstruction(
		[]*ipb.FileCheck{&ipb.FileCheck{
			FilesToCheck: []*ipb.FileSet{directory},
			CheckType:    &ipb.FileCheck_ContentEntry{ContentEntry: testContentEntry},
		}})

	config := testconfigcreator.NewBenchmarkConfig(t, "id", scanInstruction)
	checks, err := configchecks.CreateChecksFromConfig(
		context.Background(),
		&apb.ScanConfig{
			BenchmarkConfigs: []*apb.BenchmarkConfig{config},
		},
		newFakeAPI(withFileContent("VALUE1=true")),
	)
	if err != nil {
		t.Fatalf("CreateChecksFromConfig([%v]) returned an error: %v", config, err)
	}
	if len(checks) != 1 {
		t.Fatalf("Created %d checks, expected only 1", len(checks))
	}

	resultMap, err := checks[0].Exec()
	if err != nil {
		t.Fatalf("checks[0].Exec() returned an error: %v", err)
	}
	result, gotSingleton := singleComplianceResult(resultMap)
	if !gotSingleton {
		t.Fatalf("checks[0].Exec() expected to return 1 result, got %d", len(resultMap))
	}

	expected := &apb.ComplianceResult{
		Id: "id",
		ComplianceOccurrence: &gpb.ComplianceOccurrence{
			// The check passes if one of the two files in the dir had a matching entry.
			NonCompliantFiles: nil,
		},
	}
	if diff := cmp.Diff(expected, result, protocmp.Transform()); diff != "" {
		t.Errorf("check.Exec() returned unexpected diff (-want +got):\n%s", diff)
	}
}
