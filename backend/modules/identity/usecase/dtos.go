package usecase

type LoginUserDTO struct {
	Email    string
	Password string
}
type ReauthorizeDTO map[string]string
