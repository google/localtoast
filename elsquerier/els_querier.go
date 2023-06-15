// Package elsquerier provides an utility function for running ElasticSearch queries.
package elsquerier

import (
	"context"
	"io"
	"net/http"
	"strings"

	els "github.com/elastic/go-elasticsearch/v8"
)

// Query executes a ELS request and returns the JSON response as string
func Query(ctx context.Context, db *els.Client, endpoint string) (string, error) {
	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return "", err
	}

	res, err := db.Perform(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	body := new(strings.Builder)
	_, err = io.Copy(body, res.Body)
	if err != nil {
		return "", err
	}

	return body.String(), nil
}
