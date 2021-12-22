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

package configchecks

import (
	"context"
	"fmt"
	"io"

	"google.golang.org/protobuf/encoding/prototext"
	"google.golang.org/protobuf/proto"
	apb "github.com/google/localtoast/scannerlib/proto/api_go_proto"
	ipb "github.com/google/localtoast/scannerlib/proto/scan_instructions_go_proto"
)

// BenchmarkCheck is an interface representing a check to perform for one or more benchmarks
// (e.g. checking for the existence of a given file).
type BenchmarkCheck interface {
	// Exec executes the checks defined by the interface implementation.
	Exec() (ComplianceMap, error)
	// BenchmarkIDs returns the IDs of the benchmarks associated with this check.
	BenchmarkIDs() []string
	String() string
}

// scanAPIProvider is an interface that gives read access to the filesystem of
// the machine to scan and can execute SQL queries on a single database.
type scanAPIProvider interface {
	OpenFile(ctx context.Context, path string) (io.ReadCloser, error)
	FilesInDir(ctx context.Context, path string) ([]*apb.DirContent, error)
	FilePermissions(ctx context.Context, path string) (*apb.PosixPermissions, error)
	SQLQuery(ctx context.Context, query string) (int, error)
}

// ComplianceMap is returned by the checks to aggregate the results of benchmark configchecks.
// It maps a CheckAlternative ID to a compliance result associated with that alternative.
type ComplianceMap map[int]*apb.ComplianceResult

// benchmark represents a single benchmark whose compliance the scanner should check.
type benchmark struct {
	id   string // The benchmark ID as seen in the benchmark config file.
	alts []*checkAlternative
}

// checkAlternative describes a series of compliance checks to execute. A
// given benchmark is compliant if it satisfies one if its check alternatives.
type checkAlternative struct {
	id    int // Generated using a running counter, used to connect checks to their alternatives.
	proto *ipb.CheckAlternative
}

// parseCheckAlternatives deserializes the check alternatives from the benchmark config.
func parseCheckAlternatives(config *apb.BenchmarkConfig, prevAlternativeID int) ([]*checkAlternative, error) {
	serialized := config.GetComplianceNote().GetScanInstructions()
	instruction := &ipb.BenchmarkScanInstruction{}
	// The scan instructions in the Grafeas Note are serialized since they're
	// implementation-specific, so we have to deserialize them first. The
	// instructions are either in the textproto or binproto format.
	bo := proto.UnmarshalOptions{DiscardUnknown: true}
	if err := bo.Unmarshal(serialized, instruction); err != nil {
		to := &prototext.UnmarshalOptions{DiscardUnknown: true}
		if err := to.Unmarshal(serialized, instruction); err != nil {
			return nil, err
		}
	}
	if len(instruction.GetCheckAlternatives()) == 0 {
		return nil, fmt.Errorf("scan instruction %v doesn't define any checks", instruction)
	}
	result := make([]*checkAlternative, 0, len(instruction.GetCheckAlternatives()))
	for _, alt := range instruction.GetCheckAlternatives() {
		prevAlternativeID++
		result = append(result, &checkAlternative{id: prevAlternativeID, proto: alt})
	}
	return result, nil
}

// CreateChecksFromConfig parses the scan config and creates the benchmark checks defined by it.
func CreateChecksFromConfig(ctx context.Context, scanConfig *apb.ScanConfig, api scanAPIProvider) ([]BenchmarkCheck, error) {
	prevAlternativeID := 0
	benchmarks := make([]*benchmark, 0, len(scanConfig.GetBenchmarkConfigs()))
	for _, b := range scanConfig.GetBenchmarkConfigs() {
		alts, err := parseCheckAlternatives(b, prevAlternativeID)
		if err != nil {
			return nil, err
		}
		benchmarks = append(benchmarks, &benchmark{id: b.GetId(), alts: alts})
		prevAlternativeID = alts[len(alts)-1].id
	}

	fileCheckBatches, err := createFileCheckBatchesFromConfig(ctx, benchmarks, scanConfig.GetOptOutConfig(), api)
	if err != nil {
		return nil, err
	}
	sqlChecks, err := createSQLChecksFromConfig(benchmarks, api)
	if err != nil {
		return nil, err
	}

	checks := make([]BenchmarkCheck, 0, len(fileCheckBatches)+len(sqlChecks))
	for _, b := range fileCheckBatches {
		checks = append(checks, b)
	}
	for _, c := range sqlChecks {
		checks = append(checks, c)
	}
	return checks, nil
}

// ValidateScanInstructions validates the scan instructions in the given benchmark config and
// returns an error if they're invalid.
func ValidateScanInstructions(config *apb.BenchmarkConfig) error {
	alts, err := parseCheckAlternatives(config, 0)
	if err != nil {
		return err
	}
	for i, alt := range alts {
		if len(alt.proto.GetFileChecks()) == 0 && len(alt.proto.GetSqlChecks()) == 0 {
			return fmt.Errorf("alternative #%d in benchmark %s doesn't have any checks", i, config.GetId())
		}
	}
	return nil
}
