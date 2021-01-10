package achiever

// Address value object
type Address interface {
	ToString() string
}

type address string

func (a address) ToString() string {
	return string(a)
}

// NewAddress creates a address. Empty string default
func NewAddress(newAddress *string) Address {
	if newAddress == nil {
		return address("")
	}

	return address(*newAddress)
}
