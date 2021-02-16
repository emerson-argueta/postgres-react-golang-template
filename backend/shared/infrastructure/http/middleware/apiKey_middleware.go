package middleware

import (
	"emersonargueta/m/v1/modules/identity/domain/user"
	"emersonargueta/m/v1/modules/identity/infrastructure/persistence"
	"emersonargueta/m/v1/shared/infrastructure"
	"emersonargueta/m/v1/shared/infrastructure/http/authorization"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

// APIKeyMiddleware for api key validation
var APIKeyMiddleware = middleware.KeyAuthWithConfig(
	middleware.KeyAuthConfig{
		KeyLookup: "header:api_key",
		Validator: apiKeyValidator,
	},
)

var apiKeyValidator = func(key string, c echo.Context) (bool, error) {
	jwtService := authorization.NewJWTService(infrastructure.GlobalConfig)

	id, err := jwtService.VerifyTokenAndExtractIDClaim(key)
	if err != nil && err == authorization.ErrAuthorizationFailed {
		return false, nil
	}
	if err != nil {
		return false, err
	}

	_, err = persistence.IdentityRepos.User.RetrieveUserByID(id)
	if err != nil && err == user.ErrUserNotFound {
		return false, nil
	}
	if err != nil {
		return false, err
	}

	return true, nil
}
