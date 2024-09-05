package main

import (
	"encoding/csv"
	"fmt"
	"net/http"
)

func retrieveAndParseCSV(url string) ([][]string, error) {
	// Create an HTTP client and make a request
	response, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve CSV: %v", err)
	}
	defer response.Body.Close()

	// Check for a successful response
	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to download file: %v", response.Status)
	}

	reader := csv.NewReader(response.Body)

	// Reading CSV file into a slice of slice of strings
	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("failed to parse CSV: %v", err)
	}
	return records, nil
}
