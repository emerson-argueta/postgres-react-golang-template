package administrator

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

// Value marshalls Subscription.
func (s Subscription) Value() (driver.Value, error) {
	if &s == nil {
		return nil, nil
	}

	return json.Marshal(s)
}

// Scan converts raw JSON ([]byte) to Subscription.
func (s *Subscription) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(b, &s)
}
