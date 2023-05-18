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

package cli

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

// Flags contains a field for all the cli flags that can be set.
type Flags struct {
	ConfigFile              string
	ResultFile              string
	ChrootPath              string
	Database                string
	Cassandra               string
	BenchmarkOptOutIDs      string
	ContentOptOutRegexes    string
	FilenameOptOutRegexes   string
	TraversalOptOutRegexes  string
	ShowCompliantBenchmarks bool
	MaxCisProfileLevel      int
	ScanTimeout             time.Duration
	BenchmarkCheckTimeout   time.Duration
}

// ValidateFlags validates the passed command line flags.
func ValidateFlags(flags *Flags) error {
	if len(flags.ConfigFile) == 0 {
		return errors.New("--config not set")
	}
	if len(flags.ResultFile) == 0 {
		return errors.New("--result not set")
	}

	if len(flags.BenchmarkOptOutIDs) > 0 {
		for _, id := range strings.Split(flags.BenchmarkOptOutIDs, ",") {
			if len(id) == 0 {
				return errors.New("invalid --benchmark-opt-out-ids: ID cannot be left empty")
			}
		}
	}

	if err := validateRegexArg(flags.ContentOptOutRegexes, "--content-opt-out-regexes"); err != nil {
		return err
	}
	if err := validateRegexArg(flags.FilenameOptOutRegexes, "--filename-opt-out-regexes"); err != nil {
		return err
	}
	if err := validateRegexArg(flags.TraversalOptOutRegexes, "--traversal-opt-out-regexes"); err != nil {
		return err
	}

	if flags.MaxCisProfileLevel < 1 {
		return errors.New("--max-cis-profile-level must be 1 or higher")
	}

	return nil
}

func validateRegexArg(arg string, name string) error {
	if len(arg) == 0 {
		return nil
	}
	for _, regex := range strings.Split(arg, ",") {
		if len(regex) == 0 {
			return fmt.Errorf("invalid %s: Regex cannot be left empty", name)
		}
	}
	return nil
}
