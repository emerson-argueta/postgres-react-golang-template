package http

import (
	"bytes"
	"emersonargueta/m/v1/authorization"
	"emersonargueta/m/v1/communitygoaltracker"
	"emersonargueta/m/v1/communitygoaltracker/achiever"
	"emersonargueta/m/v1/communitygoaltracker/goal"
	"encoding/json"
	"net/http"
)

var _ communitygoaltracker.Processes = &communitygoaltrackerjwtservice{}

type communitygoaltrackerjwtservice struct {
	client *Client
}

// Register is an http implementation of the communitygoaltracker process.
func (cgt *communitygoaltrackerjwtservice) Register(a *achiever.Achiever) (res *achiever.Achiever, e error) {
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
func (cgt *communitygoaltrackerjwtservice) Login(email string, password string) (res *achiever.Achiever, e error) {
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

// UpdateAchiever is an http implementation of the communitygoaltracker process.
func (cgt *communitygoaltrackerjwtservice) UpdateAchiever(a *achiever.Achiever) (e error) {
	if cgt.client.Token == nil {
		return authorization.ErrAuthorizationKeyNotCreated
	}

	u := cgt.client.URL
	u.Path = RoutePrefix + AchieverLoginURL

	// Encode request body.
	reqBody, e := json.Marshal(achieverRequest{Achiever: a})
	if e != nil {
		return e
	}

	// Create a new request using http
	req, e := http.NewRequest(http.MethodPatch, u.String(), bytes.NewReader(reqBody))
	if e != nil {
		return e
	}
	// Create a Bearer string by appending string access token
	var bearer = "Bearer " + cgt.client.Token.Accesstoken
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

// UnRegister is an http implementation of the communitygoaltracker process.
func (cgt *communitygoaltrackerjwtservice) UnRegister(a *achiever.Achiever) (e error) {
	if cgt.client.Token == nil {
		return authorization.ErrAuthorizationKeyNotCreated
	}

	u := cgt.client.URL
	u.Path = RoutePrefix + AchieverURL

	// Encode request body.
	reqBody, e := json.Marshal(achieverRequest{Achiever: a})
	if e != nil {
		return e
	}

	// Create a new request using http
	req, e := http.NewRequest(http.MethodDelete, u.String(), bytes.NewReader(reqBody))
	if e != nil {
		return e
	}
	// Create a Bearer string by appending string access token
	var bearer = "Bearer " + cgt.client.Token.Accesstoken
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

// CreateGoal is an http implementation of the communitygoaltracker process.
func (cgt *communitygoaltrackerjwtservice) CreateGoal(g *goal.Goal) (res *goal.Goal, e error) {
	if cgt.client.Token == nil {
		return nil, authorization.ErrAuthorizationKeyNotCreated
	}

	u := cgt.client.URL
	u.Path = RoutePrefix + GoalURL

	// Encode request body.
	gRequest := goalToRequest(g)
	reqBody, e := json.Marshal(gRequest)
	if e != nil {
		return nil, e
	}

	// Create a new request using http
	req, e := http.NewRequest(http.MethodPost, u.String(), bytes.NewReader(reqBody))
	if e != nil {
		return nil, e
	}
	// Create a Bearer string by appending string access token
	var bearer = "Bearer " + cgt.client.Token.Accesstoken
	// add authorization header to the req
	req.Header.Add("Authorization", bearer)
	req.Header.Add("Content-Type", "application/json")

	resp, e := http.DefaultClient.Do(req)
	if e != nil {
		return nil, e
	}

	defer resp.Body.Close()

	// Decode response into JSON.
	var respBody goalResponse
	if e = json.NewDecoder(resp.Body).Decode(&respBody); e != nil {
		return nil, e
	} else if respBody.Error != "" {
		return nil, communitygoaltracker.Error(respBody.Error)
	}
	res = responseToGoal(&respBody)

	return res, e
}

// UpdateGoalProgress is an http implementation of the communitygoaltracker process.
func (cgt *communitygoaltrackerjwtservice) UpdateGoalProgress(achieverUUID string, goalID int64, progress int) (res *goal.Goal, e error) {
	if cgt.client.Token == nil {
		return nil, authorization.ErrAuthorizationKeyNotCreated
	}

	u := cgt.client.URL
	u.Path = RoutePrefix + GoalURL

	// Encode request body.
	reqBody, e := json.Marshal(goalRequest{ID: &goalID, Progress: &progress})
	if e != nil {
		return nil, e
	}

	// Create a new request using http
	req, e := http.NewRequest(http.MethodPatch, u.String(), bytes.NewReader(reqBody))
	if e != nil {
		return nil, e
	}
	// Create a Bearer string by appending string access token
	var bearer = "Bearer " + cgt.client.Token.Accesstoken
	// add authorization header to the req
	req.Header.Add("Authorization", bearer)
	req.Header.Add("Content-Type", "application/json")

	resp, e := http.DefaultClient.Do(req)
	if e != nil {
		return nil, e
	}

	defer resp.Body.Close()

	// Decode response into JSON.
	var respBody goalResponse
	if e = json.NewDecoder(resp.Body).Decode(&respBody); e != nil {
		return nil, e
	} else if respBody.Error != "" {
		return nil, communitygoaltracker.Error(respBody.Error)
	}
	res = responseToGoal(&respBody)

	return res, e
}

// AuthorizedAbandonGoal is an http implementation of the communitygoaltracker process.
func (cgt *communitygoaltrackerjwtservice) AuthorizedAbandonGoal(achieverUUID string, goalID int64) (e error) {
	return e
}

// AbandonGoal is not used
func (cgt *communitygoaltrackerjwtservice) AbandonGoal(achieverUUID string, goalID int64) (e error) {
	if cgt.client.Token == nil {
		return authorization.ErrAuthorizationKeyNotCreated
	}

	u := cgt.client.URL
	u.Path = RoutePrefix + GoalAbandonURL

	// Encode request body.
	reqBody, e := json.Marshal(goalRequest{ID: &goalID})
	if e != nil {
		return e
	}

	// Create a new request using http
	req, e := http.NewRequest(http.MethodDelete, u.String(), bytes.NewReader(reqBody))
	if e != nil {
		return e
	}
	// Create a Bearer string by appending string access token
	var bearer = "Bearer " + cgt.client.Token.Accesstoken
	// add authorization header to the req
	req.Header.Add("Authorization", bearer)
	req.Header.Add("Content-Type", "application/json")

	resp, e := http.DefaultClient.Do(req)
	if e != nil {
		return e
	}

	defer resp.Body.Close()

	// Decode response into JSON.
	var respBody goalResponse
	if e = json.NewDecoder(resp.Body).Decode(&respBody); e != nil {
		return e
	} else if respBody.Error != "" {
		return communitygoaltracker.Error(respBody.Error)
	}

	return e
}

// DeleteGoal is an http implementation of the communitygoaltracker process.
func (cgt *communitygoaltrackerjwtservice) DeleteGoal(achieverUUID string, goalID int64) (e error) {
	if cgt.client.Token == nil {
		return authorization.ErrAuthorizationKeyNotCreated
	}

	u := cgt.client.URL
	u.Path = RoutePrefix + GoalURL

	// Encode request body.
	reqBody, e := json.Marshal(goalRequest{ID: &goalID})
	if e != nil {
		return e
	}

	// Create a new request using http
	req, e := http.NewRequest(http.MethodDelete, u.String(), bytes.NewReader(reqBody))
	if e != nil {
		return e
	}
	// Create a Bearer string by appending string access token
	var bearer = "Bearer " + cgt.client.Token.Accesstoken
	// add authorization header to the req
	req.Header.Add("Authorization", bearer)
	req.Header.Add("Content-Type", "application/json")

	resp, e := http.DefaultClient.Do(req)
	if e != nil {
		return e
	}

	defer resp.Body.Close()

	// Decode response into JSON.
	var respBody goalResponse
	if e = json.NewDecoder(resp.Body).Decode(&respBody); e != nil {
		return e
	} else if respBody.Error != "" {
		return communitygoaltracker.Error(respBody.Error)
	}

	return e
}
