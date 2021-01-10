package usecase

type RetrieveAchieversDTO struct {
	AchieverUserID   string
	AchieverGoalName string
}
type RetrieveAchieverDTO struct {
	AchieverUserID string
}
type CreateAchieverDTO struct {
	UserID string
}
type CreateGoalDTO struct {
	Name           string
	AchieverUserID string
}
type RetrieveGoalsDTO struct {
	AchieverUserID string
}
type UpdateGoalProgressDTO struct {
	AchieverUserID string
	Name           string
	Progress       int
}
type AbandonGoalDTO struct {
	Name           string
	AchieverUserID string
}
type DeleteGoalDTO struct {
	Name           string
	AchieverUserID string
}
