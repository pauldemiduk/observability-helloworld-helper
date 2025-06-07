// *
// ==> sla_slo.go - [blueprint:include]
// *
// =>  sla_slo.go explained:
//
//	Description: Structures and Data routines for SLA violations, SLO compliance, and Error Budgets
package models

import (
	"fmt"
	"hello-world-go/internal/metrics"
	"time"
)

// SLAViolation represents a service-level agreement (SLA) breach
type SLAViolation struct {
	ID              string `json:"id"`
	AffectedService string `json:"affectedService"`
	BreachTime      string `json:"breachTime"`
	Reason          string `json:"reason"`
	Impact          string `json:"impact"`
}

// SLOCompliance tracks service-level objective (SLO) adherence
type SLOCompliance struct {
	ID     string `json:"id"`
	Target string `json:"target"`
	Actual string `json:"actual"`
	Status string `json:"status"`
}

// ErrorBudget tracks remaining error budget for a service
type ErrorBudget struct {
	ID                string `json:"id"`
	TotalAllocated    string `json:"totalAllocated"`
	Used              string `json:"used"`
	Remaining         string `json:"remaining"`
	DepletionForecast string `json:"depletionForecast"`
}

func GetSLAData(duration int) ([]SLAViolation, []SLOCompliance, []ErrorBudget, error) {
	durationStr := fmt.Sprintf("%dm", duration)

	// 🚀 Map metric types to queries
	queryMap := map[string]string{
		"slaViolations": fmt.Sprintf(`count_over_time(http_responses_total{http_status_code=~"5.."}[%s])`, durationStr),
		"sloCompliance": fmt.Sprintf(`avg_over_time(probe_success[%s]) * 100`, durationStr),
		"errorBudget":   fmt.Sprintf(`(1 - (sum(rate(http_requests_total{http_status_code=~"5.."}[%s])) / sum(rate(http_requests_total[%s])))) * 100`, durationStr, durationStr),
	}

	// Fetch Prometheus data
	slaViolations, _ := metrics.QueryPrometheus(queryMap["slaViolations"])
	sloCompliance, _ := metrics.QueryPrometheus(queryMap["sloCompliance"])
	errorBudget, _ := metrics.QueryPrometheus(queryMap["errorBudget"])

	// 🚀 Process SLA Violations
	var slaList []SLAViolation
	if slaViolations > 10 { // Simple threshold: More than 10 5xx errors in timeframe
		slaList = append(slaList, SLAViolation{
			ID:              "sla-violation-1",
			AffectedService: "checkout-service",
			BreachTime:      time.Now().UTC().Format(time.RFC3339),
			Reason:          "Excessive 5xx errors",
			Impact:          "High",
		})
	}

	// 🚀 Process SLO Compliance
	sloStatus := "Healthy"
	if sloCompliance < 99.0 {
		sloStatus = "At Risk"
	} else if sloCompliance < 95.0 {
		sloStatus = "Violation"
	}

	sloList := []SLOCompliance{
		{
			ID:     "slo-uptime",
			Target: "99.9% uptime",
			Actual: fmt.Sprintf("%.2f%%", sloCompliance),
			Status: sloStatus,
		},
	}

	// 🚀 Process Error Budget
	errorBudgetRemaining := fmt.Sprintf("%.2f%%", errorBudget)
	forecast := "Stable"
	if errorBudget < 20.0 {
		forecast = "Depleting fast"
	} else if errorBudget < 10.0 {
		forecast = "Critical depletion"
	}

	errorBudgetList := []ErrorBudget{
		{
			ID:                "budget-remaining",
			TotalAllocated:    "100 hours",
			Used:              fmt.Sprintf("%.2f%% used", 100-errorBudget),
			Remaining:         errorBudgetRemaining,
			DepletionForecast: forecast,
		},
	}

	return slaList, sloList, errorBudgetList, nil
}
