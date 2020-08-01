package useful

import (
	"encoding/json"
	"net/http"
	"person/internal/dto"
)

func BuildSuccess(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_, _ = w.Write(response)
}

func BuildError(w http.ResponseWriter, code int, message string) {
	response, _ := json.Marshal(buildError(message))
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_, _ = w.Write(response)
}

func buildError(message string) dto.Error {
	return dto.Error{Message: message}
}
