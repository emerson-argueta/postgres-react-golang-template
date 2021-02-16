package authorization

import (
	"emersonargueta/m/v1/shared/infrastructure"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/dgrijalva/jwt-go"
)

// JwtService to get and verify jwt tokens
type JwtService interface {
	IssueNewToken(id string) (string, error)
	VerifyTokenAndExtractIDClaim(token string) (string, error)
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

func (s *jwtService) IssueNewToken(id string) (string, error) {
	token := jwt.NewWithClaims(
		jwt.SigningMethodRS256,
		jwt.MapClaims{
			"id": id,
		},
	)
	path, err := os.Getwd()
	if err != nil {
		log.Println(err)
	}
	fmt.Println(path) // for example /home/user
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
func (s *jwtService) VerifyTokenAndExtractIDClaim(token string) (string, error) {
	verifyBytes, err := ioutil.ReadFile(s.config.Authorization.PublicKeyPath)
	if err != nil {
		return "", err
	}

	verifyKey, err := jwt.ParseRSAPublicKeyFromPEM(verifyBytes)
	if err != nil {
		return "", err
	}

	verifiedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		// since we only use the one private key to sign the tokens,
		// we also only use its public counter part to verify
		return verifyKey, nil
	})
	if err != nil {
		return "", ErrAuthorizationFailed
	}

	id, ok := verifiedToken.Claims.(jwt.MapClaims)["id"].(string)
	if !ok {
		return "", ErrAuthorizationFailed
	}

	return id, nil
}
