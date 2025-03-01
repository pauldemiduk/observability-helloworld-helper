# GraphQL Execution Flow in Hello-World-Go

## 🔍 Overview
This document describes the execution lifecycle of a GraphQL query in the `hello-world-go` service, detailing each step, related code files, and functions involved.

## 🚀 GraphQL Query Execution Flow

### **Step 1️⃣ – Client Sends Query**
- **Trigger**: Client makes an HTTP `POST` request to `/graph` with a GraphQL query.
- **Code Location**:
  - **Entry Point:** `/internal/handlers/graph.go → GraphQLHandler`
  - **Function:** `GraphQLHandler(w http.ResponseWriter, r *http.Request)`

### **Step 2️⃣ – GraphQL Parses the Query**
- **Trigger**: GraphQL engine parses and validates the request.
- **Code Location**:
  - **Schema Definition:** `/internal/graph/loader.go`
  - **Relevant Structs:** `queryType`, `contextType`, `snapshotType`

### **Step 3️⃣ – Resolver Executes Context Query**
- **Trigger**: The GraphQL engine resolves the requested fields.
- **Code Location**:
  - **GraphQL Resolver:** `/internal/graph/loader.go`
  - **Function:** `contextQuery.Resolve`
  - **Calls:** `models.GetSnapshotData(snapshotID)`

### **Step 4️⃣ – GetSnapshotData() Retrieves Contexts**
- **Trigger**: The function retrieves context data from snapshots.yaml (for now).
- **Code Location**:
  - **File:** `/internal/models/snapshot.go`
  - **Function:** `GetSnapshotData(snapshotID)`
  - **Current Behavior:** Returns contexts from **static YAML data**
  
### **Step 5️⃣ – Fetching Real Data from Logs**
- **Trigger**: The resolver attempts to replace mock data with real event logs.
- **Code Location**:
  - **File:** `/internal/models/queries.go`
  - **Function:** `GetErrorRatesFromLogs()`
  - **Calls:** `FetchRecentErrorEvents()`
  - **Issue:** Logs are not being fetched properly.

### **Step 6️⃣ – Processing Events into Error Rates**
- **Trigger**: `processErrorRates(events)` attempts to process logs into structured error rates.
- **Code Location**:
  - **File:** `/internal/models/queries.go`
  - **Function:** `processErrorRates(events)`
  - **Issue:** No real data is being processed due to missing log queries.

### **Step 7️⃣ – GraphQL Returns Final Response**
- **Trigger**: The GraphQL engine sends back the constructed response.
- **Code Location**:
  - **File:** `/internal/handlers/graph.go`
  - **Function:** `GraphQLHandler`



