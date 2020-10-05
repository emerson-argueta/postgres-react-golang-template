package middleware

import (
	"time"

	"emersonargueta/m/v1/env"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/middleware"
)

const (
	// AccestokenLimit represents the expiration time for an access token
	AccestokenLimit = time.Minute * 15
	// RefreshtokenLimit represents the expiration time for an refresh token
	RefreshtokenLimit = time.Hour * 24
)

// TokenPair represents the access token for authentication and refresh token to
// refresh access, both as strings.
type TokenPair struct {
	Accesstoken  string `json:"accesstoken,omitempty"`
	Refreshtoken string `json:"refreshtoken,omitempty"`
}

// TokenPairJWT represents the access token for authentication and refresh token to
// refresh access, both as jwt.Token.
type TokenPairJWT struct {
	Accesstoken  *jwt.Token `json:"accesstoken,omitempty"`
	Refreshtoken *jwt.Token `json:"refreshtoken,omitempty"`
}

// JwtMiddleware function
var JwtMiddleware = middleware.JWTWithConfig(middleware.JWTConfig{
	SigningKey: []byte(env.MustGet("SECRET")),
})

// GenerateTokenPair creates a token pair which contains access token and
// refresh token with exp limit for both access and refresh token.
func GenerateTokenPair(uuid string, atLimit time.Duration, rtLimit time.Duration) (*TokenPair, error) {
	secret := env.MustGet("SECRET")

	accessToken := jwt.New(jwt.SigningMethodHS256)

	atclaims := accessToken.Claims.(jwt.MapClaims)
	atclaims["uuid"] = uuid
	atclaims["exp"] = time.Now().Add(atLimit).Unix()

	// Generate encoded accessToken.
	atstr, err := accessToken.SignedString([]byte(secret))
	if err != nil {
		return nil, err
	}

	refreshToken := jwt.New(jwt.SigningMethodHS256)

	rtClaims := refreshToken.Claims.(jwt.MapClaims)
	rtClaims["uuid"] = uuid
	rtClaims["exp"] = time.Now().Add(rtLimit).Unix()

	// Generate encoded refreshToken.
	rtstr, err := refreshToken.SignedString([]byte(secret))
	if err != nil {
		return nil, err
	}

	return &TokenPair{Accesstoken: atstr, Refreshtoken: rtstr}, nil
}

// TokenPairStringToJWT converts a string access and resfresh token from a TokenPair
// to a TokenPairJWT.
func TokenPairStringToJWT(token *TokenPair) (res *TokenPairJWT, e error) {
	secret := env.MustGet("SECRET")

	atk, err := jwt.ParseWithClaims(
		token.Accesstoken,
		jwt.MapClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(secret), nil
		},
	)
	if err != nil && err.(*jwt.ValidationError).Errors != jwt.ValidationErrorExpired {
		return nil, err
	}

	rtk, err := jwt.ParseWithClaims(
		token.Refreshtoken,
		jwt.MapClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(secret), nil
		},
	)
	if err != nil && err.(*jwt.ValidationError).Errors == jwt.ValidationErrorExpired {
		return nil, ErrRefreshTokenExpired
	} else if err != nil {
		return nil, err
	}

	res = &TokenPairJWT{Accesstoken: atk, Refreshtoken: rtk}
	return res, err
}
