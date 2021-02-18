package usecase

type LoginUserDTO struct {
	Email    string
	Password string
}
type TokenDTO struct {
	AccessToken  string `json:"accesstoken,omitempty"`
	RefreshToken string `json:"refreshtoken,omitempty"`
}
