package achiever

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

// Value marshalls Goals.
func (g Goals) Value() (driver.Value, error) {
	if g == nil {
		return nil, nil
	}

	return json.Marshal(g)
}

// Scan converts raw JSON ([]byte) to Goals
func (g *Goals) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(b, &g)
}
