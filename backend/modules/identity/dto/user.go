package dto

import "emersonargueta/m/v1/modules/identity/domain/user"

// UserDTO for data transfer
type UserDTO struct {
	ID       *string `json:"id,omitempty"`
	Email    *string `json:"email,omitempty"`
	Password *string `json:"password,omitempty"`
	Role     *string `json:"role,omitempty"`
}

// ToDomain from dto
func ToDomain(dto UserDTO) (res user.User, e error) {

	role, e := user.NewRole(dto.Role)
	if e != nil {
		return nil, e
	}
	email, e := user.NewEmail(dto.Email)
	if e != nil {
		return nil, e
	}
	hashPassword, e := user.NewHashPassword(dto.Password)
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
