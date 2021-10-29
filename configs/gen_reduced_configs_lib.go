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

// Package genreducedconfigs provides a function for creating reduced versions
// of all the benchmark configs in a given directory that omit the description,
// rationale, and remediation fields.
package genreducedconfigs

import (
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	apb "github.com/google/localtoast/library/proto/api_go_proto"
	"github.com/google/localtoast/protofilehandler"
)

// Generate generates the reduced config files in the given directory.
func Generate(configDir string) error {
	if configDir == "" {
		return errors.New("--dir not set")
	}

	return filepath.Walk(configDir, func(filePath string, info os.FileInfo, err error) error {
		if info.IsDir() || path.Ext(filePath) != ".textproto" || strings.HasSuffix(filePath, "_reduced.textproto") {
			return nil
		}

		config := &apb.ScanConfig{}
		if err := protofilehandler.ReadProtoFromFile(filePath, config); err != nil {
			return fmt.Errorf("error reading %s: %v", filePath, err)
		}
		for _, b := range config.GetBenchmarkConfigs() {
			b.GetComplianceNote().Title = ""
			b.GetComplianceNote().Description = ""
			b.GetComplianceNote().Rationale = ""
			b.GetComplianceNote().Remediation = ""
		}

		baseName := strings.TrimSuffix(info.Name(), ".textproto")
		outPath := path.Join(path.Dir(filePath), baseName+"_reduced.textproto")
		if err := protofilehandler.WriteProtoToFile(outPath, config); err != nil {
			return fmt.Errorf("error writing %s: %v", outPath, err)
		}
		return nil
	})
}
