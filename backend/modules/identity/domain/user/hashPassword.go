package user

import (
	"emersonargueta/m/v1/validation"

	"golang.org/x/crypto/bcrypt"
)

// HashPassword value object
type HashPassword interface {
	ToString() string
}
type hashPassword string

func (h hashPassword) ToString() string {
	return string(h)
}

// ToHashPassword from hash password string, returns error if nil
func ToHashPassword(hashPasswordStr *string) (res HashPassword, e error) {
	if hashPasswordStr == nil {
		return nil, ErrUserIncompleteDetails
	}

	return hashPassword(*hashPasswordStr), nil
}

// NewHashPassword created from a password, returns error if nil or password is
// invalid
func NewHashPassword(password *string) (res HashPassword, e error) {
	if password == nil {
		return nil, ErrUserIncompleteDetails
	}
	if e = validation.ValidatePassword(*password); e != nil {
		return nil, e
	}

	hash, e := bcrypt.GenerateFromPassword([]byte(*password), bcrypt.DefaultCost)
	hashString := string(hash)

	return hashPassword(hashString), e
}

// CompareHashAndPassword to validate password, returns error if nil or email is
// invalid
func CompareHashAndPassword(hashPassword HashPassword, password *string) error {
	if password == nil {
		return ErrUserIncompleteDetails
	}

	hp := hashPassword.ToString()
	e := bcrypt.CompareHashAndPassword([]byte(hp), []byte(*password))

	if e != nil && e == bcrypt.ErrMismatchedHashAndPassword {
		return ErrUserIncorrectCredentials
	}
	return e

}
