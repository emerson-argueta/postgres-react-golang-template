package persistence

import (
	"emersonargueta/m/v1/modules/communitygoaltracker/domain/achiever"
	"emersonargueta/m/v1/modules/communitygoaltracker/domain/goal"
	"emersonargueta/m/v1/modules/communitygoaltracker/dto"
	"encoding/json"
)

// AchieverDTO for scanning from database
type AchieverDTO struct {
	ID        *string `db:"id" json:"id"`
	UserID    *string `db:"userid" json:"userid"`
	Role      *string `db:"role" json:"role"`
	Firstname *string `db:"firstname" json:"firstname"`
	Lastname  *string `db:"lastname" json:"lastname"`
	Address   *string `db:"address" json:"address"`
	Phone     *string `db:"phone" json:"phone"`
	Goals     *[]byte `db:"goals" json:"goals,omitempty"`
}

// AchieverPersistenceToDomain from persistence dto
func AchieverPersistenceToDomain(dto AchieverDTO) (res achiever.Achiever, e error) {
	userid, e := achiever.NewUserID(dto.UserID)
	if e != nil {
		return nil, e
	}
	role, e := achiever.NewRole(dto.Role)
	if e != nil {
		return nil, e
	}
	firstname := achiever.NewName(dto.Firstname)
	lastname := achiever.NewName(dto.Lastname)
	address := achiever.NewAddress(dto.Address)
	phone, e := achiever.NewPhone(dto.Phone)
	if e != nil {
		return nil, e
	}
	goals := achiever.NewGoals(dto.Goals)

	achieverFields := &achiever.Fields{
		UserID:    &userid,
		Role:      &role,
		Firstname: &firstname,
		Lastname:  &lastname,
		Address:   &address,
		Phone:     &phone,
		Goals:     &goals,
	}
	return achiever.Create(achieverFields, dto.ID)
}

// AchieverDomainToPersistence from achiever domain model
func AchieverDomainToPersistence(a achiever.Achiever) *AchieverDTO {
	id := a.GetID()
	rawUserid := a.GetUserID()
	var userid *string = nil
	if rawUserid != nil {
		u := a.GetUserID().ToString()
		userid = &u
	}
	role, _ := a.GetRole().String()
	firstname := a.GetFirstname().ToString()
	lastname := a.GetLastname().ToString()
	address := a.GetAddress().ToString()
	phone := a.GetPhone().ToString()
	var goals *[]byte = nil
	if a.GetGoals() != nil {
		gg, _ := json.Marshal(a.GetGoals())
		goals = &gg
	}

	return &AchieverDTO{
		ID:        &id,
		UserID:    userid,
		Role:      &role,
		Firstname: &firstname,
		Lastname:  &lastname,
		Address:   &address,
		Phone:     &phone,
		Goals:     goals,
	}

}

// GoalDTO for scanning from database.
type GoalDTO struct {
	ID        *int64  `db:"id" dbignoreinsert:"" json:"id"`
	Name      *string `db:"name" json:"name"`
	Achievers *[]byte `db:"achievers" json:"achievers"`
}

// GoalPersistenceToDomain from persistence dto
func GoalPersistenceToDomain(dto GoalDTO) (res goal.Goal, e error) {
	name, e := goal.NewName(dto.Name)
	if e != nil {
		return nil, e
	}

	achievers := NewAchieversFromByteArray(dto.Achievers)

	achieverFields := &goal.Fields{
		Name:      &name,
		Achievers: &achievers,
	}
	return goal.Create(achieverFields, dto.ID)
}

// NewAchieversFromByteArray creates new achievers. Nil default
func NewAchieversFromByteArray(newAchievers *[]byte) goal.Achievers {
	if newAchievers == nil {
		return nil
	}

	primativeTypeAchiever := make(
		map[string]*struct {
			State    *string            `json:"state,omitempty"`
			Progress *int               `json:"progress,omitempty"`
			Messages *map[string]string `json:"messages,omitempty"`
		},
	)
	json.Unmarshal(*newAchievers, &primativeTypeAchiever)
	a, _ := dto.PrimativeTypeToAchievers(&primativeTypeAchiever)

	return a

}

// GoalDomainToPersistence from goal domain model
func GoalDomainToPersistence(g goal.Goal) *GoalDTO {
	id := g.GetID()
	name := g.GetName().ToString()
	achievers := g.GetAchievers()
	var achieversByte *[]byte = nil

	if achievers != nil {
		achieversMap := make(map[string]goalAchieverDTO)
		for k, v := range achievers {
			state, _ := v.State.String()
			progress := v.Progress.ToInt()
			messages := v.Messages.ToPrimativeType()
			achieversMap[k] = goalAchieverDTO{
				State:    &state,
				Progress: &progress,
				Messages: &messages,
			}
		}
		aa, _ := json.Marshal(achieversMap)
		achieversByte = &aa
	}

	return &GoalDTO{
		ID:        &id,
		Name:      &name,
		Achievers: achieversByte,
	}

}

type goalAchieverDTO struct {
	State    *string            `json:"state,omitempty"`
	Progress *int               `json:"progress,omitempty"`
	Messages *map[string]string `json:"messages,omitempty"`
}
