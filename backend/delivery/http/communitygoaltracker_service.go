package http

import (
	"bytes"
	"emersonargueta/m/v1/communitygoaltracker"
	"emersonargueta/m/v1/communitygoaltracker/achiever"
	"emersonargueta/m/v1/communitygoaltracker/goal"
	"encoding/json"
	"net/http"
)

var _ communitygoaltracker.Service = &Communitygoaltracker{}

// Communitygoaltracker represents an HTTP implementation of communitygoaltracker.Service.
type Communitygoaltracker struct {
	client *Client
}

// Register is an http implementation of the communitygoaltracker process.
func (cgt *Communitygoaltracker) Register(a *achiever.Achiever) (res *achiever.Achiever, e error) {
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
func (cgt *Communitygoaltracker) Login(email string, password string) (res *achiever.Achiever, e error) {
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
func (cgt *Communitygoaltracker) UpdateAchiever(a *achiever.Achiever) (e error) {
	return e
}

// UnRegister is an http implementation of the communitygoaltracker process.
func (cgt *Communitygoaltracker) UnRegister(a *achiever.Achiever) (e error) {
	return e
}

// CreateGoal is an http implementation of the communitygoaltracker process.
func (cgt *Communitygoaltracker) CreateGoal(g *goal.Goal) (res *goal.Goal, e error) {
	return res, e
}

// UpdateGoalProgress is an http implementation of the communitygoaltracker process.
func (cgt *Communitygoaltracker) UpdateGoalProgress(achieverUUID string, goalID int64, progress int) (e error) {
	return e
}

// AbandonGoal is an http implementation of the communitygoaltracker process.
func (cgt *Communitygoaltracker) AbandonGoal(achieverUUID string, goalID int64) (e error) {
	return e
}

// DeleteGoal is an http implementation of the communitygoaltracker process.
func (cgt *Communitygoaltracker) DeleteGoal(achieverUUID string, goalID int64) (e error) {
	return e
}
