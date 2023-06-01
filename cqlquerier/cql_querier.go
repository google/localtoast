// Package cqlquerier provides an utility function for running CQL queries.
package cqlquerier

import (
	"context"
	"errors"

	"github.com/gocql/gocql"
)

// Query executes a CQL query and returns the number of rows in the result.
func Query(ctx context.Context, db *gocql.Session, query string) (int, error) {
	if db == nil {
		return 0, errors.New("no cassandra database specified. Please provide one using the --cassandra flag")
	}
	scanner := db.Query(query).Iter().Scanner()
	n := 0
	for scanner.Next() {
		n++
	}
	if err := scanner.Err(); err != nil {
		return 0, err
	}
	return n, nil
}
