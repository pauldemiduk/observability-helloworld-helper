// *
// ==> datastore.go - [blueprint:include]
// *
// =>  datastore.go explained:
//
//	Description: Structures and Data routines for in-memory log storage
package models

import "log"

// In-memory store for GraphQL logs
// ** ROADMAP.>, iterate from []LogEntry in-memory datastore to other optons.
var GraphQLDataStore []LogEntry

// StoreLogForGraphQL -
// Forces logs into []LogEntry in-memory datastore.
func StoreLogForGraphQL(entry LogEntry) {
	log.Printf("DEBUG: Storing Log in GraphQLDataStore -> Status: %d, Message: %s", entry.HTTPStatusCode, entry.Message)
	GraphQLDataStore = append(GraphQLDataStore, entry)
}

// fetchLogsFromStore -
// Fetch logs from the []LogEntry in-memory datastore.
func fetchLogsFromStore(limit int) []LogEntry {
	if len(GraphQLDataStore) > limit {
		return GraphQLDataStore[len(GraphQLDataStore)-limit:]
	}
	return GraphQLDataStore
}
