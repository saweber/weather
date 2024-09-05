package main

import (
	"testing"
	"time"
)

func TestConvertStrToInt(t *testing.T) {
	var tests = []struct {
		input string
		want  int
	}{
		{"", 0},
		{"100", 100},
		{"-50", -50},
		{"UNK", 0},
		{"abc", 0},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := convertStrToInt(tt.input)
			if got != tt.want {
				t.Errorf("convertStrToInt(%q) = %v; want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestGetStormReportTime(t *testing.T) {
	var tests = []struct {
		date    string
		timeStr string
		want    time.Time
	}{
		{"2019-06-02", "0100", time.Date(2019, 6, 1, 1, 0, 0, 0, time.UTC)},
		{"2019-06-02", "1200", time.Date(2019, 6, 2, 12, 0, 0, 0, time.UTC)},
		{"2019-06-02", "1300", time.Date(2019, 6, 2, 13, 0, 0, 0, time.UTC)},
		{"2019-06-02", "INVALID", time.Time{}},
	}

	for _, tt := range tests {
		t.Run(tt.date+" "+tt.timeStr, func(t *testing.T) {
			consumer := &Consumer{}
			got := consumer.getStormReportTime(tt.date, tt.timeStr)
			if !got.Equal(tt.want) {
				t.Errorf("getStormReportTime() = %v; want %v", got, tt.want)
			}
		})
	}
}
