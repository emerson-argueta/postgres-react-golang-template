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
