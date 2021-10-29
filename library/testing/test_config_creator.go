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

// Package testconfigcreator provides util functions for creating benchmark configs for testing.
package testconfigcreator

import (
	"testing"

	"google.golang.org/protobuf/encoding/prototext"
	gpb "github.com/google/localtoast/library/proto/compliance_go_proto"
	apb "github.com/google/localtoast/library/proto/api_go_proto"
	ipb "github.com/google/localtoast/library/proto/scan_instructions_go_proto"
)

// SingleFileWithPath creates a FileSet that defines a single file with the given path.
func SingleFileWithPath(path string) *ipb.FileSet {
	return &ipb.FileSet{
		FilePath: &ipb.FileSet_SingleFile_{SingleFile: &ipb.FileSet_SingleFile{Path: path}},
	}
}

// NewBenchmarkConfig creates a benchmark config with the given ID and scan instructions.
func NewBenchmarkConfig(t *testing.T, id string, scanInstruction *ipb.BenchmarkScanInstruction) *apb.BenchmarkConfig {
	t.Helper()
	serializedInstruction, err := prototext.Marshal(scanInstruction)
	if err != nil {
		t.Fatalf("error while serializing scan instructions %v: %v", scanInstruction, err)
	}
	return &apb.BenchmarkConfig{
		Id:             id,
		ComplianceNote: &gpb.ComplianceNote{ScanInstructions: serializedInstruction},
	}
}

// NewFileScanInstruction creates a scan instruction with a single alternative from the
// given file checks.
func NewFileScanInstruction(fileChecks []*ipb.FileCheck) *ipb.BenchmarkScanInstruction {
	return &ipb.BenchmarkScanInstruction{
		CheckAlternatives: []*ipb.CheckAlternative{{FileChecks: fileChecks}},
	}
}

// NewSQLScanInstruction creates a scan instruction with a single alternative from the
// given SQL checks.
func NewSQLScanInstruction(sqlChecks []*ipb.SQLCheck) *ipb.BenchmarkScanInstruction {
	return &ipb.BenchmarkScanInstruction{
		CheckAlternatives: []*ipb.CheckAlternative{{SqlChecks: sqlChecks}},
	}
}
