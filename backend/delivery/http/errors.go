package http

import (
	"encoding/json"
	"log"
	"net/http"

	"trustdonations.org/m/v2/domain"
	"trustdonations.org/m/v2/domain/administrator"
)

// http errors
const (
	ErrHTTP = Error("http error")
)

// Error represents a general http error.
type Error string

// Error returns the error message.
func (e Error) Error() string { return string(e) }

// ResponseError writes an API error message to the response and logger.
func ResponseError(w http.ResponseWriter, err error, code int, logger *log.Logger, serviceType string) error {
	// Log error.
	logger.Printf("http error: %s (code=%d)", err, code)

	// Hide error from client if it is internal.
	if code == http.StatusInternalServerError {
		switch serviceType {
		case "Administrator":
			err = administrator.ErrAdministratorInternal
		default:
			err = domain.ErrDomainInternal
		}

	}

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
