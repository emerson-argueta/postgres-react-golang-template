package transaction

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

// Value marshalls Donation.
func (d Donation) Value() (driver.Value, error) {
	if &d == nil {
		return nil, nil
	}

	return json.Marshal(d)
}

// Scan converts raw JSON ([]byte) to Donation.
func (d *Donation) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(b, &d)
}
