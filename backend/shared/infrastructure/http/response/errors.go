package response

import (
	"encoding/json"
	"log"
	"net/http"
)

// http errors
const (
	ErrHTTP        = Error("http error")
	ErrInternal    = Error("internal error")
	ErrInvalidJSON = Error("invalid json")
)

// Error represents a general http error.
type Error string

// Error returns the error message.
func (e Error) Error() string { return string(e) }

type errorResponse struct {
	Error string `json:"error,omitempty"`
}

// ErrorResponse writes an API error message to the response and logger.
func ErrorResponse(w http.ResponseWriter, err error, code int, logger *log.Logger) error {
	// Log error.
	logger.Printf("http error: %s (code=%d)", err, code)

	// Hide error from client if it is internal.
	if code == http.StatusInternalServerError {
		err = ErrInternal
	}

	// Write generic error response.
	w.WriteHeader(code)
	if err := json.NewEncoder(w).Encode(&errorResponse{Error: err.Error()}); err != nil {
		return err
	}
	return nil
}

// NotFound writes an API error message to the response.
func NotFound(w http.ResponseWriter) error {
	w.WriteHeader(http.StatusNotFound)
	_, err := w.Write([]byte(`{}` + "\n"))
	return err
}
