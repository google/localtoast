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

// The gen_reduced_configs command creates a reduced version of all the
// benchmark configs in a given directory that omit the description, rationale,
// and remediation fields.
// These fields are not used by the scanner so invoking it with the reduced configs
// produces the same behaviour but uses less memory.
package main

import (
	"flag"
	"log"

	"github.com/google/localtoast/configs/genreducedconfigs"
)

func main() {
	configDir := flag.String("dir", "", "The directory the config files are stored in")
	flag.Parse()
	if err := genreducedconfigs.Generate(*configDir); err != nil {
		log.Fatal(err)
	}
}
