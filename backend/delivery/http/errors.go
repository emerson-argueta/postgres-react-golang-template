package http

import (
	"encoding/json"
	"log"
	"net/http"
)

// http errors
const (
	ErrHTTP = Error("http error")
	// ErrInvalidJSON indicated when a request body is incorrect after decoding
	ErrInvalidJSON = Error("invalid json")
)

// Error represents a general http error.
type Error string

// Error returns the error message.
func (e Error) Error() string { return string(e) }

// ResponseError writes an API error message to the response and logger.
func ResponseError(w http.ResponseWriter, err error, code int, logger *log.Logger) error {
	// Log error.
	logger.Printf("http error: %s (code=%d)", err, code)

	// Write generic error response.
	w.WriteHeader(code)
	if err := json.NewEncoder(w).Encode(&errorResponse{Error: err.Error()}); err != nil {
		return err
	}
	return nil
}

// errorResponse is a generic response for sending a error.
type errorResponse struct {
	Error string `json:"error,omitempty"`
}
