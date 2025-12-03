// Package render provides HTTP response rendering utilities for JSON output.
// It includes functions for rendering successful responses and error responses
// with consistent formatting across the API.
package render

import (
	"cli-stat-creator/internal/pipeline"
	"encoding/json"
	"net/http"
)

// JSON writes a JSON response to the HTTP response writer with the specified status code.
// It marshals the provided data and sets the Content-Type header to application/json.
// If marshaling fails, it writes a 500 error response instead.
func JSON(w http.ResponseWriter, status int, data any) {
	jsonBytes, err := json.Marshal(data)
	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		Error(w, 500, "Failed to encode response", "PROCESSING_ERROR")
		return
	}
	w.WriteHeader(status)
	w.Write(jsonBytes)
}

// Error writes a JSON error response to the HTTP response writer with the specified status code.
// It formats the error with a message and error code in a consistent JSON structure.
// If marshaling the error fails, it writes a plain text internal server error response.
func Error(w http.ResponseWriter, status int, message string, code string) {
	jsonError := struct {
		Error string `json:"error"`
		Code  string `json:"code"`
	}{
		Error: message,
		Code:  code,
	}
	jsonBytes, err := json.Marshal(jsonError)
	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte("internal server error"))
		return
	}
	w.WriteHeader(status)
	w.Write(jsonBytes)
}

// FilterResults creates a filtered map of pipeline results based on the specified flags.
// It includes only the requested result types (overall, byLevel, byPlayer) in the returned map.
func FilterResults(results pipeline.Results, overall, levels, players bool) any {
	resultMap := make(map[string]any)
	if overall {
		resultMap["overall"] = results.Overall
	}
	if levels {
		resultMap["byLevel"] = results.ByLevel
	}
	if players {
		resultMap["byPlayer"] = results.ByPlayer
	}
	return resultMap
}
