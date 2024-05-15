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
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"google.golang.org/protobuf/testing/protocmp"
	cpb "github.com/google/localtoast/scannerlib/proto/compliance_go_proto"
	"github.com/google/localtoast/scannerlib/configchecks"
	apb "github.com/google/localtoast/scannerlib/proto/api_go_proto"
	ipb "github.com/google/localtoast/scannerlib/proto/scan_instructions_go_proto"
	"github.com/google/localtoast/scannerlib/testconfigcreator"
)

func TestSQLCheckCreation(t *testing.T) {
	testCases := []struct {
		desc     string
		sqlCheck *ipb.SQLCheck
	}{
		{
			desc: "MySQL",
			sqlCheck: &ipb.SQLCheck{
				TargetDatabase: ipb.SQLCheck_DB_MYSQL,
				Query:          "SELECT 1;",
				ExpectResults:  true,
			}},
		{
			desc: "Cassandra",
			sqlCheck: &ipb.SQLCheck{
				TargetDatabase: ipb.SQLCheck_DB_CASSANDRA,
				Query:          "SELECT 1;",
				ExpectResults:  true,
			}},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			scanInstruction := testconfigcreator.NewSQLScanInstruction([]*ipb.SQLCheck{tc.sqlCheck})
			config := testconfigcreator.NewBenchmarkConfig(t, "id", scanInstruction)

			newchecks, err := configchecks.CreateChecksFromConfig(
				context.Background(),
				&apb.ScanConfig{
					BenchmarkConfigs: []*apb.BenchmarkConfig{config},
				},
				newFakeAPI(withSupportedDatabase(tc.sqlCheck.TargetDatabase)))
			if err != nil {
				t.Fatalf("configchecks.CreateChecksFromConfig([%v]) returned an error: %v", config, err)
			}
			if len(newchecks) != 1 {
				t.Fatalf("Expected 1 check to be created, got %d", len(newchecks))
			}
			expectedIDs := []string{"id"}
			if diff := cmp.Diff(expectedIDs, newchecks[0].BenchmarkIDs()); diff != "" {
				t.Errorf("%v.BenchmarkIDs() returned unexpected diff (-want +got):\n%s", newchecks[0], diff)
			}
		})
	}
}

func TestSQLCheckWithEmptyInstructionsReturnsError(t *testing.T) {
	scanInstruction := &ipb.BenchmarkScanInstruction{
		CheckAlternatives: []*ipb.CheckAlternative{},
	}
	config := testconfigcreator.NewBenchmarkConfig(t, "id", scanInstruction)

	if _, err := configchecks.CreateChecksFromConfig(
		context.Background(),
		&apb.ScanConfig{
			BenchmarkConfigs: []*apb.BenchmarkConfig{config},
		},
		newFakeAPI()); err == nil {
		t.Errorf("configchecks.CreateChecksFromConfig([%v]) didn't return an error", config)
	}
}

func TestSQLCheckUnsupportedTypeReturnsError(t *testing.T) {
	testCases := []struct {
		desc        string
		sqlCheck    *ipb.SQLCheck
		supportedDB ipb.SQLCheck_SQLDatabase
	}{
		{
			desc: "Unspecified",
			sqlCheck: &ipb.SQLCheck{
				TargetDatabase: ipb.SQLCheck_DB_UNSPECIFIED,
				Query:          "SELECT 1;",
				ExpectResults:  true,
			},
			supportedDB: ipb.SQLCheck_DB_MYSQL,
		},
		{
			desc: "Wrong DB type",
			sqlCheck: &ipb.SQLCheck{
				TargetDatabase: ipb.SQLCheck_DB_CASSANDRA,
				Query:          "SELECT 1;",
				ExpectResults:  true,
			},
			supportedDB: ipb.SQLCheck_DB_MYSQL,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			scanInstruction := testconfigcreator.NewSQLScanInstruction([]*ipb.SQLCheck{tc.sqlCheck})
			config := testconfigcreator.NewBenchmarkConfig(t, "id", scanInstruction)

			if _, err := configchecks.CreateChecksFromConfig(
				context.Background(),
				&apb.ScanConfig{
					BenchmarkConfigs: []*apb.BenchmarkConfig{config},
				},
				newFakeAPI()); err == nil {
				t.Errorf("configchecks.CreateChecksFromConfig([%v]) didn't return an error", config)
			}
		})
	}
}

func createMySQLCheck(t *testing.T, id string, sqlChecks []*ipb.SQLCheck, api *fakeAPI) configchecks.BenchmarkCheck {
	t.Helper()
	scanInstruction := testconfigcreator.NewSQLScanInstruction(sqlChecks)
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

func TestMySQLCheckComplianceResults(t *testing.T) {
	testCases := []struct {
		desc           string
		sqlCheck       *ipb.SQLCheck
		expectedResult *apb.ComplianceResult
	}{
		{
			desc: "expect rows, get one row",
			sqlCheck: &ipb.SQLCheck{
				TargetDatabase: ipb.SQLCheck_DB_MYSQL,
				Query:          fakeQueryOneRow,
				ExpectResults:  true,
			},
			expectedResult: &apb.ComplianceResult{
				Id:                   "id",
				ComplianceOccurrence: &cpb.ComplianceOccurrence{},
			},
		},
		{
			desc: "expect no rows, get no rows",
			sqlCheck: &ipb.SQLCheck{
				TargetDatabase: ipb.SQLCheck_DB_MYSQL,
				Query:          fakeQueryNoRows,
				ExpectResults:  false,
			},
			expectedResult: &apb.ComplianceResult{
				Id:                   "id",
				ComplianceOccurrence: &cpb.ComplianceOccurrence{},
			},
		},
		{
			desc: "expect rows, get no rows",
			sqlCheck: &ipb.SQLCheck{
				TargetDatabase: ipb.SQLCheck_DB_MYSQL,
				Query:          fakeQueryNoRows,
				ExpectResults:  true,
			},
			expectedResult: &apb.ComplianceResult{
				Id: "id",
				ComplianceOccurrence: &cpb.ComplianceOccurrence{
					NonComplianceReason: fmt.Sprintf("Expected results for query %q, but got none.", fakeQueryNoRows),
				},
			},
		},
		{
			desc: "expect no rows, get one row",
			sqlCheck: &ipb.SQLCheck{
				TargetDatabase: ipb.SQLCheck_DB_MYSQL,
				Query:          fakeQueryOneRow,
				ExpectResults:  false,
			},
			expectedResult: &apb.ComplianceResult{
				Id: "id",
				ComplianceOccurrence: &cpb.ComplianceOccurrence{
					NonComplianceReason: fmt.Sprintf("Expected no results for query %q, but got some.", fakeQueryOneRow),
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			check := createMySQLCheck(t, "id", []*ipb.SQLCheck{tc.sqlCheck}, newFakeAPI())

			var pVal string
			resultMap, _, err := check.Exec(pVal)
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

func TestMySQLCustomNonComplianceMessage(t *testing.T) {
	reason := "custom reason"
	check := &ipb.SQLCheck{
		TargetDatabase:   ipb.SQLCheck_DB_MYSQL,
		Query:            fakeQueryOneRow,
		ExpectResults:    false,
		NonComplianceMsg: reason,
	}
	scanInstruction := testconfigcreator.NewSQLScanInstruction([]*ipb.SQLCheck{check})
	config := testconfigcreator.NewBenchmarkConfig(t, "id", scanInstruction)
	newchecks, err := configchecks.CreateChecksFromConfig(
		context.Background(),
		&apb.ScanConfig{
			BenchmarkConfigs: []*apb.BenchmarkConfig{config},
		},
		newFakeAPI())
	if err != nil {
		t.Fatalf("configchecks.CreateChecksFromConfig([%v]) returned an error: %v", config, err)
	}
	if len(newchecks) != 1 {
		t.Fatalf("Expected 1 check to be created, got %d", len(newchecks))
	}

	var pVal string
	newcheck := newchecks[0]
	resultMap, _, err := newcheck.Exec(pVal)
	if err != nil {
		t.Fatalf("newcheck.Exec() returned an error: %v", err)
	}
	result, gotSingleton := singleComplianceResult(resultMap)
	if !gotSingleton {
		t.Fatalf("newcheck.Exec() expected to return 1 result, got %d", len(resultMap))
	}

	if diff := cmp.Diff(&apb.ComplianceResult{
		Id: "id",
		ComplianceOccurrence: &cpb.ComplianceOccurrence{
			NonComplianceReason: reason,
		},
	}, result, protocmp.Transform()); diff != "" {
		t.Errorf("check.Exec() returned unexpected diff (-want +got):\n%s", diff)
	}
}

func TestMySQLCheckPropagatesError(t *testing.T) {
	check := createMySQLCheck(t, "id", []*ipb.SQLCheck{{
		TargetDatabase: ipb.SQLCheck_DB_MYSQL,
		Query:          fakeQueryError,
		ExpectResults:  true,
	}}, newFakeAPI())
	var pVal string
	_, _, err := check.Exec(pVal)
	if err == nil {
		t.Errorf("check.Exec() didn't return an error")
	}
	if !strings.Contains(err.Error(), queryErrorMsg) {
		t.Errorf("check.Exec returned the wrong error: want %q, got %v", queryErrorMsg, err)
	}
}
