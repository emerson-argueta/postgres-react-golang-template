package jwt

import (
	"emersonargueta/m/v1/authorization"
	"emersonargueta/m/v1/delivery/middleware"

	"github.com/dgrijalva/jwt-go"
)

var _ authorization.Service = &Jwt{}

// Jwt represents an jwt implementation of authorization.Actions.
type Jwt struct {
	client     *Client
	middleware *middleware.Middleware
}

// Authorize a user by checking uuid parameter to the uuid extracted from the
// token cliams, if matching then generates a valid token and issues a new token
// pair to that user.
func (s *Jwt) Authorize(token *middleware.TokenPair, uuid string) (res *middleware.TokenPair, e error) {

	// Validate refresh token
	if jwtTokenPair, err := s.middleware.TokenPairStringToJWT(token); err != nil {
		return nil, ErrJWTAuth
	} else if uuidClaim, ok := jwtTokenPair.Refreshtoken.Claims.(jwt.MapClaims)["uuid"].(string); !ok {
		return nil, ErrJWTAuth
	} else if uuidClaim != uuid {
		return nil, ErrJWTAuth
	} else if tokenPair, err := s.middleware.GenerateTokenPair(uuid, middleware.AccestokenLimit, middleware.RefreshtokenLimit); err != nil {
		return nil, err
	} else {

		res = tokenPair
	}
	return res, nil
}

// ExtractUUIDFromToken from a domain model using JWT for authentication.
func (s *Jwt) ExtractUUIDFromToken(token interface{}) (res *string, e error) {
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
