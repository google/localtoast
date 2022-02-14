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

package genfullconfiglib_test

import (
	"fmt"
	"os"
	"path"
	"testing"

	"github.com/google/go-cmp/cmp"
	"google.golang.org/protobuf/testing/protocmp"
	cpb "github.com/google/localtoast/scannerlib/proto/compliance_go_proto"
	gpb "github.com/google/localtoast/scannerlib/proto/compliance_go_proto"
	"github.com/google/localtoast/configs/genfullconfig/genfullconfiglib"
	"github.com/google/localtoast/protofilehandler"
	apb "github.com/google/localtoast/scannerlib/proto/api_go_proto"
)

type testDirs struct {
	testDirPath    string
	reducedDirPath string
	defDirPath     string
	outDirPath     string
}

// Create directories for the input and output config files.
func createTestDirs(t *testing.T) testDirs {
	testDirPath := t.TempDir()
	dirs := testDirs{
		testDirPath:    testDirPath,
		reducedDirPath: path.Join(testDirPath, "config"),
		defDirPath:     path.Join(testDirPath, "defs"),
		outDirPath:     path.Join(testDirPath, "out"),
	}
	for _, d := range []string{dirs.reducedDirPath, dirs.defDirPath, dirs.outDirPath} {
		if err := os.Mkdir(d, 0744); err != nil {
			panic(fmt.Sprintf("error while creating directory %s: %v", d, err))
		}
	}
	return dirs
}

func getDefaultConfigPaths(dirs testDirs) (reducedPath string, defPath string, outPath string) {
	reducedPath = path.Join(dirs.reducedDirPath, "instance_scanning.textproto")
	defPath = path.Join(dirs.defDirPath, "def.textproto")
	outPath = path.Join(dirs.outDirPath, "config_instance_scanning.textproto")
	return
}

func TestNoInputPaths(t *testing.T) {
	createTestDirs(t)
	if err := genfullconfiglib.Generate([]string{}, []string{"out.textproto"}, false); err == nil {
		t.Errorf("genfullconfiglib.Generate({}, {'out.textproto'}, false) didn't return an error")
	}
}

func TestNoOutputPaths(t *testing.T) {
	createTestDirs(t)
	if err := genfullconfiglib.Generate([]string{"in.textproto"}, []string{}, false); err == nil {
		t.Errorf("genfullconfiglib.Generate({'in.textproto'}, {}, false) didn't return an error")
	}
}

func TestTooFewInputPaths(t *testing.T) {
	createTestDirs(t)
	out := []string{"full_config_1.textproto", "full_config_2.textproto"}
	in := []string{"reduced_config.textproto", "config_def.textproto"}
	if err := genfullconfiglib.Generate(in, out, false); err == nil {
		t.Errorf("genfullconfiglib.Generate(%v, %v, false) didn't return an error", in, out)
	}
}

func createTestScanConfig(id string, versions []*gpb.ComplianceVersion, instructions string, profileLevel int32) *apb.ScanConfig {
	return &apb.ScanConfig{BenchmarkConfigs: []*apb.BenchmarkConfig{&apb.BenchmarkConfig{
		Id: id,
		ComplianceNote: &gpb.ComplianceNote{
			Version: versions,
			ComplianceType: &cpb.ComplianceNote_CisBenchmark_{
				CisBenchmark: &cpb.ComplianceNote_CisBenchmark{ProfileLevel: profileLevel},
			},
			ScanInstructions: []byte(instructions),
		},
	}}}
}

func writeReducedConfigToFile(t *testing.T, path string, id string, version *gpb.ComplianceVersion) {
	reduced := &apb.PerOsBenchmarkConfig{Version: version, BenchmarkId: []string{id}}
	if err := protofilehandler.WriteProtoToFile(path, reduced); err != nil {
		t.Errorf("protofilehandler.WriteProtoToFile(%s, %v) returned an error: %v", path, reduced, err)
	}
}

func writeConfigDefToFile(t *testing.T, path string, id string, versions []*gpb.ComplianceVersion, instructions string) {
	def := createTestScanConfig(id, versions, instructions, 1)
	if err := protofilehandler.WriteProtoToFile(path, def); err != nil {
		t.Errorf("protofilehandler.WriteProtoToFile(%s, %v) returned an error: %v", path, def, err)
	}
}

func TestCreateSingleConfig(t *testing.T) {
	dirs := createTestDirs(t)
	reducedPath, defPath, outPath := getDefaultConfigPaths(dirs)
	version := &gpb.ComplianceVersion{CpeUri: "cpe", Version: "1.0.0"}
	writeReducedConfigToFile(t, reducedPath, "id", version)
	writeConfigDefToFile(t, defPath, "id", []*gpb.ComplianceVersion{version}, "generic:{check_alternatives:{}}")

	if err := genfullconfiglib.Generate([]string{reducedPath, defPath}, []string{outPath}, false); err != nil {
		t.Errorf("genfullconfiglib.Generate([%v, %v], [%v], false) returned an error: %v", reducedPath, defPath, outPath, err)
	}

	got := &apb.ScanConfig{}
	if err := protofilehandler.ReadProtoFromFile(outPath, got); err != nil {
		t.Errorf("protofilehandler.ReadProtoFromFile(%s, %v) returned an error: %v", outPath, got, err)
	}

	want := createTestScanConfig("id", []*gpb.ComplianceVersion{version}, "check_alternatives:{}", 1)
	if diff := cmp.Diff(want, got, protocmp.Transform()); diff != "" {
		t.Errorf("genfullconfiglib.Generate([%v, %v], [%v], false) returned unexpected diff (-want +got):\n%s",
			reducedPath, defPath, outPath, diff)
	}
}

func TestMissingDefinition(t *testing.T) {
	dirs := createTestDirs(t)
	reducedPath, defPath, outPath := getDefaultConfigPaths(dirs)
	version := &gpb.ComplianceVersion{CpeUri: "cpe", Version: "1.0.0"}
	writeReducedConfigToFile(t, reducedPath, "id", version)
	writeConfigDefToFile(t, defPath, "different-id", []*gpb.ComplianceVersion{version}, "generic:{check_alternatives:{}}")

	if err := genfullconfiglib.Generate([]string{reducedPath, defPath}, []string{outPath}, false); err == nil {
		t.Errorf("genfullconfiglib.Generate([%v, %v], [%v], false) didn't return an error", reducedPath, defPath, outPath)
	}
}

func TestInstructionsDifferPerScanType(t *testing.T) {
	dirs := createTestDirs(t)
	version := &gpb.ComplianceVersion{CpeUri: "cpe", Version: "1.0.0"}
	instancePath := path.Join(dirs.reducedDirPath, "instance_scanning.textproto")
	containerPath := path.Join(dirs.reducedDirPath, "container_image_scanning.textproto")
	vmPath := path.Join(dirs.reducedDirPath, "vm_image_scanning.textproto")
	defPath := path.Join(dirs.defDirPath, "def.textproto")
	writeReducedConfigToFile(t, instancePath, "id", version)
	writeReducedConfigToFile(t, containerPath, "id", version)
	writeReducedConfigToFile(t, vmPath, "id", version)
	writeConfigDefToFile(t, defPath, "id", []*gpb.ComplianceVersion{version},
		"scan_type_specific:{"+
			"  instance_scanning:{check_alternatives:{sql_checks:{}}}"+
			"  image_scanning:{check_alternatives:{file_checks:{}}}"+
			"}",
	)

	instanceOutPath := path.Join(dirs.reducedDirPath, "instance_out.textproto")
	containerOutPath := path.Join(dirs.reducedDirPath, "container_out.textproto")
	vmOutPath := path.Join(dirs.reducedDirPath, "vm_out.textproto")
	inPaths := []string{instancePath, containerPath, vmPath, defPath}
	outPaths := []string{instanceOutPath, containerOutPath, vmOutPath}
	if err := genfullconfiglib.Generate(inPaths, outPaths, false); err != nil {
		t.Errorf("genfullconfiglib.Generate(%v, %v, false) returned an error: %v", inPaths, outPaths, err)
	}

	testCases := []struct {
		description          string
		path                 string
		expectedInstructions string
	}{
		{
			description:          "instance scanning",
			path:                 instanceOutPath,
			expectedInstructions: "check_alternatives:{sql_checks:{}}",
		},
		{
			description:          "container image scanning",
			path:                 containerOutPath,
			expectedInstructions: "check_alternatives:{file_checks:{}}",
		},
		{
			description:          "VM image scanning",
			path:                 vmOutPath,
			expectedInstructions: "check_alternatives:{file_checks:{}}",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			got := &apb.ScanConfig{}
			if err := protofilehandler.ReadProtoFromFile(tc.path, got); err != nil {
				t.Errorf("protofilehandler.ReadProtoFromFile(%s, %v) returned an error: %v", tc.path, got, err)
			}

			want := createTestScanConfig("id", []*gpb.ComplianceVersion{version}, tc.expectedInstructions, 1)
			if diff := cmp.Diff(want, got, protocmp.Transform()); diff != "" {
				t.Errorf("genfullconfiglib.Generate(, false) returned unexpected diff for %s (-want +got):\n%s", tc.path, diff)
			}
		})
	}
}

func TestSameBenchmarkWithDifferentVersions(t *testing.T) {
	dirs := createTestDirs(t)
	version1 := &gpb.ComplianceVersion{CpeUri: "cpe", Version: "1.0.0"}
	version2 := &gpb.ComplianceVersion{CpeUri: "cpe", Version: "2.0.0"}
	version1ReducedPath := path.Join(dirs.reducedDirPath, "instance_scanning.textproto")
	version2ReducedPath := path.Join(dirs.reducedDirPath, "container_image_scanning.textproto")
	version1DefPath := path.Join(dirs.defDirPath, "def1.textproto")
	version2DefPath := path.Join(dirs.defDirPath, "def2.textproto")
	writeReducedConfigToFile(t, version1ReducedPath, "id", version1)
	writeReducedConfigToFile(t, version2ReducedPath, "id", version2)
	writeConfigDefToFile(t, version1DefPath, "id", []*gpb.ComplianceVersion{version1}, "generic:{check_alternatives:{file_checks:{}}}")
	writeConfigDefToFile(t, version2DefPath, "id", []*gpb.ComplianceVersion{version2}, "generic:{check_alternatives:{sql_checks:{}}}")

	version1OutPath := path.Join(dirs.reducedDirPath, "version1_out.textproto")
	version2OutPath := path.Join(dirs.reducedDirPath, "version2_out.textproto")
	inPaths := []string{version1ReducedPath, version2ReducedPath, version1DefPath, version2DefPath}
	outPaths := []string{version1OutPath, version2OutPath}
	if err := genfullconfiglib.Generate(inPaths, outPaths, false); err != nil {
		t.Errorf("genfullconfiglib.Generate(%v, %v, false) returned an error: %v", inPaths, outPaths, err)
	}

	testCases := []struct {
		description          string
		path                 string
		expectedVersion      *gpb.ComplianceVersion
		expectedInstructions string
	}{
		{
			description:          "version 1",
			path:                 version1OutPath,
			expectedVersion:      version1,
			expectedInstructions: "check_alternatives:{file_checks:{}}",
		},
		{
			description:          "version 2",
			path:                 version2OutPath,
			expectedVersion:      version2,
			expectedInstructions: "check_alternatives:{sql_checks:{}}",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			got := &apb.ScanConfig{}
			if err := protofilehandler.ReadProtoFromFile(tc.path, got); err != nil {
				t.Errorf("protofilehandler.ReadProtoFromFile(%s, %v) returned an error: %v", tc.path, got, err)
			}

			want := createTestScanConfig("id", []*gpb.ComplianceVersion{tc.expectedVersion}, tc.expectedInstructions, 1)
			if diff := cmp.Diff(want, got, protocmp.Transform()); diff != "" {
				t.Errorf("genfullconfiglib.Generate(, false) returned unexpected diff for %s (-want +got):\n%s", tc.path, diff)
			}
		})
	}
}

func TestOmitDescriptionFields(t *testing.T) {
	dirs := createTestDirs(t)
	reducedPath, defPath, outPath := getDefaultConfigPaths(dirs)
	version := &gpb.ComplianceVersion{CpeUri: "cpe", Version: "1.0.0"}
	writeReducedConfigToFile(t, reducedPath, "id", version)

	def := &apb.ScanConfig{BenchmarkConfigs: []*apb.BenchmarkConfig{&apb.BenchmarkConfig{
		Id: "id",
		ComplianceNote: &gpb.ComplianceNote{
			Title:            "Title",
			Description:      "Description",
			Version:          []*gpb.ComplianceVersion{version},
			Rationale:        "Rationale",
			Remediation:      "Remediation",
			ScanInstructions: []byte("generic:{check_alternatives:{}}"),
		},
	}}}
	if err := protofilehandler.WriteProtoToFile(defPath, def); err != nil {
		t.Errorf("protofilehandler.WriteProtoToFile(%s, %v) returned an error: %v", defPath, def, err)
	}

	if err := genfullconfiglib.Generate([]string{reducedPath, defPath}, []string{outPath}, true); err != nil {
		t.Errorf("genfullconfiglib.Generate([%v, %v], [%v], true) returned an error: %v", reducedPath, defPath, outPath, err)
	}

	got := &apb.ScanConfig{}
	if err := protofilehandler.ReadProtoFromFile(outPath, got); err != nil {
		t.Errorf("protofilehandler.ReadProtoFromFile(%s, %v) returned an error: %v", outPath, got, err)
	}

	want := &apb.ScanConfig{BenchmarkConfigs: []*apb.BenchmarkConfig{&apb.BenchmarkConfig{
		Id: "id",
		ComplianceNote: &gpb.ComplianceNote{
			Version:          []*gpb.ComplianceVersion{version},
			ScanInstructions: []byte("check_alternatives:{}"),
		},
	}}}
	if diff := cmp.Diff(want, got, protocmp.Transform()); diff != "" {
		t.Errorf("genfullconfiglib.Generate([%v, %v], [%v], true) returned unexpected diff (-want +got):\n%s",
			reducedPath, defPath, outPath, diff)
	}
}

func TestOverrideProfileLevel(t *testing.T) {
	dirs := createTestDirs(t)
	reducedPath, defPath, outPath := getDefaultConfigPaths(dirs)
	version := &gpb.ComplianceVersion{CpeUri: "cpe", Version: "1.0.0"}
	testCases := []struct {
		description   string
		id            string
		level         int32
		overrideID    string
		overrideLevel int32
		wantLevel     int32
	}{
		{
			description:   "override level if ID matches",
			id:            "id1",
			level:         1,
			overrideID:    "id1",
			overrideLevel: 2,
			wantLevel:     2,
		},
		{
			description:   "don't override level if ID doesn't match",
			id:            "id1",
			level:         1,
			overrideID:    "id2",
			overrideLevel: 2,
			wantLevel:     1,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			reduced := &apb.PerOsBenchmarkConfig{
				Version: version, BenchmarkId: []string{tc.id},
				ProfileLevelOverride: []*apb.ProfileLevelOverride{{Level: tc.overrideLevel, BenchmarkId: []string{tc.overrideID}}},
			}
			if err := protofilehandler.WriteProtoToFile(reducedPath, reduced); err != nil {
				t.Errorf("protofilehandler.WriteProtoToFile(%s, %v) returned an error: %v", reducedPath, reduced, err)
			}
			def := createTestScanConfig(tc.id, []*gpb.ComplianceVersion{version}, "generic:{check_alternatives:{}}", 1)
			if err := protofilehandler.WriteProtoToFile(defPath, def); err != nil {
				t.Errorf("protofilehandler.WriteProtoToFile(%s, %v) returned an error: %v", defPath, def, err)
			}

			if err := genfullconfiglib.Generate([]string{reducedPath, defPath}, []string{outPath}, false); err != nil {
				t.Errorf("genfullconfiglib.Generate([%v, %v], [%v], false) returned an error: %v", reducedPath, defPath, outPath, err)
			}

			got := &apb.ScanConfig{}
			if err := protofilehandler.ReadProtoFromFile(outPath, got); err != nil {
				t.Errorf("protofilehandler.ReadProtoFromFile(%s, %v) returned an error: %v", outPath, got, err)
			}

			want := createTestScanConfig(tc.id, []*gpb.ComplianceVersion{version}, "check_alternatives:{}", tc.wantLevel)
			if diff := cmp.Diff(want, got, protocmp.Transform()); diff != "" {
				t.Errorf("genfullconfiglib.Generate([%v, %v], [%v], false) returned unexpected diff (-want +got):\n%s",
					reducedPath, defPath, outPath, diff)
			}
		})
	}
}
