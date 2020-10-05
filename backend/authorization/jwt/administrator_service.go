package jwt

import (
	"emersonargueta/m/v1/authorization"
	"emersonargueta/m/v1/delivery/middleware"
	"emersonargueta/m/v1/domain/administrator"

	"github.com/dgrijalva/jwt-go"
)

var _ authorization.Actions = &Administrator{}

// Administrator represents an jwt implementation of authorization.Actions.
type Administrator struct {
	Usecase administrator.Usecase
	client  *Client
}

// Authorize an administrator with a valid token by issuing a new token pair to
// that administrator.
func (s *Administrator) Authorize(token *middleware.TokenPair) (res *middleware.TokenPair, e error) {

	// Validate refresh token
	if jwtTokenPair, err := middleware.TokenPairStringToJWT(token); err != nil {
		return nil, ErrJWTAuth
	} else if uuidClaim, ok := jwtTokenPair.Refreshtoken.Claims.(jwt.MapClaims)["uuid"].(string); !ok {
		return nil, ErrJWTAuth
	} else if exists, err := s.Usecase.Services.Administrator.Read(&administrator.Administrator{UUID: &uuidClaim}); err != nil {
		return nil, err
	} else if exists == nil {
		return nil, administrator.ErrAdministratorNotFound
	} else if tokenPair, err := middleware.GenerateTokenPair(*exists.UUID, middleware.AccestokenLimit, middleware.RefreshtokenLimit); err != nil {
		return nil, err
	} else {

		res = tokenPair
	}
	return res, nil
}

// ExtractUUIDFromToken from an administrator using JWT for authentication.
func (s *Administrator) ExtractUUIDFromToken(token interface{}) (res *string, e error) {
	jwtToken, ok := token.(*jwt.Token)
	if !ok {
		return nil, ErrJWTAuth
	}
	uuid, ok := jwtToken.Claims.(jwt.MapClaims)["uuid"].(string)
	if !ok {
		return nil, ErrJWTAuth
	}

	res = &uuid

	return res, nil
}
