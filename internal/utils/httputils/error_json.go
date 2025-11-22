package httputils

import (
	"encoding/json"
	"net/http"
)

func SendJSONError(w http.ResponseWriter, msg any, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(msg)
}
