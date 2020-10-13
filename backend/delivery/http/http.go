package http

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"
)

// RoutePrefix to use for REST api
const RoutePrefix = "/api/v1"

// NotFound writes an API error message to the response.
func NotFound(w http.ResponseWriter) error {
	w.WriteHeader(http.StatusNotFound)
	_, err := w.Write([]byte(`{}` + "\n"))
	return err
}

// EncodeJSON encodes the jsonFormatObject as an httpResponseWriter. Error() is
// called if encoding fails.
func EncodeJSON(httpResponseWriter http.ResponseWriter, jsonFormatObject interface{}, logger *log.Logger) {
	if err := json.NewEncoder(httpResponseWriter).Encode(jsonFormatObject); err != nil {
		ResponseError(httpResponseWriter, err, http.StatusInternalServerError, logger)
	}
}

// DecodeUnixTime function decodes a unix time string to time.time.
func DecodeUnixTime(unixTime string) (time.Time, error) {
	i, err := strconv.ParseInt(unixTime, 10, 64)
	if err != nil {
		panic(err)
	}
	tm := time.Unix(i, 0)
	return tm, err
}

// EncodeUnixTime function encodes a time.Time to a unix string.
func EncodeUnixTime(t time.Time) string {
	timeString := strconv.FormatInt(t.Unix(), 10)
	return timeString
}
