// Package cqlquerier provides an utility function for running CQL queries.
package cqlquerier

import (
	"context"
	"errors"

	"github.com/gocql/gocql"
)

// Query executes a CQL query and returns the first row in the result.
func Query(ctx context.Context, db *gocql.Session, query string) (string, error) {
	if db == nil {
		return "", errors.New("no cassandra database specified. Please provide one using the --cassandra flag")
	}
	scanner := db.Query(query).Iter().Scanner()
	result := ""
	if scanner.Next() {
		if err := scanner.Scan(&result); err != nil {
			return "", err
		}
		if err := scanner.Err(); err != nil {
			return "", err
		}
	}
	return result, nil
}
