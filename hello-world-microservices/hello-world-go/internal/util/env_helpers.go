// http_helpers.go -
package util

import (
	"fmt"
	"hello-world-go/internal/logging"
	"hello-world-go/internal/models"

	"log"
	"os"
	"strconv"
	"sync"
)

var warningLogged sync.Map

// util.GetEnvOrDefault -
// Retrieves an environment variable or returns a default value.
// GetEnvOrDefault - Retrieves an environment variable or returns a default value.
func GetEnvOrDefault(envVar, defaultValue string) string {
	value, exists := os.LookupEnv(envVar)
	if exists {
		return value
	}

	// ✅ Check log level properly (debug-only logging)
	if models.HelloInstance.Config.LogLevel == string(logging.LogLevelDebug) {
		if _, loaded := warningLogged.LoadOrStore(envVar, true); !loaded {
			log.Printf("DEBUG: %s is not set, using default: %s", envVar, defaultValue)
		}
	}

	return defaultValue
}

// util.GetEnvIntValidated -
// Retrieves an environment variable as an int, ensuring it's within a valid range.
func GetEnvIntValidated(envVar string, defaultValue, min, max int) int {
	valueStr := GetEnvOrDefault(envVar, fmt.Sprintf("%d", defaultValue))
	value, err := strconv.Atoi(valueStr)

	if err != nil || value < min || value > max {
		log.Printf("WARNING: Invalid %s value '%s', resetting to default: %d", envVar, valueStr, defaultValue)
		return defaultValue
	}

	return value
}

// util.GetEnvFloatValidated -
func GetEnvFloatValidated(envVar string, defaultValue, min, max float64) float64 {
	valueStr := GetEnvOrDefault(envVar, fmt.Sprintf("%.1f", defaultValue))
	value, err := strconv.ParseFloat(valueStr, 64)

	if err != nil || value < min || value > max {
		log.Printf("WARNING: Invalid %s value '%s', resetting to default: %.1f", envVar, valueStr, defaultValue)
		return defaultValue
	}

	return value
}

// util.GetEnvOrDefaultBool -
// Retrieves a boolean environment variable or returns a default value.
func GetEnvOrDefaultBool(envVar string, defaultValue bool) bool {
	valueStr := GetEnvOrDefault(envVar, "")
	if valueStr == "" {
		return defaultValue
	}

	parsedValue, err := strconv.ParseBool(valueStr)
	if err != nil {
		log.Printf("WARNING: Invalid boolean value for %s: '%s', defaulting to %t", envVar, valueStr, defaultValue)
		return defaultValue
	}

	return parsedValue
}
