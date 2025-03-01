// reload.go -
package handlers

import (
	"hello-world-go/internal/config"
	"net/http"
)

// ReloadConfigHandler - HTTP handler to reload config dynamically
func ReloadConfigHandler(w http.ResponseWriter, r *http.Request) {
	config.ReloadConfig() // ✅ Reload config at runtime
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "Configuration reloaded successfully"}`))
}
