// *
// ==> saturation.go - [blueprint:include]
// *
// =>  saturation.go explained:
//
//	Description: Structures and Data routines for CPU, Memory, Disk, Network Utilization
package models

// ServiceSaturation -
type ServiceSaturation struct {
	CPUUtilization     float64 `json:"cpu_utilization"`    // CPU % usage
	MemoryUtilization  float64 `json:"memory_utilization"` // Memory % usage
	DiskUtilization    float64 `json:"disk_utilization"`
	NetworkUtilization float64 `json:"network_utilization"`
}
