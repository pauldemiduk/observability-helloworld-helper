// *
// ==> prom_helper.go - [blueprint:include]
// *
// => prom_helper.go explained:
//
//	Description: Enables simple query integration witih Prometheus.
package metrics

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/graphql-go/graphql"
)

// NameValuePair -
// Generic struct for Prometheus vector results
type NameValuePair struct {
	Name  string  `json:"name"`  // Represents request_url, http_status_code, etc.
	Value float64 `json:"value"` // Represents latency, error count, etc.
}

// PrometheusQueryResult -
// Represents a standard response structure from Prometheus.
type PrometheusQueryResult struct {
	Status string `json:"status"`
	Data   struct {
		ResultType string `json:"resultType"`
		Result     []struct {
			Metric map[string]string `json:"metric"`
			Value  []interface{}     `json:"value"`
		} `json:"result"`
	} `json:"data"`
}

// GetMetricStats -
// Wrappler to abstract and fetch stat metrics from Prometheus given a metric name.
func GetMetricStats(duration int, metricType string) (interface{}, error) {
	durationStr := fmt.Sprintf("%dm", duration)

	// Map metric types to queries
	queryMap := map[string]string{
		"p95Latency":         fmt.Sprintf("histogram_quantile(0.95, sum(rate(http_response_duration_seconds_bucket[%s])) by (le))", durationStr),
		"p99Latency":         fmt.Sprintf("histogram_quantile(0.99, sum(rate(http_response_duration_seconds_bucket[%s])) by (le))", durationStr),
		"totalRequests":      fmt.Sprintf("sum(increase(http_requests_total[%s]))", durationStr),
		"requestRate":        fmt.Sprintf("sum(rate(http_requests_total[%s]))", durationStr),
		"slowRequests":       fmt.Sprintf("sum(rate(http_response_duration_seconds_count[%s])) by (request_url, operation, request_method)", durationStr),
		"failedRequests":     fmt.Sprintf(`sum(increase(http_responses_total{http_status_code=~"4..|5.."}[%s]))`, durationStr),
		"errorRate":          fmt.Sprintf(`sum(rate(http_responses_total{http_status_code=~"4..|5.."}[%s])) / sum(rate(http_requests_total[%s])) * 100`, durationStr, durationStr), // ✅ Converts to %
		"cpuUtilization":     fmt.Sprintf(`100 - (avg(rate(node_cpu_seconds_total{mode="idle"}[%s])) * 100)`, durationStr),
		"memoryUtilization":  "node_memory_active_bytes / node_memory_total_bytes * 100",
		"diskUtilization":    `100 - ((node_filesystem_avail_bytes{mountpoint="/"} / node_filesystem_size_bytes{mountpoint="/"}) * 100)`,
		"networkUtilization": fmt.Sprintf(`sum(rate(node_network_receive_bytes_total[%s])) + sum(rate(node_network_transmit_bytes_total[%s]))`, durationStr, durationStr),
		"slaViolations":      fmt.Sprintf("count_over_time(http_requests_total{http_status_code=~'5..'}[%s])", durationStr),
		"sloCompliance":      fmt.Sprintf("avg_over_time(probe_success[%s]) * 100", durationStr),
		"errorBudget":        fmt.Sprintf("(1 - (sum(rate(http_requests_total{http_status_code=~'5..'}[%s])) / sum(rate(http_requests_total[%s])))) * 100", durationStr, durationStr), // ✅ Uses http_status_code
	}

	// Validate metric type
	query, exists := queryMap[metricType]
	if !exists {
		return nil, fmt.Errorf("invalid metric type: %s", metricType)
	}

	// Execute query
	result, err := QueryPrometheus(query)
	if err != nil {
		log.Printf("WARNING: No data returned for metric '%s' (Query: %s), defaulting to 0", metricType, query)
		return 0, nil
	}

	// Convert errorRate to percentage if needed
	if metricType == "errorRate" {
		result = result * 100.0 // ✅ Convert to percentage
	}

	return result, nil
}

// QueryPrometheus -
// Fetches stat metrics from Prometheus given a query string.
func QueryPrometheus(query string) (float64, error) {
	baseURL := "http://localhost:9090/api/v1/query"
	fullURL := fmt.Sprintf("%s?query=%s", baseURL, url.QueryEscape(query))

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Get(fullURL)
	if err != nil {
		return 0, fmt.Errorf("failed to fetch Prometheus data: %w", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, fmt.Errorf("failed to read Prometheus response: %w", err)
	}

	var result PrometheusQueryResult
	if err := json.Unmarshal(body, &result); err != nil {
		return 0, fmt.Errorf("failed to parse Prometheus response: %w", err)
	}

	// Ensure we have data
	if len(result.Data.Result) == 0 || len(result.Data.Result[0].Value) < 2 {
		log.Printf("WARNING: No valid data returned from Prometheus for query: %s", query)
		return 0, nil
	}

	// Extract the first metric value
	valueStr, ok := result.Data.Result[0].Value[1].(string)
	if !ok {
		return 0, fmt.Errorf("value is not a string")
	}

	var parsedValue float64
	_, err = fmt.Sscanf(valueStr, "%f", &parsedValue)
	if err != nil {
		return 0, fmt.Errorf("failed to parse Prometheus value: %s", valueStr)
	}

	return parsedValue, nil
}

// ConvertVectorToNameValuePairs -
// Fetches Prometheus vector data dynamically
func ConvertVectorToNameValuePairs(duration int, metricType string) ([]NameValuePair, error) {
	durationStr := fmt.Sprintf("%dm", duration)

	// 🚀 Map vector metric types to PromQL queries
	queryMap := map[string]string{
		"topErrorCodes":   fmt.Sprintf(`topk(5, sum(increase(http_responses_total{http_status_code=~"4..|5.."}[%s])) by (http_status_code))`, durationStr),
		"failedEndpoints": fmt.Sprintf(`topk(5, sum(increase(http_responses_total{http_status_code=~"4..|5.."}[%s])) by (request_url))`, durationStr),
		"topEndpoints":    fmt.Sprintf(`topk(5, histogram_quantile(0.95, sum(rate(http_response_duration_seconds_bucket{request_url!='/metrics'}[%s])) by (le, request_url)))`, durationStr),
		"slowEndpoints":   fmt.Sprintf(`topk(5, histogram_quantile(0.95, sum(rate(http_response_duration_seconds_bucket{request_url!='/metrics'}[%s])) by (le, request_url)))`, durationStr),
		"topConsumers":    fmt.Sprintf(`topk(5, sum(increase(http_requests_total[%s])) by (client_id))`, durationStr),
	}

	// 🚀 Validate requested metric type
	query, exists := queryMap[metricType]
	if !exists {
		return nil, fmt.Errorf("invalid vector metric type: %s", metricType)
	}

	// 🚀 Execute Prometheus query
	results, err := QueryPrometheusVector(query)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch Prometheus vector data: %w", err)
	}

	// 🚀 Determine label key for Prometheus dimension
	labelKeyMap := map[string]string{
		"topErrorCodes":   "http_status_code",
		"failedEndpoints": "request_url",
		"topEndpoints":    "request_url",
		"slowEndpoints":   "request_url",
		"topConsumers":    "client_id",
	}

	labelKey, exists := labelKeyMap[metricType]
	if !exists {
		return nil, fmt.Errorf("invalid label mapping for metric type: %s", metricType)
	}

	// 🚀 Convert Prometheus response to NameValuePair list
	var entries []NameValuePair
	for _, entry := range results {
		if dimension, exists := entry.Metric[labelKey]; exists {
			// Ensure value is a string before parsing
			value, ok := entry.Value[1].(string)
			if !ok {
				continue
			}

			var parsedValue float64
			if _, err := fmt.Sscanf(value, "%f", &parsedValue); err != nil {
				continue
			}

			entries = append(entries, NameValuePair{
				Name:  dimension,
				Value: parsedValue,
			})
		}
	}

	return entries, nil
}

// GetMetricData -
// GraphQL helper to fetch Prometheus metrics.
func GetMetricData(params graphql.ResolveParams) (interface{}, error) {
	metricType := params.Info.FieldName // Auto-detect query type
	duration, _ := params.Args["duration"].(int)
	if duration <= 0 {
		duration = 5
	}

	// ✅ Capture both return values
	data, err := GetMetricStats(duration, metricType)
	if err != nil {
		return nil, err // ✅ Return error properly
	}

	return data, nil // ✅ Return the actual data
}

// QueryPrometheusVector -
// Fetches vector metrics from Prometheus given a query string.
func QueryPrometheusVector(query string) ([]struct {
	Metric map[string]string `json:"metric"`
	Value  []interface{}     `json:"value"`
}, error) {
	baseURL := "http://localhost:9090/api/v1/query"
	fullURL := fmt.Sprintf("%s?query=%s", baseURL, url.QueryEscape(query))

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Get(fullURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch Prometheus data: %w", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read Prometheus response: %w", err)
	}

	var result PrometheusQueryResult
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to parse Prometheus response: %w", err)
	}

	// ✅ Explicitly return the existing slice type
	return result.Data.Result, nil
}
