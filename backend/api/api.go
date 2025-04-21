// Package api contains basic error handling methods.
package api

import (
	"encoding/json"
	"net/http"
)

// Error represents an error message and status code to be returned.
type Error struct {
	Code    int
	Message string
}

// writeError writes an error response with a given message and HTTP status code.
func writeError(w http.ResponseWriter, message string, code int) {
	resp := Error{
		Code:    code,
		Message: message,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	err := json.NewEncoder(w).Encode(resp)
	if err != nil {
		return
	}
}

// RequestErrorhandler handles client-side (bad request) errors.
var RequestErrorhandler = func(w http.ResponseWriter, err error) {
	writeError(w, err.Error(), http.StatusBadRequest)
}

// InternalErrorHandler handles server-side (internal) errors.
var InternalErrorHandler = func(w http.ResponseWriter) {
	writeError(w, "An Unexpected Error occured", http.StatusInternalServerError)
}
