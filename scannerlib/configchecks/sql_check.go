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
	"time"

	cpb "github.com/google/localtoast/scannerlib/proto/compliance_go_proto"
	apb "github.com/google/localtoast/scannerlib/proto/api_go_proto"
	ipb "github.com/google/localtoast/scannerlib/proto/scan_instructions_go_proto"
)

// SQLQuerier is an interface that supports SQL queries to a target
// SQL database.
type SQLQuerier interface {
	SQLQuery(ctx context.Context, query string) (int, error)
}

// MySQLCheck is an an implementation of scanner.Check that executes
// SQL queries on a MySQL/MariaDB database.
type MySQLCheck struct {
	ctx              context.Context
	benchmarkID      string
	alternativeID    int
	checkInstruction *ipb.SQLCheck
	querier          SQLQuerier
}

// Exec executes the SQL checks and returns the compliance status.
func (c *MySQLCheck) Exec() (ComplianceMap, error) {
	query := c.checkInstruction.GetQuery()
	rows, err := c.querier.SQLQuery(c.ctx, query)
	if err != nil {
		return nil, err
	}
	var reason string
	if rows > 0 && !c.checkInstruction.GetExpectResults() {
		reason = fmt.Sprintf("Expected no results for query %q, but got %d rows.", query, rows)
	} else if rows == 0 && c.checkInstruction.GetExpectResults() {
		reason = fmt.Sprintf("Expected results for query %q, but got none.", query)
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
	return ComplianceMap{c.alternativeID: r}, nil
}

// BenchmarkIDs returns the IDs of the benchmarks associated with this check.
func (c *MySQLCheck) BenchmarkIDs() []string {
	// We don't do batching for MySQL checks, so we will always have exactly one ID.
	return []string{c.benchmarkID}
}

func (c *MySQLCheck) String() string {
	return fmt.Sprintf("[MySQL check with id %q]", c.benchmarkID)
}

// createSQLChecksFromConfig parses the benchmark config and creates the executable
// SQL checks that it defines.
func createSQLChecksFromConfig(ctx context.Context, benchmarks []*benchmark, scanTimeout time.Time, sq SQLQuerier) ([]*MySQLCheck, error) {
	// TODO(b/235991635): Use scanTimeout.
	checks := []*MySQLCheck{}
	for _, b := range benchmarks {
		for _, alt := range b.alts {
			for _, sqlCheckInstruction := range alt.proto.GetSqlChecks() {
				if sqlCheckInstruction.GetTargetDatabase() != ipb.SQLCheck_DB_MYSQL {
					return nil, errors.New("only MySQL/MariaDB database checks are supported")
				}
				checks = append(checks, &MySQLCheck{
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
