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
	"path/filepath"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"google.golang.org/protobuf/encoding/prototext"
	"google.golang.org/protobuf/testing/protocmp"
	"bitbucket.org/creachadair/stringset"
	cpb "github.com/google/localtoast/scannerlib/proto/compliance_go_proto"
	spb "github.com/google/localtoast/scannerlib/proto/compliance_go_proto"
	apb "github.com/google/localtoast/scannerlib/proto/api_go_proto"
	sipb "github.com/google/localtoast/scannerlib/proto/scan_instructions_go_proto"
)

var configFileNames = stringset.New(
	"vm_image_scanning.textproto", "container_image_scanning.textproto", "instance_scanning.textproto")

// Validate behavior across all the configs.
func TestRequiredAttributes(t *testing.T) {
	fullNoteIDmap := make(map[string]*cpb.ComplianceNote)
	cpeVersionMap := make(map[string]*stringset.Set)
	for filePath, configBytes := range scanConfigs {
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
				if version.GetCpeUri() != "fallback" && version.GetCpeUri() != "cpe:/distribution_independent_linux" && len(strings.Split(version.GetCpeUri(), ":")) != 5 {
					t.Errorf("%s version.getCpeUri(): got %s, want a valid CPE", noteID, version.GetCpeUri())
				}
			}
			if note.GetCisBenchmark().GetProfileLevel() < 1 || note.GetCisBenchmark().GetProfileLevel() > 3 {
				t.Errorf("%s GetCisBenchmark().GetProfilelevel(): got %d, want 1 <= level <= 3", noteID, note.GetCisBenchmark().GetProfileLevel())
			}
			if note.GetCisBenchmark().GetSeverity() == spb.Severity_SEVERITY_UNSPECIFIED {
				t.Errorf("%s GetCisBenchmark().GetSeverity(): got %s, want any specified severity", noteID, note.GetCisBenchmark().GetSeverity())
			}
			scanInstructions := &sipb.BenchmarkScanInstruction{}
			if err := prototext.Unmarshal(note.GetScanInstructions(), scanInstructions); err != nil {
				t.Errorf("%s could not parse scan instructions: %v", noteID, err)
			}
			// Ensure that each note ID has an identical note, excluding version number, and scan instructions.
			// Image and instance scanning notes will check for the same thing, so they have mostly the same
			// information (note ID, description, etc...).
			// However, image notes have shorter scan instructions because we can't scan for everything that
			// instance scanning does.
			if val, ok := fullNoteIDmap[noteID]; ok {
				if diff := cmp.Diff(val, note, protocmp.Transform(), protocmp.IgnoreFields(&cpb.ComplianceNote{}, "version", "scan_instructions")); diff != "" {
					t.Errorf("%s checking note uniqueness: got %v, wanted no diff", noteID, diff)
				}
			} else {
				fullNoteIDmap[noteID] = note
			}
			for _, noteVer := range note.GetVersion() {
				if val, ok := cpeVersionMap[noteVer.GetCpeUri()]; ok {
					val.Add(noteVer.GetVersion())
				} else {
					s := stringset.New(noteVer.GetVersion())
					cpeVersionMap[noteVer.GetCpeUri()] = &s
				}
			}
		}
	}
	for cpe, versions := range cpeVersionMap {
		if versions.Len() > 1 {
			t.Errorf("CPE %s checking number of versions, got %v, wanted a single version", cpe, versions)
		}
	}
}

func TestFilesHaveSupportedName(t *testing.T) {
	for filePath := range scanConfigs {
		fileName := filepath.Base(filePath)
		if !configFileNames.Contains(fileName) {
			t.Errorf("checking file names: got %q, want image_scanning.textproto or instance_scanning.textproto", fileName)
		}
	}
}

func TestScanInstructionsSameForFileSets(t *testing.T) {
	for _, fileName := range configFileNames.Elements() {
		scanInstructionForNoteID := make(map[string]*sipb.BenchmarkScanInstruction)
		for filePath, configBytes := range scanConfigs {
			if filepath.Base(filePath) != fileName {
				continue
			}
			config := &apb.ScanConfig{}
			if err := prototext.Unmarshal(configBytes, config); err != nil {
				t.Errorf("error reading %s: %v", filePath, err)
			}
			for _, b := range config.GetBenchmarkConfigs() {
				noteID := b.GetId()
				note := b.GetComplianceNote()
				// Ensure that each note ID for images have identical scan instructions.
				scanInstructions := &sipb.BenchmarkScanInstruction{}
				if err := prototext.Unmarshal(note.GetScanInstructions(), scanInstructions); err != nil {
					t.Errorf("%s could not parse scan instructions: %v", noteID, err)
				}
				if val, ok := scanInstructionForNoteID[noteID]; ok {
					if diff := cmp.Diff(val, scanInstructions, protocmp.Transform()); diff != "" {
						t.Errorf("%s checking scan instructions unique: got %v, wanted no diff", noteID, diff)
					}
				} else {
					scanInstructionForNoteID[noteID] = scanInstructions
				}
			}
		}
	}
}

func TestFallbackNotesHaveExpectedIdFormat(t *testing.T) {
	for filePath, configBytes := range scanConfigs {
		if !strings.Contains(filePath, "fallback") {
			continue
		}
		config := &apb.ScanConfig{}
		if err := prototext.Unmarshal(configBytes, config); err != nil {
			t.Errorf("error reading %s: %v", filePath, err)
		}
		for _, b := range config.GetBenchmarkConfigs() {
			noteID := b.GetId()
			if !strings.HasSuffix(noteID, "-fallback") {
				t.Errorf("Fallback note ID %q should end with -fallback", noteID)
			}
		}
	}
}
