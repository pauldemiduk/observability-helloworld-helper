// log_rotation.go - Optimized structured logging for the Hello World microservice
package logging

import (
	"fmt"
	"hello-world-go/internal/models"
	"log"
	"os"
	"path/filepath"
	"time"
)

// ManageLogRotation -
// Rotate log files when they exceed 10MB
func ManageLogRotation(logFile string) {
	// ✅ If log rotation is disabled, enforce `LOG_MAX_SIZE_MB`
	if !models.HelloInstance.Config.LogRotationEnabled {
		EnforceMaxLogSize(logFile)
		return
	}

	// ✅ Standard log rotation logic (if enabled)
	fileInfo, err := os.Stat(logFile)
	if err == nil && fileInfo.Size() > int64(models.HelloInstance.Config.LogRotationSizeMB)*1024*1024 {
		backupName := fmt.Sprintf("%s.%s", logFile, time.Now().Format("20060102-150405"))
		err := os.Rename(logFile, backupName)
		if err != nil {
			log.Printf("ERROR: Failed to rotate log file: %v", err)
			return
		}
		log.Printf("INFO: Log file rotated: %s -> %s", logFile, backupName)

		// ✅ Apply log retention policy
		CleanupOldLogs(logFile)
	}
}

// EnforceMaxLogSize -
// ✅ Enforce a Hard Limit on Log File Size (No Rotation)
func EnforceMaxLogSize(logFile string) {
	fileInfo, err := os.Stat(logFile)
	if err == nil && fileInfo.Size() > int64(models.HelloInstance.Config.LogMaxSizeMB)*1024*1024 {
		log.Printf("INFO: Log file exceeded %dMB, clearing and starting over.", models.HelloInstance.Config.LogMaxSizeMB)
		os.Truncate(logFile, 0) // ✅ Clears the file without deleting it
	}
}

// CleanupOldLogs -
// Keeps only the last `N` log backups, or deletes all if `LOG_RETENTION_COUNT=0`
func CleanupOldLogs(logFile string) {
	logDir := filepath.Dir(logFile)
	logBase := filepath.Base(logFile)
	retentionCount := models.HelloInstance.Config.LogRetentionCount

	// ✅ Find all rotated log backups
	matches, _ := filepath.Glob(fmt.Sprintf("%s/%s.*", logDir, logBase))

	// ✅ If retention count is 0, delete all backups
	if retentionCount == 0 {
		for _, oldLog := range matches {
			err := os.Remove(oldLog)
			if err == nil {
				log.Printf("INFO: Deleted log file (retention disabled): %s", oldLog)
			} else {
				log.Printf("ERROR: Failed to delete log file: %s", oldLog)
			}
		}
		return
	}

	// ✅ Keep only the last `N` log backups
	if len(matches) > retentionCount {
		toDelete := matches[:len(matches)-retentionCount]
		for _, oldLog := range toDelete {
			err := os.Remove(oldLog)
			if err == nil {
				log.Printf("INFO: Deleted old log file: %s", oldLog)
			} else {
				log.Printf("ERROR: Failed to delete log file: %s", oldLog)
			}
		}
	}
}
