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

// Package genfullconfiglib creates full per-OS scan config files by combining
// the config definition and the reduced config files.
package genfullconfiglib

import (
	"errors"
	"fmt"
	"path"

	"google.golang.org/protobuf/encoding/prototext"
	"google.golang.org/protobuf/proto"
	cpb "github.com/google/localtoast/scannerlib/proto/compliance_go_proto"
	"github.com/google/localtoast/protofilehandler"
	apb "github.com/google/localtoast/scannerlib/proto/api_go_proto"
	ipb "github.com/google/localtoast/scannerlib/proto/scan_instructions_go_proto"
)

type configDefKey struct {
	id      string
	version string
}

type configDefMap map[configDefKey]*apb.BenchmarkConfig

type scanTypeEnum int

const (
	scanTypeUnknown scanTypeEnum = iota
	scanTypeInstance
	scanTypeVMImage
	scanTypeContainerImage
)

// Generate creates the full OS-specific config values from the config
// definitions and per-OS configs and writes them to the specified output path.
func Generate(inPaths []string, outPaths []string, omitDescriptions bool) error {
	if len(inPaths) == 0 {
		return errors.New("--in not set")
	}
	if len(outPaths) == 0 {
		return errors.New("--out not set")
	}
	reducedConfigCount := len(outPaths)
	defConfigCount := len(inPaths) - reducedConfigCount
	if defConfigCount <= 0 {
		return errors.New("no definition files specified")
	}

	// Split the input paths into the reduced configs and the config defs.
	// The input paths list the reduced config paths first and then the
	// config def paths.
	reducedConfigPaths := make([]string, 0, reducedConfigCount)
	configDefPaths := make([]string, 0, defConfigCount)
	for i, p := range inPaths {
		if i < reducedConfigCount {
			reducedConfigPaths = append(reducedConfigPaths, p)
		} else {
			configDefPaths = append(configDefPaths, p)
		}
	}

	configDefs, err := createConfigDefMap(configDefPaths)
	if err != nil {
		return fmt.Errorf("error fetching config definitions: %v", err)
	}
	for i, p := range reducedConfigPaths {
		outPath := outPaths[i]

		reduced := &apb.PerOsBenchmarkConfig{}
		if err := protofilehandler.ReadProtoFromFile(p, reduced); err != nil {
			return err
		}
		scanType, err := getScanTypeFromFileName(path.Base(p))
		if err != nil {
			return err
		}
		config, err := getFullConfig(reduced, configDefs, scanType)
		if omitDescriptions {
			removeDescriptionFields(config)
		}
		if err != nil {
			return fmt.Errorf("error fetching full config: %v", err)
		}
		if err = protofilehandler.WriteProtoToFile(outPath, config); err != nil {
			return fmt.Errorf("error writing %s: %v", outPath, err)
		}
	}
	return nil
}

func getScanTypeFromFileName(name string) (scanTypeEnum, error) {
	switch name {
	case "instance_scanning.textproto":
		return scanTypeInstance, nil
	case "vm_image_scanning.textproto":
		return scanTypeVMImage, nil
	case "container_image_scanning.textproto":
		return scanTypeContainerImage, nil
	default:
		return scanTypeUnknown, fmt.Errorf("unknown scan type %s", name)
	}
}

func createConfigDefMap(configDefPaths []string) (configDefMap, error) {
	defs := make(configDefMap)
	for _, p := range configDefPaths {
		benchmarkDefs := &apb.ScanConfig{}
		if err := protofilehandler.ReadProtoFromFile(p, benchmarkDefs); err != nil {
			return nil, err
		}
		for _, b := range benchmarkDefs.BenchmarkConfigs {
			for _, v := range b.ComplianceNote.Version {
				key, err := createConfigDefKey(b.Id, v)
				if err != nil {
					return nil, err
				}
				defs[*key] = b
			}
		}
	}
	return defs, nil
}

func createConfigDefKey(id string, version *cpb.ComplianceVersion) (*configDefKey, error) {
	versionAsBytes, err := proto.MarshalOptions{Deterministic: true}.Marshal(version)
	if err != nil {
		return nil, err
	}
	return &configDefKey{id: id, version: string(versionAsBytes)}, nil
}

func getFullConfig(reduced *apb.PerOsBenchmarkConfig, configDefs configDefMap, scanType scanTypeEnum) (*apb.ScanConfig, error) {
	fullConfigs := make([]*apb.BenchmarkConfig, 0, len(reduced.BenchmarkId))
	for _, id := range reduced.BenchmarkId {
		key, err := createConfigDefKey(id, reduced.Version)
		if err != nil {
			return nil, err
		}
		config, ok := configDefs[*key]
		if !ok {
			return nil, fmt.Errorf("no benchmark definition for %s %v", id, reduced.Version)
		}
		if config, err = selectInstructionForScanType(config, scanType); err != nil {
			return nil, err
		}
		config.ComplianceNote.Version = []*cpb.ComplianceVersion{reduced.Version}
		fullConfigs = append(fullConfigs, config)
	}
	return &apb.ScanConfig{BenchmarkConfigs: fullConfigs}, nil
}

func selectInstructionForScanType(config *apb.BenchmarkConfig, scanType scanTypeEnum) (*apb.BenchmarkConfig, error) {
	instructionDef := &ipb.BenchmarkScanInstructionDef{}
	if err := prototext.Unmarshal(config.ComplianceNote.ScanInstructions, instructionDef); err != nil {
		return nil, fmt.Errorf("error extracting %s: %v", config.Id, err)
	}

	var instructions *ipb.BenchmarkScanInstruction
	if instructionDef.GetGeneric() != nil {
		instructions = instructionDef.GetGeneric()
	} else if instructionDef.GetScanTypeSpecific() != nil {
		switch scanType {
		case scanTypeInstance:
			instructions = instructionDef.GetScanTypeSpecific().InstanceScanning
		case scanTypeVMImage:
			instructions = instructionDef.GetScanTypeSpecific().ImageScanning
		case scanTypeContainerImage:
			instructions = instructionDef.GetScanTypeSpecific().ImageScanning
		default:
			return nil, fmt.Errorf("unknown scan type %v", scanType)
		}
	} else {
		return nil, fmt.Errorf("no scan instructions provided in %v", instructionDef)
	}
	instructionBytes, err := prototext.Marshal(instructions)
	if err != nil {
		return nil, err
	}

	result := proto.Clone(config).(*apb.BenchmarkConfig)
	result.ComplianceNote.ScanInstructions = instructionBytes
	return result, nil
}

func removeDescriptionFields(config *apb.ScanConfig) {
	for _, b := range config.GetBenchmarkConfigs() {
		b.GetComplianceNote().Title = ""
		b.GetComplianceNote().Description = ""
		b.GetComplianceNote().Rationale = ""
		b.GetComplianceNote().Remediation = ""
	}
}
