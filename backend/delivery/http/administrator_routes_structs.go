package http

import (
	"emersonargueta/m/v1/delivery/middleware"
	"emersonargueta/m/v1/domain/administrator"
	"emersonargueta/m/v1/domain/church"
	"emersonargueta/m/v1/domain/donator"
	"emersonargueta/m/v1/domain/transaction"
)

type administratorRequest struct {
	Administrator *administrator.Administrator `json:"administrator,omitempty"`
	Token         *middleware.TokenPair        `json:"token,omitempty"`
}

type administratorResponse struct {
	Administrator *administrator.Administrator `json:"administrator,omitempty"`
	Token         *middleware.TokenPair        `json:"token,omitempty"`
	Error         string                       `json:"error,omitempty"`
}

type authorizationAdministratorRequest struct {
	Token *middleware.TokenPair `json:"token,omitempty"`
}
type authorizationAdministratorResponse struct {
	Token *middleware.TokenPair `json:"token,omitempty"`
	Error string                `json:"error,omitempty"`
}

type churchRequest struct {
	Church *church.Church `json:"church,omitempty"`
}
type churchResponse struct {
	Church   *church.Church `json:"church,omitempty"`
	ChurchID *int64         `json:"churchid,omitempty"`
	Error    string         `json:"error,omitempty"`
}

type churchAdministratorRequest struct {
	Administrator *church.Administrator `json:"administrator,omitempty"`
	Token         *middleware.TokenPair `json:"token,omitempty"`
}
type churchAdministratorResponse struct {
	AdministratorUUID *string               `json:"administratoruuid,omitempty"`
	Administrator     *church.Administrator `json:"administrator,omitempty"`
	Error             string                `json:"error,omitempty"`
}

type donatorRequest struct {
	Church  *church.Church   `json:"church,omitempty"`
	Donator *donator.Donator `json:"donator,omitempty"`
}

type donatorResponse struct {
	Donator   *donator.Donator `json:"donator,omitempty"`
	DonatorID *int64           `json:"donatorid,omitempty"`
	Error     string           `json:"error,omitempty"`
}

type donationRequest struct {
	Donation *transaction.Donation `json:"donation,omitempty"`
}

type donationResponse struct {
	Donation *transaction.Donation `json:"donation,omitempty"`
	Error    string                `json:"error,omitempty"`
}
