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

package genreducedconfigs_test

import (
	"errors"
	"io/ioutil"
	"os"
	"path"
	"testing"

	"github.com/google/go-cmp/cmp"
	"google.golang.org/protobuf/encoding/prototext"
	"google.golang.org/protobuf/testing/protocmp"
	gpb "google.golang.org/genproto/googleapis/grafeas/v1"
	"github.com/google/localtoast/configs/genreducedconfigs"
	apb "github.com/google/localtoast/library/proto/api_go_proto"
)

func createTestFile(t *testing.T, testDir, filename string, config *apb.ScanConfig) error {
	content, err := prototext.Marshal(config)
	if err != nil {
		t.Fatalf("prototext.Marshal(%v) returned an error: %v", config, err)
	}

	if err := ioutil.WriteFile(path.Join(testDir, filename), content, 0644); err != nil {
		t.Fatalf("createTestFile(%s, %s) returned an error: %v", filename, content, err)
	}
	return nil
}

func fileExists(testDir, filename string) bool {
	_, err := os.Stat(path.Join(testDir, filename))
	return !errors.Is(err, os.ErrNotExist)
}

func TestDirectoryNotSet(t *testing.T) {
	if err := genreducedconfigs.Generate(""); err == nil {
		t.Errorf("genreducedconfigs.Generate('') didn't return an error")
	}
}

func TestDescriptionsRemoved(t *testing.T) {
	testDir := t.TempDir()
	config := &apb.ScanConfig{
		BenchmarkConfigs: []*apb.BenchmarkConfig{
			&apb.BenchmarkConfig{
				Id: "id1",
				ComplianceNote: &gpb.ComplianceNote{
					Title:       "Title1",
					Description: "Description1",
					Rationale:   "Rationale1",
					Remediation: "Remediation1",
				},
			},
			&apb.BenchmarkConfig{
				Id: "id2",
				ComplianceNote: &gpb.ComplianceNote{
					Title:       "Title2",
					Description: "Description2",
					Rationale:   "Rationale2",
					Remediation: "Remediation2",
				},
			},
		},
	}
	want := &apb.ScanConfig{
		BenchmarkConfigs: []*apb.BenchmarkConfig{
			&apb.BenchmarkConfig{
				Id:             "id1",
				ComplianceNote: &gpb.ComplianceNote{},
			},
			&apb.BenchmarkConfig{
				Id:             "id2",
				ComplianceNote: &gpb.ComplianceNote{},
			}},
	}

	createTestFile(t, testDir, "config.textproto", config)
	if err := genreducedconfigs.Generate(testDir); err != nil {
		t.Fatalf("genreducedconfigs.Generate(%s) returned an error %v", testDir, err)
	}
	reducedFile := "config_reduced.textproto"

	content, err := ioutil.ReadFile(path.Join(testDir, reducedFile))
	if err != nil {
		t.Fatalf("error while reading %s: %v", reducedFile, err)
	}
	got := &apb.ScanConfig{}
	if err := prototext.Unmarshal(content, got); err != nil {
		t.Fatalf("prototext.Unmarshal(%s) returned error: %v", content, err)
	}
	if diff := cmp.Diff(want, got, protocmp.Transform()); diff != "" {
		t.Errorf("reduced proto contains unexpected diff (-want +got):\n%s", diff)
	}
}

func TestAllTextprotosProcessed(t *testing.T) {
	testDir := t.TempDir()
	createTestFile(t, testDir, "config1.textproto", &apb.ScanConfig{})
	createTestFile(t, testDir, "config2.textproto", &apb.ScanConfig{})
	if err := genreducedconfigs.Generate(testDir); err != nil {
		t.Fatalf("genreducedconfigs.Generate(%s) returned an error %v", testDir, err)
	}
	if !fileExists(testDir, "config1_reduced.textproto") {
		t.Errorf("genreducedconfigs.Generate(%s) didn't create config1_reduced.textproto", testDir)
	}
	if !fileExists(testDir, "config2_reduced.textproto") {
		t.Errorf("genreducedconfigs.Generate(%s) didn't create config2_reduced.textproto", testDir)
	}
}

func TestNonTextprotosNotProcessed(t *testing.T) {
	testDir := t.TempDir()
	createTestFile(t, testDir, "config.txt", &apb.ScanConfig{})
	if err := genreducedconfigs.Generate(testDir); err != nil {
		t.Fatalf("genreducedconfigs.Generate(%s) returned an error %v", testDir, err)
	}
	if fileExists(testDir, "config_reduced.txt") {
		t.Errorf("genreducedconfigs.Generate(%s) didn't ignore config.txt", testDir)
	}
}

func TestAlreadyReducedTextprotosNotProcessedAgain(t *testing.T) {
	testDir := t.TempDir()
	createTestFile(t, testDir, "config_reduced.textproto", &apb.ScanConfig{})
	if err := genreducedconfigs.Generate(testDir); err != nil {
		t.Fatalf("genreducedconfigs.Generate(%s) returned an error %v", testDir, err)
	}
	if fileExists(testDir, "config_reduced_reduced.textproto") {
		t.Errorf("genreducedconfigs.Generate(%s) didn't ignore config_reduced.textproto", testDir)
	}
}
