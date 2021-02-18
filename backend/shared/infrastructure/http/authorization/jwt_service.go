package authorization

import (
	"emersonargueta/m/v1/shared/infrastructure"
	"io/ioutil"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// JwtService to get and verify jwt tokens
type JwtService interface {
	IssueNewToken(id string, expiresAt *time.Duration) (string, error)
	VerifyTokenAndExtractIDClaim(token string) (string, error)
	IssueTokenPair(
		uuid string,
		accessTokenLimit *time.Duration,
		refreshTokenLimit *time.Duration,
	) (*TokenPair, error)
}

type jwtService struct {
	config *infrastructure.Config
}

// NewJWTService jwt service
func NewJWTService(config *infrastructure.Config) JwtService {
	s := &jwtService{
		config: config,
	}

	return s
}

func (s *jwtService) IssueNewToken(id string, expiresAt *time.Duration) (string, error) {
	claims := jwt.MapClaims{
		"id": id,
	}
	if expiresAt != nil {
		claims["exp"] = time.Now().Add(*expiresAt).Unix()
	}

	token := jwt.NewWithClaims(
		jwt.SigningMethodRS256,
		claims,
	)

	signBytes, err := ioutil.ReadFile(s.config.Authorization.PrivateKeyPath)
	if err != nil {
		return "", err
	}

	signKey, err := jwt.ParseRSAPrivateKeyFromPEM(signBytes)
	if err != nil {
		return "", err
	}

	return token.SignedString(signKey)
}

// IssueTokenPair where the token pair contains an access token and a
// refresh token with exp limit for both.
// Defualt accessToken limit is AccessTokenLimit.
// Defualt refreshToken limit is RefreshTokenLimit.
func (s *jwtService) IssueTokenPair(
	uuid string,
	accessTokenLimit *time.Duration,
	refreshTokenLimit *time.Duration,
) (*TokenPair, error) {
	tokenPair := &TokenPair{}
	if accessTokenLimit == nil {
		defaultATKLimit := AccestokenLimit
		accessTokenLimit = &defaultATKLimit
	}
	if refreshTokenLimit == nil {
		defaultRTKLimit := RefreshtokenLimit
		refreshTokenLimit = &defaultRTKLimit
	}

	atk, err := s.IssueNewToken(uuid, accessTokenLimit)
	if err != nil {
		return nil, err
	}
	tokenPair.Accesstoken = atk

	rtk, err := s.IssueNewToken(uuid, refreshTokenLimit)
	if err != nil {
		return nil, err
	}
	tokenPair.Refreshtoken = rtk

	return tokenPair, nil
}
func (s *jwtService) VerifyTokenAndExtractIDClaim(token string) (string, error) {
	verifiedToken, err := s.verifyToken(token)
	if err != nil {
		return "", err
	}

	id, err := s.extractIDFromToken(verifiedToken)
	if err != nil {
		return "", err
	}

	err = s.verifyExpiryTime(verifiedToken)
	if err != nil {
		return id, err
	}

	return id, nil
}

func (s *jwtService) verifyToken(token string) (*jwt.Token, error) {
	verifyBytes, err := ioutil.ReadFile(s.config.Authorization.PublicKeyPath)
	if err != nil {
		return nil, err
	}

	verifyKey, err := jwt.ParseRSAPublicKeyFromPEM(verifyBytes)
	if err != nil {
		return nil, err
	}

	verifiedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		// since we only use the one private key to sign the tokens,
		// we also only use its public counter part to verify
		return verifyKey, nil
	})
	if err != nil {
		return nil, ErrAuthorizationFailed
	}
	return verifiedToken, nil
}

func (s *jwtService) extractIDFromToken(verifiedToken *jwt.Token) (string, error) {
	id, hasID := verifiedToken.Claims.(jwt.MapClaims)["id"].(string)
	if !hasID {
		return "", ErrAuthorizationFailed
	}
	return id, nil
}

func (s *jwtService) verifyExpiryTime(verifiedToken *jwt.Token) error {
	mapClaims := verifiedToken.Claims.(jwt.MapClaims)
	_, hasExpiryTime := mapClaims["exp"].(int64)
	if !hasExpiryTime {
		return nil
	}

	verifiedTime := mapClaims.VerifyExpiresAt(time.Now().Unix(), false)
	if !verifiedTime {
		return ErrAuthorizationFailed
	}

	return nil
}
