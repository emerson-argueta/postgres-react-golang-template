package achiever

// UserID value object
type UserID interface {
	ToString() string
}

type userid string

func (u userid) ToString() string {
	return string(u)
}

// NewUserID creates a userid, returns error if newUserID is nil
func NewUserID(newUserID *string) (res UserID, e error) {
	if newUserID == nil {
		return nil, ErrAchieverIncompleteDetails
	}

	return userid(*newUserID), nil
}
