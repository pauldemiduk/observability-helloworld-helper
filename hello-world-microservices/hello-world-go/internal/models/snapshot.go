package models

import (
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
)

// ✅ Define Data Structures Matching the YAML Schema
type Snapshot struct {
	ID        string    `yaml:"id"`
	Name      string    `yaml:"name"`
	CreatedAt string    `yaml:"createdAt"`
	Duration  string    `yaml:"duration"`
	Contexts  []Context `yaml:"contexts"`
}

type Context struct {
	ID              string           `yaml:"id"`
	Name            string           `yaml:"name"`
	Description     string           `yaml:"description"`
	Queries         []QueryExecution `yaml:"queries"`
	ErrorRates      []ErrorRate      `yaml:"errorRates"`
	SlowestAPICalls []SlowAPICall    `yaml:"slowestAPICalls"`
	FailedEndpoints []FailedEndpoint `yaml:"failedEndpoints"`
}

type QueryExecution struct {
	ID   string `yaml:"id"`
	Name string `yaml:"name"`
	Data string `yaml:"data"`
}

type ErrorRate struct {
	ID              string           `yaml:"id"`
	Name            string           `yaml:"name"`
	TotalRequests   int              `yaml:"totalRequests"`
	FailedRequests  int              `yaml:"failedRequests"`
	ErrorRate       float64          `yaml:"errorRate"`
	TopErrorCodes   []ErrorCode      `yaml:"topErrorCodes"`
	FailedEndpoints []FailedEndpoint `yaml:"failedEndpoints"`
}

type ErrorCode struct {
	Code  int `yaml:"code"`
	Count int `yaml:"count"`
}

type FailedEndpoint struct {
	Endpoint    string  `yaml:"endpoint"`
	Failures    int     `yaml:"failures"`
	FailureRate float64 `yaml:"failureRate"`
}

type SlowAPICall struct {
	ID           string         `yaml:"id"`
	Name         string         `yaml:"name"`
	SlowRequests int            `yaml:"slowRequests"`
	P95Latency   float64        `yaml:"p95Latency"`
	P99Latency   float64        `yaml:"p99Latency"`
	TopEndpoints []SlowEndpoint `yaml:"topEndpoints"`
}

type SlowEndpoint struct {
	Endpoint string  `yaml:"endpoint"`
	Latency  float64 `yaml:"latency"`
}

// ✅ Load YAML File and Parse into Go Structs
func LoadSnapshotsFromYAML(filename string) ([]Snapshot, error) {

	log.Printf("INFO: Attempting to load YAML file: %s", filename) // ✅ Add this log

	data, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Printf("ERROR: Failed to read YAML file: %v", err)
		return nil, err
	}

	log.Println("INFO: YAML file read successfully, attempting to parse...") // ✅ Add this log

	var snapshots []Snapshot
	err = yaml.Unmarshal(data, &snapshots)
	if err != nil {
		log.Printf("ERROR: Failed to parse YAML: %v", err)
		return nil, err
	}

	log.Printf("INFO: Loaded %d snapshots from YAML", len(snapshots)) // ✅ Add this log
	return snapshots, nil
}

func GetSnapshots() []Snapshot {
	log.Println("INFO: Fetching snapshots...") // ✅ Add log before calling the loader

	filename := "configs/snapshots.yaml" // ✅ Correct file path
	snapshots, err := LoadSnapshotsFromYAML(filename)

	if err != nil {
		log.Println("WARNING: Falling back to empty snapshot list")
		return []Snapshot{}
	}

	log.Printf("INFO: Returning %d snapshots", len(snapshots)) // ✅ Confirm how many snapshots are returned
	return snapshots
}

// ✅ Retrieve a Single Snapshot by ID
func GetSnapshotData(snapshotID string) []Context {
	snapshots := GetSnapshots()
	for _, snapshot := range snapshots {
		if snapshot.ID == snapshotID {
			log.Printf("INFO: Found snapshot %s with %d contexts", snapshotID, len(snapshot.Contexts))
			return snapshot.Contexts
		}
	}
	log.Printf("WARNING: Snapshot %s not found", snapshotID)
	return []Context{}
}
