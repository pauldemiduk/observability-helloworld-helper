// *
// ==> availability.go - [blueprint:include]
// *
// =>  availability.go explained:
//
//	Description: Structures and Data routines for service availability calculations
package models

import (
	"fmt"
	"hello-world-go/internal/metrics"
	"time"
)

// ServiceAvailability - Tracks uptime and downtime
type ServiceAvailability struct {
	UptimeSeconds   float64 `json:"uptime_seconds"`   // Float64 in seconds (native)
	UptimeHuman     string  `json:"uptime_human"`     // Readable format: "2d 3h 15m"
	DowntimeSeconds float64 `json:"downtime_seconds"` // Float64 in seconds (native)
	DowntimeHuman   string  `json:"downtime_human"`   // Readable format: "0d 2h 42m"
}

func GetUptime() (float64, string, error) {
	uptimeSeconds, err := metrics.QueryPrometheus("service_uptime_seconds")
	if err != nil {
		return 0, "", err
	}

	// Convert to human-readable format
	_, uptimeHuman := formatDuration(uptimeSeconds)

	return uptimeSeconds, uptimeHuman, nil
}

func GetDowntimeSeconds() float64 {
	logs := GraphQLDataStore
	downtimeCount := 0

	for _, log := range logs {
		if log.HTTPStatusCode >= 500 {
			downtimeCount++
		}
	}

	return float64(downtimeCount * 60) // Convert to seconds
}

func UpdateUptime(startTime time.Time) {
	// ✅ Correctly computes elapsed uptime since service start
	uptime := time.Since(startTime).Seconds()

	// ✅ Ensures Prometheus metric correctly reflects elapsed time
	metrics.ServiceUptime.Set(uptime)
}

func formatDuration(seconds float64) (float64, string) {
	totalMinutes := int(seconds / 60)
	hours := totalMinutes / 60
	minutes := totalMinutes % 60
	days := hours / 24
	hours = hours % 24

	humanReadable := fmt.Sprintf("%dd %dh %dm", days, hours, minutes)
	return seconds, humanReadable
}
