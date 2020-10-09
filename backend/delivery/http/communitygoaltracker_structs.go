package http

import (
	"emersonargueta/m/v1/communitygoaltracker/achiever"
	"emersonargueta/m/v1/delivery/middleware"
)

type achieverRequest struct {
	Achiever *achiever.Achiever    `json:"achiever,omitempty"`
	Token    *middleware.TokenPair `json:"token,omitempty"`
}

type achieverResponse struct {
	Achiever *achiever.Achiever    `json:"achiever,omitempty"`
	Token    *middleware.TokenPair `json:"token,omitempty"`
	Error    string                `json:"error,omitempty"`
}
