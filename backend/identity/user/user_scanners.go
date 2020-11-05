package user

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

// Value to store Role in database.
func (r Role) Value() (driver.Value, error) {
	return r.String()
}

// Scan string stored in database to Role type.
func (r *Role) Scan(value interface{}) (e error) {
	if value == nil {
		return nil
	}

	i, ok := value.(string)
	if !ok {
		return errors.New("type assertion to string failed")
	}

	*r, e = ToRole(i)

	return e
}
