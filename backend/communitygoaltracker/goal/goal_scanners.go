package goal

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

// Value marshalls Achievers.
func (a Achievers) Value() (driver.Value, error) {
	if a == nil {
		return nil, nil
	}

	return json.Marshal(a)
}

// Scan converts raw JSON ([]byte) to Achievers.
func (a *Achievers) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(b, &a)
}

// Value marshalls Messages.
func (m Messages) Value() (driver.Value, error) {
	if m == nil {
		return nil, nil
	}

	return json.Marshal(m)
}

// Scan converts raw JSON ([]byte) to Messages.
func (m *Messages) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(b, &m)
}
