package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

func parseTimeParam(param string) (*time.Time, error) {
	if param == "" {
		return nil, nil
	}
	paramTime, err := time.Parse(time.DateOnly, param)
	if err != nil {
		return nil, err
	}
	return &paramTime, nil
}

func handleGetStormReports(repository *Repository) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		start := r.URL.Query().Get("windowStart")
		startTime, err := parseTimeParam(start)

		end := r.URL.Query().Get("windowEnd")
		endTime, err := parseTimeParam(end)

		opts := GetStormReportsOptions{
			ReportDateStart: startTime,
			ReportDateEnd:   endTime,
			Location:        r.URL.Query().Get("location"),
		}
		reports := repository.GetStormReports(opts)

		bytes, err := json.Marshal(reports)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
		w.Header().Set("Content-Type", "application/bytes")
		w.WriteHeader(http.StatusOK)
		_, err = w.Write(bytes)
		if err != nil {
			fmt.Printf("Error writing json response: %s\n", err)
		}
		return
	}
}
