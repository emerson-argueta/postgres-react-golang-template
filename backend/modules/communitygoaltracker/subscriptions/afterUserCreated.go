package subscriptions

import (
	"emersonargueta/m/v1/modules/communitygoaltracker/usecase"
	"emersonargueta/m/v1/modules/identity/domain/user"
	"emersonargueta/m/v1/shared/domain"
	"fmt"
)

type AfterUserCreated interface {
	domain.Handle
}

type afterUserCreated struct {
	createAchiever *usecase.CreateAchieverUsecase // TODO replace type with usecase, createAchiever
}

// TODO replace type with usecase, createAchiever
func NewAfterUserCreated(createAchiever *usecase.CreateAchieverUsecase) AfterUserCreated {
	subscription := &afterUserCreated{createAchiever: createAchiever}
	subscription.SetupSubscriptions()

	return subscription
}

func (s *afterUserCreated) SetupSubscriptions() {

	domain.DomainEvents.Register(s.onUserCreated, user.UserCreatedEventName)
}

func (s *afterUserCreated) onUserCreated(event domain.Event) {
	userID := event.(user.UserCreated).GetUser().GetID()

	e := s.createAchiever.Execute(&usecase.CreateAchieverDTO{UserID: userID})
	if e != nil {
		fmt.Println("[AfterUserCreated]: Failed to execute createAchiever use case AfterUserCreated ==>", e)
	} else {
		fmt.Println("[AfterUserCreated]: Successfully executed createAchiever use case AfterUserCreated ==>", userID)
	}

}
