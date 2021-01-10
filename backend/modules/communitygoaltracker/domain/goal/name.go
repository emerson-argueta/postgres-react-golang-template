package goal

// Name value object
type Name interface {
	ToString() string
}

type name string

func (n name) ToString() string {
	return string(n)
}

// NewName creates a name. Returns error if newName is nil.
func NewName(newName *string) (res Name, e error) {
	if newName == nil {
		return nil, ErrGoalIncompleteDetails
	}

	return name(*newName), nil
}
