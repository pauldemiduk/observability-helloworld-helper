package models

import (
	"bufio"
	"log"
	"os"
	"regexp"
	"strings"
	"time"
)

// ✅ Fetches recent error events from structured logs
func FetchRecentErrorEvents() []EventPayload {
	log.Println("DEBUG: Fetching recent error events from logs...")

	// ✅ Query logs for errors (Replace with actual logic)
	events := []EventPayload{}

	// ✅ Scan service.log for errors
	file, err := os.Open("service.log")
	if err != nil {
		log.Println("ERROR: Unable to open log file:", err)
		return events
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, `"log_level":"error"`) {
			// ✅ Extract structured log details (Replace with actual parsing logic)
			event := EventPayload{
				EventID:   "parsed-event", // Replace with extracted event ID
				Category:  "error",
				Source:    "api",
				Severity:  "high",
				Message:   line,
				Timestamp: time.Now(),
			}
			events = append(events, event)
		}
	}

	log.Printf("INFO: Fetched %d error events from logs", len(events))
	return events
}

// ✅ Extracts a field from JSON-like log messages
func extractField(logLine, fieldName string) string {
	re := regexp.MustCompile(`"` + fieldName + `":"([^"]+)"`)
	matches := re.FindStringSubmatch(logLine)
	if len(matches) > 1 {
		return matches[1]
	}
	return "unknown"
}

// ✅ Extracts timestamps (assumes logs use ISO format)
func extractTimestamp(logLine string) time.Time {
	re := regexp.MustCompile(`"timestamp":"([^"]+)"`)
	matches := re.FindStringSubmatch(logLine)
	if len(matches) > 1 {
		parsedTime, err := time.Parse(time.RFC3339, matches[1])
		if err == nil {
			return parsedTime
		}
	}
	return time.Now()
}
