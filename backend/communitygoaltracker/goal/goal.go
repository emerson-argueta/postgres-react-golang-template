package goal

// Goal represents a goal that an achiever is trying to complete.
type Goal struct {
	ID        *int64     `db:"id" dbignoreinsert:"" json:"id"`
	Name      *string    `db:"name" json:"name"`
	Achievers *Achievers `db:"achievers" json:"achievers"`
}

// Achievers represensts a map of achievers within a goal where the
// key is an achiever UUID and the value contains the adminstrator.
type Achievers map[string]*Achiever

// Keys represent the achiever uuid
func (c *Achievers) Keys() []string {
	keys := make([]string, len(*c))

	i := 0
	for k := range *c {
		keys[i] = k
		i++
	}
	return keys
}

// Achiever represents an Achievers within a goal.
type Achiever struct {
	State    *State    `json:"state,omitempty"`
	Progress *int      `json:"progress,omitempty"`
	Messages *Messages `json:"messages,omitempty"`
}

// State represents the state of the goal for a particular achiever.
type State int

const (
	// InProgress when a goal is below 100 in Progress.
	InProgress State = iota
	// Abondoned when a goal is no long inprogress.
	Abondoned
	// Completed when a goal is 100 in Progress.
	Completed
)

func (s State) String() string {
	return [...]string{"inprogress", "abandoned", "completed"}[s]

}

// Messages represents a map of messages for an achiever within a goal
// where the key is a timestamp and the value is the message
type Messages map[string]string

// Keys represent the timestamps
func (m *Messages) Keys() []string {
	keys := make([]string, len(*m))

	i := 0
	for k := range *m {
		keys[i] = k
		i++
	}
	return keys
}

// Service provides processes provided by the identity domain.
type Service interface {
	CreateGoal(*Goal) (*Goal, error)
	RetrieveGoal(id int64) (*Goal, error)
	UpdateGoal(*Goal) error
	DeleteGoal(id int64) error
	RetrieveGoals(ids []int64) ([]*Goal, error)
	UpdateGoals(goals []*Goal) error
}
