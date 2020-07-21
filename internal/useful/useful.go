package useful

import (
	"encoding/json"
	"net/http"
)

func BuildSuccess(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_, _ = w.Write(response)
}

func BuildError(w http.ResponseWriter, code int, message string) {
	response, _ := json.Marshal(buildJsonMessage(message))
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_, _ = w.Write(response)
}

func buildJsonMessage(message string) map[string]string {
	return map[string]string{"message":message}
}