package utils

import (
	"encoding/json"
	"net/http"
)

// Status represents the status part of the API response.
type Status struct {
	Success bool `json:"success"`
	Code    int  `json:"code"`
}

// APIResponse is the standardized JSON response structure.
type APIResponse struct {
	Status  Status      `json:"status"`
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message,omitempty"`
}

// respond is the base function for sending all JSON responses.
func respond(w http.ResponseWriter, statusCode int, success bool, data interface{}, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	response := APIResponse{
		Status: Status{
			Success: success,
			Code:    statusCode,
		},
		Data:    data,
		Message: message,
	}

	json.NewEncoder(w).Encode(response)
}

// RespondJSON sends a standard success response with data.
func RespondJSON(w http.ResponseWriter, statusCode int, data interface{}) {
	respond(w, statusCode, true, data, "")
}

// RespondError sends a standard error response with a message.
func RespondError(w http.ResponseWriter, statusCode int, message string) {
	respond(w, statusCode, false, nil, message)
}
