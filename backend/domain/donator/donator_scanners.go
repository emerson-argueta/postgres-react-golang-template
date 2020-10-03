package donator

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

// Value marshalls Churches.
func (c Churches) Value() (driver.Value, error) {
	if c == nil {
		return nil, nil
	}

	return json.Marshal(c)
}

// Scan converts raw JSON ([]byte) to Churches
func (c *Churches) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(b, &c)
}
