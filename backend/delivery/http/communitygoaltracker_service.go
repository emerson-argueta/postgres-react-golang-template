package http

import (
	"bytes"
	"emersonargueta/m/v1/communitygoaltracker"
	"emersonargueta/m/v1/communitygoaltracker/achiever"
	"emersonargueta/m/v1/communitygoaltracker/goal"
	"emersonargueta/m/v1/delivery/middleware"
	"encoding/json"
	"net/http"
)

var _ communitygoaltracker.Processes = &communitygoaltrackerservice{}

// communitygoaltrackerservice represents an HTTP implementation of communitygoaltracker.Service.
type communitygoaltrackerservice struct {
	client *Client
}

// Register is an http implementation of the communitygoaltracker process.
func (cgt *communitygoaltrackerservice) Register(a *achiever.Achiever) (res *achiever.Achiever, e error) {
	u := cgt.client.URL
	u.Path = RoutePrefix + AchieverURL

	// Encode request body.
	reqBody, e := json.Marshal(achieverRequest{Achiever: a})
	if e != nil {
		return nil, e
	}

	// Execute request.
	url := u.String()
	resp, e := http.Post(url, "application/json", bytes.NewReader(reqBody))
	if e != nil {
		return nil, e
	}
	defer resp.Body.Close()

	// Decode response into JSON.
	var respBody achieverResponse
	if e = json.NewDecoder(resp.Body).Decode(&respBody); e != nil {
		return nil, e
	} else if respBody.Error != "" {
		return nil, communitygoaltracker.Error(respBody.Error)
	}
	res = respBody.Achiever

	return res, e
}

// Login is an http implementation of the communitygoaltracker process.
func (cgt *communitygoaltrackerservice) Login(email string, password string) (res *achiever.Achiever, e error) {
	u := cgt.client.URL
	u.Path = RoutePrefix + AchieverLoginURL

	// Encode request body.
	reqBody, e := json.Marshal(achieverRequest{Achiever: &achiever.Achiever{Email: &email, Password: &password}})
	if e != nil {
		return nil, e
	}

	// Execute request.
	url := u.String()
	resp, e := http.Post(url, "application/json", bytes.NewReader(reqBody))
	if e != nil {
		return nil, e
	}
	defer resp.Body.Close()

	// Decode response into JSON.
	var respBody achieverResponse
	if e = json.NewDecoder(resp.Body).Decode(&respBody); e != nil {
		return nil, e
	} else if respBody.Error != "" {
		return nil, communitygoaltracker.Error(respBody.Error)
	}
	res = respBody.Achiever

	return res, e
}

// AuthorizedUpdateAchiever is an http implementation of the communitygoaltracker process.
func (cgt *communitygoaltrackerservice) AuthorizedUpdateAchiever(a *achiever.Achiever, token *middleware.TokenPair) (e error) {
	u := cgt.client.URL
	u.Path = RoutePrefix + AchieverLoginURL

	// Encode request body.
	reqBody, e := json.Marshal(achieverRequest{Achiever: a})
	if e != nil {
		return e
	}

	// Create a new request using http
	req, e := http.NewRequest("PATCH", u.String(), bytes.NewReader(reqBody))
	if e != nil {
		return e
	}
	// Create a Bearer string by appending string access token
	var bearer = "Bearer " + token.Accesstoken
	// add authorization header to the req
	req.Header.Add("Authorization", bearer)
	req.Header.Add("Content-Type", "application/json")

	resp, e := http.DefaultClient.Do(req)
	if e != nil {
		return e
	}

	defer resp.Body.Close()

	// Decode response into JSON.
	var respBody achieverResponse
	if e = json.NewDecoder(resp.Body).Decode(&respBody); e != nil {
		return e
	} else if respBody.Error != "" {
		return communitygoaltracker.Error(respBody.Error)
	}
	a = respBody.Achiever

	return e
}

// UpdateAchiever not used
func (cgt *communitygoaltrackerservice) UpdateAchiever(a *achiever.Achiever) (e error) {
	return e
}

// AuthorizedUnRegister is an http implementation of the communitygoaltracker process.
func (cgt *communitygoaltrackerservice) AuthorizedUnRegister(a *achiever.Achiever, token *middleware.TokenPair) (e error) {
	return e
}

// UnRegister is not used
func (cgt *communitygoaltrackerservice) UnRegister(a *achiever.Achiever) (e error) {
	return e
}

// AuthorizedCreateGoal is an http implementation of the communitygoaltracker process.
func (cgt *communitygoaltrackerservice) AuthorizedCreateGoal(g *goal.Goal) (res *goal.Goal, e error) {
	return res, e
}

// CreateGoal is not used.
func (cgt *communitygoaltrackerservice) CreateGoal(g *goal.Goal) (res *goal.Goal, e error) {
	return res, e
}

// AuthorizedUpdateGoalProgress is an http implementation of the communitygoaltracker process.
func (cgt *communitygoaltrackerservice) AuthorizedUpdateGoalProgress(achieverUUID string, goalID int64, progress int) (res *goal.Goal, e error) {
	return res, e
}

// UpdateGoalProgress is not used
func (cgt *communitygoaltrackerservice) UpdateGoalProgress(achieverUUID string, goalID int64, progress int) (res *goal.Goal, e error) {
	return res, e
}

// AuthorizedAbandonGoal is an http implementation of the communitygoaltracker process.
func (cgt *communitygoaltrackerservice) AuthorizedAbandonGoal(achieverUUID string, goalID int64) (e error) {
	return e
}

// AbandonGoal is not used
func (cgt *communitygoaltrackerservice) AbandonGoal(achieverUUID string, goalID int64) (e error) {
	return e
}

// AuthorizedDeleteGoal is an http implementation of the communitygoaltracker process.
func (cgt *communitygoaltrackerservice) AuthorizedDeleteGoal(achieverUUID string, goalID int64) (e error) {
	return e
}

// DeleteGoal is not used
func (cgt *communitygoaltrackerservice) DeleteGoal(achieverUUID string, goalID int64) (e error) {
	return e
}
