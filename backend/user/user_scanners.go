package user

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

// Value marshalls the Services.
func (s Services) Value() (driver.Value, error) {
	if s == nil {
		return nil, nil
	}
	return json.Marshal(s)
}

// Scan converts raw JSONB ([]byte) from postgres results and transforms it to
// Services
func (s *Services) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(b, &s)
}
