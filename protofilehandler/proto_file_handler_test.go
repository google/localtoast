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

package protofilehandler_test

import (
	"io/ioutil"
	"path/filepath"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"google.golang.org/protobuf/testing/protocmp"
	"github.com/google/localtoast/protofilehandler"
	apb "github.com/google/localtoast/scannerlib/proto/api_go_proto"
)

func TestReadProtoFromFile(t *testing.T) {
	testDirPath := t.TempDir()
	var expectedConfig = &apb.ScanConfig{
		BenchmarkConfigs: []*apb.BenchmarkConfig{
			&apb.BenchmarkConfig{Id: "test-benchmark"},
		},
	}
	testPaths := []string{"config.textproto", "config.binproto", "config.textproto.gz"}

	for _, path := range testPaths {
		fullPath := filepath.Join(testDirPath, path)
		if err := protofilehandler.WriteProtoToFile(fullPath, expectedConfig); err != nil {
			t.Fatalf("protofilehandler.WriteProtoToFile(%s, %v) returned an error: %v", fullPath, expectedConfig, err)
		}

		config := &apb.ScanConfig{}
		if err := protofilehandler.ReadProtoFromFile(fullPath, config); err != nil {
			t.Fatalf("protofilehandler.ReadProtoFromFile(%s) returned an error: %v", fullPath, err)
		}
		if diff := cmp.Diff(expectedConfig, config, protocmp.Transform()); diff != "" {
			t.Errorf("protofilehandler.ReadProtoFromFile(%s) returned unexpected diff (-want +got):\n%s",
				fullPath, diff)
		}
	}
}

func TestReadProtoFromFileInvalidData(t *testing.T) {
	testDirPath := t.TempDir()
	testCases := []struct {
		desc    string
		path    string
		content string
	}{
		{
			desc:    "textproto",
			path:    "config.textproto",
			content: "invalid textproto",
		},
		{
			desc:    "binproto",
			path:    "config.binproto",
			content: "invalid binproto",
		},
		{
			desc:    "gzipped file",
			path:    "config.textproto.gz",
			content: "invalid gzip",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			fullPath := filepath.Join(testDirPath, tc.path)
			if err := ioutil.WriteFile(fullPath, []byte(tc.content), 0644); err != nil {
				t.Fatalf("failed to create config file %s: %v", fullPath, err)
			}
			config := &apb.ScanConfig{}
			if err := protofilehandler.ReadProtoFromFile(fullPath, config); err == nil {
				t.Errorf("protofilehandler.ReadProtoFromFile(%s) didn't return an error", fullPath)
			}
		})
	}
}

func TestWriteResultToFile(t *testing.T) {
	testDirPath := t.TempDir()
	var result = &apb.ScanResults{ScannerVersion: "1.0.0"}
	testCases := []struct {
		desc           string
		path           string
		expectedPrefix string
	}{
		{
			desc:           "textproto",
			path:           "output.textproto",
			expectedPrefix: "scanner_version:",
		},
		{
			desc:           "binproto",
			path:           "output.binproto",
			expectedPrefix: "\x1a\x051.0.0",
		},
		{
			desc:           "gzipped file",
			path:           "output.textproto.gz",
			expectedPrefix: "\x1f\x8b",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			fullPath := filepath.Join(testDirPath, tc.path)
			err := protofilehandler.WriteProtoToFile(fullPath, result)
			if err != nil {
				t.Fatalf("protofilehandler.WriteProtoToFile(%s, %v) returned an error: %v", fullPath, result, err)
			}

			content, err := ioutil.ReadFile(fullPath)
			if err != nil {
				t.Fatalf("error while reading %s: %v", fullPath, err)
			}
			prefix := content[:len(tc.expectedPrefix)]
			if diff := cmp.Diff(tc.expectedPrefix, string(prefix)); diff != "" {
				t.Errorf("%s contains unexpected prefix, diff (-want +got):\n%s", fullPath, diff)
			}
		})
	}
}

func TestInvalidProtoFileName(t *testing.T) {
	testDirPath := t.TempDir()
	testPaths := []string{
		"config.invalid-extension",
		"config.invalid-extension.gz",
		"no-extension",
		"no-extension.gz",
	}
	for _, p := range testPaths {
		fullPath := filepath.Join(testDirPath, p)
		config := &apb.ScanConfig{}
		if err := protofilehandler.ReadProtoFromFile(fullPath, config); err == nil || !strings.HasPrefix(err.Error(), "invalid filename") {
			t.Errorf("protofilehandler.ReadProtoFromFile(%s) didn't return an invalid file error: %v", fullPath, err)
		}
		if err := protofilehandler.WriteProtoToFile(fullPath, &apb.ScanResults{}); err == nil ||
			!strings.HasPrefix(err.Error(), "invalid filename") {
			t.Errorf("protofilehandler.WriteProtoToFile(%s) didn't return an invalid file error: %v", fullPath, err)
		}
	}
}
