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

// Package scannerlib provides an interface for running compliance and config
// security checks on a machine.
package scannerlib

import (
	"context"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	cpb "github.com/google/localtoast/scannerlib/proto/compliance_go_proto"
	"github.com/google/localtoast/scanapi"
	"github.com/google/localtoast/scannerlib/configchecks"

	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
	apb "github.com/google/localtoast/scannerlib/proto/api_go_proto"
)

// Scanner is the main entry point of the scanner.
type Scanner struct{}

// Scan executes the scan for benchmark compliance using the provided scan
// config and API for accessing the scanned machine.
func (Scanner) Scan(ctx context.Context, config *apb.ScanConfig, api scanapi.ScanAPI) (*apb.ScanResults, error) {
	benchmarkConfigs := config.GetBenchmarkConfigs()
	if err := validateBenchmarkConfigs(benchmarkConfigs); err != nil {
		return nil, err
	}
	scanStartTime := time.Now()
	checks, err := configchecks.CreateChecksFromConfig(ctx, config, &apiErrorWrapper{api: api})
	if err != nil {
		return nil, err
	}

	checkResults, benchmarkErrors := executeChecks(checks)
	configchecks.AddBenchmarkVersionToResults(checkResults, benchmarkConfigs)
	complianceResults := determineBenchmarkCompliance(checks, checkResults, benchmarkErrors)

	benchmarkVersion, err := oldestBenchmarkVersion(config.GetBenchmarkConfigs())
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error while determining oldest benchmark version: %v\n", err)
		benchmarkVersion = "0.0.0"
	}

	options := newScanResultsOptions{
		startTime:              scanStartTime,
		benchmarkVersion:       benchmarkVersion,
		benchmarkDocument:      getBenchmarkDocument(config.GetBenchmarkConfigs()),
		compliantBenchmarks:    complianceResults.compliantBenchmarks,
		nonCompliantBenchmarks: complianceResults.nonCompliantBenchmarks,
	}
	if len(complianceResults.unknownBenchmarks) == 0 {
		options.status = apb.ScanStatus_SUCCEEDED
	} else {
		options.status = apb.ScanStatus_FAILED
		errorStrings := ""
		for _, errs := range benchmarkErrors {
			for _, err := range errs {
				errorStrings += err.Error() + "\n"
			}
		}
		options.failureReason = fmt.Sprintf(
			"Compliance state of the following benchmarks couldn't be determined: [%s]\n"+
				"The following errors were encountered while running the checks:\n%s",
			strings.Join(complianceResults.unknownBenchmarks, ","), errorStrings)
	}
	return newScanResults(options), nil
}

func validateBenchmarkConfigs(configs []*apb.BenchmarkConfig) error {
	ids := map[string]bool{}
	for _, config := range configs {
		id := config.GetId()
		if ids[id] {
			return fmt.Errorf("duplicate benchmark ID %s", id)
		}
		ids[id] = true

		if err := configchecks.ValidateScanInstructions(config); err != nil {
			return err
		}
	}
	return nil
}

// executeChecks runs the given benchmarkChecks and returns their findings.
// It also returns a map of benchmark IDs to errors that the benchmark's checks
// produced while running, or an empty slice if no errors occurred for a benchmark.
func executeChecks(checks []configchecks.BenchmarkCheck) ([]*apb.ComplianceResult, map[string][]error) {
	compliancePerAlternative := make(configchecks.ComplianceMap)
	benchmarkErrors := make(map[string][]error)
	var prvRes string
	for _, check := range checks {
		checkResults, res, err := check.Exec(prvRes)

		if len(res) > 0 {
			prvRes = res
		}

		if err != nil {
			fmt.Fprintf(os.Stderr, "Error while executing check %s: %v\n", check, err)
		}
		for _, id := range check.BenchmarkIDs() {
			if _, ok := benchmarkErrors[id]; !ok {
				benchmarkErrors[id] = []error{}
			}
			if err != nil {
				benchmarkErrors[id] = append(benchmarkErrors[id], fmt.Errorf("%s: %s", check, err))
			}
		}
		// Merge the check results into a unified compliance map.
		for altNum, compliance := range checkResults {
			if prev, ok := compliancePerAlternative[altNum]; ok {
				appendToComplianceResult(prev, compliance)
			} else {
				compliancePerAlternative[altNum] = compliance
			}
		}
	}

	result := make([]*apb.ComplianceResult, 0, len(compliancePerAlternative))
	for _, c := range compliancePerAlternative {
		result = append(result, c)
	}
	return result, benchmarkErrors
}

type benchmarkCompliance struct {
	compliantBenchmarks    []*apb.ComplianceResult
	nonCompliantBenchmarks []*apb.ComplianceResult
	unknownBenchmarks      []string
}

// determineBenchmarkCompliance takes the results of the check runs and aggregates them to figure
// out which benchmarks were compliant, which weren't, and which benchmarks' compliance couldn't be
// determined.
func determineBenchmarkCompliance(checks []configchecks.BenchmarkCheck, checkResults []*apb.ComplianceResult, benchmarkErrors map[string][]error) *benchmarkCompliance {
	complianceResultForBenchmark := make(map[string]*apb.ComplianceResult)
	for _, c := range checkResults {
		// A benchmark is compliant if any of its check alternatives is compliant.
		if prev, ok := complianceResultForBenchmark[c.GetId()]; ok {
			if isResultCompliant(prev) {
				// Result is already compliant.
				continue
			}
			if isResultCompliant(c) {
				// Overwrite any previous non-compliant results.
				complianceResultForBenchmark[c.GetId()] = c
			} else {
				// Merge the two non-compliant results.
				appendToComplianceResult(prev, c)
			}
		} else {
			complianceResultForBenchmark[c.GetId()] = c
		}
	}

	result := &benchmarkCompliance{
		compliantBenchmarks:    []*apb.ComplianceResult{},
		nonCompliantBenchmarks: []*apb.ComplianceResult{},
		unknownBenchmarks:      []string{},
	}
	for id, err := range benchmarkErrors {
		if len(err) > 0 {
			// Some checks haven't completed successfully, compliance can't be determined.
			result.unknownBenchmarks = append(result.unknownBenchmarks, id)
			continue
		}
		c := complianceResultForBenchmark[id]
		if !isResultCompliant(c) {
			// The benchmark is not compliant.
			result.nonCompliantBenchmarks = append(result.nonCompliantBenchmarks, c)
		} else {
			// The benchmark is compliant.
			result.compliantBenchmarks = append(result.compliantBenchmarks, c)
		}
	}
	deduplicateNonCompliantFiles(result.nonCompliantBenchmarks)

	return result
}

func isResultCompliant(c *apb.ComplianceResult) bool {
	return len(c.GetComplianceOccurrence().NonCompliantFiles) == 0 && c.GetComplianceOccurrence().NonComplianceReason == ""
}

func appendToComplianceResult(currentResult, newResult *apb.ComplianceResult) {
	current := currentResult.GetComplianceOccurrence()
	new := newResult.GetComplianceOccurrence()
	current.NonCompliantFiles = append(current.NonCompliantFiles, new.NonCompliantFiles...)
	if new.NonComplianceReason != "" {
		if current.NonComplianceReason != "" {
			current.NonComplianceReason = current.NonComplianceReason + "\n" + new.NonComplianceReason
		} else {
			current.NonComplianceReason = new.NonComplianceReason
		}
	}
}

// deduplicateNonCompliantFiles removes duplicate non-compliant file entries from the
// given compliance findings.
func deduplicateNonCompliantFiles(compliances []*apb.ComplianceResult) {
	type mapKey struct {
		path           string
		displayCommand string
		reason         string
	}
	for _, c := range compliances {
		occ := c.GetComplianceOccurrence()
		fileMap := make(map[mapKey]bool)
		for _, f := range occ.NonCompliantFiles {
			fileMap[mapKey{
				path:           f.Path,
				displayCommand: f.DisplayCommand,
				reason:         f.Reason,
			}] = true
		}
		occ.NonCompliantFiles = make([]*cpb.NonCompliantFile, 0, len(fileMap))
		for key, _ := range fileMap {
			f := &cpb.NonCompliantFile{
				Path:           key.path,
				Reason:         key.reason,
				DisplayCommand: key.displayCommand,
			}
			occ.NonCompliantFiles = append(occ.NonCompliantFiles, f)
		}
		// Keep the file list sorted for a better overview of non-compliance reports.
		sort.Slice(occ.NonCompliantFiles, func(i, j int) bool {
			f1 := occ.NonCompliantFiles[i]
			f2 := occ.NonCompliantFiles[j]
			if f1.Path != f2.Path {
				return f1.Path < f2.Path
			}
			if f1.Reason != f2.Reason {
				return f1.Reason < f2.Reason
			}
			return f1.DisplayCommand < f2.DisplayCommand
		})
	}
}

type newScanResultsOptions struct {
	startTime              time.Time
	benchmarkVersion       string
	benchmarkDocument      string
	compliantBenchmarks    []*apb.ComplianceResult
	nonCompliantBenchmarks []*apb.ComplianceResult
	status                 apb.ScanStatus_ScanStatusEnum
	failureReason          string
}

func newScanResults(options newScanResultsOptions) *apb.ScanResults {
	return &apb.ScanResults{
		StartTime:         timestamppb.New(options.startTime),
		EndTime:           timestamppb.New(time.Now()),
		ScannerVersion:    ScannerVersion,
		BenchmarkVersion:  options.benchmarkVersion,
		BenchmarkDocument: options.benchmarkDocument,
		Status: &apb.ScanStatus{
			Status:        options.status,
			FailureReason: options.failureReason,
		},
		CompliantBenchmarks:    options.compliantBenchmarks,
		NonCompliantBenchmarks: options.nonCompliantBenchmarks,
	}
}

// oldestBenchmarkVersion returns the version of the oldest benchmark in the list. This value
// can be compared to the benchmark  document version to figure out if a scan is up-to-date, i.e.
// if it was only using benchmarks defined in the latest document.
func oldestBenchmarkVersion(benchmarks []*apb.BenchmarkConfig) (string, error) {
	oldest := ""
	for _, benchmark := range benchmarks {
		// A single benchmark can be defined in several documents. In this case we're interested
		// in the version of the latest doc it's defined in.
		newestInBenchmark := "0.0.0"
		for _, version := range benchmark.GetComplianceNote().GetVersion() {
			isNewer, err := isVersionNewer(version.GetVersion(), newestInBenchmark)
			if err != nil {
				return "", err
			}
			if isNewer {
				newestInBenchmark = version.GetVersion()
			}
		}
		if len(oldest) == 0 {
			oldest = newestInBenchmark
		} else {
			isNewer, err := isVersionNewer(oldest, newestInBenchmark)
			if err != nil {
				return "", err
			}
			if isNewer {
				oldest = newestInBenchmark
			}
		}
	}
	return oldest, nil
}

// getBenchmarkDocument returns name of the benchmark document used in the scan.
// This value is expected to be the same in all benchmarks from the config.
func getBenchmarkDocument(benchmarks []*apb.BenchmarkConfig) string {
	if len(benchmarks) == 0 {
		return ""
	}
	if len(benchmarks[0].GetComplianceNote().Version) == 0 {
		return ""
	}
	return benchmarks[0].GetComplianceNote().Version[0].BenchmarkDocument
}

// isVersionNewer compares version strings that follow a num.num.num... pattern.
// It returns true if v1 is newer (a higher version) than v2, and false otherwise.
func isVersionNewer(v1, v2 string) (bool, error) {
	v1Nums := strings.Split(v1, ".")
	v2Nums := strings.Split(v2, ".")
	if len(v1Nums) == 0 {
		return false, fmt.Errorf("Invalid version string: %s", v1)
	}
	if len(v2Nums) == 0 {
		return false, fmt.Errorf("Invalid version string: %s", v2)
	}

	numLength := len(v1Nums)
	if numLength > len(v2Nums) {
		numLength = len(v2Nums)
	}
	for i := 0; i < numLength; i++ {
		v1Num, err := strconv.Atoi(v1Nums[i])
		if err != nil {
			return false, fmt.Errorf("Invalid version string: %s", v1)
		}
		v2Num, err := strconv.Atoi(v2Nums[i])
		if err != nil {
			return false, fmt.Errorf("Invalid version string: %s", v2)
		}
		if v1Num > v2Num {
			return true, nil
		}
		if v1Num < v2Num {
			return false, nil
		}
	}
	return len(v1Nums) > len(v2Nums), nil
}
