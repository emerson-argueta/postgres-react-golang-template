package achiever

// Name value object
type Name interface {
	ToString() string
}

type name string

func (n name) ToString() string {
	return string(n)
}

// NewName creates a name. Empty string default
func NewName(newName *string) Name {
	if newName == nil {
		return name("")
	}

	return name(*newName)
}
