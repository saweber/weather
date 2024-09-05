package main

import (
	"testing"
	"time"
)

func TestParseTimeParam(t *testing.T) {
	tests := []struct {
		name     string
		param    string
		wantTime *time.Time
		wantErr  bool
	}{
		{
			name:     "EmptyString",
			param:    "",
			wantTime: nil,
			wantErr:  false,
		},
		{
			name:     "ValidRFC3339",
			param:    "2020-12-31T23:59:59Z",
			wantTime: func() *time.Time { t, _ := time.Parse(time.RFC3339, "2020-12-31T23:59:59Z"); return &t }(),
			wantErr:  false,
		},
		{
			name:     "InvalidRFC3339",
			param:    "invalid RFC3339 time string",
			wantTime: nil,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotTime, err := parseTimeParam(tt.param)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseTimeParam() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantTime == nil && gotTime != nil {
				t.Errorf("parseTimeParam() = %v, want %v", gotTime, tt.wantTime)
				return
			}
			if tt.wantTime != nil && (gotTime == nil || !gotTime.Equal(*tt.wantTime)) {
				t.Errorf("parseTimeParam() = %v, want %v", gotTime, tt.wantTime)
			}
		})
	}
}
