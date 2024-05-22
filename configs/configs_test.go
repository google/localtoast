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

package configs_test

import (
	"embed"
	"io/fs"
	"log"
	"path"
	"path/filepath"
	"strings"
	"testing"

	"google.golang.org/protobuf/encoding/prototext"
	"google.golang.org/protobuf/proto"
	"bitbucket.org/creachadair/stringset"
	spb "github.com/google/localtoast/scannerlib/proto/severity_go_proto"
	apb "github.com/google/localtoast/scannerlib/proto/api_go_proto"
	ipb "github.com/google/localtoast/scannerlib/proto/scan_instructions_go_proto"
)

//go:embed defs/*
var defEfs embed.FS

//go:embed reduced/*/*
var configEfs embed.FS

var configFileNames = stringset.New(
	"vm_image_scanning.textproto", "container_image_scanning.textproto", "instance_scanning.textproto")

var scanConfigDefs, reducedScanConfigs = readConfigFiles()

func readConfigFiles() (map[string][]byte, map[string][]byte) {
	return readFilesInDir(defEfs), readFilesInDir(configEfs)
}

func readFilesInDir(efs embed.FS) map[string][]byte {
	result := make(map[string][]byte)
	err := fs.WalkDir(efs, ".", func(filePath string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if path.Ext(filePath) != ".textproto" {
			return nil
		}
		content, err := efs.ReadFile(filePath)
		if err != nil {
			return err
		}
		result[filePath] = content
		return nil
	})
	if err != nil {
		log.Fatalf("Error reading scan config: %v\n", err)
	}
	return result
}

// Validate behavior across all the configs.
func TestRequiredAttributes(t *testing.T) {
	for filePath, configBytes := range scanConfigDefs {
		config := &apb.ScanConfig{}
		if err := prototext.Unmarshal(configBytes, config); err != nil {
			t.Errorf("error reading %s: %v", filePath, err)
		}
		for _, b := range config.GetBenchmarkConfigs() {
			noteID := b.GetId()
			note := b.GetComplianceNote()
			if note.GetTitle() == "" {
				t.Errorf(`%s GetTitle(): got "", want non-empty string`, noteID)
			}
			if note.GetDescription() == "" {
				t.Errorf(`%s GetDescription(): got "", want non-empty string`, noteID)
			}
			if note.GetRationale() == "" {
				t.Errorf(`%s GetRationale(): got "", want non-empty string`, noteID)
			}
			if note.GetRemediation() == "" {
				t.Errorf(`%s GetRemediation(): got "", want non-empty string`, noteID)
			}
			if len(note.GetVersion()) == 0 {
				t.Errorf(`%s len(GetVersion()): got 0, want at least 1`, noteID)
			}
			for _, version := range note.GetVersion() {
				if version.GetVersion() == "" {
					t.Errorf(`%s version.GetVersion(): got "", want non-empty string`, noteID)
				}
				if version.GetCpeUri() != "fallback" && version.GetCpeUri() != "cpe:/distribution_independent_linux" && len(strings.Split(version.GetCpeUri(), ":")) > 5 {
					t.Errorf("%s version.getCpeUri(): got %s, want a valid CPE", noteID, version.GetCpeUri())
				}
			}
			if note.GetCisBenchmark().GetProfileLevel() < 1 || note.GetCisBenchmark().GetProfileLevel() > 3 {
				t.Errorf("%s GetCisBenchmark().GetProfilelevel(): got %d, want 1 <= level <= 3", noteID, note.GetCisBenchmark().GetProfileLevel())
			}
			if note.GetCisBenchmark().GetSeverity() == spb.Severity_SEVERITY_UNSPECIFIED {
				t.Errorf("%s GetCisBenchmark().GetSeverity(): got %s, want any specified severity", noteID, note.GetCisBenchmark().GetSeverity())
			}
			scanInstructions := &ipb.BenchmarkScanInstructionDef{}
			if err := prototext.Unmarshal(note.GetScanInstructions(), scanInstructions); err != nil {
				t.Errorf("%s could not parse scan instructions: %v", noteID, err)
			}
		}
	}
}

func TestFilesHaveSupportedName(t *testing.T) {
	for filePath := range reducedScanConfigs {
		fileName := filepath.Base(filePath)
		if !configFileNames.Contains(fileName) {
			t.Errorf("checking file names: got %q, want image_scanning.textproto or instance_scanning.textproto", fileName)
		}
	}
}

func TestScanInstructionsHaveDisplayCommandAndNonComplianceReason(t *testing.T) {
	for filePath, configBytes := range scanConfigDefs {
		config := &apb.ScanConfig{}
		if err := prototext.Unmarshal(configBytes, config); err != nil {
			t.Fatalf("error reading %s: %v", filePath, err)
		}
		for _, b := range config.GetBenchmarkConfigs() {
			noteID := b.GetId()
			instructionDef := &ipb.BenchmarkScanInstructionDef{}
			if err := prototext.Unmarshal(b.GetComplianceNote().GetScanInstructions(), instructionDef); err != nil {
				t.Errorf("%s could not parse scan instructions: %v", noteID, err)
			}
			var scanInstructions []*ipb.BenchmarkScanInstruction
			if instructionDef.GetGeneric() != nil {
				scanInstructions = []*ipb.BenchmarkScanInstruction{instructionDef.GetGeneric()}
			} else if instructionDef.GetScanTypeSpecific() != nil {
				scanInstructions = []*ipb.BenchmarkScanInstruction{
					instructionDef.GetScanTypeSpecific().InstanceScanning,
					instructionDef.GetScanTypeSpecific().ImageScanning,
				}
			} else {
				t.Fatalf("benchmark %s has invalid instruction def %v", noteID, instructionDef)
			}
			for _, i := range scanInstructions {
				for _, a := range i.GetCheckAlternatives() {
					for _, f := range a.GetFileChecks() {
						if len(f.GetFileDisplayCommand()) > 0 && len(f.GetNonComplianceMsg()) == 0 {
							t.Errorf("check for benchmark %s has a file display command set but no non-compliance message: %v", noteID, f)
						}
					}
				}
			}
		}
	}
}

func TestFallbackBenchmarkDefsHaveExpectedIdFormat(t *testing.T) {
	for filePath, configBytes := range scanConfigDefs {
		if filePath != "fallback.textproto" {
			continue
		}
		config := &apb.ScanConfig{}
		if err := prototext.Unmarshal(configBytes, config); err != nil {
			t.Errorf("error reading %s: %v", filePath, err)
		}
		for _, b := range config.GetBenchmarkConfigs() {
			noteID := b.GetId()
			if !strings.HasSuffix(noteID, "-fallback") {
				t.Errorf("Fallback benchmark ID %q should end with -fallback", noteID)
			}
		}
	}
}

func TestFallbackPerOsBenchmarksHaveExpectedIdFormat(t *testing.T) {
	for filePath, configBytes := range reducedScanConfigs {
		if !strings.Contains(filePath, "/fallback/") {
			continue
		}
		config := &apb.PerOsBenchmarkConfig{}
		if err := prototext.Unmarshal(configBytes, config); err != nil {
			t.Errorf("error reading %s: %v", filePath, err)
		}
		for _, id := range config.BenchmarkId {
			if !strings.HasSuffix(id, "-fallback") {
				t.Errorf("Fallback benchmark ID %q should end with -fallback", id)
			}
		}
	}
}

func TestProfileLevelOverridesUseExistingIDs(t *testing.T) {
	for filePath, configBytes := range reducedScanConfigs {
		config := &apb.PerOsBenchmarkConfig{}
		if err := prototext.Unmarshal(configBytes, config); err != nil {
			t.Errorf("error reading %s: %v", filePath, err)
		}
		ids := make(map[string]bool)
		for _, id := range config.BenchmarkId {
			ids[id] = true
		}
		for _, o := range config.ProfileLevelOverride {
			for _, id := range o.BenchmarkId {
				if _, ok := ids[id]; !ok {
					t.Errorf("%s: overridden benchmark ID %q isn't used by the config", filePath, id)
				}
			}
		}
	}
}

func TestConfigDefContainsUniqueNoteDefinitions(t *testing.T) {
	notePresent := make(map[string]bool)
	for filePath, configBytes := range scanConfigDefs {
		config := &apb.ScanConfig{}
		if err := prototext.Unmarshal(configBytes, config); err != nil {
			t.Errorf("error reading %s: %v", filePath, err)
		}
		for _, benchmarkConfig := range config.GetBenchmarkConfigs() {
			noteID := benchmarkConfig.GetId()
			if _, present := notePresent[noteID]; present {
				t.Errorf("%s: Benchmark ID %q specified multiple times", filePath, noteID)
			} else {
				notePresent[noteID] = true
			}
		}
	}
}

func TestOnlyOneDelimiterForGivenFile(t *testing.T) {
	fileToDelimiter := make(map[string]string)
	for filePath, configBytes := range scanConfigDefs {
		config := &apb.ScanConfig{}
		if err := prototext.Unmarshal(configBytes, config); err != nil {
			t.Fatalf("error reading %s: %v", filePath, err)
		}
		for _, b := range config.GetBenchmarkConfigs() {
			noteID := b.GetId()
			instructionDef := &ipb.BenchmarkScanInstructionDef{}
			if err := prototext.Unmarshal(b.GetComplianceNote().GetScanInstructions(), instructionDef); err != nil {
				t.Errorf("%s could not parse scan instructions: %v", noteID, err)
			}
			var scanInstructions []*ipb.BenchmarkScanInstruction
			if instructionDef.GetGeneric() != nil {
				scanInstructions = []*ipb.BenchmarkScanInstruction{instructionDef.GetGeneric()}
			} else if instructionDef.GetScanTypeSpecific() != nil {
				scanInstructions = []*ipb.BenchmarkScanInstruction{
					instructionDef.GetScanTypeSpecific().InstanceScanning,
					instructionDef.GetScanTypeSpecific().ImageScanning,
				}
			} else {
				t.Fatalf("benchmark %s has invalid instruction def %v", noteID, instructionDef)
			}
			for _, i := range scanInstructions {
				for _, a := range i.GetCheckAlternatives() {
					for _, f := range a.GetFileChecks() {
						if f.GetContentEntry() != nil {
							delimiter := string(f.GetContentEntry().GetDelimiter())
							if delimiter == "" {
								delimiter = "\n" // Default value.
							}
							for _, fs := range f.GetFilesToCheck() {
								fileToCheck, err := proto.MarshalOptions{Deterministic: true}.Marshal(fs)
								if err != nil {
									t.Fatalf("Error marshaling %v: %v", f.GetFilesToCheck(), err)
								}
								if prev, ok := fileToDelimiter[string(fileToCheck)]; ok {
									if prev != delimiter {
										t.Fatalf("file %v has content entry checks with different delimiters: %q vs. %q",
											fs, prev, delimiter)
									}
								} else {
									fileToDelimiter[string(fileToCheck)] = delimiter
								}
							}
						}
					}
				}
			}
		}
	}
}
