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

// Package scannercommon provides the common functions used by the scanner binaries.
package scannercommon

import (
	"context"
	"flag"
	"log"
	"strings"

	"github.com/google/localtoast/cli"
	"github.com/google/localtoast/protofilehandler"
	apb "github.com/google/localtoast/scannerlib/proto/api_go_proto"
	"github.com/google/localtoast/scannerlib"
)

// ParseFlags parses the scanner binary's cli flags.
func ParseFlags() *cli.Flags {
	configFile := flag.String("config", "", "The path of the scan config file")
	resultFile := flag.String("result", "", "The path of the output scan result file")
	chrootPath := flag.String("chroot", "",
		"A path that will be prefixed to the paths of the files to be checked. "+
			"To be used when scanning a container/VM whose filesystem mounted to a disk")
	database := flag.String("database", "", "The ODBC data source name of the SQL database connection")
	benchmarkOptOutIDs := flag.String("benchmark-opt-out-ids", "",
		"A comma-separated list of benchmark IDs to exclude from scanning")
	contentOptOutRegexes := flag.String("content-opt-out-regexes", "",
		"A comma-separated list of file path regexes whose content shouldn't be displayed")
	filenameOptOutRegexes := flag.String("filename-opt-out-regexes", "",
		"A comma-separated list of file path regexes whose filename shouldn't be displayed")
	traversalOptOutRegexes := flag.String("traversal-opt-out-regexes", "",
		"A comma-separated list of file path regexes that should be omitted when traversing the filesystem recursively")
	showCompliantBenchmarks := flag.Bool("show-compliant-benchmarks", true,
		"Whether to show compliant benchmarks in the scan results.")
	maxCisProfileLevel := flag.Int("max-cis-profile-level", 3,
		"Don't scan for any CIS benchmarks with a higher level than this")

	flag.Parse()
	flags := &cli.Flags{
		ConfigFile:              *configFile,
		ResultFile:              *resultFile,
		ChrootPath:              *chrootPath,
		Database:                *database,
		BenchmarkOptOutIDs:      *benchmarkOptOutIDs,
		ContentOptOutRegexes:    *contentOptOutRegexes,
		FilenameOptOutRegexes:   *filenameOptOutRegexes,
		TraversalOptOutRegexes:  *traversalOptOutRegexes,
		ShowCompliantBenchmarks: *showCompliantBenchmarks,
		MaxCisProfileLevel:      *maxCisProfileLevel,
	}
	if err := cli.ValidateFlags(flags); err != nil {
		log.Fatalf("Error parsing CLI args: %v\n", err)
	}
	return flags
}

// RunScan executes the scan with the given CLI flags and API provider.
// Returns the exit code that the main binary should exit with.
func RunScan(flags *cli.Flags, provider scannerlib.ScanAPIProvider) int {
	log.Printf("Reading scan config from %s\n", flags.ConfigFile)
	config := &apb.ScanConfig{}
	if err := protofilehandler.ReadProtoFromFile(flags.ConfigFile, config); err != nil {
		log.Fatalf("Error reading config file: %v\n", err)
	}
	ApplyCLIFlagsToConfig(config, flags)

	log.Printf("Running scan of %d benchmarks\n", len(config.GetBenchmarkConfigs()))
	scanner := scannerlib.Scanner{}
	result, err := scanner.Scan(context.Background(), config, provider)
	if err != nil {
		log.Fatalf("Error while scanning: %v\n", err)
	}
	log.Printf("Scan status: %s\n", result.GetStatus().GetStatus().String())

	log.Printf("Found %d non-compliant benchmarks\n", len(result.GetNonCompliantBenchmarks()))
	if !flags.ShowCompliantBenchmarks {
		result.CompliantBenchmarks = []*apb.ComplianceResult{}
	}

	log.Printf("Writing scan results to %s\n", flags.ResultFile)
	if err := protofilehandler.WriteProtoToFile(flags.ResultFile, result); err != nil {
		log.Fatalf("Error writing scan results: %v\n", err)
	}

	if result.GetStatus().GetStatus() != apb.ScanStatus_SUCCEEDED {
		log.Printf("Not all checks completed successfully: %s\n", result.GetStatus().GetFailureReason())
		return 1
	}
	if len(result.GetNonCompliantBenchmarks()) > 0 {
		return 2
	}
	return 0
}

// ApplyCLIFlagsToConfig applies the given CLI flags to the scan config.
func ApplyCLIFlagsToConfig(config *apb.ScanConfig, flags *cli.Flags) {
	config.BenchmarkConfigs = removeOptedOutBenchmarks(config.GetBenchmarkConfigs(), strings.Split(flags.BenchmarkOptOutIDs, ","))
	config.BenchmarkConfigs = removeHighLevelBenchmarks(config.GetBenchmarkConfigs(), flags.MaxCisProfileLevel)
	addOptOutRegexes(config, flags.ContentOptOutRegexes, flags.FilenameOptOutRegexes, flags.TraversalOptOutRegexes)
}

func removeOptedOutBenchmarks(configs []*apb.BenchmarkConfig, optOutBenchmarks []string) []*apb.BenchmarkConfig {
	if len(optOutBenchmarks) == 0 {
		return configs
	}
	result := make([]*apb.BenchmarkConfig, 0, len(configs))
	optOutMap := make(map[string]bool)
	for _, id := range optOutBenchmarks {
		optOutMap[id] = true
	}
	for _, config := range configs {
		if optOutMap[config.GetId()] {
			continue
		}
		result = append(result, config)
	}
	return result
}

func removeHighLevelBenchmarks(configs []*apb.BenchmarkConfig, maxLevel int) []*apb.BenchmarkConfig {
	result := make([]*apb.BenchmarkConfig, 0, len(configs))
	for _, c := range configs {
		if c.ComplianceNote.GetCisBenchmark() != nil && int(c.ComplianceNote.GetCisBenchmark().ProfileLevel) > maxLevel {
			continue
		}
		result = append(result, c)
	}
	return result
}

func addOptOutRegexes(config *apb.ScanConfig, contentOptOutRegexes string, filenameOptOutRegexes string, traversalOptOutRegexes string) {
	content := config.GetOptOutConfig().GetContentOptoutRegexes()
	if len(contentOptOutRegexes) > 0 {
		content = append(content, strings.Split(contentOptOutRegexes, ",")...)
	}
	filename := config.GetOptOutConfig().GetFilenameOptoutRegexes()
	if len(filenameOptOutRegexes) > 0 {
		filename = append(filename, strings.Split(filenameOptOutRegexes, ",")...)
	}
	traversal := config.GetOptOutConfig().GetTraversalOptoutRegexes()
	if len(traversalOptOutRegexes) > 0 {
		traversal = append(traversal, strings.Split(traversalOptOutRegexes, ",")...)
	}

	if len(content) > 0 || len(filename) > 0 || len(traversal) > 0 {
		config.OptOutConfig = &apb.OptOutConfig{
			ContentOptoutRegexes:   content,
			FilenameOptoutRegexes:  filename,
			TraversalOptoutRegexes: traversal,
		}
	}
}
