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
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"google.golang.org/protobuf/testing/protocmp"
	apb "github.com/google/localtoast/library/proto/api_go_proto"
	"github.com/google/localtoast/protofilehandler"
)

var testDirPath = os.Getenv("TEST_TMPDIR")

func TestReadProtoFromFile(t *testing.T) {
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
	testCases := []struct {
		path    string
		content string
	}{
		{
			path:    "config.textproto",
			content: "invalid textproto",
		},
		{
			path:    "config.binproto",
			content: "invalid binproto",
		},
		{
			path:    "config.textproto.gz",
			content: "invalid gzip",
		},
	}

	for _, tc := range testCases {
		fullPath := filepath.Join(testDirPath, tc.path)
		if err := ioutil.WriteFile(fullPath, []byte(tc.content), 0644); err != nil {
			t.Fatalf("failed to create config file %s: %v", fullPath, err)
		}
		config := &apb.ScanConfig{}
		if err := protofilehandler.ReadProtoFromFile(fullPath, config); err == nil {
			t.Errorf("protofilehandler.ReadProtoFromFile(%s) didn't return an error", fullPath)
		}
	}
}

func TestWriteResultToFile(t *testing.T) {
	var result = &apb.ScanResults{ScannerVersion: "1.0.0"}
	testCases := []struct {
		path           string
		expectedPrefix string
	}{
		{
			path:           "output.textproto",
			expectedPrefix: "scanner_version:",
		},
		{
			path:           "output.binproto",
			expectedPrefix: "\x1a\x051.0.0",
		},
		{
			path:           "output.textproto.gz",
			expectedPrefix: "\x1f\x8b",
		},
	}

	for _, tc := range testCases {
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
	}
}

func TestInvalidProtoFileName(t *testing.T) {
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
