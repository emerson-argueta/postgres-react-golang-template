package persistence

import (
	"emersonargueta/m/v1/modules/identity/domain/user"
)

// UserDTO for scanning from database
type UserDTO struct {
	ID           *string `db:"id" json:"id"`
	Role         *string `db:"role" json:"role"`
	Email        *string `db:"email" json:"email"`
	HashPassword *string `db:"hashpassword" json:"hashpassword"`
}

// UserPersistenceToDomain from persistence
func UserPersistenceToDomain(dto UserDTO) (res user.User, e error) {
	role, e := user.NewRole(dto.Role)
	if e != nil {
		return nil, e
	}
	email, e := user.NewEmail(dto.Email)
	if e != nil {
		return nil, e
	}
	hashPassword, e := user.ToHashPassword(dto.HashPassword)
	if e != nil {
		return nil, e
	}
	userFields := &user.Fields{
		Role:         &role,
		Email:        &email,
		HashPassword: &hashPassword,
	}
	return user.Create(userFields, dto.ID)
}

// UserDomainToPersistence from domain
func UserDomainToPersistence(u user.User) *UserDTO {
	id := u.GetID()
	role, _ := u.GetRole().String()
	email := u.GetEmail().ToString()
	hashPassword := u.GetHashPassword().ToString()

	return &UserDTO{
		ID:           &id,
		Role:         &role,
		Email:        &email,
		HashPassword: &hashPassword,
	}

}
