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

package configchecks

import (
	"context"
	"errors"
	"fmt"
	"regexp"

	cpb "github.com/google/localtoast/scannerlib/proto/compliance_go_proto"
	"github.com/google/localtoast/scanapi"
	apb "github.com/google/localtoast/scannerlib/proto/api_go_proto"
	ipb "github.com/google/localtoast/scannerlib/proto/scan_instructions_go_proto"
)

// SQLCheck is an implementation of configchecks.BenchmarkCheck
// It runs queries on the database specified by the check.
type SQLCheck struct {
	ctx              context.Context
	benchmarkID      string
	alternativeID    int
	checkInstruction *ipb.SQLCheck
	querier          scanapi.SQLQuerier
}

// Exec executes the SQL checks and returns the compliance status.
func (c *SQLCheck) Exec(prvRes string) (ComplianceMap, string, error) {
	query := c.checkInstruction.GetQuery()
	var resVal string = ""

	var reason string
	if c.checkInstruction.TargetDatabase == ipb.SQLCheck_DB_MYSQL || c.checkInstruction.TargetDatabase == ipb.SQLCheck_DB_CASSANDRA {
		// Check number of returned rows for MySQL and Cassandra
		resVal, err := c.querier.SQLQuery(c.ctx, query)
		if err != nil {
			return nil, "", err
		}
		if len(resVal) > 0 && !c.checkInstruction.GetExpectResults() {
			reason = fmt.Sprintf("Expected no results for query %q, but got some.", query)
		} else if len(resVal) == 0 && c.checkInstruction.GetExpectResults() {
			reason = fmt.Sprintf("Expected results for query %q, but got none.", query)
		}

	} else if c.checkInstruction.TargetDatabase == ipb.SQLCheck_DB_ELASTICSEARCH {
		// Perform regex match on result string for ElasticSearch
		filterRegex, err := regexp.Compile("^" + c.checkInstruction.FilterRegex + "$")
		if err != nil {
			return nil, "", err
		}
		// Execute ElasticSearch query
		resVal, err = c.querier.SQLQuery(c.ctx, query)
		if err != nil {
			return nil, "", err
		}
		// Check if regex obtains results and compare with expected result
		if !filterRegex.MatchString(resVal) && c.checkInstruction.GetExpectResults() {
			reason = fmt.Sprintf("ElasticSearch response %q does not match the Filter Regex %q and it should.", resVal, c.checkInstruction.FilterRegex)
		} else if filterRegex.MatchString(resVal) && !c.checkInstruction.GetExpectResults() {
			reason = fmt.Sprintf("ElasticSearch response %q matches the Filter Regex %q and it should not.", resVal, c.checkInstruction.FilterRegex)
		}
	} else {
		// Return error for unsupported database
		return nil, "", errors.New("unsupported database for SQLCheck")
	}

	if reason != "" && c.checkInstruction.GetNonComplianceMsg() != "" {
		reason = c.checkInstruction.GetNonComplianceMsg()
	}
	r := &apb.ComplianceResult{
		Id: c.benchmarkID,
		ComplianceOccurrence: &cpb.ComplianceOccurrence{
			NonComplianceReason: reason,
		},
	}
	return ComplianceMap{c.alternativeID: r}, resVal, nil
}

// BenchmarkIDs returns the IDs of the benchmarks associated with this check.
func (c *SQLCheck) BenchmarkIDs() []string {
	// We don't do batching for SQL checks, so we will always have exactly one ID.
	return []string{c.benchmarkID}
}

func (c *SQLCheck) String() string {
	return fmt.Sprintf("[SQL check with id %q]", c.benchmarkID)
}

// createSQLChecksFromConfig parses the benchmark config and creates the executable
// SQL checks that it defines.
func createSQLChecksFromConfig(ctx context.Context, benchmarks []*benchmark, timeout *timeoutOptions, sq scanapi.SQLQuerier) ([]*SQLCheck, error) {
	// TODO(b/235991635): Use timeout.
	checks := []*SQLCheck{}
	for _, b := range benchmarks {
		for _, alt := range b.alts {
			for _, sqlCheckInstruction := range alt.proto.GetSqlChecks() {
				dbtype, err := sq.SupportedDatabase()
				if err != nil {
					return nil, err
				}
				if dbtype != sqlCheckInstruction.GetTargetDatabase() {
					return nil, fmt.Errorf("sql check %v does not match the connected database type %v", sqlCheckInstruction.GetTargetDatabase(), dbtype)
				}
				if sqlCheckInstruction.GetTargetDatabase() == ipb.SQLCheck_DB_ELASTICSEARCH && sqlCheckInstruction.GetFilterRegex() == "" {
					return nil, errors.New("no regex provided for ElasticSearch database SQLCheck")
				}
				checks = append(checks, &SQLCheck{
					ctx:              ctx,
					benchmarkID:      b.id,
					alternativeID:    alt.id,
					checkInstruction: sqlCheckInstruction,
					querier:          sq,
				})
			}
		}
	}
	return checks, nil
}
