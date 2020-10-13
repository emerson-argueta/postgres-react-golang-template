package http

import (
	"emersonargueta/m/v1/communitygoaltracker/achiever"
	"emersonargueta/m/v1/communitygoaltracker/goal"
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

func requestToGoal(request *goalRequest, achieverUUID string) (res *goal.Goal) {
	achiever := make(goal.Achievers)

	var state goal.State
	var message goal.Messages
	if request.State != nil {
		state = goal.ToState(*request.State)
	}
	if request.Message != nil && request.Timestamp != nil {
		message = make(goal.Messages)
		message[*request.Timestamp] = *request.Message
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

func goalToResponse(g *goal.Goal) (res *goalResponse) {
	var achievers Achievers

	if g.Achievers != nil {
		achievers = make(Achievers)
		for uuid, achiever := range *g.Achievers {
			state := achiever.State.String()
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
