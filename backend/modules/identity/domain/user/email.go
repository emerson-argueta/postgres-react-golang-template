package user

import (
	"emersonargueta/m/v1/validation"
)

// Email value object
type Email interface {
	ToString() string
}

type email string

func (e email) ToString() string {
	return string(e)
}

// NewEmail created a user, returns error if email is nil
func NewEmail(newEmail *string) (res Email, e error) {
	if newEmail == nil {
		return nil, ErrUserIncompleteDetails
	}
	if e = validation.ValidateUserEmail(*newEmail); e != nil {
		return nil, e
	}

	return email(*newEmail), nil
}
