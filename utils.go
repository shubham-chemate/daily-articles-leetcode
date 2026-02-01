package main

import (
	"fmt"
	"os"
	"strings"
	"time"
)

const lastTimestampFile = "last_processed_timestamp.txt"

// readLastProcessedTimestamp reads the last processed timestamp from file
func readLastProcessedTimestamp() (time.Time, error) {
	data, err := os.ReadFile(lastTimestampFile)
	if err != nil {
		if os.IsNotExist(err) {
			// File doesn't exist, return zero time
			return time.Time{}, nil
		}
		return time.Time{}, fmt.Errorf("failed to read timestamp file: %w", err)
	}

	timestampStr := strings.TrimSpace(string(data))
	if timestampStr == "" {
		return time.Time{}, nil
	}

	t, err := time.Parse(time.RFC3339, timestampStr)
	if err != nil {
		return time.Time{}, fmt.Errorf("failed to parse timestamp: %w", err)
	}

	return t, nil
}

// writeLastProcessedTimestamp writes the last processed timestamp to file
func writeLastProcessedTimestamp(t time.Time) error {
	return os.WriteFile(lastTimestampFile, []byte(t.Format(time.RFC3339)), 0644)
}

// formatStringTimestamp formats an ISO timestamp string to IST
func formatStringTimestamp(ts string) string {
	t, err := time.Parse(time.RFC3339, ts)
	if err != nil {
		return ts
	}
	ist := time.FixedZone("IST", 5*3600+30*60)
	return t.In(ist).Format("2006-01-02 15:04:05 MST")
}
