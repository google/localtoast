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
	"fmt"
	"io"
	"os"
	"sort"
	"testing"

	"github.com/google/go-cmp/cmp"
	"google.golang.org/protobuf/testing/protocmp"
	gpb "github.com/google/localtoast/library/proto/compliance_go_proto"
	"github.com/google/localtoast/library/configchecks"
	"github.com/google/localtoast/library/fileset"
	apb "github.com/google/localtoast/library/proto/api_go_proto"
	ipb "github.com/google/localtoast/library/proto/scan_instructions_go_proto"
	"github.com/google/localtoast/library/testing/testconfigcreator"
)

func TestChecksOfSameTypeOnSameFileGroupedTogether(t *testing.T) {
	scanInstruction := testconfigcreator.NewFileScanInstruction([]*ipb.FileCheck{
		&ipb.FileCheck{
			FilesToCheck: []*ipb.FileSet{testconfigcreator.SingleFileWithPath("/path")},
			CheckType:    &ipb.FileCheck_Existence{Existence: &ipb.ExistenceCheck{ShouldExist: true}},
		},
	})
	config1 := testconfigcreator.NewBenchmarkConfig(t, "id1", scanInstruction)
	config2 := testconfigcreator.NewBenchmarkConfig(t, "id2", scanInstruction)

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

func TestDifferentChecksGroupedSeparately(t *testing.T) {
	type fileCheckPair struct {
		fileCheck1 *ipb.FileCheck
		fileCheck2 *ipb.FileCheck
	}
	pairsToTest := []*fileCheckPair{
		// Same check types but different files to check.
		&fileCheckPair{
			fileCheck1: &ipb.FileCheck{
				FilesToCheck: []*ipb.FileSet{testconfigcreator.SingleFileWithPath("/path1")},
				CheckType:    &ipb.FileCheck_Existence{Existence: &ipb.ExistenceCheck{ShouldExist: true}},
			},
			fileCheck2: &ipb.FileCheck{
				FilesToCheck: []*ipb.FileSet{testconfigcreator.SingleFileWithPath("/path2")},
				CheckType:    &ipb.FileCheck_Existence{Existence: &ipb.ExistenceCheck{ShouldExist: false}},
			},
		},
		// Same files to check but different check types.
		&fileCheckPair{
			fileCheck1: &ipb.FileCheck{
				FilesToCheck: []*ipb.FileSet{testconfigcreator.SingleFileWithPath("/path")},
				CheckType:    &ipb.FileCheck_Existence{Existence: &ipb.ExistenceCheck{ShouldExist: true}},
			},
			fileCheck2: &ipb.FileCheck{
				FilesToCheck: []*ipb.FileSet{testconfigcreator.SingleFileWithPath("/path")},
				CheckType:    &ipb.FileCheck_Content{Content: &ipb.ContentCheck{Content: "content"}},
			},
		},
		// Same files and check types but different delimiters.
		&fileCheckPair{
			&ipb.FileCheck{
				FilesToCheck: []*ipb.FileSet{testconfigcreator.SingleFileWithPath("/path")},
				CheckType: &ipb.FileCheck_ContentEntry{ContentEntry: &ipb.ContentEntryCheck{
					Delimiter: []byte{'\n'},
				}},
			},
			&ipb.FileCheck{
				FilesToCheck: []*ipb.FileSet{testconfigcreator.SingleFileWithPath("/path")},
				CheckType: &ipb.FileCheck_ContentEntry{ContentEntry: &ipb.ContentEntryCheck{
					Delimiter: []byte{';'},
				}},
			},
		},
	}

	for _, pair := range pairsToTest {
		scanInstruction1 := testconfigcreator.NewFileScanInstruction([]*ipb.FileCheck{pair.fileCheck1})
		scanInstruction2 := testconfigcreator.NewFileScanInstruction([]*ipb.FileCheck{pair.fileCheck2})
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
}

func TestIdsDeduplicated(t *testing.T) {
	type fileCheckPair struct {
		fileCheck1 *ipb.FileCheck
		fileCheck2 *ipb.FileCheck
	}
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

func createFileCheckBatch(t *testing.T, id string, fileChecks []*ipb.FileCheck, api *fakeAPI) configchecks.BenchmarkCheck {
	t.Helper()
	scanInstruction := testconfigcreator.NewFileScanInstruction(fileChecks)
	config := testconfigcreator.NewBenchmarkConfig(t, "id", scanInstruction)

	checks, err := configchecks.CreateChecksFromConfig(
		context.Background(),
		&apb.ScanConfig{
			BenchmarkConfigs: []*apb.BenchmarkConfig{config},
		},
		api)
	if err != nil {
		t.Fatalf("configchecks.CreateChecksFromConfig([%v]) returned an error: %v", config, err)
	}
	if len(checks) != 1 {
		t.Fatalf("Created %d checks, expected only 1", len(checks))
	}
	return checks[0]
}

func TestFileCustomNonComplianceMessage(t *testing.T) {
	testCases := []struct {
		fileCheck      *ipb.FileCheck
		expectedResult *apb.ComplianceResult
	}{
		{
			fileCheck: &ipb.FileCheck{
				FilesToCheck:     []*ipb.FileSet{testconfigcreator.SingleFileWithPath(testFilePath)},
				CheckType:        &ipb.FileCheck_Existence{Existence: &ipb.ExistenceCheck{ShouldExist: false}},
				NonComplianceMsg: "custom reason",
			},
			expectedResult: &apb.ComplianceResult{
				Id: "id",
				ComplianceOccurrence: &gpb.ComplianceOccurrence{
					NonCompliantFiles: []*gpb.NonCompliantFile{
						&gpb.NonCompliantFile{
							Path:   testFilePath,
							Reason: "custom reason",
						},
					},
				},
			},
		},
		{
			fileCheck: &ipb.FileCheck{
				FilesToCheck:       []*ipb.FileSet{testconfigcreator.SingleFileWithPath(testFilePath)},
				CheckType:          &ipb.FileCheck_Existence{Existence: &ipb.ExistenceCheck{ShouldExist: false}},
				NonComplianceMsg:   "custom reason",
				FileDisplayCommand: "display command",
			},
			expectedResult: &apb.ComplianceResult{
				Id: "id",
				ComplianceOccurrence: &gpb.ComplianceOccurrence{
					NonCompliantFiles: []*gpb.NonCompliantFile{
						&gpb.NonCompliantFile{
							DisplayCommand: "display command",
							Reason:         "custom reason",
						},
					},
				},
			},
		},
		{
			fileCheck: &ipb.FileCheck{
				FilesToCheck:       []*ipb.FileSet{testconfigcreator.SingleFileWithPath(testFilePath)},
				CheckType:          &ipb.FileCheck_Existence{Existence: &ipb.ExistenceCheck{ShouldExist: true}},
				FileDisplayCommand: "display command",
			},
			expectedResult: &apb.ComplianceResult{
				Id:                   "id",
				ComplianceOccurrence: &gpb.ComplianceOccurrence{},
			},
		},
	}

	for _, tc := range testCases {
		check := createFileCheckBatch(t, "id", []*ipb.FileCheck{tc.fileCheck}, newFakeAPI())
		resultMap, err := check.Exec()
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
	}
}

func TestFileDisplayCommandWithoutCustomMessage(t *testing.T) {
	fileCheck := &ipb.FileCheck{
		FilesToCheck:       []*ipb.FileSet{testconfigcreator.SingleFileWithPath(testFilePath)},
		CheckType:          &ipb.FileCheck_Existence{Existence: &ipb.ExistenceCheck{ShouldExist: false}},
		FileDisplayCommand: "display command",
	}

	check := createFileCheckBatch(t, "id", []*ipb.FileCheck{fileCheck}, newFakeAPI())
	if _, err := check.Exec(); err == nil {
		t.Fatalf("check.Exec() didn't return an error")
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

	resultMap, err := check.Exec()
	if err != nil {
		t.Fatalf("check.Exec() returned an error: %v", err)
	}

	if len(resultMap) != 2 {
		t.Errorf("Expected results to be present for 2 check alternatives, found %d: %v", len(resultMap), resultMap)
	}
}

func TestFileExistenceCheckComplianceResults(t *testing.T) {
	testCases := []struct {
		fileCheck      *ipb.FileCheck
		expectedResult *apb.ComplianceResult
	}{
		{
			fileCheck: &ipb.FileCheck{
				FilesToCheck: []*ipb.FileSet{testconfigcreator.SingleFileWithPath(testFilePath)},
				CheckType:    &ipb.FileCheck_Existence{Existence: &ipb.ExistenceCheck{ShouldExist: true}},
			},
			expectedResult: &apb.ComplianceResult{
				Id:                   "id",
				ComplianceOccurrence: &gpb.ComplianceOccurrence{},
			},
		},
		{
			fileCheck: &ipb.FileCheck{
				FilesToCheck: []*ipb.FileSet{testconfigcreator.SingleFileWithPath(nonExistentFilePath)},
				CheckType:    &ipb.FileCheck_Existence{Existence: &ipb.ExistenceCheck{ShouldExist: false}},
			},
			expectedResult: &apb.ComplianceResult{
				Id:                   "id",
				ComplianceOccurrence: &gpb.ComplianceOccurrence{},
			},
		},
		{
			fileCheck: &ipb.FileCheck{
				FilesToCheck: []*ipb.FileSet{testconfigcreator.SingleFileWithPath(nonExistentFilePath)},
				CheckType:    &ipb.FileCheck_Existence{Existence: &ipb.ExistenceCheck{ShouldExist: true}},
			},
			expectedResult: &apb.ComplianceResult{
				Id: "id",
				ComplianceOccurrence: &gpb.ComplianceOccurrence{
					NonCompliantFiles: []*gpb.NonCompliantFile{
						&gpb.NonCompliantFile{
							Path:   fileset.FileSetToString(testconfigcreator.SingleFileWithPath(nonExistentFilePath)),
							Reason: "File doesn't exist but it should",
						},
					},
				},
			},
		},
		{
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
				ComplianceOccurrence: &gpb.ComplianceOccurrence{
					NonCompliantFiles: []*gpb.NonCompliantFile{
						&gpb.NonCompliantFile{
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
			fileCheck: &ipb.FileCheck{
				FilesToCheck: []*ipb.FileSet{testconfigcreator.SingleFileWithPath(testFilePath)},
				CheckType:    &ipb.FileCheck_Existence{Existence: &ipb.ExistenceCheck{ShouldExist: false}},
			},
			expectedResult: &apb.ComplianceResult{
				Id: "id",
				ComplianceOccurrence: &gpb.ComplianceOccurrence{
					NonCompliantFiles: []*gpb.NonCompliantFile{
						&gpb.NonCompliantFile{
							Path:   testFilePath,
							Reason: "File exists but it shouldn't",
						},
					},
				},
			},
		},
	}

	for _, tc := range testCases {
		check := createFileCheckBatch(t, "id", []*ipb.FileCheck{tc.fileCheck}, newFakeAPI())
		resultMap, err := check.Exec()
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
	}
}

func TestFileExistenceCheckPropagatesError(t *testing.T) {
	check := createFileCheckBatch(t, "id", []*ipb.FileCheck{&ipb.FileCheck{
		FilesToCheck: []*ipb.FileSet{testconfigcreator.SingleFileWithPath(unreadableFilePath)},
		CheckType:    &ipb.FileCheck_Existence{Existence: &ipb.ExistenceCheck{ShouldExist: true}},
	}}, newFakeAPI())
	if _, err := check.Exec(); err == nil {
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
	expectedResult := &apb.ComplianceResult{Id: "id", ComplianceOccurrence: &gpb.ComplianceOccurrence{}}

	check := createFileCheckBatch(t, "id", []*ipb.FileCheck{fileCheck}, newFakeAPI(withOpenFileFunc(openFileFunc)))
	resultMap, err := check.Exec()

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
				ComplianceOccurrence: &gpb.ComplianceOccurrence{},
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
				ComplianceOccurrence: &gpb.ComplianceOccurrence{
					NonCompliantFiles: []*gpb.NonCompliantFile{
						&gpb.NonCompliantFile{
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
				ComplianceOccurrence: &gpb.ComplianceOccurrence{
					NonCompliantFiles: []*gpb.NonCompliantFile{
						&gpb.NonCompliantFile{
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
				ComplianceOccurrence: &gpb.ComplianceOccurrence{
					NonCompliantFiles: []*gpb.NonCompliantFile{
						&gpb.NonCompliantFile{
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
				ComplianceOccurrence: &gpb.ComplianceOccurrence{},
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
				ComplianceOccurrence: &gpb.ComplianceOccurrence{
					NonCompliantFiles: []*gpb.NonCompliantFile{
						&gpb.NonCompliantFile{
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
				ComplianceOccurrence: &gpb.ComplianceOccurrence{
					NonCompliantFiles: []*gpb.NonCompliantFile{
						&gpb.NonCompliantFile{
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
				ComplianceOccurrence: &gpb.ComplianceOccurrence{
					NonCompliantFiles: []*gpb.NonCompliantFile{
						&gpb.NonCompliantFile{
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
				ComplianceOccurrence: &gpb.ComplianceOccurrence{
					NonCompliantFiles: []*gpb.NonCompliantFile{
						&gpb.NonCompliantFile{
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
				ComplianceOccurrence: &gpb.ComplianceOccurrence{
					NonCompliantFiles: []*gpb.NonCompliantFile{
						&gpb.NonCompliantFile{
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
				ComplianceOccurrence: &gpb.ComplianceOccurrence{
					NonCompliantFiles: []*gpb.NonCompliantFile{
						&gpb.NonCompliantFile{
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
				ComplianceOccurrence: &gpb.ComplianceOccurrence{
					NonCompliantFiles: []*gpb.NonCompliantFile{
						&gpb.NonCompliantFile{
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

			resultMap, err := check.Exec()
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
				ComplianceOccurrence: &gpb.ComplianceOccurrence{},
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
				ComplianceOccurrence: &gpb.ComplianceOccurrence{
					NonCompliantFiles: []*gpb.NonCompliantFile{
						&gpb.NonCompliantFile{
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
				ComplianceOccurrence: &gpb.ComplianceOccurrence{
					NonCompliantFiles: []*gpb.NonCompliantFile{
						&gpb.NonCompliantFile{
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
				ComplianceOccurrence: &gpb.ComplianceOccurrence{
					NonCompliantFiles: []*gpb.NonCompliantFile{
						&gpb.NonCompliantFile{
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
				ComplianceOccurrence: &gpb.ComplianceOccurrence{
					NonCompliantFiles: []*gpb.NonCompliantFile{
						&gpb.NonCompliantFile{
							Path:   testFilePath,
							Reason: fmt.Sprintf("Got content %q, expected \"Different content 1\"", testFileContent),
						},
						&gpb.NonCompliantFile{
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
			resultMap, err := check.Exec()
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
	if _, err := checks[0].Exec(); err != nil {
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
			ComplianceOccurrence: &gpb.ComplianceOccurrence{
				NonCompliantFiles: []*gpb.NonCompliantFile{
					&gpb.NonCompliantFile{
						Path:   filePath,
						Reason: fmt.Sprintf("Got content %q, expected \"Different content\"", fileContent),
					},
				},
			},
		}

	check := createFileCheckBatch(t, "id", fileChecks, newFakeAPI(withFileContent(gzipFileContent)))
	resultMap, err := check.Exec()
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
			ComplianceOccurrence: &gpb.ComplianceOccurrence{
				NonCompliantFiles: []*gpb.NonCompliantFile{
					&gpb.NonCompliantFile{
						Path:   "/home/user1/file.txt",
						Reason: "Owner is root, expected it to be user1",
					},
				},
			},
		}
	expectedResult2 :=
		&apb.ComplianceResult{
			Id: "id",
			ComplianceOccurrence: &gpb.ComplianceOccurrence{
				NonCompliantFiles: []*gpb.NonCompliantFile{
					&gpb.NonCompliantFile{
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

	resultMap1, err := checks[0].Exec()
	if err != nil {
		t.Fatalf("checks[0].Exec() returned an error: %v", err)
	}
	result1, gotSingleton := singleComplianceResult(resultMap1)
	if !gotSingleton {
		t.Fatalf("checks[0].Exec() expected to return 1 result, got %d", len(resultMap1))
	}

	resultMap2, err := checks[1].Exec()
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
		ComplianceOccurrence: &gpb.ComplianceOccurrence{
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

	resultMap, err := checks[0].Exec()
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
