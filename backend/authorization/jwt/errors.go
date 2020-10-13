package jwt

// jwt authorization errors.
const (
	ErrJWTAuth             = Error("invalid token, could not authorize")
	ErrRefreshTokenExpired = Error("refresh token is expired")
	ErrAcccessTokenExpired = Error("access token is expired")
)

// Error represents a general middleware error.
type Error string

// Error returns the error message.
func (e Error) Error() string { return string(e) }
