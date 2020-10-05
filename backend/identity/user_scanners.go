package identity

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

// Value marshalls the Domains.
func (d Domains) Value() (driver.Value, error) {
	if d == nil {
		return nil, nil
	}
	return json.Marshal(d)
}

// Scan converts raw JSONB ([]byte) from postgres results and transforms it to
// Domains
func (d *Domains) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(b, &d)
}
