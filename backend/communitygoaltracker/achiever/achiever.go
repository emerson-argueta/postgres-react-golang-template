package achiever

// Role represents access role for an achiever within the communitygoaltracker domain
type Role int

const (
	// UserRole has access to manage their own information
	UserRole Role = iota
	// AdministratorRole has access to manage users within their own goals
	AdministratorRole
)

func (s Role) String() string {
	return [...]string{"user", "administrator"}[s]

}

// Achiever model in the communinty_goal_tracker domain
type Achiever struct {
	UUID      *string `db:"uuid" json:"uuid"`
	Firstname *string `db:"firstname" json:"firstname"`
	Lastname  *string `db:"lastname" json:"lastname"`
	Address   *string `db:"address" json:"address"`
	Phone     *string `db:"phone" json:"phone"`
	Goals     *Goals  `db:"goals" json:"goals,omitempty"`
	Email     *string `json:"email,omitempty"`
	Password  *string `json:"password,omitempty"`
}

// Goals represents a slice goals ids for an achiever
type Goals []int64

// Service provides processes that can be achieved by an achiever.
type Service interface {
	CreateAchiever(*Achiever) (*Achiever, error)
	RetrieveAchiever(email string, password string) (*Achiever, error)
	UpdateAchiever(*Achiever) error
	DeleteAchiever(*Achiever) error
}
