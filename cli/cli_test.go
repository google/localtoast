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

package cli_test

import (
	"testing"

	"github.com/google/localtoast/cli"
)

func TestValidateFlags(t *testing.T) {
	for _, tc := range []struct {
		desc        string
		flags       *cli.Flags
		expectError bool
	}{
		{
			desc: "Valid config",
			flags: &cli.Flags{
				ConfigFile:             "config.textproto",
				ResultFile:             "result.textproto",
				BenchmarkOptOutIDs:     "id1,id2,id3",
				ContentOptOutRegexes:   "regex1,regex2,regex3",
				FilenameOptOutRegexes:  "regex1,regex2,regex3",
				TraversalOptOutRegexes: "regex1,regex2,regex3",
				MaxCisProfileLevel:     3,
				MySQLDatabase:          "127.0.0.1",
			},
			expectError: false,
		},
		{
			desc: "Config missing",
			flags: &cli.Flags{
				ConfigFile:         "",
				ResultFile:         "result.textproto",
				MaxCisProfileLevel: 3,
			},
			expectError: true,
		},
		{
			desc: "Result missing",
			flags: &cli.Flags{
				ConfigFile:         "config.textproto",
				ResultFile:         "",
				MaxCisProfileLevel: 3,
			},
			expectError: true,
		},
		{
			desc: "Empty benchmark opt-out ID",
			flags: &cli.Flags{
				ConfigFile:         "config.textproto",
				ResultFile:         "result.textproto",
				BenchmarkOptOutIDs: "id1,,id2",
				MaxCisProfileLevel: 3,
			},
			expectError: true,
		},
		{
			desc: "Empty content opt-out regex",
			flags: &cli.Flags{
				ConfigFile:           "config.textproto",
				ResultFile:           "result.textproto",
				ContentOptOutRegexes: "regex1,,regex2",
				MaxCisProfileLevel:   3,
			},
			expectError: true,
		},
		{
			desc: "Empty filename opt-out regex",
			flags: &cli.Flags{
				ConfigFile:            "config.textproto",
				ResultFile:            "result.textproto",
				FilenameOptOutRegexes: "regex1,,regex2",
				MaxCisProfileLevel:    3,
			},
			expectError: true,
		},
		{
			desc: "Empty traversal opt-out regex",
			flags: &cli.Flags{
				ConfigFile:             "config.textproto",
				ResultFile:             "result.textproto",
				TraversalOptOutRegexes: "regex1,regex2,",
				MaxCisProfileLevel:     3,
			},
			expectError: true,
		},
		{
			desc: "Invalid profile level",
			flags: &cli.Flags{
				ConfigFile:         "config.textproto",
				ResultFile:         "result.textproto",
				MaxCisProfileLevel: 0,
			},
			expectError: true,
		},
		{
			desc: "Multiple database set",
			flags: &cli.Flags{
				ConfigFile:        "config.textproto",
				ResultFile:        "result.textproto",
				MySQLDatabase:     "127.0.0.1",
				CassandraDatabase: "127.0.0.1",
			},
			expectError: true,
		},
		{
			desc: "Multiple database set",
			flags: &cli.Flags{
				ConfigFile:            "config.textproto",
				ResultFile:            "result.textproto",
				MySQLDatabase:         "127.0.0.1",
				ElasticSearchDatabase: "https://elastic:test@localhost:9200",
			},
			expectError: true,
		},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			err := cli.ValidateFlags(tc.flags)
			if err == nil && tc.expectError {
				t.Errorf("validation passed, expected it to fail")
			} else if err != nil && !tc.expectError {
				t.Errorf("validation failed with %v, expected it to pass", err)
			}
		})
	}
}
