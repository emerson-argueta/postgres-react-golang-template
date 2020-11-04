package http

import (
	"emersonargueta/m/v1/communitygoaltracker/achiever"
	"emersonargueta/m/v1/communitygoaltracker/goal"
	"time"
)

type achieverRequest struct {
	Achiever      *achiever.Achiever `json:"achiever,omitempty"`
	Authorization map[string]string  `json:"authorization,omitempty"`
}

type achieverResponse struct {
	Achiever      *achiever.Achiever `json:"achiever,omitempty"`
	Authorization map[string]string  `json:"authorization,omitempty"`
	Error         string             `json:"error,omitempty"`
}

type goalRequest struct {
	ID            *int64            `json:"id,omitempty"`
	Name          *string           `json:"name,omitempty"`
	State         *string           `json:"state,omitempty"`
	Progress      *int              `json:"progress,omitempty"`
	Message       *string           `json:"message,omitempty"`
	Timestamp     *string           `json:"timestamp,omitempty"`
	Authorization map[string]string `json:"authorization,omitempty"`
}

type goalResponse struct {
	ID        *int64     `json:"id,omitempty"`
	Name      *string    `json:"name,omitempty"`
	Achievers *Achievers `json:"achievers,omitempty"`
	Error     string     `json:"error,omitempty"`
}

// Achievers is a wrapper for goal.Achievers where the key is an achiever UUID and the value contains the achiever.
type Achievers map[string]struct {
	State    *string        `json:"state,omitempty"`
	Progress *int           `json:"progress,omitempty"`
	Messages *goal.Messages `json:"messages,omitempty"`
}

// Creates a goal from a goalRequest with the goal's the achiever's UUID
// as the first entry in the goal's Achiever field
func requestToGoal(request *goalRequest, achieverUUID string) (res *goal.Goal) {
	achiever := make(goal.Achievers)

	state := goal.InProgress
	if request.State != nil {
		state, _ = goal.ToState(*request.State)
	}

	message := make(goal.Messages)
	if request.Message != nil && request.Timestamp != nil {
		message[*request.Timestamp] = *request.Message
	} else if request.Message != nil && request.Timestamp == nil {
		currentTimestamp := time.Now().UTC().Format(time.RFC3339)
		message[currentTimestamp] = *request.Message
	}

	achiever[achieverUUID] = &goal.Achiever{
		State:    &state,
		Progress: request.Progress,
		Messages: &message,
	}

	res = &goal.Goal{
		Name:      request.Name,
		Achievers: &achiever,
	}
	return res
}

func goalToRequest(g *goal.Goal) (res *goalRequest) {
	if g.Achievers != nil {
		achieverUUID := g.Achievers.Keys()[0]
		achieverGoal := (*g.Achievers)[achieverUUID]
		res.ID = g.ID
		res.Name = g.Name
		currentTimestamp := time.Now().UTC().Format(time.RFC3339)
		res.Timestamp = &currentTimestamp
		if achieverGoal.Messages != nil {
			messageTime := achieverGoal.Messages.Keys()[0]
			firstMessage := (*achieverGoal.Messages)[messageTime]
			res.Message = &firstMessage
			res.Timestamp = &messageTime
		}
		res.Progress = achieverGoal.Progress
		state, _ := achieverGoal.State.String()
		res.State = &state
	}

	return res
}

func goalToResponse(g *goal.Goal) (res *goalResponse) {
	var achievers Achievers

	if g.Achievers != nil {
		achievers = make(Achievers)
		for uuid, achiever := range *g.Achievers {
			state, _ := achiever.State.String()
			achievers[uuid] = struct {
				State    *string        "json:\"state,omitempty\""
				Progress *int           "json:\"progress,omitempty\""
				Messages *goal.Messages "json:\"messages,omitempty\""
			}{
				State:    &state,
				Progress: achiever.Progress,
				Messages: achiever.Messages,
			}

		}
	}

	res = &goalResponse{
		ID:        g.ID,
		Name:      g.Name,
		Achievers: &achievers,
	}
	return res
}

func responseToGoal(g *goalResponse) (res *goal.Goal) {
	achievers := make(goal.Achievers)
	for achieverUUID, achieverGoal := range *g.Achievers {
		achievers[achieverUUID].Messages = achieverGoal.Messages
		achievers[achieverUUID].Progress = achieverGoal.Progress
		stringState := achieverGoal.State
		state, _ := goal.ToState(*stringState)
		achievers[achieverUUID].State = &state
	}
	res.Achievers = &achievers

	return res
}
