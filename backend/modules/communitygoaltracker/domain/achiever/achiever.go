package achiever

import (
	"emersonargueta/m/v1/shared/domain"

	"github.com/google/uuid"
)

// Fields used to create achiever
type Fields struct {
	UserID    *UserID
	Role      *Role
	Firstname *Name
	Lastname  *Name
	Address   *Address
	Phone     *Phone
	Goals     *Goals
}

// Achiever model in the communinty_goal_tracker domain
type Achiever interface {
	GetID() string
	GetUserID() UserID
	GetRole() Role
	GetFirstname() Name
	GetLastname() Name
	GetAddress() Address
	GetPhone() Phone
	GetGoals() Goals
	SetFirstname(Name)
	SetLastname(Name)
	SetAddress(Address)
	SetPhone(Phone)
	SetGoals(Goals)
}

type achiever struct {
	ID        string
	UserID    UserID
	Role      Role
	Firstname Name
	Lastname  Name
	Address   Address
	Phone     Phone
	Goals     Goals
	aggregate *domain.AbstractAggregateRoot
}

// Create user with role, email and password
func Create(achieverFields *Fields, id *string) (res Achiever, e error) {
	if achieverFields.UserID == nil {
		return nil, ErrAchieverIncompleteDetails
	}

	achiever := &achiever{
		UserID:    *achieverFields.UserID,
		Role:      *achieverFields.Role,
		Firstname: *achieverFields.Firstname,
		Lastname:  *achieverFields.Lastname,
		Address:   *achieverFields.Address,
		Phone:     *achieverFields.Phone,
		Goals:     *achieverFields.Goals,
	}

	isNewAchiever := (id == nil)
	if isNewAchiever {
		achiever.ID = uuid.New().String()
	}
	if !isNewAchiever {
		achiever.ID = *id
	}

	achiever.aggregate = &domain.AbstractAggregateRoot{}
	achiever.aggregate.DomainEvents = make([]domain.Event, 0)
	achiever.aggregate.Name = "Achiever"
	achiever.aggregate.ID = achiever.ID

	if isNewAchiever {
		achiever.aggregate.AddDomainEvent(NewAchieverCreated(achiever))
	}

	return achiever, nil
}

func (a *achiever) GetID() string {
	return a.ID
}
func (a *achiever) GetUserID() UserID {
	return a.UserID
}
func (a *achiever) GetRole() Role {
	return a.Role
}
func (a *achiever) GetFirstname() Name {
	return a.Firstname
}
func (a *achiever) GetLastname() Name {
	return a.Lastname
}
func (a *achiever) GetAddress() Address {
	return a.Address
}
func (a *achiever) GetPhone() Phone {
	return a.Phone
}
func (a *achiever) GetGoals() Goals {
	return a.Goals
}
func (a *achiever) SetFirstname(firstname Name) {
	a.Firstname = firstname
}
func (a *achiever) SetLastname(lastname Name) {
	a.Lastname = lastname
}
func (a *achiever) SetAddress(address Address) {
	a.Address = address
}
func (a *achiever) SetPhone(phone Phone) {
	a.Phone = phone
}
func (a *achiever) SetGoals(goals Goals) {
	a.Goals = goals
}
