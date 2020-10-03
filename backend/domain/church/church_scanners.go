package church

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

// Value marshalls Administrators.
func (a Administrators) Value() (driver.Value, error) {
	if a == nil {
		return nil, nil
	}

	return json.Marshal(a)
}

// Scan converts raw JSON ([]byte) to Administrators.
func (a *Administrators) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(b, &a)
}

// Value marshalls Donators.
func (d Donators) Value() (driver.Value, error) {
	if d == nil {
		return nil, nil
	}

	return json.Marshal(d)
}

// Scan converts raw JSON ([]byte) to Donators.
func (d *Donators) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(b, &d)
}

// Value marshalls DonationCategories.
func (dc DonationCategories) Value() (driver.Value, error) {
	if dc == nil {
		return nil, nil
	}

	return json.Marshal(dc)
}

// Scan converts raw JSON ([]byte) to DonationCategories.
func (dc *DonationCategories) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(b, &dc)
}
