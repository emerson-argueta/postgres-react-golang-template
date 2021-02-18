package authorization

import "time"

const (
	// AccestokenLimit represents the expiration time for an access token
	AccestokenLimit = time.Minute * 15
	// RefreshtokenLimit represents the expiration time for an refresh token
	RefreshtokenLimit = time.Hour * 24
)

// TokenPair with access token for jwt authentication and refresh token to
// refresh access.
type TokenPair struct {
	Accesstoken  string `json:"accesstoken,omitempty"`
	Refreshtoken string `json:"refreshtoken,omitempty"`
}
