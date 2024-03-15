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
	cpb "github.com/google/localtoast/scannerlib/proto/compliance_go_proto"
	"github.com/google/localtoast/scannerlib/configchecks"
	apb "github.com/google/localtoast/scannerlib/proto/api_go_proto"
	ipb "github.com/google/localtoast/scannerlib/proto/scan_instructions_go_proto"
	"github.com/google/localtoast/scannerlib/testconfigcreator"
)

func TestGroupCriteria(t *testing.T) {
	testCases := []struct {
		description               string
		fileContent               string
		check                     *ipb.ContentEntryCheck
		expectedError             bool
		expectedNonCompliantFiles []*cpb.NonCompliantFile
	}{
		{
			description: "LESS_THAN does not match if group cannot be parsed",
			fileContent: "VALUE=abc\n",
			check: &ipb.ContentEntryCheck{
				MatchType: ipb.ContentEntryCheck_ALL_MATCH_STRICT_ORDER,
				MatchCriteria: []*ipb.MatchCriterion{&ipb.MatchCriterion{
					FilterRegex:   "VALUE=.*",
					ExpectedRegex: `VALUE=(\w+)`,
					GroupCriteria: []*ipb.GroupCriterion{
						&ipb.GroupCriterion{
							GroupIndex:      1,
							Type:            ipb.GroupCriterion_LESS_THAN,
							ComparisonValue: &ipb.GroupCriterion_Const{Const: 300},
						},
					},
				}},
			},
			expectedError: false,
			expectedNonCompliantFiles: []*cpb.NonCompliantFile{
				&cpb.NonCompliantFile{
					Path:   testFilePath,
					Reason: `File contains entry "VALUE=abc", expected "(?s)^VALUE=(\\w+)$ with group criteria {[group#1 < 300]}"`,
				},
			},
		},
		{
			description: "LESS_THAN matches const",
			fileContent: "VALUE=100\n" +
				"VALUE=200\n",
			check: &ipb.ContentEntryCheck{
				MatchType: ipb.ContentEntryCheck_ALL_MATCH_ANY_ORDER,
				MatchCriteria: []*ipb.MatchCriterion{&ipb.MatchCriterion{
					FilterRegex:   "VALUE=.*",
					ExpectedRegex: `VALUE=(\d+)`,
					GroupCriteria: []*ipb.GroupCriterion{
						&ipb.GroupCriterion{
							GroupIndex:      1,
							Type:            ipb.GroupCriterion_LESS_THAN,
							ComparisonValue: &ipb.GroupCriterion_Const{Const: 300},
						},
					},
				}},
			},
			expectedError:             false,
			expectedNonCompliantFiles: []*cpb.NonCompliantFile{},
		},
		{
			description: "LESS_THAN does not match const",
			fileContent: "VALUE=100\n" +
				"VALUE=200\n",
			check: &ipb.ContentEntryCheck{
				MatchType: ipb.ContentEntryCheck_ALL_MATCH_ANY_ORDER,
				MatchCriteria: []*ipb.MatchCriterion{&ipb.MatchCriterion{
					FilterRegex:   "VALUE=.*",
					ExpectedRegex: `VALUE=(\d+)`,
					GroupCriteria: []*ipb.GroupCriterion{
						&ipb.GroupCriterion{
							GroupIndex:      1,
							Type:            ipb.GroupCriterion_LESS_THAN,
							ComparisonValue: &ipb.GroupCriterion_Const{Const: 100},
						},
					},
				}},
			},
			expectedError: false,
			expectedNonCompliantFiles: []*cpb.NonCompliantFile{
				&cpb.NonCompliantFile{
					Path:   testFilePath,
					Reason: `File contains entry "VALUE=100", expected "(?s)^VALUE=(\\d+)$ with group criteria {[group#1 < 100]}"`,
				},
				&cpb.NonCompliantFile{
					Path:   testFilePath,
					Reason: `File contains entry "VALUE=200", expected "(?s)^VALUE=(\\d+)$ with group criteria {[group#1 < 100]}"`,
				},
			},
		},
		{
			description: "LESS_THAN matches before Today()",
			fileContent: "DAYS_SINCE_EPOCH=0\n",
			check: &ipb.ContentEntryCheck{
				MatchType: ipb.ContentEntryCheck_ALL_MATCH_ANY_ORDER,
				MatchCriteria: []*ipb.MatchCriterion{&ipb.MatchCriterion{
					FilterRegex:   "DAYS_SINCE_EPOCH=.*",
					ExpectedRegex: `DAYS_SINCE_EPOCH=(\d+)`,
					GroupCriteria: []*ipb.GroupCriterion{
						&ipb.GroupCriterion{
							GroupIndex:      1,
							Type:            ipb.GroupCriterion_LESS_THAN,
							ComparisonValue: &ipb.GroupCriterion_Today_{Today: &ipb.GroupCriterion_Today{}},
						},
					},
				}},
			},
			expectedError:             false,
			expectedNonCompliantFiles: []*cpb.NonCompliantFile{},
		},
		{
			description: "LESS_THAN does not match on Today()",
			fileContent: fmt.Sprintf("DAYS_SINCE_EPOCH=%d\n", today),
			check: &ipb.ContentEntryCheck{
				MatchType: ipb.ContentEntryCheck_ALL_MATCH_ANY_ORDER,
				MatchCriteria: []*ipb.MatchCriterion{&ipb.MatchCriterion{
					FilterRegex:   "DAYS_SINCE_EPOCH=.*",
					ExpectedRegex: `DAYS_SINCE_EPOCH=(\d+)`,
					GroupCriteria: []*ipb.GroupCriterion{
						&ipb.GroupCriterion{
							GroupIndex:      1,
							Type:            ipb.GroupCriterion_LESS_THAN,
							ComparisonValue: &ipb.GroupCriterion_Today_{Today: &ipb.GroupCriterion_Today{}},
						},
					},
				}},
			},
			expectedError: false,
			expectedNonCompliantFiles: []*cpb.NonCompliantFile{
				&cpb.NonCompliantFile{
					Path:   testFilePath,
					Reason: fmt.Sprintf(`File contains entry "DAYS_SINCE_EPOCH=%d", expected "(?s)^DAYS_SINCE_EPOCH=(\\d+)$ with group criteria {[group#1 < today]}"`, today),
				},
			},
		},
		{
			description: "LESS_THAN does not match after Today()",
			fileContent: fmt.Sprintf("DAYS_SINCE_EPOCH=%d\n", today+1),
			check: &ipb.ContentEntryCheck{
				MatchType: ipb.ContentEntryCheck_ALL_MATCH_ANY_ORDER,
				MatchCriteria: []*ipb.MatchCriterion{&ipb.MatchCriterion{
					FilterRegex:   "DAYS_SINCE_EPOCH=.*",
					ExpectedRegex: `DAYS_SINCE_EPOCH=(\d+)`,
					GroupCriteria: []*ipb.GroupCriterion{
						&ipb.GroupCriterion{
							GroupIndex:      1,
							Type:            ipb.GroupCriterion_LESS_THAN,
							ComparisonValue: &ipb.GroupCriterion_Today_{Today: &ipb.GroupCriterion_Today{}},
						},
					},
				}},
			},
			expectedError: false,
			expectedNonCompliantFiles: []*cpb.NonCompliantFile{
				&cpb.NonCompliantFile{
					Path:   testFilePath,
					Reason: fmt.Sprintf(`File contains entry "DAYS_SINCE_EPOCH=%d", expected "(?s)^DAYS_SINCE_EPOCH=(\\d+)$ with group criteria {[group#1 < today]}"`, today+1),
				},
			},
		},
		{
			description: "GREATER_THAN does not match if group cannot be parsed",
			fileContent: "VALUE=abc\n",
			check: &ipb.ContentEntryCheck{
				MatchType: ipb.ContentEntryCheck_ALL_MATCH_STRICT_ORDER,
				MatchCriteria: []*ipb.MatchCriterion{&ipb.MatchCriterion{
					FilterRegex:   "VALUE=.*",
					ExpectedRegex: `VALUE=(\w+)`,
					GroupCriteria: []*ipb.GroupCriterion{
						&ipb.GroupCriterion{
							GroupIndex:      1,
							Type:            ipb.GroupCriterion_GREATER_THAN,
							ComparisonValue: &ipb.GroupCriterion_Const{Const: 0},
						},
					},
				}},
			},
			expectedError: false,
			expectedNonCompliantFiles: []*cpb.NonCompliantFile{
				&cpb.NonCompliantFile{
					Path:   testFilePath,
					Reason: `File contains entry "VALUE=abc", expected "(?s)^VALUE=(\\w+)$ with group criteria {[group#1 > 0]}"`,
				},
			},
		},
		{
			description: "GREATER_THAN matches const",
			fileContent: "VALUE=100\n" +
				"VALUE=200\n",
			check: &ipb.ContentEntryCheck{
				MatchType: ipb.ContentEntryCheck_ALL_MATCH_ANY_ORDER,
				MatchCriteria: []*ipb.MatchCriterion{&ipb.MatchCriterion{
					FilterRegex:   "VALUE=.*",
					ExpectedRegex: `VALUE=(\d+)`,
					GroupCriteria: []*ipb.GroupCriterion{
						&ipb.GroupCriterion{
							GroupIndex:      1,
							Type:            ipb.GroupCriterion_GREATER_THAN,
							ComparisonValue: &ipb.GroupCriterion_Const{Const: 0},
						},
					},
				}},
			},
			expectedError:             false,
			expectedNonCompliantFiles: []*cpb.NonCompliantFile{},
		},
		{
			description: "GREATER_THAN does not match const",
			fileContent: "VALUE=100\n" +
				"VALUE=200\n",
			check: &ipb.ContentEntryCheck{
				MatchType: ipb.ContentEntryCheck_ALL_MATCH_ANY_ORDER,
				MatchCriteria: []*ipb.MatchCriterion{&ipb.MatchCriterion{
					FilterRegex:   "VALUE=.*",
					ExpectedRegex: `VALUE=(\d+)`,
					GroupCriteria: []*ipb.GroupCriterion{
						&ipb.GroupCriterion{
							GroupIndex:      1,
							Type:            ipb.GroupCriterion_GREATER_THAN,
							ComparisonValue: &ipb.GroupCriterion_Const{Const: 200},
						},
					},
				}},
			},
			expectedError: false,
			expectedNonCompliantFiles: []*cpb.NonCompliantFile{
				&cpb.NonCompliantFile{
					Path:   testFilePath,
					Reason: `File contains entry "VALUE=100", expected "(?s)^VALUE=(\\d+)$ with group criteria {[group#1 > 200]}"`,
				},
				&cpb.NonCompliantFile{
					Path:   testFilePath,
					Reason: `File contains entry "VALUE=200", expected "(?s)^VALUE=(\\d+)$ with group criteria {[group#1 > 200]}"`,
				},
			},
		},
		{
			description: "GREATER_THAN matches after Today()",
			fileContent: fmt.Sprintf("DAYS_SINCE_EPOCH=%d\n", today+1),
			check: &ipb.ContentEntryCheck{
				MatchType: ipb.ContentEntryCheck_ALL_MATCH_ANY_ORDER,
				MatchCriteria: []*ipb.MatchCriterion{&ipb.MatchCriterion{
					FilterRegex:   "DAYS_SINCE_EPOCH=.*",
					ExpectedRegex: `DAYS_SINCE_EPOCH=(\d+)`,
					GroupCriteria: []*ipb.GroupCriterion{
						&ipb.GroupCriterion{
							GroupIndex:      1,
							Type:            ipb.GroupCriterion_GREATER_THAN,
							ComparisonValue: &ipb.GroupCriterion_Today_{Today: &ipb.GroupCriterion_Today{}},
						},
					},
				}},
			},
			expectedError:             false,
			expectedNonCompliantFiles: []*cpb.NonCompliantFile{},
		},
		{
			description: "GREATER_THAN does not match on Today()",
			fileContent: fmt.Sprintf("DAYS_SINCE_EPOCH=%d\n", today),
			check: &ipb.ContentEntryCheck{
				MatchType: ipb.ContentEntryCheck_ALL_MATCH_ANY_ORDER,
				MatchCriteria: []*ipb.MatchCriterion{&ipb.MatchCriterion{
					FilterRegex:   "DAYS_SINCE_EPOCH=.*",
					ExpectedRegex: `DAYS_SINCE_EPOCH=(\d+)`,
					GroupCriteria: []*ipb.GroupCriterion{
						&ipb.GroupCriterion{
							GroupIndex:      1,
							Type:            ipb.GroupCriterion_GREATER_THAN,
							ComparisonValue: &ipb.GroupCriterion_Today_{Today: &ipb.GroupCriterion_Today{}},
						},
					},
				}},
			},
			expectedError: false,
			expectedNonCompliantFiles: []*cpb.NonCompliantFile{
				&cpb.NonCompliantFile{
					Path:   testFilePath,
					Reason: fmt.Sprintf(`File contains entry "DAYS_SINCE_EPOCH=%d", expected "(?s)^DAYS_SINCE_EPOCH=(\\d+)$ with group criteria {[group#1 > today]}"`, today),
				},
			},
		},
		{
			description: "GREATER_THAN does not match before Today()",
			fileContent: "DAYS_SINCE_EPOCH=0\n",
			check: &ipb.ContentEntryCheck{
				MatchType: ipb.ContentEntryCheck_ALL_MATCH_ANY_ORDER,
				MatchCriteria: []*ipb.MatchCriterion{&ipb.MatchCriterion{
					FilterRegex:   "DAYS_SINCE_EPOCH=.*",
					ExpectedRegex: `DAYS_SINCE_EPOCH=(\d+)`,
					GroupCriteria: []*ipb.GroupCriterion{
						&ipb.GroupCriterion{
							GroupIndex:      1,
							Type:            ipb.GroupCriterion_GREATER_THAN,
							ComparisonValue: &ipb.GroupCriterion_Today_{Today: &ipb.GroupCriterion_Today{}},
						},
					},
				}},
			},
			expectedError: false,
			expectedNonCompliantFiles: []*cpb.NonCompliantFile{
				&cpb.NonCompliantFile{
					Path:   testFilePath,
					Reason: `File contains entry "DAYS_SINCE_EPOCH=0", expected "(?s)^DAYS_SINCE_EPOCH=(\\d+)$ with group criteria {[group#1 > today]}"`,
				},
			},
		},
		{
			description: "NO_LESS_RESTRICTIVE_UMASK requires const",
			fileContent: "",
			check: &ipb.ContentEntryCheck{
				MatchType: ipb.ContentEntryCheck_ALL_MATCH_ANY_ORDER,
				MatchCriteria: []*ipb.MatchCriterion{&ipb.MatchCriterion{
					FilterRegex:   "UMASK=.*",
					ExpectedRegex: `UMASK=(.*)`,
					GroupCriteria: []*ipb.GroupCriterion{
						&ipb.GroupCriterion{
							GroupIndex:      1,
							Type:            ipb.GroupCriterion_NO_LESS_RESTRICTIVE_UMASK,
							ComparisonValue: &ipb.GroupCriterion_Today_{Today: &ipb.GroupCriterion_Today{}},
						},
					},
				}},
			},
			expectedError: true,
		},
		{
			description: "NO_LESS_RESTRICTIVE_UMASK does not match non-umask",
			fileContent: "UMASK=888\n",
			check: &ipb.ContentEntryCheck{
				MatchType: ipb.ContentEntryCheck_ALL_MATCH_STRICT_ORDER,
				MatchCriteria: []*ipb.MatchCriterion{
					&ipb.MatchCriterion{
						FilterRegex:   "UMASK=.*",
						ExpectedRegex: "UMASK=(.*)",
						GroupCriteria: []*ipb.GroupCriterion{
							&ipb.GroupCriterion{
								GroupIndex:      1,
								Type:            ipb.GroupCriterion_NO_LESS_RESTRICTIVE_UMASK,
								ComparisonValue: &ipb.GroupCriterion_Const{Const: 027},
							},
						},
					},
				},
			},
			expectedError: false,
			expectedNonCompliantFiles: []*cpb.NonCompliantFile{
				&cpb.NonCompliantFile{
					Path:   testFilePath,
					Reason: `File contains entry "UMASK=888", expected "(?s)^UMASK=(.*)$ with group criteria {[group#1 not less restrictive than 23]}"`,
				},
			},
		},
		{
			description: "NO_LESS_RESTRICTIVE_UMASK does not match",
			fileContent: "UMASK=000\n",
			check: &ipb.ContentEntryCheck{
				MatchType: ipb.ContentEntryCheck_ALL_MATCH_STRICT_ORDER,
				MatchCriteria: []*ipb.MatchCriterion{
					&ipb.MatchCriterion{
						FilterRegex:   "UMASK=.*",
						ExpectedRegex: "UMASK=(.*)",
						GroupCriteria: []*ipb.GroupCriterion{
							&ipb.GroupCriterion{
								GroupIndex:      1,
								Type:            ipb.GroupCriterion_NO_LESS_RESTRICTIVE_UMASK,
								ComparisonValue: &ipb.GroupCriterion_Const{Const: 027},
							},
						},
					},
				},
			},
			expectedError: false,
			expectedNonCompliantFiles: []*cpb.NonCompliantFile{
				&cpb.NonCompliantFile{
					Path:   testFilePath,
					Reason: `File contains entry "UMASK=000", expected "(?s)^UMASK=(.*)$ with group criteria {[group#1 not less restrictive than 23]}"`,
				},
			},
		},
		{
			description: "NO_LESS_RESTRICTIVE_UMASK matches exact",
			fileContent: "UMASK=027\n",
			check: &ipb.ContentEntryCheck{
				MatchType: ipb.ContentEntryCheck_ALL_MATCH_ANY_ORDER,
				MatchCriteria: []*ipb.MatchCriterion{
					&ipb.MatchCriterion{
						FilterRegex:   "UMASK=.*",
						ExpectedRegex: "UMASK=(.*)",
						GroupCriteria: []*ipb.GroupCriterion{
							&ipb.GroupCriterion{
								GroupIndex:      1,
								Type:            ipb.GroupCriterion_NO_LESS_RESTRICTIVE_UMASK,
								ComparisonValue: &ipb.GroupCriterion_Const{Const: 027},
							},
						},
					},
				},
			},
			expectedError:             false,
			expectedNonCompliantFiles: []*cpb.NonCompliantFile{},
		},
		{
			description: "NO_LESS_RESTRICTIVE_UMASK matches",
			fileContent: "UMASK=0077\n",
			check: &ipb.ContentEntryCheck{
				MatchType: ipb.ContentEntryCheck_ALL_MATCH_ANY_ORDER,
				MatchCriteria: []*ipb.MatchCriterion{
					&ipb.MatchCriterion{
						FilterRegex:   "UMASK=.*",
						ExpectedRegex: "UMASK=(.*)",
						GroupCriteria: []*ipb.GroupCriterion{
							&ipb.GroupCriterion{
								GroupIndex:      1,
								Type:            ipb.GroupCriterion_NO_LESS_RESTRICTIVE_UMASK,
								ComparisonValue: &ipb.GroupCriterion_Const{Const: 027},
							},
						},
					},
				},
			},
			expectedError:             false,
			expectedNonCompliantFiles: []*cpb.NonCompliantFile{},
		},
		{
			description: "UNIQUE does not match",
			fileContent: "USER=sergey\n" +
				"USER=larry\n" +
				"USER=sergey\n",
			check: &ipb.ContentEntryCheck{
				MatchType: ipb.ContentEntryCheck_ALL_MATCH_ANY_ORDER,
				MatchCriteria: []*ipb.MatchCriterion{&ipb.MatchCriterion{
					FilterRegex:   "USER=.*",
					ExpectedRegex: `USER=(.*)`,
					GroupCriteria: []*ipb.GroupCriterion{
						&ipb.GroupCriterion{
							GroupIndex: 1,
							Type:       ipb.GroupCriterion_UNIQUE,
						},
					},
				}},
			},
			expectedNonCompliantFiles: []*cpb.NonCompliantFile{
				&cpb.NonCompliantFile{
					Path:   testFilePath,
					Reason: `File contains entry "USER=sergey", expected "(?s)^USER=(.*)$ with group criteria {[group#1 is unique]}"`,
				},
			},
		},
		{
			description: "UNIQUE matches",
			fileContent: "USER=sergey\n" +
				"USER=larry\n" +
				"USER=sundar\n",
			check: &ipb.ContentEntryCheck{
				MatchType: ipb.ContentEntryCheck_ALL_MATCH_ANY_ORDER,
				MatchCriteria: []*ipb.MatchCriterion{&ipb.MatchCriterion{
					FilterRegex:   "USER=.*",
					ExpectedRegex: `USER=(.*)`,
					GroupCriteria: []*ipb.GroupCriterion{
						&ipb.GroupCriterion{
							GroupIndex: 1,
							Type:       ipb.GroupCriterion_UNIQUE,
						},
					},
				}},
			},
		},
		{
			description: "index out of bounds",
			fileContent: "",
			check: &ipb.ContentEntryCheck{
				MatchType: ipb.ContentEntryCheck_ALL_MATCH_ANY_ORDER,
				MatchCriteria: []*ipb.MatchCriterion{&ipb.MatchCriterion{
					FilterRegex:   "USER=.*",
					ExpectedRegex: `USER=(.*)`,
					GroupCriteria: []*ipb.GroupCriterion{
						&ipb.GroupCriterion{
							GroupIndex: 12,
							Type:       ipb.GroupCriterion_UNIQUE,
						},
					},
				}},
			},
			expectedError: true,
		},
		{
			description: "UNIQUE and NONE_MATCH are incompatible",
			fileContent: "",
			check: &ipb.ContentEntryCheck{
				MatchType: ipb.ContentEntryCheck_NONE_MATCH,
				MatchCriteria: []*ipb.MatchCriterion{&ipb.MatchCriterion{
					FilterRegex:   "USER=.*",
					ExpectedRegex: `USER=(.*)`,
					GroupCriteria: []*ipb.GroupCriterion{
						&ipb.GroupCriterion{
							GroupIndex: 1,
							Type:       ipb.GroupCriterion_UNIQUE,
						},
					},
				}},
			},
			expectedError: true,
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

			if tc.expectedError {
				if err == nil {
					t.Errorf("CreateChecksFromConfig([%v]) did not return an error", config)
				}
				return
			}

			if err != nil {
				t.Fatalf("CreateChecksFromConfig([%v]) returned an error: %v", config, err)
			}
			if len(checks) != 1 {
				t.Fatalf("Created %d checks, expected only 1", len(checks))
			}

			var pVal string
			resultMap, _, err := checks[0].Exec(pVal)
			if err != nil {
				t.Fatalf("check.Exec() returned an error: %v", err)
			}
			result, gotSingleton := singleComplianceResult(resultMap)
			if !gotSingleton {
				t.Fatalf("check.Exec() expected to return 1 result, got %d", len(resultMap))
			}

			want := &apb.ComplianceResult{
				Id: "id",
				ComplianceOccurrence: &cpb.ComplianceOccurrence{
					NonCompliantFiles: tc.expectedNonCompliantFiles,
				},
			}
			if diff := cmp.Diff(want, result, protocmp.Transform()); diff != "" {
				t.Errorf("check.Exec() returned unexpected diff (-want +got):\n%s", diff)
			}
		})
	}
}
