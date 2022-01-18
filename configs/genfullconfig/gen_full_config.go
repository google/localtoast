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

// The gen_full_config command creates full per-OS scan config files by
// combining the config definition and the reduced config files.
package main

import (
	"flag"
	"log"
	"strings"

	"github.com/google/localtoast/configs/genfullconfig/genfullconfig"
)

func main() {
	inPaths := flag.String("in", "", "Comma-separated list of the reduced per-OS configs, followed by a list of the config definition paths")
	outPaths := flag.String("out", "", "Comma-separated list of the output paths for the produced full configs")
	flag.Parse()

	// Remove trailing commas.
	if strings.HasSuffix(*inPaths, ",") {
		*inPaths = (*inPaths)[:len(*inPaths)-1]
	}
	if strings.HasSuffix(*outPaths, ",") {
		*outPaths = (*outPaths)[:len(*outPaths)-1]
	}

	if err := genfullconfig.Generate(strings.Split(*inPaths, ","), strings.Split(*outPaths, ",")); err != nil {
		log.Fatal(err)
	}
}
