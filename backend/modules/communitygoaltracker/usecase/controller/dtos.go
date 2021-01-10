package controller

import "emersonargueta/m/v1/modules/communitygoaltracker/dto"

type achieverRequest struct {
	Achiever *dto.AchieverDTO `json:"achiever,omitempty"`
}
type achieverResponse struct {
	Achiever *dto.AchieverDTO `json:"achiever,omitempty"`
}

type retrieveAchieversResponse struct {
	Achievers *achieversDTO `json:"achievers,omitempty"`
}
type achieversDTO map[string]*achieverDTO
type achieverDTO struct {
	Firstname *string `json:"firstname,omitempty"`
	Lastname  *string `json:"lastname,omitempty"`
	UserID    *string `json:"userid,omitempty"`
}
type updateAchieverResponse struct {
	Message string `json:"message,omitempty"`
}

type createGoalRequest struct {
	Name *string `json:"name,omitempty"`
}
type createGoalResponse struct {
	Message string `json:"message,omitempty"`
}

type retrieveGoalResponse struct {
	Achievers *goalAchievers `json:"achievers,omitempty"`
}
type goalAchievers map[string]*goalAchiever
type goalMessages map[string]string
type goalAchiever struct {
	State    *string       `json:"state,omitempty"`
	Progress *int          `json:"progress,omitempty"`
	Messages *goalMessages `json:"messages,omitempty"`
}

type retrieveGoalsResponse struct {
	Goals *[]*dto.GoalDTO `json:"goals,omitempty"`
}

type updateGoalProgressRequest struct {
	Name     *string `json:"name,omitempty"`
	Progress *int    `json:"progress,omitempty"`
}
type updateGoalProgressResponse struct {
	Message string `json:"message,omitempty"`
}
type abandonGoalRequest struct {
	Name *string `json:"name,omitempty"`
}
type abandonGoalResponse struct {
	Message string `json:"message,omitempty"`
}
type deletedGoalResponse struct {
	Message string `json:"message,omitempty"`
}
