package authorization

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)

const (
	// AccestokenLimit represents the expiration time for an access token
	AccestokenLimit = time.Minute * 15
	// RefreshtokenLimit represents the expiration time for an refresh token
	RefreshtokenLimit = time.Hour * 24
)

var _ Processes = &jwtservice{}

type jwtservice struct {
	client *Client
}

// TokenPair with access token for jwt authentication and refresh token to
// refresh access.
type TokenPair struct {
	Accesstoken  string `json:"accesstoken,omitempty"`
	Refreshtoken string `json:"refreshtoken,omitempty"`
}

// TokenPairJWT with access token for jwt authentication and refresh token to
// refresh access.
type TokenPairJWT struct {
	Accesstoken  *jwt.Token `json:"accesstoken,omitempty"`
	Refreshtoken *jwt.Token `json:"refreshtoken,omitempty"`
}

// Authorize a user where the key is a jwt token. Returns the user's uuid if succesful,or ErrAuthorizationFailed
func (s *jwtservice) Authorize(key interface{}) (res *string, e error) {
	jwtToken, ok := key.(*jwt.Token)
	if !ok {
		return nil, ErrAuthorizationInvalidKey
	}
	uuid, ok := jwtToken.Claims.(jwt.MapClaims)["uuid"].(string)
	if !ok {
		return nil, ErrAuthorizationFailed
	}

	res = &uuid

	return res, nil
}
func (s *jwtservice) NewKey(uuid string) (res map[string]string, e error) {
	tokenPair, e := s.generateTokenPair(uuid)
	if e != nil {
		return nil, ErrAuthorizationKeyNotCreated
	}

	res = make(map[string]string)
	res["accesstoken"] = tokenPair.Accesstoken
	res["refreshtoken"] = tokenPair.Refreshtoken

	return res, e
}

// ReAuthorize a user when their accesstoken is expired. If the refresh token is
// not expired then generates a valid token pair and issues to the user.
func (s *jwtservice) ReAuthorize(key interface{}) (res map[string]string, e error) {
	token, ok := key.(*TokenPair)
	if !ok {
		return nil, ErrAuthorizationInvalidKey
	}

	if jwtTokenPair, err := s.tokenPairStringToJWT(token); err != nil {
		return nil, ErrAuthorizationInvalidKey
	} else if uuidClaim, ok := jwtTokenPair.Refreshtoken.Claims.(jwt.MapClaims)["uuid"].(string); !ok {
		return nil, ErrAuthorizationInvalidKey
	} else if newKey, err := s.NewKey(uuidClaim); err != nil {
		return nil, err
	} else {
		res = newKey
	}

	return res, nil
}

// tokenPairStringToJWT converts a string access and resfresh token from a TokenPair
// to a TokenPairJWT. It will fail if the refresh token is expired.
func (s *jwtservice) tokenPairStringToJWT(token *TokenPair) (res *TokenPairJWT, e error) {
	accessToken, e := jwt.ParseWithClaims(
		token.Accesstoken,
		jwt.MapClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(s.client.config.Authorization.Secret), nil
		},
	)
	if e != nil && e.(*jwt.ValidationError).Errors != jwt.ValidationErrorExpired {
		return nil, e
	}

	refreshToken, e := jwt.ParseWithClaims(
		token.Refreshtoken,
		jwt.MapClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(s.client.config.Authorization.Secret), nil
		},
	)
	if e != nil && e.(*jwt.ValidationError).Errors == jwt.ValidationErrorExpired {
		return nil, ErrAuthorizationKeyNotCreated
	} else if e != nil {
		return nil, e
	}

	res = &TokenPairJWT{Accesstoken: accessToken, Refreshtoken: refreshToken}
	return res, e
}

// GenerateTokenPair creates a token pair which contains access token and
// refresh token with exp limit for both access and refresh token.
func (s *jwtservice) generateTokenPair(uuid string) (*TokenPair, error) {
	accessToken := jwt.New(jwt.SigningMethodHS256)

	atclaims := accessToken.Claims.(jwt.MapClaims)
	atclaims["uuid"] = uuid
	atclaims["exp"] = time.Now().Add(AccestokenLimit).Unix()

	// Generate encoded accessToken.
	atstr, err := accessToken.SignedString([]byte(s.client.config.Authorization.Secret))
	if err != nil {
		return nil, err
	}

	refreshToken := jwt.New(jwt.SigningMethodHS256)

	rtClaims := refreshToken.Claims.(jwt.MapClaims)
	rtClaims["uuid"] = uuid
	rtClaims["exp"] = time.Now().Add(RefreshtokenLimit).Unix()

	// Generate encoded refreshToken.
	rtstr, err := refreshToken.SignedString([]byte(s.client.config.Authorization.Secret))
	if err != nil {
		return nil, err
	}

	return &TokenPair{Accesstoken: atstr, Refreshtoken: rtstr}, nil
}
