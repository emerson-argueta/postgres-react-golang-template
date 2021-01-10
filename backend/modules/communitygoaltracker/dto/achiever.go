package dto

import "emersonargueta/m/v1/modules/communitygoaltracker/domain/achiever"

// AchieverDTO for data transfer
type AchieverDTO struct {
	ID        *string          `json:"id,omitempty"`
	UserID    *string          `json:"userid,omitempty"`
	Role      *string          `json:"role,omitempty"`
	Firstname *string          `json:"firstname,omitempty"`
	Lastname  *string          `json:"lastname,omitempty"`
	Address   *string          `json:"address,omitempty"`
	Phone     *string          `json:"phone,omitempty"`
	Goals     *map[string]bool `json:"goals,omitempty"`
}

// ToDomain from achiever dto
func ToDomain(dto AchieverDTO) (res achiever.Achiever, e error) {
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
	phone, e := achiever.NewPhone(dto.Address)
	if e != nil {
		return nil, e
	}
	goals := achiever.ToGoals(dto.Goals)

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

// ToDTO from achiever domain model
func ToDTO(a achiever.Achiever) *AchieverDTO {
	id := a.GetID()
	userID := a.GetUserID().ToString()

	role := a.GetRole()
	roleStr, _ := role.String()

	firstname := a.GetFirstname().ToString()
	lastname := a.GetLastname().ToString()
	address := a.GetAddress().ToString()
	phone := a.GetPhone().ToString()

	goals := a.GetGoals()
	goalsPrimitive := goals.ToPrimativeType()

	dto := &AchieverDTO{
		ID:        &id,
		UserID:    &userID,
		Role:      &roleStr,
		Firstname: &firstname,
		Lastname:  &lastname,
		Address:   &address,
		Phone:     &phone,
		Goals:     &goalsPrimitive,
	}

	return dto
}
