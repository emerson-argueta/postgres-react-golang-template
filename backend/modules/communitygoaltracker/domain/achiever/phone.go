package achiever

import "emersonargueta/m/v1/validation"

// Phone value object
type Phone interface {
	ToString() string
}

type phone string

func (p phone) ToString() string {
	return string(p)
}

// NewPhone creates a phone. Empty string default. Returns error if invalid
// phone.
func NewPhone(newPhone *string) (res Phone, e error) {
	if newPhone == nil {
		return phone(""), nil
	}
	if e = validation.ValidatePhoneNumber(*newPhone); e != nil {
		return nil, e
	}

	return phone(*newPhone), nil
}
