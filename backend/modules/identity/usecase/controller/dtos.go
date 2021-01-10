package controller

import (
	"emersonargueta/m/v1/modules/identity/dto"
	"emersonargueta/m/v1/modules/identity/usecase"
)

type userRequest struct {
	User          *dto.UserDTO            `json:"user,omitempty"`
	Authorization *usecase.ReauthorizeDTO `json:"authorization,omitempty"`
}
type userResponse struct {
	User          *dto.UserDTO       `json:"user,omitempty"`
	Authorization *map[string]string `json:"authorization,omitempty"`
}
type registerResponse struct {
	Message string `json:"message,omitempty"`
}
type updateUserResponse struct {
	Message string `json:"message,omitempty"`
}
type loginRequest struct {
	Email    *string `json:"email,omitempty"`
	Password *string `json:"password,omitempty"`
}
type loginResponse struct {
	Authorization *map[string]string `json:"authorization,omitempty"`
}
