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
	"bytes"
	"compress/gzip"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"sort"
	"strings"
	"testing"

	dpb "google.golang.org/protobuf/types/known/durationpb"
	"github.com/google/go-cmp/cmp"
	"google.golang.org/protobuf/testing/protocmp"
	cpb "github.com/google/localtoast/scannerlib/proto/compliance_go_proto"
	"github.com/google/localtoast/scanapi"
	"github.com/google/localtoast/scannerlib/configchecks"
	"github.com/google/localtoast/scannerlib/fileset"
	apb "github.com/google/localtoast/scannerlib/proto/api_go_proto"
	ipb "github.com/google/localtoast/scannerlib/proto/scan_instructions_go_proto"
	"github.com/google/localtoast/scannerlib/testconfigcreator"
)

func TestChecksOnSameFileGroupedTogether(t *testing.T) {
	testCases := []struct {
		desc       string
		fileCheck1 *ipb.FileCheck
		fileCheck2 *ipb.FileCheck
	}{
		{
			desc: "Same check types",
			fileCheck1: &ipb.FileCheck{
				FilesToCheck: []*ipb.FileSet{testconfigcreator.SingleFileWithPath("/path")},
				CheckType:    &ipb.FileCheck_Existence{Existence: &ipb.ExistenceCheck{ShouldExist: true}},
			},
			fileCheck2: &ipb.FileCheck{
				FilesToCheck: []*ipb.FileSet{testconfigcreator.SingleFileWithPath("/path")},
				CheckType:    &ipb.FileCheck_Existence{Existence: &ipb.ExistenceCheck{ShouldExist: false}},
			},
		},
		{
			desc: "Different check types",
			fileCheck1: &ipb.FileCheck{
				FilesToCheck: []*ipb.FileSet{testconfigcreator.SingleFileWithPath("/path")},
				CheckType:    &ipb.FileCheck_Existence{Existence: &ipb.ExistenceCheck{ShouldExist: true}},
			},
			fileCheck2: &ipb.FileCheck{
				FilesToCheck: []*ipb.FileSet{testconfigcreator.SingleFileWithPath("/path")},
				CheckType:    &ipb.FileCheck_Content{Content: &ipb.ContentCheck{Content: "content"}},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			scanInstruction1 := testconfigcreator.NewFileScanInstruction([]*ipb.FileCheck{tc.fileCheck1})
			scanInstruction2 := testconfigcreator.NewFileScanInstruction([]*ipb.FileCheck{tc.fileCheck2})
			config1 := testconfigcreator.NewBenchmarkConfig(t, "id1", scanInstruction1)
			config2 := testconfigcreator.NewBenchmarkConfig(t, "id2", scanInstruction2)
			checks, err := configchecks.CreateChecksFromConfig(
				context.Background(),
				&apb.ScanConfig{
					BenchmarkConfigs: []*apb.BenchmarkConfig{config1, config2},
				},
				&fakeAPI{fileContent: testFileContent},
			)
			if err != nil {
				t.Fatalf("configchecks.CreateChecksFromConfig([%v %v]) returned an error: %v", config1, config2, err)
			}
			if len(checks) != 1 {
				t.Fatalf("Expected 1 check to be created, got %d", len(checks))
			}

			expectedIDs := []string{"id1", "id2"}
			actualIDs := checks[0].BenchmarkIDs()
			sort.Strings(actualIDs)
			if diff := cmp.Diff(expectedIDs, actualIDs); diff != "" {
				t.Errorf("%v.BenchmarkIDs() returned unexpected diff (-want +got):\n%s", checks[0], diff)
			}
		})
	}
}

func TestSameChecksOnDifferentAlternativesGroupedTogether(t *testing.T) {
	check := &ipb.FileCheck{
		FilesToCheck: []*ipb.FileSet{testconfigcreator.SingleFileWithPath("/path")},
		CheckType:    &ipb.FileCheck_Existence{Existence: &ipb.ExistenceCheck{ShouldExist: true}},
	}
	instruction := &ipb.BenchmarkScanInstruction{
		CheckAlternatives: []*ipb.CheckAlternative{
			&ipb.CheckAlternative{
				FileChecks: []*ipb.FileCheck{check},
			},
			&ipb.CheckAlternative{
				FileChecks: []*ipb.FileCheck{check},
			},
		},
	}

	config := testconfigcreator.NewBenchmarkConfig(t, "id", instruction)

	checks, err := configchecks.CreateChecksFromConfig(
		context.Background(),
		&apb.ScanConfig{
			BenchmarkConfigs: []*apb.BenchmarkConfig{config},
		},
		newFakeAPI())
	if err != nil {
		t.Fatalf("configchecks.CreateChecksFromConfig([%v]) returned an error: %v", config, err)
	}
	if len(checks) != 1 {
		t.Errorf("configchecks.CreateChecksFromConfig([%v]) expected to create 1 check, got %d", config, len(checks))
	}
}

func TestChecksOnDifferentFilesGroupedSeparately(t *testing.T) {
	fileCheck1 := &ipb.FileCheck{
		FilesToCheck: []*ipb.FileSet{testconfigcreator.SingleFileWithPath("/path1")},
		CheckType:    &ipb.FileCheck_Existence{Existence: &ipb.ExistenceCheck{ShouldExist: true}},
	}
	fileCheck2 := &ipb.FileCheck{
		FilesToCheck: []*ipb.FileSet{testconfigcreator.SingleFileWithPath("/path2")},
		CheckType:    &ipb.FileCheck_Existence{Existence: &ipb.ExistenceCheck{ShouldExist: false}},
	}

	scanInstruction1 := testconfigcreator.NewFileScanInstruction([]*ipb.FileCheck{fileCheck1})
	scanInstruction2 := testconfigcreator.NewFileScanInstruction([]*ipb.FileCheck{fileCheck2})
	config1 := testconfigcreator.NewBenchmarkConfig(t, "id1", scanInstruction1)
	config2 := testconfigcreator.NewBenchmarkConfig(t, "id2", scanInstruction2)

	checks, err := configchecks.CreateChecksFromConfig(
		context.Background(),
		&apb.ScanConfig{
			BenchmarkConfigs: []*apb.BenchmarkConfig{config1, config2},
		},
		newFakeAPI())
	if err != nil {
		t.Fatalf("configchecks.CreateChecksFromConfig([%v %v]) returned an error: %v", config1, config2, err)
	}
	if len(checks) != 2 {
		t.Errorf("configchecks.CreateChecksFromConfig([%v %v]) expected to create 2 checks,got %d", config1, config2, len(checks))
	}
}

func TestInvalidCheckCombination(t *testing.T) {
	testCases := []struct {
		desc       string
		fileCheck1 *ipb.FileCheck
		fileCheck2 *ipb.FileCheck
	}{
		{
			desc: "Different delimiters",
			fileCheck1: &ipb.FileCheck{
				FilesToCheck: []*ipb.FileSet{testconfigcreator.SingleFileWithPath("/path")},
				CheckType: &ipb.FileCheck_ContentEntry{ContentEntry: &ipb.ContentEntryCheck{
					Delimiter: []byte{'\n'},
				}},
			},
			fileCheck2: &ipb.FileCheck{
				FilesToCheck: []*ipb.FileSet{testconfigcreator.SingleFileWithPath("/path")},
				CheckType: &ipb.FileCheck_ContentEntry{ContentEntry: &ipb.ContentEntryCheck{
					Delimiter: []byte{';'},
				}},
			},
		},
		{
			desc: "Both content and content entry checks present",
			fileCheck1: &ipb.FileCheck{
				FilesToCheck: []*ipb.FileSet{testconfigcreator.SingleFileWithPath("/path")},
				CheckType:    &ipb.FileCheck_Content{Content: &ipb.ContentCheck{Content: "content"}},
			},
			fileCheck2: &ipb.FileCheck{
				FilesToCheck: []*ipb.FileSet{testconfigcreator.SingleFileWithPath("/path")},
				CheckType: &ipb.FileCheck_ContentEntry{ContentEntry: &ipb.ContentEntryCheck{
					Delimiter: []byte{'\n'},
				}},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			scanInstruction1 := testconfigcreator.NewFileScanInstruction([]*ipb.FileCheck{tc.fileCheck1})
			scanInstruction2 := testconfigcreator.NewFileScanInstruction([]*ipb.FileCheck{tc.fileCheck2})
			config1 := testconfigcreator.NewBenchmarkConfig(t, "id1", scanInstruction1)
			config2 := testconfigcreator.NewBenchmarkConfig(t, "id2", scanInstruction2)
			_, err := configchecks.CreateChecksFromConfig(
				context.Background(),
				&apb.ScanConfig{
					BenchmarkConfigs: []*apb.BenchmarkConfig{config1, config2},
				},
				&fakeAPI{fileContent: testFileContent},
			)
			if err == nil {
				t.Fatalf("configchecks.CreateChecksFromConfig([%v %v]) didn't return an error", config1, config2)
			}
		})
	}
}

func TestIdsDeduplicated(t *testing.T) {
	fileCheck1 := &ipb.FileCheck{
		FilesToCheck: []*ipb.FileSet{testconfigcreator.SingleFileWithPath("/path1")},
		CheckType:    &ipb.FileCheck_Existence{Existence: &ipb.ExistenceCheck{ShouldExist: true}},
	}
	fileCheck2 := &ipb.FileCheck{
		FilesToCheck: []*ipb.FileSet{testconfigcreator.SingleFileWithPath("/path1")},
		CheckType:    &ipb.FileCheck_Existence{Existence: &ipb.ExistenceCheck{ShouldExist: false}},
	}

	scanInstruction1 := testconfigcreator.NewFileScanInstruction([]*ipb.FileCheck{fileCheck1})
	scanInstruction2 := testconfigcreator.NewFileScanInstruction([]*ipb.FileCheck{fileCheck2})
	config1 := testconfigcreator.NewBenchmarkConfig(t, "id", scanInstruction1)
	config2 := testconfigcreator.NewBenchmarkConfig(t, "id", scanInstruction2)

	checks, err := configchecks.CreateChecksFromConfig(
		context.Background(),
		&apb.ScanConfig{
			BenchmarkConfigs: []*apb.BenchmarkConfig{config1, config2},
		},
		newFakeAPI())
	if err != nil {
		t.Fatalf("configchecks.CreateChecksFromConfig([%v %v]) returned an error: %v", config1, config2, err)
	}
	if len(checks) != 1 {
		t.Errorf("configchecks.CreateChecksFromConfig([%v %v]) expected to create one check ,got %d", config1, config2, len(checks))
	}

	wantIDs := []string{"id"}
	if diff := cmp.Diff(wantIDs, checks[0].BenchmarkIDs()); diff != "" {
		t.Errorf("%v.BenchmarkIDs() returned unexpected diff (-want +got):\n%s", checks[0], diff)
	}
}

func TestMultipleChecksCreatedForMultipleFileSets(t *testing.T) {
	file1 := testconfigcreator.SingleFileWithPath("/path1")
	file2 := testconfigcreator.SingleFileWithPath("/path2")
	file3 := testconfigcreator.SingleFileWithPath("/path3")
	scanInstruction1 := testconfigcreator.NewFileScanInstruction([]*ipb.FileCheck{
		&ipb.FileCheck{
			FilesToCheck: []*ipb.FileSet{file1, file2},
			CheckType:    &ipb.FileCheck_Existence{Existence: &ipb.ExistenceCheck{ShouldExist: true}},
		},
	})
	scanInstruction2 := testconfigcreator.NewFileScanInstruction([]*ipb.FileCheck{
		&ipb.FileCheck{
			FilesToCheck: []*ipb.FileSet{file2, file3},
			CheckType:    &ipb.FileCheck_Existence{Existence: &ipb.ExistenceCheck{ShouldExist: true}},
		},
	})
	config1 := testconfigcreator.NewBenchmarkConfig(t, "id1", scanInstruction1)
	config2 := testconfigcreator.NewBenchmarkConfig(t, "id2", scanInstruction2)

	checks, err := configchecks.CreateChecksFromConfig(
		context.Background(),
		&apb.ScanConfig{
			BenchmarkConfigs: []*apb.BenchmarkConfig{config1, config2},
		},
		newFakeAPI())
	if err != nil {
		t.Fatalf("configchecks.CreateChecksFromConfig([%v %v]) returned an error: %v", config1, config2, err)
	}
	if len(checks) != 3 {
		t.Fatalf("Expected 3 check to be created, got %d", len(checks))
	}
}

func TestFileCheckWithEmptyInstructionsReturnsError(t *testing.T) {
	scanInstruction := &ipb.BenchmarkScanInstruction{
		CheckAlternatives: []*ipb.CheckAlternative{},
	}
	config := testconfigcreator.NewBenchmarkConfig(t, "id", scanInstruction)

	scanConfig := &apb.ScanConfig{BenchmarkConfigs: []*apb.BenchmarkConfig{config}}
	if _, err := configchecks.CreateChecksFromConfig(context.Background(), scanConfig, newFakeAPI()); err == nil {
		t.Errorf("configchecks.CreateChecksFromConfig([%v]) didn't return an error", config)
	}
}

func TestCheckCreation(t *testing.T) {
	filesToCheck := []*ipb.FileSet{testconfigcreator.SingleFileWithPath("/path")}
	// Test the creation of all check types.
	checks := []*ipb.FileCheck{
		&ipb.FileCheck{
			FilesToCheck: filesToCheck,
			CheckType:    &ipb.FileCheck_Existence{Existence: &ipb.ExistenceCheck{ShouldExist: true}},
		},
		&ipb.FileCheck{
			FilesToCheck: filesToCheck,
			CheckType:    &ipb.FileCheck_Permission{Permission: &ipb.PermissionCheck{SetBits: 0b0101}},
		},
		&ipb.FileCheck{
			FilesToCheck: filesToCheck,
			CheckType:    &ipb.FileCheck_Content{Content: &ipb.ContentCheck{Content: "content"}},
		},
		&ipb.FileCheck{
			FilesToCheck: filesToCheck,
			CheckType: &ipb.FileCheck_ContentEntry{ContentEntry: &ipb.ContentEntryCheck{
				MatchCriteria: []*ipb.MatchCriterion{&ipb.MatchCriterion{
					FilterRegex:   "config_field = .*",
					ExpectedRegex: "config_field = 'value'",
				}},
			}},
		},
	}

	for _, check := range checks {
		scanInstruction := testconfigcreator.NewFileScanInstruction([]*ipb.FileCheck{check})
		config := testconfigcreator.NewBenchmarkConfig(t, "id", scanInstruction)

		checks, err := configchecks.CreateChecksFromConfig(
			context.Background(),
			&apb.ScanConfig{
				BenchmarkConfigs: []*apb.BenchmarkConfig{config},
			},
			newFakeAPI())
		if err != nil {
			t.Fatalf("configchecks.CreateChecksFromConfig([%v]) returned an error: %v", config, err)
		}
		if len(checks) != 1 {
			t.Fatalf("Expected 1 check to be created, got %d", len(checks))
		}
		expectedIDs := []string{"id"}
		if diff := cmp.Diff(expectedIDs, checks[0].BenchmarkIDs()); diff != "" {
			t.Errorf("%v.BenchmarkIDs() returned unexpected diff (-want +got):\n%s", checks[0], diff)
		}
	}
}

func createFileCheckBatch(t *testing.T, id string, fileChecks []*ipb.FileCheck, api scanapi.ScanAPI) configchecks.BenchmarkCheck {
	t.Helper()
	scanInstruction := testconfigcreator.NewFileScanInstruction(fileChecks)
	config := testconfigcreator.NewBenchmarkConfig(t, "id", scanInstruction)
	return createFileCheckBatchFromScanConfig(t, id, &apb.ScanConfig{BenchmarkConfigs: []*apb.BenchmarkConfig{config}}, api)
}

func createFileCheckBatchFromScanConfig(t *testing.T, id string, scanConfig *apb.ScanConfig, api scanapi.ScanAPI) configchecks.BenchmarkCheck {
	t.Helper()
	checks, err := configchecks.CreateChecksFromConfig(context.Background(), scanConfig, api)
	if err != nil {
		t.Fatalf("configchecks.CreateChecksFromConfig(%v) returned an error: %v", scanConfig, err)
	}
	if len(checks) != 1 {
		t.Fatalf("Created %d checks, expected only 1", len(checks))
	}
	return checks[0]
}

func TestFileCustomNonComplianceMessage(t *testing.T) {
	testCases := []struct {
		desc           string
		fileCheck      *ipb.FileCheck
		expectedResult *apb.ComplianceResult
	}{
		{
			desc: "Custom message",
			fileCheck: &ipb.FileCheck{
				FilesToCheck:     []*ipb.FileSet{testconfigcreator.SingleFileWithPath(testFilePath)},
				CheckType:        &ipb.FileCheck_Existence{Existence: &ipb.ExistenceCheck{ShouldExist: false}},
				NonComplianceMsg: "custom reason",
			},
			expectedResult: &apb.ComplianceResult{
				Id: "id",
				ComplianceOccurrence: &cpb.ComplianceOccurrence{
					NonCompliantFiles: []*cpb.NonCompliantFile{
						&cpb.NonCompliantFile{
							Path:   testFilePath,
							Reason: "custom reason",
						},
					},
				},
			},
		},
		{
			desc: "Custom message + command",
			fileCheck: &ipb.FileCheck{
				FilesToCheck:       []*ipb.FileSet{testconfigcreator.SingleFileWithPath(testFilePath)},
				CheckType:          &ipb.FileCheck_Existence{Existence: &ipb.ExistenceCheck{ShouldExist: false}},
				NonComplianceMsg:   "custom reason",
				FileDisplayCommand: "display command",
			},
			expectedResult: &apb.ComplianceResult{
				Id: "id",
				ComplianceOccurrence: &cpb.ComplianceOccurrence{
					NonCompliantFiles: []*cpb.NonCompliantFile{
						&cpb.NonCompliantFile{
							DisplayCommand: "display command",
							Reason:         "custom reason",
						},
					},
				},
			},
		},
		{
			desc: "Compliant check",
			fileCheck: &ipb.FileCheck{
				FilesToCheck:       []*ipb.FileSet{testconfigcreator.SingleFileWithPath(testFilePath)},
				CheckType:          &ipb.FileCheck_Existence{Existence: &ipb.ExistenceCheck{ShouldExist: true}},
				NonComplianceMsg:   "custom reason",
				FileDisplayCommand: "display command",
			},
			expectedResult: &apb.ComplianceResult{
				Id:                   "id",
				ComplianceOccurrence: &cpb.ComplianceOccurrence{},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			check := createFileCheckBatch(t, "id", []*ipb.FileCheck{tc.fileCheck}, newFakeAPI())
			resultMap, _, err := check.Exec("")
			if err != nil {
				t.Fatalf("check.Exec() returned an error: %v", err)
			}
			result, gotSingleton := singleComplianceResult(resultMap)
			if !gotSingleton {
				t.Fatalf("check.Exec() expected to return 1 result, got %d", len(resultMap))
			}

			if diff := cmp.Diff(tc.expectedResult, result, protocmp.Transform()); diff != "" {
				t.Errorf("check.Exec() returned unexpected diff (-want +got):\n%s", diff)
			}
		})
	}
}

func TestFileDisplayCommandWithoutCustomMessage(t *testing.T) {
	fileCheck := &ipb.FileCheck{
		FilesToCheck:       []*ipb.FileSet{testconfigcreator.SingleFileWithPath(testFilePath)},
		CheckType:          &ipb.FileCheck_Existence{Existence: &ipb.ExistenceCheck{ShouldExist: false}},
		FileDisplayCommand: "display command",
	}

	scanInstruction := testconfigcreator.NewFileScanInstruction([]*ipb.FileCheck{fileCheck})
	config := testconfigcreator.NewBenchmarkConfig(t, "id", scanInstruction)

	if _, err := configchecks.CreateChecksFromConfig(context.Background(), &apb.ScanConfig{BenchmarkConfigs: []*apb.BenchmarkConfig{config}}, newFakeAPI()); err == nil {
		t.Fatalf("configchecks.CreateChecksFromConfig(%v) didn't return an error", config)
	}
}

func TestFilesInOptOutConfigRedacted(t *testing.T) {
	testPath := []*ipb.FileSet{testconfigcreator.SingleFileWithPath(testFilePath)}
	testCases := []struct {
		desc           string
		fileCheck      *ipb.FileCheck
		optOutConfig   *apb.OptOutConfig
		expectedResult *apb.ComplianceResult
	}{
		{
			desc: "content not displayed",
			fileCheck: &ipb.FileCheck{
				FilesToCheck: testPath,
				CheckType:    &ipb.FileCheck_Content{Content: &ipb.ContentCheck{Content: "content"}},
			},
			optOutConfig: &apb.OptOutConfig{
				ContentOptoutRegexes: []string{".*"},
			},
			expectedResult: &apb.ComplianceResult{
				Id: "id",
				ComplianceOccurrence: &cpb.ComplianceOccurrence{
					NonCompliantFiles: []*cpb.NonCompliantFile{
						&cpb.NonCompliantFile{
							Path:   testFilePath,
							Reason: "Got content \"[redacted due to opt-out config]\", expected \"content\"",
						},
					},
				},
			},
		},
		{
			desc: "filename not displayed",
			fileCheck: &ipb.FileCheck{
				FilesToCheck: testPath,
				CheckType:    &ipb.FileCheck_Content{Content: &ipb.ContentCheck{Content: "content"}},
			},
			optOutConfig: &apb.OptOutConfig{
				FilenameOptoutRegexes: []string{".*"},
			},
			expectedResult: &apb.ComplianceResult{
				Id: "id",
				ComplianceOccurrence: &cpb.ComplianceOccurrence{
					NonCompliantFiles: []*cpb.NonCompliantFile{
						&cpb.NonCompliantFile{
							Path:   "[redacted due to opt-out config]",
							Reason: fmt.Sprintf("Got content %q, expected \"content\"", testFileContent),
						},
					},
				},
			},
		},
		{
			desc: "don't redact if regex doesn't match",
			fileCheck: &ipb.FileCheck{
				FilesToCheck: testPath,
				CheckType:    &ipb.FileCheck_Content{Content: &ipb.ContentCheck{Content: "content"}},
			},
			optOutConfig: &apb.OptOutConfig{
				ContentOptoutRegexes: []string{"no match"},
			},
			expectedResult: &apb.ComplianceResult{
				Id: "id",
				ComplianceOccurrence: &cpb.ComplianceOccurrence{
					NonCompliantFiles: []*cpb.NonCompliantFile{
						&cpb.NonCompliantFile{
							Path:   testFilePath,
							Reason: fmt.Sprintf("Got content %q, expected \"content\"", testFileContent),
						},
					},
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			scanInstruction := testconfigcreator.NewFileScanInstruction([]*ipb.FileCheck{tc.fileCheck})
			config := testconfigcreator.NewBenchmarkConfig(t, "id", scanInstruction)
			scanConfig := &apb.ScanConfig{
				BenchmarkConfigs: []*apb.BenchmarkConfig{config},
				OptOutConfig:     tc.optOutConfig,
			}
			check := createFileCheckBatchFromScanConfig(t, "id", scanConfig, newFakeAPI())
			resultMap, _, err := check.Exec("")
			if err != nil {
				t.Fatalf("check.Exec() returned an error: %v", err)
			}
			result, gotSingleton := singleComplianceResult(resultMap)
			if !gotSingleton {
				t.Fatalf("check.Exec() expected to return 1 result, got %d", len(resultMap))
			}

			if diff := cmp.Diff(tc.expectedResult, result, protocmp.Transform()); diff != "" {
				t.Errorf("check.Exec() returned unexpected diff (-want +got):\n%s", diff)
			}
		})
	}
}

func TestReplacementConfigApplied(t *testing.T) {
	testPath := []*ipb.FileSet{&ipb.FileSet{
		FilePath: &ipb.FileSet_SingleFile_{SingleFile: &ipb.FileSet_SingleFile{Path: unreadableFilePath}},
	}}
	fileCheck := &ipb.FileCheck{
		FilesToCheck: testPath,
		CheckType:    &ipb.FileCheck_Content{Content: &ipb.ContentCheck{Content: testFileContent}},
	}
	// Replace file from instruction with a readable file.
	prefix := strings.TrimSuffix(unreadableFilePath, path.Base(unreadableFilePath))
	replacement := strings.TrimSuffix(testFilePath, path.Base(testFilePath))
	replacementConfig := &apb.ReplacementConfig{PathPrefixReplacements: map[string]string{prefix: replacement}}
	// The check should succeed after the file gets replaced.
	wantResult := &apb.ComplianceResult{
		Id:                   "id",
		ComplianceOccurrence: &cpb.ComplianceOccurrence{},
	}

	scanInstruction := testconfigcreator.NewFileScanInstruction([]*ipb.FileCheck{fileCheck})
	config := testconfigcreator.NewBenchmarkConfig(t, "id", scanInstruction)
	scanConfig := &apb.ScanConfig{
		BenchmarkConfigs:  []*apb.BenchmarkConfig{config},
		ReplacementConfig: replacementConfig,
	}
	check := createFileCheckBatchFromScanConfig(t, "id", scanConfig, newFakeAPI())
	resultMap, _, err := check.Exec("")
	if err != nil {
		t.Fatalf("check.Exec() returned an error: %v", err)
	}
	result, gotSingleton := singleComplianceResult(resultMap)
	if !gotSingleton {
		t.Fatalf("check.Exec() expected to return 1 result, got %d", len(resultMap))
	}

	if diff := cmp.Diff(wantResult, result, protocmp.Transform()); diff != "" {
		t.Errorf("check.Exec() returned unexpected diff (-want +got):\n%s", diff)
	}
}

func TestResultsForDifferentAlternativesAggregatedSeparately(t *testing.T) {
	fileCheck := &ipb.FileCheck{
		FilesToCheck: []*ipb.FileSet{testconfigcreator.SingleFileWithPath(testFilePath)},
		CheckType:    &ipb.FileCheck_Existence{Existence: &ipb.ExistenceCheck{ShouldExist: true}},
	}
	scanInstruction := &ipb.BenchmarkScanInstruction{
		CheckAlternatives: []*ipb.CheckAlternative{
			&ipb.CheckAlternative{
				FileChecks: []*ipb.FileCheck{fileCheck},
			},
			&ipb.CheckAlternative{
				FileChecks: []*ipb.FileCheck{fileCheck},
			},
		},
	}

	config := testconfigcreator.NewBenchmarkConfig(t, "id", scanInstruction)
	checks, err := configchecks.CreateChecksFromConfig(
		context.Background(),
		&apb.ScanConfig{
			BenchmarkConfigs: []*apb.BenchmarkConfig{config},
		}, newFakeAPI())
	if err != nil {
		t.Fatalf("configchecks.CreateChecksFromConfig([%v]) returned an error: %v", config, err)
	}
	if len(checks) != 1 {
		t.Fatalf("Created %d checks, expected only 1", len(checks))
	}
	check := checks[0]

	resultMap, _, err := check.Exec("")
	if err != nil {
		t.Fatalf("check.Exec() returned an error: %v", err)
	}

	if len(resultMap) != 2 {
		t.Errorf("Expected results to be present for 2 check alternatives, found %d: %v", len(resultMap), resultMap)
	}
}

type dirWithUnreadableFile struct{}

func (dirWithUnreadableFile) OpenFile(ctx context.Context, filePath string) (io.ReadCloser, error) {
	return nil, os.ErrNotExist
}

func (dirWithUnreadableFile) OpenDir(ctx context.Context, filePath string) (scanapi.DirReader, error) {
	return scanapi.SliceToDirReader([]*apb.DirContent{
		{Name: path.Base(testFilePath), IsDir: false},
	}), nil
}

func (dirWithUnreadableFile) FilePermissions(ctx context.Context, filePath string) (*apb.PosixPermissions, error) {
	return nil, errors.New("not implemented")
}

func (dirWithUnreadableFile) SQLQuery(ctx context.Context, query string) (string, error) {
	return "", errors.New("not implemented")
}

func (dirWithUnreadableFile) SupportedDatabase() (ipb.SQLCheck_SQLDatabase, error) {
	return ipb.SQLCheck_DB_UNSPECIFIED, errors.New("not implemented")
}

func TestFileExistenceCheckComplianceResults(t *testing.T) {
	testCases := []struct {
		desc           string
		api            scanapi.ScanAPI
		fileCheck      *ipb.FileCheck
		expectedResult *apb.ComplianceResult
	}{
		{
			desc: "File exists and should",
			api:  newFakeAPI(),
			fileCheck: &ipb.FileCheck{
				FilesToCheck: []*ipb.FileSet{testconfigcreator.SingleFileWithPath(testFilePath)},
				CheckType:    &ipb.FileCheck_Existence{Existence: &ipb.ExistenceCheck{ShouldExist: true}},
			},
			expectedResult: &apb.ComplianceResult{
				Id:                   "id",
				ComplianceOccurrence: &cpb.ComplianceOccurrence{},
			},
		},
		{
			desc: "File doesn't exist and shouldn't",
			api:  newFakeAPI(),
			fileCheck: &ipb.FileCheck{
				FilesToCheck: []*ipb.FileSet{testconfigcreator.SingleFileWithPath(nonExistentFilePath)},
				CheckType:    &ipb.FileCheck_Existence{Existence: &ipb.ExistenceCheck{ShouldExist: false}},
			},
			expectedResult: &apb.ComplianceResult{
				Id:                   "id",
				ComplianceOccurrence: &cpb.ComplianceOccurrence{},
			},
		},
		{
			desc: "File doesn't exist but should",
			api:  newFakeAPI(),
			fileCheck: &ipb.FileCheck{
				FilesToCheck: []*ipb.FileSet{testconfigcreator.SingleFileWithPath(nonExistentFilePath)},
				CheckType:    &ipb.FileCheck_Existence{Existence: &ipb.ExistenceCheck{ShouldExist: true}},
			},
			expectedResult: &apb.ComplianceResult{
				Id: "id",
				ComplianceOccurrence: &cpb.ComplianceOccurrence{
					NonCompliantFiles: []*cpb.NonCompliantFile{
						&cpb.NonCompliantFile{
							Path:   fileset.FileSetToString(testconfigcreator.SingleFileWithPath(nonExistentFilePath)),
							Reason: "File doesn't exist but it should",
						},
					},
				},
			},
		},
		{
			desc: "File in directory doesn't exist but it should",
			api:  newFakeAPI(),
			fileCheck: &ipb.FileCheck{
				FilesToCheck: []*ipb.FileSet{&ipb.FileSet{
					FilePath: &ipb.FileSet_FilesInDir_{FilesInDir: &ipb.FileSet_FilesInDir{
						DirPath:       testDirPath,
						Recursive:     true,
						FilenameRegex: "non-existent",
					}},
				}},
				CheckType: &ipb.FileCheck_Existence{Existence: &ipb.ExistenceCheck{ShouldExist: true}},
			},
			expectedResult: &apb.ComplianceResult{
				Id: "id",
				ComplianceOccurrence: &cpb.ComplianceOccurrence{
					NonCompliantFiles: []*cpb.NonCompliantFile{
						&cpb.NonCompliantFile{
							Path: fileset.FileSetToString(&ipb.FileSet{
								FilePath: &ipb.FileSet_FilesInDir_{FilesInDir: &ipb.FileSet_FilesInDir{
									DirPath:       testDirPath,
									Recursive:     true,
									FilenameRegex: "non-existent",
								}},
							}),
							Reason: "File doesn't exist but it should",
						},
					},
				},
			},
		},
		{
			desc: "File exists but it shouldn't",
			api:  newFakeAPI(),
			fileCheck: &ipb.FileCheck{
				FilesToCheck: []*ipb.FileSet{testconfigcreator.SingleFileWithPath(testFilePath)},
				CheckType:    &ipb.FileCheck_Existence{Existence: &ipb.ExistenceCheck{ShouldExist: false}},
			},
			expectedResult: &apb.ComplianceResult{
				Id: "id",
				ComplianceOccurrence: &cpb.ComplianceOccurrence{
					NonCompliantFiles: []*cpb.NonCompliantFile{
						&cpb.NonCompliantFile{
							Path:   testFilePath,
							Reason: "File exists but it shouldn't",
						},
					},
				},
			},
		},
		{
			desc: "File exists but unreadable",
			api:  &dirWithUnreadableFile{},
			fileCheck: &ipb.FileCheck{
				FilesToCheck: []*ipb.FileSet{&ipb.FileSet{
					FilePath: &ipb.FileSet_FilesInDir_{FilesInDir: &ipb.FileSet_FilesInDir{
						DirPath:   testDirPath,
						Recursive: false,
					}},
				}},
				CheckType: &ipb.FileCheck_Existence{Existence: &ipb.ExistenceCheck{ShouldExist: true}},
			},
			expectedResult: &apb.ComplianceResult{
				Id:                   "id",
				ComplianceOccurrence: &cpb.ComplianceOccurrence{},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			check := createFileCheckBatch(t, "id", []*ipb.FileCheck{tc.fileCheck}, tc.api)
			resultMap, _, err := check.Exec("")
			if err != nil {
				t.Fatalf("check.Exec() returned an error: %v", err)
			}
			result, gotSingleton := singleComplianceResult(resultMap)
			if !gotSingleton {
				t.Fatalf("check.Exec() expected to return 1 result, got %d", len(resultMap))
			}

			if diff := cmp.Diff(tc.expectedResult, result, protocmp.Transform()); diff != "" {
				t.Errorf("check.Exec() returned unexpected diff (-want +got):\n%s", diff)
			}
		})
	}
}

func TestFileWithPipelining(t *testing.T) {
	testCases := []struct {
		desc           string
		fileCheck      *ipb.FileCheck
		expectedResult *apb.ComplianceResult
		fileToCheck    string
	}{
		{
			desc: "file exists with pipelined input",
			fileCheck: &ipb.FileCheck{
				FilesToCheck: []*ipb.FileSet{testconfigcreator.SingleFileWithPath(pipelineFileToken)},
				CheckType:    &ipb.FileCheck_Existence{Existence: &ipb.ExistenceCheck{ShouldExist: true}},
			},
			expectedResult: &apb.ComplianceResult{
				Id:                   "id",
				ComplianceOccurrence: &cpb.ComplianceOccurrence{},
			},
			fileToCheck: testFilePath,
		},
		{
			desc: "file not exists with pipelined input",
			fileCheck: &ipb.FileCheck{
				FilesToCheck: []*ipb.FileSet{testconfigcreator.SingleFileWithPath(pipelineFileToken)},
				CheckType:    &ipb.FileCheck_Existence{Existence: &ipb.ExistenceCheck{ShouldExist: false}},
			},
			expectedResult: &apb.ComplianceResult{
				Id:                   "id",
				ComplianceOccurrence: &cpb.ComplianceOccurrence{},
			},
			fileToCheck: nonExistentFilePath,
		},
		{
			desc: "pipeline input expected but not provided",
			fileCheck: &ipb.FileCheck{
				FilesToCheck: []*ipb.FileSet{testconfigcreator.SingleFileWithPath(pipelineFileToken)},
				CheckType:    &ipb.FileCheck_Existence{Existence: &ipb.ExistenceCheck{ShouldExist: false}},
			},
			expectedResult: &apb.ComplianceResult{
				Id:                   "id",
				ComplianceOccurrence: &cpb.ComplianceOccurrence{},
			},
			fileToCheck: "",
		},
		{
			desc: "pipeline input not expected but provided",
			fileCheck: &ipb.FileCheck{
				FilesToCheck: []*ipb.FileSet{testconfigcreator.SingleFileWithPath(testFileContent)},
				CheckType:    &ipb.FileCheck_Existence{Existence: &ipb.ExistenceCheck{ShouldExist: true}},
			},
			expectedResult: &apb.ComplianceResult{
				Id:                   "id",
				ComplianceOccurrence: &cpb.ComplianceOccurrence{},
			},
			fileToCheck: nonExistentFilePath,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			check := createFileCheckBatch(t, "id", []*ipb.FileCheck{tc.fileCheck}, newFakeAPI())
			resultMap, _, err := check.Exec(tc.fileToCheck)
			if err != nil {
				t.Fatalf("check.Exec() returned an error: %v", err)
			}
			result, gotSingleton := singleComplianceResult(resultMap)
			if !gotSingleton {
				t.Fatalf("check.Exec() expected to return 1 result, got %d", len(resultMap))
			}

			if diff := cmp.Diff(tc.expectedResult, result, protocmp.Transform()); diff != "" {
				t.Errorf("check.Exec() returned unexpected diff (-want +got):\n%s", diff)
			}
		})
	}
}

func TestFileExistenceCheckPropagatesError(t *testing.T) {
	check := createFileCheckBatch(t, "id", []*ipb.FileCheck{&ipb.FileCheck{
		FilesToCheck: []*ipb.FileSet{testconfigcreator.SingleFileWithPath(unreadableFilePath)},
		CheckType:    &ipb.FileCheck_Existence{Existence: &ipb.ExistenceCheck{ShouldExist: true}},
	}}, newFakeAPI())
	if _, _, err := check.Exec(""); err == nil {
		t.Errorf("check.Exec() didn't return an error")
	}
}

func TestFileExistenceWithWrappedError(t *testing.T) {
	fileCheck := &ipb.FileCheck{
		FilesToCheck: []*ipb.FileSet{testconfigcreator.SingleFileWithPath(nonExistentFilePath)},
		CheckType:    &ipb.FileCheck_Existence{Existence: &ipb.ExistenceCheck{ShouldExist: false}},
	}
	openFileFunc := func(ctx context.Context, filePath string) (io.ReadCloser, error) {
		return nil, fmt.Errorf("Error: %w", os.ErrNotExist)
	}
	expectedResult := &apb.ComplianceResult{Id: "id", ComplianceOccurrence: &cpb.ComplianceOccurrence{}}

	check := createFileCheckBatch(t, "id", []*ipb.FileCheck{fileCheck}, newFakeAPI(withOpenFileFunc(openFileFunc)))
	resultMap, _, err := check.Exec("")

	if err != nil {
		t.Fatalf("check.Exec() returned an error: %v", err)
	}
	result, gotSingleton := singleComplianceResult(resultMap)
	if !gotSingleton {
		t.Fatalf("check.Exec() expected to return 1 result, got %d", len(resultMap))
	}

	if diff := cmp.Diff(expectedResult, result, protocmp.Transform()); diff != "" {
		t.Errorf("check.Exec() returned unexpected diff (-want +got):\n%s", diff)
	}
}

func TestPermissionCheckComplianceResults(t *testing.T) {
	testCases := []struct {
		description     string
		file            string
		permissionCheck *ipb.PermissionCheck
		expectedResult  *apb.ComplianceResult
	}{
		{
			description: "compliant permission check",
			file:        testFilePath,
			permissionCheck: &ipb.PermissionCheck{
				SetBits:         0644,
				ClearBits:       0133,
				BitsShouldMatch: ipb.PermissionCheck_BOTH_SET_AND_CLEAR,
				User:            &ipb.PermissionCheck_OwnerCheck{Name: "root", ShouldOwn: true},
				Group:           &ipb.PermissionCheck_OwnerCheck{Name: "root", ShouldOwn: true},
			},
			expectedResult: &apb.ComplianceResult{
				Id:                   "id",
				ComplianceOccurrence: &cpb.ComplianceOccurrence{},
			},
		},
		{
			description: "file doesn't exist",
			file:        nonExistentFilePath,
			permissionCheck: &ipb.PermissionCheck{
				SetBits: 0644,
			},
			expectedResult: &apb.ComplianceResult{
				Id: "id",
				ComplianceOccurrence: &cpb.ComplianceOccurrence{
					NonCompliantFiles: []*cpb.NonCompliantFile{
						&cpb.NonCompliantFile{
							Path:   nonExistentFilePath,
							Reason: "File doesn't exist",
						},
					},
				},
			},
		},
		{
			description:     "set bits don't match",
			file:            testFilePath,
			permissionCheck: &ipb.PermissionCheck{SetBits: 0643},
			expectedResult: &apb.ComplianceResult{
				Id: "id",
				ComplianceOccurrence: &cpb.ComplianceOccurrence{
					NonCompliantFiles: []*cpb.NonCompliantFile{
						&cpb.NonCompliantFile{
							Path:   testFilePath,
							Reason: "File permission is 0644, expected the following bits to be set: 0643",
						},
					},
				},
			},
		},
		{
			description:     "clear bits don't match",
			file:            testFilePath,
			permissionCheck: &ipb.PermissionCheck{ClearBits: 0135},
			expectedResult: &apb.ComplianceResult{
				Id: "id",
				ComplianceOccurrence: &cpb.ComplianceOccurrence{
					NonCompliantFiles: []*cpb.NonCompliantFile{
						&cpb.NonCompliantFile{
							Path:   testFilePath,
							Reason: "File permission is 0644, expected the following bits to be clear: 0135",
						},
					},
				},
			},
		},
		{
			description: "set bits don't match but clear ones do, expected either",
			file:        testFilePath,
			permissionCheck: &ipb.PermissionCheck{
				SetBits:         0643,
				ClearBits:       0133,
				BitsShouldMatch: ipb.PermissionCheck_EITHER_SET_OR_CLEAR,
			},
			expectedResult: &apb.ComplianceResult{
				Id:                   "id",
				ComplianceOccurrence: &cpb.ComplianceOccurrence{},
			},
		},
		{
			description: "set bits don't match but clear ones do, expected both",
			file:        testFilePath,
			permissionCheck: &ipb.PermissionCheck{
				SetBits:         0643,
				ClearBits:       0133,
				BitsShouldMatch: ipb.PermissionCheck_BOTH_SET_AND_CLEAR,
			},
			expectedResult: &apb.ComplianceResult{
				Id: "id",
				ComplianceOccurrence: &cpb.ComplianceOccurrence{
					NonCompliantFiles: []*cpb.NonCompliantFile{
						&cpb.NonCompliantFile{
							Path:   testFilePath,
							Reason: "File permission is 0644, expected the following bits to be set: 0643 and the following bits to be clear: 0133",
						},
					},
				},
			},
		},
		{
			description: "neither set nor clear bits match, expected either",
			file:        testFilePath,
			permissionCheck: &ipb.PermissionCheck{
				SetBits:         0643,
				ClearBits:       0135,
				BitsShouldMatch: ipb.PermissionCheck_EITHER_SET_OR_CLEAR,
			},
			expectedResult: &apb.ComplianceResult{
				Id: "id",
				ComplianceOccurrence: &cpb.ComplianceOccurrence{
					NonCompliantFiles: []*cpb.NonCompliantFile{
						&cpb.NonCompliantFile{
							Path:   testFilePath,
							Reason: "File permission is 0644, expected the following bits to be set: 0643 or the following bits to be clear: 0135",
						},
					},
				},
			},
		},
		{
			description: "neither set nor clear bits match, expected both",
			file:        testFilePath,
			permissionCheck: &ipb.PermissionCheck{
				SetBits:         0643,
				ClearBits:       0135,
				BitsShouldMatch: ipb.PermissionCheck_BOTH_SET_AND_CLEAR,
			},
			expectedResult: &apb.ComplianceResult{
				Id: "id",
				ComplianceOccurrence: &cpb.ComplianceOccurrence{
					NonCompliantFiles: []*cpb.NonCompliantFile{
						&cpb.NonCompliantFile{
							Path:   testFilePath,
							Reason: "File permission is 0644, expected the following bits to be set: 0643 and the following bits to be clear: 0135",
						},
					},
				},
			},
		},
		{
			description: "user doesn't match",
			file:        testFilePath,
			permissionCheck: &ipb.PermissionCheck{
				User: &ipb.PermissionCheck_OwnerCheck{Name: "not-root", ShouldOwn: true},
			},
			expectedResult: &apb.ComplianceResult{
				Id: "id",
				ComplianceOccurrence: &cpb.ComplianceOccurrence{
					NonCompliantFiles: []*cpb.NonCompliantFile{
						&cpb.NonCompliantFile{
							Path:   testFilePath,
							Reason: "Owner is root, expected it to be not-root",
						},
					},
				},
			},
		},
		{
			description: "user matches but shouldn't",
			file:        testFilePath,
			permissionCheck: &ipb.PermissionCheck{
				User: &ipb.PermissionCheck_OwnerCheck{Name: "root", ShouldOwn: false},
			},
			expectedResult: &apb.ComplianceResult{
				Id: "id",
				ComplianceOccurrence: &cpb.ComplianceOccurrence{
					NonCompliantFiles: []*cpb.NonCompliantFile{
						&cpb.NonCompliantFile{
							Path:   testFilePath,
							Reason: "Owner is root, expected it to be a different user",
						},
					},
				},
			},
		},
		{
			description: "group doesn't match",
			file:        testFilePath,
			permissionCheck: &ipb.PermissionCheck{
				Group: &ipb.PermissionCheck_OwnerCheck{Name: "not-root", ShouldOwn: true},
			},
			expectedResult: &apb.ComplianceResult{
				Id: "id",
				ComplianceOccurrence: &cpb.ComplianceOccurrence{
					NonCompliantFiles: []*cpb.NonCompliantFile{
						&cpb.NonCompliantFile{
							Path:   testFilePath,
							Reason: "Group is root, expected it to be not-root",
						},
					},
				},
			},
		},
		{
			description: "group matches but shouldn't",
			file:        testFilePath,
			permissionCheck: &ipb.PermissionCheck{
				Group: &ipb.PermissionCheck_OwnerCheck{Name: "root", ShouldOwn: false},
			},
			expectedResult: &apb.ComplianceResult{
				Id: "id",
				ComplianceOccurrence: &cpb.ComplianceOccurrence{
					NonCompliantFiles: []*cpb.NonCompliantFile{
						&cpb.NonCompliantFile{
							Path:   testFilePath,
							Reason: "Group is root, expected it to be a different group",
						},
					},
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			check := createFileCheckBatch(t, "id", []*ipb.FileCheck{&ipb.FileCheck{
				FilesToCheck: []*ipb.FileSet{testconfigcreator.SingleFileWithPath(tc.file)},
				CheckType:    &ipb.FileCheck_Permission{Permission: tc.permissionCheck},
			}}, newFakeAPI())

			resultMap, _, err := check.Exec("")
			if err != nil {
				t.Fatalf("check.Exec() returned an error: %v", err)
			}
			result, gotSingleton := singleComplianceResult(resultMap)
			if !gotSingleton {
				t.Fatalf("check.Exec() expected to return 1 result, got %d", len(resultMap))
			}

			if diff := cmp.Diff(tc.expectedResult, result, protocmp.Transform()); diff != "" {
				t.Errorf("check.Exec() returned unexpected diff (-want +got):\n%s", diff)
			}
		})
	}
}

func TestFileContentCheckComplianceResults(t *testing.T) {
	testCases := []struct {
		description    string
		fileChecks     []*ipb.FileCheck
		expectedResult *apb.ComplianceResult
	}{
		{
			description: "content matches",
			fileChecks: []*ipb.FileCheck{&ipb.FileCheck{
				FilesToCheck: []*ipb.FileSet{testconfigcreator.SingleFileWithPath(testFilePath)},
				CheckType:    &ipb.FileCheck_Content{Content: &ipb.ContentCheck{Content: testFileContent}},
			}},
			expectedResult: &apb.ComplianceResult{
				Id:                   "id",
				ComplianceOccurrence: &cpb.ComplianceOccurrence{},
			},
		},
		{
			description: "content doesn't match",
			fileChecks: []*ipb.FileCheck{&ipb.FileCheck{
				CheckType:    &ipb.FileCheck_Content{Content: &ipb.ContentCheck{Content: "Different content"}},
				FilesToCheck: []*ipb.FileSet{testconfigcreator.SingleFileWithPath(testFilePath)},
			}},
			expectedResult: &apb.ComplianceResult{
				Id: "id",
				ComplianceOccurrence: &cpb.ComplianceOccurrence{
					NonCompliantFiles: []*cpb.NonCompliantFile{
						&cpb.NonCompliantFile{
							Path:   testFilePath,
							Reason: fmt.Sprintf("Got content %q, expected \"Different content\"", testFileContent),
						},
					},
				},
			},
		},
		{
			description: "file doesn't exist",
			fileChecks: []*ipb.FileCheck{&ipb.FileCheck{
				FilesToCheck: []*ipb.FileSet{testconfigcreator.SingleFileWithPath(nonExistentFilePath)},
				CheckType:    &ipb.FileCheck_Content{Content: &ipb.ContentCheck{Content: testFileContent}},
			}},
			expectedResult: &apb.ComplianceResult{
				Id: "id",
				ComplianceOccurrence: &cpb.ComplianceOccurrence{
					NonCompliantFiles: []*cpb.NonCompliantFile{
						&cpb.NonCompliantFile{
							Path:   nonExistentFilePath,
							Reason: "File doesn't exist",
						},
					},
				},
			},
		},
		{
			description: "directory doesn't exist",
			fileChecks: []*ipb.FileCheck{&ipb.FileCheck{
				FilesToCheck: []*ipb.FileSet{&ipb.FileSet{
					FilePath: &ipb.FileSet_FilesInDir_{FilesInDir: &ipb.FileSet_FilesInDir{
						DirPath: nonExistentFilePath,
					}},
				}},
				CheckType: &ipb.FileCheck_Content{Content: &ipb.ContentCheck{Content: testFileContent}},
			}},
			expectedResult: &apb.ComplianceResult{
				Id: "id",
				ComplianceOccurrence: &cpb.ComplianceOccurrence{
					NonCompliantFiles: []*cpb.NonCompliantFile{
						&cpb.NonCompliantFile{
							Path:   nonExistentFilePath,
							Reason: "File doesn't exist",
						},
					},
				},
			},
		},
		{
			// The results of several non-compliant checks are aggregated.
			fileChecks: []*ipb.FileCheck{
				&ipb.FileCheck{
					FilesToCheck: []*ipb.FileSet{testconfigcreator.SingleFileWithPath(testFilePath)},
					CheckType:    &ipb.FileCheck_Content{Content: &ipb.ContentCheck{Content: testFileContent}},
				},
				&ipb.FileCheck{
					FilesToCheck: []*ipb.FileSet{testconfigcreator.SingleFileWithPath(testFilePath)},
					CheckType:    &ipb.FileCheck_Content{Content: &ipb.ContentCheck{Content: "Different content 1"}},
				},
				&ipb.FileCheck{
					FilesToCheck: []*ipb.FileSet{testconfigcreator.SingleFileWithPath(testFilePath)},
					CheckType:    &ipb.FileCheck_Content{Content: &ipb.ContentCheck{Content: "Different content 2"}},
				},
			},
			expectedResult: &apb.ComplianceResult{
				Id: "id",
				ComplianceOccurrence: &cpb.ComplianceOccurrence{
					NonCompliantFiles: []*cpb.NonCompliantFile{
						&cpb.NonCompliantFile{
							Path:   testFilePath,
							Reason: fmt.Sprintf("Got content %q, expected \"Different content 1\"", testFileContent),
						},
						&cpb.NonCompliantFile{
							Path:   testFilePath,
							Reason: fmt.Sprintf("Got content %q, expected \"Different content 2\"", testFileContent),
						},
					},
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			check := createFileCheckBatch(t, "id", tc.fileChecks, newFakeAPI())
			resultMap, _, err := check.Exec("")
			if err != nil {
				t.Fatalf("check.Exec() returned an error: %v", err)
			}
			result, gotSingleton := singleComplianceResult(resultMap)
			if !gotSingleton {
				t.Fatalf("check.Exec() expected to return 1 result, got %d", len(resultMap))
			}

			if diff := cmp.Diff(tc.expectedResult, result, protocmp.Transform()); diff != "" {
				t.Errorf("check.Exec() returned unexpected diff (-want +got):\n%s", diff)
			}
		})
	}
}

func TestTraversalOptOut(t *testing.T) {
	dirPath := "/non/existent/dir"
	optOutRegex := "/non/existe.*"

	fileChecks := []*ipb.FileCheck{&ipb.FileCheck{
		FilesToCheck: []*ipb.FileSet{&ipb.FileSet{
			FilePath: &ipb.FileSet_FilesInDir_{FilesInDir: &ipb.FileSet_FilesInDir{DirPath: dirPath}},
		}},
		CheckType: &ipb.FileCheck_Content{Content: &ipb.ContentCheck{Content: "content"}},
	}}

	scanInstruction := testconfigcreator.NewFileScanInstruction(fileChecks)
	config := testconfigcreator.NewBenchmarkConfig(t, "id", scanInstruction)
	checks, err := configchecks.CreateChecksFromConfig(
		context.Background(),
		&apb.ScanConfig{
			BenchmarkConfigs: []*apb.BenchmarkConfig{config},
			OptOutConfig: &apb.OptOutConfig{
				TraversalOptoutRegexes: []string{optOutRegex},
			},
		},
		newFakeAPI())
	if err != nil {
		t.Fatalf("configchecks.CreateChecksFromConfig([%v]) returned an error: %v", config, err)
	}
	if len(checks) != 1 {
		t.Fatalf("configchecks.CreateChecksFromConfig([%v]) created %d checks, expected 1",
			config, len(checks))
	}

	// The non-existent directory should be skipped and not cause an error.
	var pVal string
	if _, _, err := checks[0].Exec(pVal); err != nil {
		t.Errorf("check.Exec() returned an error: %v", err)
	}
}

func gzipString(str string) (string, error) {
	var b bytes.Buffer
	gz := gzip.NewWriter(&b)
	if _, err := gz.Write([]byte(str)); err != nil {
		return "", err
	}
	gz.Close()
	return b.String(), nil
}

func TestGzippedFileUnzipped(t *testing.T) {
	filePath := "file.gz"
	fileContent := "File content"
	gzipFileContent, err := gzipString(fileContent)
	if err != nil {
		t.Fatalf("gzipString(%s) returned an error: %v", fileContent, err)
	}
	fileChecks := []*ipb.FileCheck{&ipb.FileCheck{
		FilesToCheck: []*ipb.FileSet{testconfigcreator.SingleFileWithPath(filePath)},
		CheckType:    &ipb.FileCheck_Content{Content: &ipb.ContentCheck{Content: "Different content"}},
	}}
	expectedResult :=
		&apb.ComplianceResult{
			Id: "id",
			ComplianceOccurrence: &cpb.ComplianceOccurrence{
				NonCompliantFiles: []*cpb.NonCompliantFile{
					&cpb.NonCompliantFile{
						Path:   filePath,
						Reason: fmt.Sprintf("Got content %q, expected \"Different content\"", fileContent),
					},
				},
			},
		}

	check := createFileCheckBatch(t, "id", fileChecks, newFakeAPI(withFileContent(gzipFileContent)))
	resultMap, _, err := check.Exec("")
	if err != nil {
		t.Fatalf("check.Exec() returned an error: %v", err)
	}
	result, gotSingleton := singleComplianceResult(resultMap)
	if !gotSingleton {
		t.Fatalf("check.Exec() expected to return 1 result, got %d", len(resultMap))
	}

	if diff := cmp.Diff(expectedResult, result, protocmp.Transform()); diff != "" {
		t.Errorf("check.Exec() returned unexpected diff (-want +got):\n%s", diff)
	}
}

func TestRepeatConfigApplied(t *testing.T) {
	passwdContent := "user1:x:1337:1338::/home/user1:/bin/bash\n" +
		"user2:x:1337:1338::/home/user2:/bin/bash"
	fileChecks := []*ipb.FileCheck{&ipb.FileCheck{
		FilesToCheck: []*ipb.FileSet{testconfigcreator.SingleFileWithPath("$home/file.txt")},
		CheckType: &ipb.FileCheck_Permission{Permission: &ipb.PermissionCheck{
			User: &ipb.PermissionCheck_OwnerCheck{Name: "$user", ShouldOwn: true},
		}},
		RepeatConfig: &ipb.RepeatConfig{
			Type: ipb.RepeatConfig_FOR_EACH_USER_WITH_LOGIN,
		},
	}}
	expectedResult1 :=
		&apb.ComplianceResult{
			Id: "id",
			ComplianceOccurrence: &cpb.ComplianceOccurrence{
				NonCompliantFiles: []*cpb.NonCompliantFile{
					&cpb.NonCompliantFile{
						Path:   "/home/user1/file.txt",
						Reason: "Owner is root, expected it to be user1",
					},
				},
			},
		}
	expectedResult2 :=
		&apb.ComplianceResult{
			Id: "id",
			ComplianceOccurrence: &cpb.ComplianceOccurrence{
				NonCompliantFiles: []*cpb.NonCompliantFile{
					&cpb.NonCompliantFile{
						Path:   "/home/user2/file.txt",
						Reason: "Owner is root, expected it to be user2",
					},
				},
			},
		}

	scanInstruction := testconfigcreator.NewFileScanInstruction(fileChecks)
	config := testconfigcreator.NewBenchmarkConfig(t, "id", scanInstruction)
	checks, err := configchecks.CreateChecksFromConfig(
		context.Background(),
		&apb.ScanConfig{
			BenchmarkConfigs: []*apb.BenchmarkConfig{config},
		},
		newFakeAPI(withFileContent(passwdContent)))
	if err != nil {
		t.Fatalf("configchecks.CreateChecksFromConfig([%v]) returned an error: %v", config, err)
	}
	if len(checks) != 2 {
		t.Fatalf("configchecks.CreateChecksFromConfig([%v]) created %d checks, expected 2",
			config, len(checks))
	}

	var pVal string
	resultMap1, _, err := checks[0].Exec(pVal)
	if err != nil {
		t.Fatalf("checks[0].Exec() returned an error: %v", err)
	}
	result1, gotSingleton := singleComplianceResult(resultMap1)
	if !gotSingleton {
		t.Fatalf("checks[0].Exec() expected to return 1 result, got %d", len(resultMap1))
	}

	var pVal2 string
	resultMap2, _, err := checks[1].Exec(pVal2)
	if err != nil {
		t.Fatalf("checks[1].Exec() returned an error: %v", err)
	}
	result2, gotSingleton := singleComplianceResult(resultMap2)
	if !gotSingleton {
		t.Fatalf("checks[1].Exec() expected to return 1 result, got %d", len(resultMap2))
	}

	// The checks are created in an arbitrary order, check both results for a match.
	diff1 := cmp.Diff(expectedResult1, result1, protocmp.Transform())
	diff2 := cmp.Diff(expectedResult2, result1, protocmp.Transform())
	if diff1 != "" && diff2 != "" {
		t.Errorf("checks[0].Exec() returned unexpected diff (-want +got):\n%s%s", diff1, diff2)
	}
	diff1 = cmp.Diff(expectedResult1, result2, protocmp.Transform())
	diff2 = cmp.Diff(expectedResult2, result2, protocmp.Transform())
	if diff1 != "" && diff2 != "" {
		t.Errorf("checks[1].Exec() returned unexpected diff (-want +got):\n%s%s", diff1, diff2)
	}
}

func TestRepeatConfigCreationFails(t *testing.T) {
	passwdContent := "invalid\n"
	fileChecks := []*ipb.FileCheck{&ipb.FileCheck{
		FilesToCheck: []*ipb.FileSet{testconfigcreator.SingleFileWithPath("$home/file.txt")},
		CheckType:    &ipb.FileCheck_Existence{Existence: &ipb.ExistenceCheck{ShouldExist: true}},
		RepeatConfig: &ipb.RepeatConfig{Type: ipb.RepeatConfig_FOR_EACH_USER_WITH_LOGIN},
	}}
	want := &apb.ComplianceResult{
		Id: "id",
		ComplianceOccurrence: &cpb.ComplianceOccurrence{
			NonComplianceReason: "error creating RepeatConfig: can't parse line 1 in /etc/passwd: expected at least 7 tokens, got 1",
		},
	}

	scanInstruction := testconfigcreator.NewFileScanInstruction(fileChecks)
	config := testconfigcreator.NewBenchmarkConfig(t, "id", scanInstruction)
	checks, err := configchecks.CreateChecksFromConfig(
		context.Background(),
		&apb.ScanConfig{
			BenchmarkConfigs: []*apb.BenchmarkConfig{config},
		},
		newFakeAPI(withFileContent(passwdContent)))
	if err != nil {
		t.Fatalf("configchecks.CreateChecksFromConfig([%v]) returned an error: %v", config, err)
	}
	if len(checks) != 1 {
		t.Fatalf("configchecks.CreateChecksFromConfig([%v]) created %d checks, expected 2",
			config, len(checks))
	}

	var pVal string
	resultMap, _, err := checks[0].Exec(pVal)
	if err != nil {
		t.Fatalf("checks[0].Exec() returned an error: %v", err)
	}
	result, gotSingleton := singleComplianceResult(resultMap)
	if !gotSingleton {
		t.Fatalf("checks[0].Exec() expected to return 1 result, got %d", len(resultMap))
	}
	if diff := cmp.Diff(want, result, protocmp.Transform()); diff != "" {
		t.Errorf("checks[0].Exec() returned unexpected diff (-want +got):\n%s", diff)
	}
}

// A fake ScanAPIProvider implementation that returns a lot of files that are owned by root.
type manyFilesAPI struct{}

func (r *manyFilesAPI) OpenFile(ctx context.Context, filePath string) (io.ReadCloser, error) {
	return nil, errors.New("not implemented")
}

func (manyFilesAPI) OpenDir(ctx context.Context, filePath string) (scanapi.DirReader, error) {
	files := make([]*apb.DirContent, 0, 20)
	for i := 0; i < 20; i++ {
		files = append(files, &apb.DirContent{Name: fmt.Sprintf("file%d", i), IsDir: false})
	}
	return scanapi.SliceToDirReader(files), nil
}

func (manyFilesAPI) FilePermissions(ctx context.Context, filePath string) (*apb.PosixPermissions, error) {
	return &apb.PosixPermissions{User: "root"}, nil
}

func (manyFilesAPI) SQLQuery(ctx context.Context, query string) (string, error) {
	return "", errors.New("not implemented")
}

func (manyFilesAPI) SupportedDatabase() (ipb.SQLCheck_SQLDatabase, error) {
	return 0, errors.New("not implemented")
}

func TestLongCheckResultsPruned(t *testing.T) {
	// Set up a check that returns many non-compliant files.
	fileCheck := &ipb.FileCheck{
		FilesToCheck: []*ipb.FileSet{&ipb.FileSet{
			FilePath: &ipb.FileSet_FilesInDir_{FilesInDir: &ipb.FileSet_FilesInDir{DirPath: "/"}},
		}},
		CheckType: &ipb.FileCheck_Permission{Permission: &ipb.PermissionCheck{
			User: &ipb.PermissionCheck_OwnerCheck{Name: "root", ShouldOwn: false},
		}},
	}
	check := createFileCheckBatch(t, "id", []*ipb.FileCheck{fileCheck}, &manyFilesAPI{})

	resultMap, _, err := check.Exec("")
	if err != nil {
		t.Fatalf("check.Exec() returned an error: %v", err)
	}
	result, gotSingleton := singleComplianceResult(resultMap)
	if !gotSingleton {
		t.Fatalf("check.Exec() expected to return 1 result, got %d", len(resultMap))
	}
	count := len(result.GetComplianceOccurrence().GetNonCompliantFiles())
	if count != configchecks.MaxNonCompliantFiles {
		t.Errorf("want %d non-compliant files, check.Exec() returned %d:\n%s", configchecks.MaxNonCompliantFiles, count, result)
	}
}

func TestTimeout(t *testing.T) {
	smallTimeout := &dpb.Duration{Nanos: 1}
	bigTimeout := &dpb.Duration{Seconds: 60 * 60}
	zeroTimeout := &dpb.Duration{}
	testCases := []struct {
		description           string
		scanTimeout           *dpb.Duration
		benchmarkCheckTimeout *dpb.Duration
		expectTimeout         bool
	}{
		{
			description:           "Both timeouts small should timeout",
			scanTimeout:           smallTimeout,
			benchmarkCheckTimeout: smallTimeout,
			expectTimeout:         true,
		},
		{
			description:           "Scan timeout small should timeout",
			scanTimeout:           smallTimeout,
			benchmarkCheckTimeout: bigTimeout,
			expectTimeout:         true,
		},
		{
			description:           "Benchmark timeout small should timeout",
			scanTimeout:           bigTimeout,
			benchmarkCheckTimeout: smallTimeout,
			expectTimeout:         true,
		},
		{
			description:           "Both timeouts big shouldn't timeout",
			scanTimeout:           bigTimeout,
			benchmarkCheckTimeout: bigTimeout,
			expectTimeout:         false,
		},
		{
			description:           "Scan timeout at 0 should be ignored",
			scanTimeout:           zeroTimeout,
			benchmarkCheckTimeout: bigTimeout,
			expectTimeout:         false,
		},
		{
			description:           "Benchmark timeout at 0 should be ignored",
			scanTimeout:           bigTimeout,
			benchmarkCheckTimeout: zeroTimeout,
			expectTimeout:         false,
		},
		{
			description:           "Both timeouts at zero should be ignored",
			scanTimeout:           zeroTimeout,
			benchmarkCheckTimeout: zeroTimeout,
			expectTimeout:         false,
		},
	}

	benchmarkConfigs := []*apb.BenchmarkConfig{
		testconfigcreator.NewBenchmarkConfig(t, "id", testconfigcreator.NewFileScanInstruction([]*ipb.FileCheck{
			&ipb.FileCheck{
				FilesToCheck: []*ipb.FileSet{testconfigcreator.SingleFileWithPath("/path")},
				CheckType:    &ipb.FileCheck_Content{Content: &ipb.ContentCheck{Content: "content"}},
			},
		})),
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			config := &apb.ScanConfig{
				BenchmarkConfigs:      benchmarkConfigs,
				ScanTimeout:           tc.scanTimeout,
				BenchmarkCheckTimeout: tc.benchmarkCheckTimeout,
			}
			checks, err := configchecks.CreateChecksFromConfig(context.Background(), config, newFakeAPI(withFileContent("content")))
			if err != nil {
				t.Fatalf("configchecks.CreateChecksFromConfig([%v]) returned an error: %v", config, err)
			}
			var pVal string
			_, _, err = checks[0].Exec(pVal)
			if err != nil && err.Error() != "scan timed out" {
				t.Fatalf("check.Exec() with {ScanTimeout: %v, BenchmarkCheckTimeout: %v} returned an unexpected error: %v",
					tc.scanTimeout.AsDuration(), tc.benchmarkCheckTimeout.AsDuration(), err)
			}
			if tc.expectTimeout && err == nil {
				t.Fatalf("check.Exec() with {ScanTimeout: %v, BenchmarkCheckTimeout: %v} was expected to timeout but didn't",
					tc.scanTimeout.AsDuration(), tc.benchmarkCheckTimeout.AsDuration())
			}
			if !tc.expectTimeout && err != nil {
				t.Fatalf("check.Exec() with {ScanTimeout: %v, BenchmarkCheckTimeout: %v} was expected to succeed but timed out",
					tc.scanTimeout.AsDuration(), tc.benchmarkCheckTimeout.AsDuration())
			}
		})
	}
}
