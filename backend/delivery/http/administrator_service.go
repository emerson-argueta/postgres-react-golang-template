package http

import (
	"bytes"
	"encoding/json"

	"emersonargueta/m/v1/authorization"
	"emersonargueta/m/v1/delivery/middleware"
	"emersonargueta/m/v1/domain/administrator"
	"emersonargueta/m/v1/domain/church"
	"emersonargueta/m/v1/domain/donator"
	"emersonargueta/m/v1/domain/transaction"

	"net/http"
)

var _ AdministratorActions = &Administrator{}

// AdministratorActions has additional actions for http service
type AdministratorActions interface {
	administrator.UnauthorizedManagementActions
	authorization.Actions
}

// Administrator represents an HTTP implementation of administrator.Service.
type Administrator struct {
	client *Client
}

// Register a new administrator.
func (s *Administrator) Register(a *administrator.Administrator) error {
	u := s.client.URL
	u.Path = "/api/administrator"

	// Encode request body.
	reqBody, err := json.Marshal(administratorRequest{Administrator: a})
	if err != nil {
		return err
	}

	// Execute request.
	url := u.String()
	resp, err := http.Post(url, "application/json", bytes.NewReader(reqBody))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Decode response into JSON.
	var respBody administratorResponse
	if err := json.NewDecoder(resp.Body).Decode(&respBody); err != nil {
		return err
	} else if respBody.Error != "" {
		return administrator.Error(respBody.Error)
	}

	// Copy returned administrator.
	*a = *respBody.Administrator

	return nil
}

// Login an administrator to give access them accesss to the church_fund_management service
func (s *Administrator) Login(a *administrator.Administrator) (*administrator.Administrator, error) {
	u := s.client.URL
	u.Path = "/api/administrator/login"

	// Encode request body.
	reqBody, err := json.Marshal(administratorRequest{Administrator: a})
	if err != nil {
		return nil, err
	}

	// Execute request.
	url := u.String()
	resp, err := http.Post(url, "application/json", bytes.NewReader(reqBody))
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	// Decode response into JSON.
	var respBody administratorResponse
	if err := json.NewDecoder(resp.Body).Decode(&respBody); err != nil {
		return nil, err
	} else if respBody.Error != "" {
		return nil, administrator.Error(respBody.Error)
	}
	return respBody.Administrator, nil

}

// Authorize an administrator's by refreshing a valid jwt token that has expired.
func (s *Administrator) Authorize(t *middleware.TokenPair) (*middleware.TokenPair, error) {
	u := s.client.URL
	u.Path = "/api/administrator/token"

	// Encode request body.
	reqBody, err := json.Marshal(authorizationAdministratorRequest{Token: t})
	if err != nil {
		return nil, err
	}

	// Execute request.
	url := u.String()
	resp, err := http.Post(url, "application/json", bytes.NewReader(reqBody))
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	// Decode response into JSON.
	var respBody authorizationAdministratorResponse
	if err := json.NewDecoder(resp.Body).Decode(&respBody); err != nil {
		return nil, err
	} else if respBody.Error != "" {
		return nil, administrator.Error(respBody.Error)
	}
	return respBody.Token, nil
}

// Read an administrator given the administrators token. The UUID is extracted
// from the token to find the administrator by UUID.
func (s *Administrator) Read(token *middleware.TokenPair) (*administrator.Administrator, error) {
	u := s.client.URL
	u.Path = "/api/administrator"

	// Create a new request using http
	url := u.String()
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	// Create a Bearer string by appending string access token
	var bearer = "Bearer " + token.Accesstoken
	// add authorization header to the req
	req.Header.Add("Authorization", bearer)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	// Decode response into JSON.
	var respBody administratorResponse
	if err := json.NewDecoder(resp.Body).Decode(&respBody); err != nil {
		return nil, err
	} else if respBody.Error != "" {
		return nil, administrator.Error(respBody.Error)
	}
	return respBody.Administrator, nil

}

// AddChurch adds a church for an administrator with a valid token.
func (s *Administrator) AddChurch(token *middleware.TokenPair, a *administrator.Administrator, c *church.Church) (res *church.Church, e error) {
	u := s.client.URL
	u.Path = "/api/administrator/church/add"

	// Encode request body.
	reqBody, err := json.Marshal(churchRequest{Church: c})
	if err != nil {
		return nil, err
	}

	// Create a new request using http
	url := u.String()
	req, err := http.NewRequest("POST", url, bytes.NewReader(reqBody))
	if err != nil {
		return nil, err
	}
	// Create a Bearer string by appending string access token
	var bearer = "Bearer " + token.Accesstoken
	// add authorization header to the req
	req.Header.Add("Authorization", bearer)
	req.Header.Add("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	// Decode response into JSON.
	var respBody churchResponse
	if err := json.NewDecoder(resp.Body).Decode(&respBody); err != nil {
		return nil, err
	} else if respBody.Error != "" {
		return nil, church.Error(respBody.Error)
	}
	return respBody.Church, nil
}

// CreateChurch adds a church for an administrator with a valid token.
func (s *Administrator) CreateChurch(token *middleware.TokenPair, a *administrator.Administrator, c *church.Church) (res *church.Church, e error) {
	var resp *http.Response
	u := s.client.URL

	u.Path = "/api/administrator/church"

	// Encode request body.
	reqBody, err := json.Marshal(churchRequest{Church: c})
	if err != nil {
		return nil, err
	}

	// Create a new request using http
	req, err := http.NewRequest("POST", u.String(), bytes.NewReader(reqBody))
	if err != nil {
		return nil, err
	}
	// Create a Bearer string by appending string access token
	var bearer = "Bearer " + token.Accesstoken
	// add authorization header to the req
	req.Header.Add("Authorization", bearer)
	req.Header.Add("Content-Type", "application/json")

	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	// Decode response into JSON.
	var respBody churchResponse
	if err := json.NewDecoder(resp.Body).Decode(&respBody); err != nil {
		return nil, err
	} else if respBody.Error != "" {
		return nil, church.Error(respBody.Error)
	}
	return respBody.Church, nil
}

// CreateDonator adds a church for an administrator with a valid token.
func (s *Administrator) CreateDonator(token *middleware.TokenPair, c *church.Church, d *donator.Donator) (res *donator.Donator, e error) {
	var resp *http.Response
	u := s.client.URL

	u.Path = "/api/administrator/donator"

	// Encode request body.
	reqBody, err := json.Marshal(donatorRequest{Church: c, Donator: d})
	if err != nil {
		return nil, err
	}

	// Create a new request using http
	req, err := http.NewRequest("POST", u.String(), bytes.NewReader(reqBody))
	if err != nil {
		return nil, err
	}
	// Create a Bearer string by appending string access token
	var bearer = "Bearer " + token.Accesstoken
	// add authorization header to the req
	req.Header.Add("Authorization", bearer)
	req.Header.Add("Content-Type", "application/json")

	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	// Decode response into JSON.
	var respBody donatorResponse
	if err := json.NewDecoder(resp.Body).Decode(&respBody); err != nil {
		return nil, err
	} else if respBody.Error != "" {
		return nil, church.Error(respBody.Error)
	}
	return respBody.Donator, nil
}

// CreateDonation adds a church for an administrator with a valid token.
func (s *Administrator) CreateDonation(token *middleware.TokenPair, c *church.Church, d *donator.Donator, donation *transaction.Donation) (res *transaction.Donation, e error) {
	var resp *http.Response
	u := s.client.URL

	u.Path = "/api/administrator/donation"

	// Encode request body.
	reqBody, err := json.Marshal(donationRequest{Donation: donation})
	if err != nil {
		return nil, err
	}

	// Create a new request using http
	req, err := http.NewRequest("POST", u.String(), bytes.NewReader(reqBody))
	if err != nil {
		return nil, err
	}
	// Create a Bearer string by appending string access token
	var bearer = "Bearer " + token.Accesstoken
	// add authorization header to the req
	req.Header.Add("Authorization", bearer)
	req.Header.Add("Content-Type", "application/json")

	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	// Decode response into JSON.
	var respBody donationResponse
	if err := json.NewDecoder(resp.Body).Decode(&respBody); err != nil {
		return nil, err
	} else if respBody.Error != "" {
		return nil, church.Error(respBody.Error)
	}
	return respBody.Donation, nil
}

// // Update updates an admin's information.
// // Searches for admin by email or id if byEmail flag is false
// func (s *AdministratorService) Update(a *administrator.Administrator, byEmail bool) error {
// 	u := *s.URL
// 	u.Path = "/api/admin"

// 	// Encode request body.
// 	reqBody, err := json.Marshal(patchAdministratorRequest{Administrator: a})
// 	if err != nil {
// 		return err
// 	}

// 	// Create request.
// 	req, err := http.NewRequest("PATCH", u.String(), bytes.NewReader(reqBody))
// 	if err != nil {
// 		return err
// 	}

// 	// Create a Bearer string by appending string access token
// 	var bearer = "Bearer " + a.Token.Accesstoken
// 	// add authorization header to the req
// 	req.Header.Add("Authorization", bearer)
// 	req.Header.Add("Content-Type", "application/json")

// 	// Execute request.
// 	resp, err := http.DefaultClient.Do(req)
// 	if err != nil {
// 		return err
// 	}
// 	defer resp.Body.Close()

// 	// Decode response into JSON.
// 	var respBody postAdministratorResponse
// 	if err := json.NewDecoder(resp.Body).Decode(&respBody); err != nil {
// 		return err
// 	} else if respBody.Err != "" {
// 		return administrator.Error(respBody.Err)
// 	}

// 	return nil
// }

// // Delete scenarios:
// //  1. If the administrator to be deleted has non-restricted access and is the only administrator of the church,
// // he will be asked whether he wants to delete the church information or keep church information on the system.
// //      1a.If administrator decides to delete the church, all information will be deleted by trigerring the process to delete a church.
// //		1b.If administrator decides to keep the church, another administrator account should be created and church should be transferred to the new administrator.
// //	2.If a deleted administrator has non-restricted access but the church has multiple administrators:
// //      2a.Before deleting his account, administrator needs to transfer the church to one of the other administrators with restricted or non-restricted access.
// func (s *AdministratorService) Delete(a *administrator.Administrator, byEmail bool) error {
// 	u := *s.URL
// 	u.Path = "/api/admin"

// 	// Create request.
// 	req, err := http.NewRequest("DELETE", u.String(), nil)
// 	// Create a Bearer string by appending string access token
// 	var bearer = "Bearer " + a.Token.Accesstoken
// 	// add authorization header to the req
// 	req.Header.Add("Authorization", bearer)
// 	// Execute request.
// 	resp, err := http.DefaultClient.Do(req)
// 	if err != nil {
// 		return err
// 	}
// 	defer resp.Body.Close()

// 	// Decode response into JSON.
// 	var respBody deleteAdministratorResponse
// 	if err := json.NewDecoder(resp.Body).Decode(&respBody); err != nil {
// 		return err
// 	} else if respBody.Err != "" {
// 		return administrator.Error(respBody.Err)
// 	}
// 	return nil
// }

// // Login authorizes a user with correct email and password
// func (s *AdministratorService) Login(a *administrator.Administrator) error {

// 	_, err := s.Read(a, true)
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }
