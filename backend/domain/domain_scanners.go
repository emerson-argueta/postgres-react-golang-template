package domain

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

// Value marshalls AccountStatement.
func (a AccountStatement) Value() (driver.Value, error) {
	if &a == nil {
		return nil, nil
	}

	return json.Marshal(a)
}

// Scan converts raw JSON ([]byte) to AccountStatement.
func (a *AccountStatement) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(b, &a)
}
