package http

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"trustdonations.org/m/v2/domain"
)

// ErrInvalidJSON indicated when a request body is incorrect after decoding
const ErrInvalidJSON = domain.Error("invalid json")

// Handler is a collection of all the service handlers.
type Handler struct {
	AdministratorHandler *AdministratorHandler
}

// ServeHTTP delegates a request to the appropriate subhandler.
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if strings.HasPrefix(r.URL.Path, "/api/administrator") {
		h.AdministratorHandler.ServeHTTP(w, r)
	} else {
		http.NotFound(w, r)
	}

}

// NotFound writes an API error message to the response.
func NotFound(w http.ResponseWriter) error {
	w.WriteHeader(http.StatusNotFound)
	_, err := w.Write([]byte(`{}` + "\n"))
	return err
}

// encodeJSON encodes v to w in JSON format. Error() is called if encoding fails.
func encodeJSON(w http.ResponseWriter, v interface{}, logger *log.Logger, serviceType string) {
	if err := json.NewEncoder(w).Encode(v); err != nil {
		ResponseError(w, err, http.StatusInternalServerError, logger, serviceType)
	}
}

/*
DecodeUnixTime function decodes a unix time string to time.time.
*/
func DecodeUnixTime(unixTime string) (time.Time, error) {
	i, err := strconv.ParseInt(unixTime, 10, 64)
	if err != nil {
		panic(err)
	}
	tm := time.Unix(i, 0)
	return tm, err
}

/*
EncodeUnixTime function encodes a time.Time to a unix string.
*/
func EncodeUnixTime(t time.Time) string {
	timeString := strconv.FormatInt(t.Unix(), 10)
	return timeString
}
