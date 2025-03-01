// outbound.go -
package handlers

import (
	"encoding/json"
	"hello-world-go/internal/util"
	"net/http"
)

// outbound route
func OutboundHandler(w http.ResponseWriter, r *http.Request) {
	util.SetSecurityHeaders(w) // ✅ Secure responses
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Outbound request processed successfully."})
}
