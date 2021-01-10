package identity

const (
	// DomainName of this package
	DomainName = "identity"
)

// // Processes an interface for identity domain processes.
// type Processes interface {
// 	UnRegisterUser(*user.User) error
// }

// // UnRegisterUser using the following business logic
// // Validate email and password.
// // Delete the user.
// func (s *Service) UnRegisterUser(u *user.User) (e error) {
// 	if u.UUID == nil || u.Email == nil || u.Password == nil {
// 		return ErrUserIncompleteDetails
// 	}

// 	if exists, e := s.User.RetrieveUser(*u.UUID); e != nil {
// 		return e
// 	} else if *exists.Email != *u.Email || bcrypt.CompareHashAndPassword([]byte(*exists.Password), []byte(*u.Password)) == bcrypt.ErrMismatchedHashAndPassword {
// 		return ErrUserIncorrectCredentials
// 	}

// 	return s.User.DeleteUser(*u.UUID)
// }
