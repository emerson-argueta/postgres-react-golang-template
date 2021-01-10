package repository

import "emersonargueta/m/v1/modules/communitygoaltracker/domain/achiever"

// AchieverRepo used to modify the achiever model.
type AchieverRepo interface {
	// CreateAchiever implementation must return ErrAchieverExists if achiever
	// exists.
	CreateAchiever(achiever.Achiever) error
	// RetrieveAchiever implementation must return ErrAchieverNotFound if the
	// achiever is not found.
	RetrieveAchiever(uuid string) (achiever.Achiever, error)
	// RetrieveAchiever implementation must return ErrAchieverNotFound if the
	// achiever is not found.
	RetrieveAchieverByUserID(userID string) (achiever.Achiever, error)
	// RetrieveGoals implementation must return ErrAchieverNotFound if the
	// none of the achievers are found.
	RetrieveAchieversByUserIDs(userIDs []string) ([]achiever.Achiever, error)
	// UpdateAchiever implementation must search achiever by uuid and return
	// ErrAchieverNotFound if achiever is not found.
	UpdateAchiever(achiever.Achiever) error
	// DeleteAchiever implementation should search the achiever by uuid before
	// deleting the achiever and must return ErrAchieverNotFound if the achiever
	// does not exists.
	DeleteAchiever(uuid string) error
}
