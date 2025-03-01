# 🔄 GraphQL Cross-Reference: Queries, Contexts & Implementation Status

## **📌 Purpose**
This document maps GraphQL touchpoints in the codebase, links them to relevant context use cases, outlines required data points, and tracks implementation status.

---

## **1️⃣ GraphQL Query Mapping in Code**
| **GraphQL Query**   | **Code Location**               | **Resolver Function**                    | **Notes** |
|---------------------|--------------------------------|-----------------------------------------|----------|
| `snapshots`        | `/internal/graph/loader.go`   | `snapshotQuery.Resolve`                | Retrieves all snapshots from YAML. |
| `contexts`         | `/internal/graph/loader.go`   | `contextQuery.Resolve`                 | Retrieves contexts for a snapshot. |
| `errorRates`       | `/internal/models/queries.go` | `GetErrorRatesFromLogs()`               | Fetches real error logs. **(Incomplete)** |
| `slowestAPICalls`  | `/internal/models/queries.go` | `GetSlowestAPICallsFromLogs()`          | Fetches slow API logs. **(To be added)** |
| `slaViolations`    | `/internal/models/queries.go` | `GetSLAViolationsFromLogs()`            | Retrieves SLA violations. **(To be added)** |

---

## **2️⃣ Context Use Cases & Required Data**
| **Context Use Case**       | **Required Data**                                      | **Current Status** |
|---------------------------|-------------------------------------------------------|----------------|
| **Error Rates**           | Total requests, failed requests, error codes         | Mock data only ⚠️ |
| **Slowest API Calls**      | Response times (P95, P99), API endpoints             | Not implemented ❌ |
| **Database Performance**   | Query execution time, transaction volume, slow queries | Not implemented ❌ |
| **API Traffic & Scaling**  | Request volume, autoscaling events                   | Not implemented ❌ |
| **SLA Violations**         | Downtime incidents, impact, compensation             | Not implemented ❌ |
| **Error Budget Tracking**  | Burn rate, remaining budget, forecasted depletion    | Not implemented ❌ |

---

## **3️⃣ Sample SQL-Like Queries for Data Retrieval**
_(These pseudo-queries represent what we need to fetch in Go code)_

🔹 **Error Rates Query:**
```sql
SELECT COUNT(*) AS total_requests,
       SUM(CASE WHEN status_code >= 400 THEN 1 ELSE 0 END) AS failed_requests,
       status_code, COUNT(status_code) AS error_count
FROM logs
WHERE timestamp >= NOW() - INTERVAL '1 hour'
GROUP BY status_code;
```

🔹 **Slowest API Calls Query:**
```sql
SELECT endpoint, PERCENTILE_CONT(0.95) WITHIN GROUP (ORDER BY response_time) AS p95_latency,
       PERCENTILE_CONT(0.99) WITHIN GROUP (ORDER BY response_time) AS p99_latency
FROM logs
WHERE timestamp >= NOW() - INTERVAL '1 hour'
GROUP BY endpoint;
```

🔹 **SLA Violations Query:**
```sql
SELECT incident_id, downtime_minutes, compensation_amount
FROM sla_violations
WHERE timestamp >= NOW() - INTERVAL '30 days';
```

---

## **4️⃣ Implementation Plan**
1️⃣ **Fix `GetErrorRatesFromLogs()` to actually fetch logs.** ✅ _(In progress)_
2️⃣ **Implement `GetSlowestAPICallsFromLogs()` for slow API tracking.** ⚠️ _(Pending)_
3️⃣ **Implement `GetSLAViolationsFromLogs()` to retrieve SLA data.** ⚠️ _(Pending)_
4️⃣ **Ensure GraphQL queries are correctly resolving real data.** 🚀 _(Next priority)_

