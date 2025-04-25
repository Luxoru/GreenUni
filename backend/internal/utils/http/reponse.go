package response

import (
	"encoding/json"
	"net/http"
)

func WriteJson(w http.ResponseWriter, obj interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(obj)
}

type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message,omitempty"`
}

func SuccessResponse(data interface{}, message string) *Response {
	return &Response{
		Success: true,
		Data:    data,
	}
}

func ErrorResponse(message string) *Response {
	return &Response{
		Success: false,
		Message: message,
	}
}
