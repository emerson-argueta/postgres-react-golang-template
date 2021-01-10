package response

import (
	"encoding/json"
	"log"
	"net/http"
)

// EncodeJSON encodes v to w in JSON format. Error() is called if encoding fails.
func EncodeJSON(w http.ResponseWriter, v interface{}, logger *log.Logger) {
	if err := json.NewEncoder(w).Encode(v); err != nil {
		ErrorResponse(w, err, http.StatusInternalServerError, logger)
	}
}
