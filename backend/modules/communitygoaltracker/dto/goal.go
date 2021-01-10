package dto

import "emersonargueta/m/v1/modules/communitygoaltracker/domain/goal"

// GoalDTO for data transfer
type GoalDTO struct {
	ID        *int64  `json:"id,omitempty"`
	Name      *string `json:"name,omitempty"`
	Achievers *map[string]*struct {
		State    *string            `json:"state,omitempty"`
		Progress *int               `json:"progress,omitempty"`
		Messages *map[string]string `json:"messages,omitempty"`
	} `json:"achievers,omitempty"`
}

// GoalDtoToDomain from goal dto to goal domain model
func GoalDtoToDomain(dto GoalDTO) (res goal.Goal, e error) {
	name, e := goal.NewName(dto.Name)
	if e != nil {
		return nil, e
	}
	achievers, e := PrimativeTypeToAchievers(dto.Achievers)

	goalFields := &goal.Fields{
		Name:      &name,
		Achievers: &achievers,
	}
	return goal.Create(goalFields, dto.ID)
}

// PrimativeTypeToAchievers given a primitive type. Empty map if primitive type is nil.
// Returns error if map values contain invalid fields.
func PrimativeTypeToAchievers(
	a *map[string]*struct {
		State    *string            `json:"state,omitempty"`
		Progress *int               `json:"progress,omitempty"`
		Messages *map[string]string `json:"messages,omitempty"`
	}) (res goal.Achievers, e error) {
	if a == nil {
		return make(map[string]*goal.Achiever), nil
	}
	res = make(map[string]*goal.Achiever)
	for k, v := range *a {
		state, e := goal.NewState(v.State)
		if e != nil {
			return nil, e
		}
		progress, e := goal.NewProgress(v.Progress)
		if e != nil {
			return nil, e
		}
		messages := goal.NewMessages(v.Messages)
		achiever := &goal.Achiever{
			State:    &state,
			Progress: &progress,
			Messages: &messages,
		}
		res[k] = achiever
	}

	return goal.Achievers(res), nil
}

// GoalToDTO from goal domain model to goal dto
func GoalToDTO(g goal.Goal) *GoalDTO {
	id := g.GetID()

	name := g.GetName().ToString()

	achievers := g.GetAchievers()
	achieversPrimitive := AchieversToPrimativeType(&achievers)

	dto := &GoalDTO{
		ID:        &id,
		Name:      &name,
		Achievers: &achieversPrimitive,
	}

	return dto
}

// AchieversToPrimativeType from Achievers type
func AchieversToPrimativeType(a *goal.Achievers) map[string]*struct {
	State    *string            `json:"state,omitempty"`
	Progress *int               `json:"progress,omitempty"`
	Messages *map[string]string `json:"messages,omitempty"`
} {
	primativeType := make(map[string]*struct {
		State    *string            `json:"state,omitempty"`
		Progress *int               `json:"progress,omitempty"`
		Messages *map[string]string `json:"messages,omitempty"`
	})

	for k, v := range *a {
		state, _ := v.State.String()
		progress := v.Progress.ToInt()
		messages := v.Messages.ToPrimativeType()
		pt := &struct {
			State    *string            `json:"state,omitempty"`
			Progress *int               `json:"progress,omitempty"`
			Messages *map[string]string `json:"messages,omitempty"`
		}{
			State:    &state,
			Progress: &progress,
			Messages: &messages,
		}
		primativeType[k] = pt
	}
	return primativeType
}
