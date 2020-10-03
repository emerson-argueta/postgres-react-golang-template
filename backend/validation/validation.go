package validation

import "regexp"

// ValidatePassword will validate a password whose length is greater or equal than 8
// characters long, returns nil if validation is successful.
func ValidatePassword(password string) error {
	if len(password) < 8 {
		return ErrValidationPassword
	}
	return nil
}

// ValidateUserEmail will validate emails. Returns nil if validation is successful.
func ValidateUserEmail(email string) error {
	emailRegexp := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

	if !emailRegexp.MatchString(email) {
		return ErrValidationUserEmail
	}
	return nil

}
