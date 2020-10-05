package administrator

import (
	"strconv"
	"strings"
	"time"

	"emersonargueta/m/v1/domain"
	"emersonargueta/m/v1/domain/church"
	"emersonargueta/m/v1/domain/donator"
	"emersonargueta/m/v1/domain/transaction"
	"emersonargueta/m/v1/user"
	"emersonargueta/m/v1/validation"

	"golang.org/x/crypto/bcrypt"
)

//TODO: refactor function parameters that are domain structs with only id field value present to int64 type

// SubscriptionActions provides functions that can be used to manage
// subscriptions for administrator.
type SubscriptionActions interface {
	NewSubscription(*Subscription, *Administrator) error
}

// UnauthorizedManagementActions provides functions that an administrator can
// use when administrator has not been authorized.
type UnauthorizedManagementActions interface {
	Register(*Administrator) error
	Login(*Administrator) (*Administrator, error)
}

// AuthorizedManagementActions provides functions that an authorized
// administrator can use to manage administrators,churches,donators,and
// donations
type AuthorizedManagementActions interface {
	ReadAdministrators(a *Administrator, c *church.Church) (map[string]*Administrator, error)
	EditAdministrator(a *Administrator) (*Administrator, error)
	EditAdministrators(
		c *church.Church,
		a *Administrator,
		aToEdit *Administrator,
		newChurchAdministrator *church.Administrator,
	) (*church.Administrator, error)

	AddChurch(*Administrator, *church.Church) (*church.Church, error)
	CreateChurch(*Administrator, *church.Church) (*church.Church, error)
	ReadChurches(*Administrator) (map[int64]*church.Church, error)
	EditChurch(a *Administrator, c *church.Church) (*church.Church, error)
	// RemoveChurch(a *Administrator, c *church.Church) error

	CreateDonator(*Administrator, *church.Church, *donator.Donator) (*donator.Donator, error)
	ReadDonators(*Administrator, *church.Church) (map[int64]*donator.Donator, error)
	EditDonator(*Administrator, *church.Church, *donator.Donator) (*donator.Donator, error)
	// RemoveDonator(a *Administrator, d *donator.Donator) error

	CreateDonation(*Administrator, *transaction.Donation) (*transaction.Donation, error)
	ReadDonations(a *Administrator, c *church.Church) (map[int64][]*transaction.Donation, error)
	// EditDonation(a *Administrator, c *church.Church, d *donator.Donator, t *transaction.Transaction) error

}

// Usecase contains the business logic.
type Usecase struct {
	Services Services
}

// Services used by administrator usecase
type Services struct {
	Administrator Service
	User          user.Service
	Church        church.Service
	Donator       donator.Service
	Transaction   transaction.Service
}

// Register an administrator using the following business logic: Search user by
// email and administrator by user UUID
//  case0:If user does not exists and administrator does not exists --> create user and administrator
//		with free subscription type and freeusagelimitcount set to 0
//  case1:If user exists and administrator does not exists
//          If password matches --> set user.UUID to administrator.UUID, then create administrator
//				with free subscription type and freeusagelimitcount set 0
//          If password does not match --> return error user incorrect credentials
//  case2:If user exists and administrator exists --> return error administrator already registered
func (uc *Usecase) Register(a *Administrator) (e error) {

	if a.Email == nil || a.Password == nil {
		return ErrAdministratorIncompleteDetails
	}

	byEmail := true
	u := user.User{Email: a.Email, Password: a.Password}
	userRead, _ := uc.Services.User.Retrieve(&u, byEmail)
	// case0
	if userRead == nil {
		if err := validation.ValidateUserEmail(*a.Email); err != nil {
			return validation.ErrValidationUserEmail
		} else if err := validation.ValidatePassword(*a.Password); err != nil {
			return validation.ErrValidationPassword
		}
		return uc.registerUserNF(&u, a)
	}

	a.UUID = userRead.UUID
	administratorRead, _ := uc.Services.Administrator.Read(a)
	// case1
	if userRead != nil && administratorRead == nil {
		return uc.registerUserFadminNF(userRead, a)
	}
	// case2
	if userRead != nil && administratorRead != nil {
		return ErrAdministratorExists
	}

	return ErrAdministratorRegister

}

// Register Case0: When user not found administrator must not exist.
func (uc *Usecase) registerUserNF(u *user.User, a *Administrator) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(*a.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	hashstr := string(hash)
	u.Password = &hashstr
	a.Password = &hashstr

	serviceRole := role
	u.Domains = &user.Domains{domain.ServiceName: {Role: &serviceRole}}
	if err := uc.Services.User.Register(u); err != nil {
		return err
	}
	a.UUID = u.UUID

	freeusagelimitcount := int64(0)
	subscriptionType := FreePlan
	subscription := &Subscription{
		Freeusagelimitcount: &freeusagelimitcount,
		Customeremail:       a.Email,
		Type:                &subscriptionType,
	}
	a.Subscription = subscription
	return uc.Services.Administrator.Create(a)
}

// Register Case1: When user is found and administrator is not found.
func (uc *Usecase) registerUserFadminNF(u *user.User, a *Administrator) error {
	if err := bcrypt.CompareHashAndPassword([]byte(*u.Password), []byte(*a.Password)); err == nil {
		hash, err := bcrypt.GenerateFromPassword([]byte(*a.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		hashstr := string(hash)
		a.Password = &hashstr

		freeusagelimitcount := int64(0)
		subscriptionType := FreePlan
		subscription := &Subscription{
			Freeusagelimitcount: &freeusagelimitcount,
			Customeremail:       a.Email,
			Type:                &subscriptionType,
		}
		a.Subscription = subscription
		return uc.Services.Administrator.Create(a)
	} else if err == bcrypt.ErrMismatchedHashAndPassword {
		return ErrAdministratorIncorrectCredentials
	} else {
		return err
	}

}

// Login an administrator using the following business logic:
// Search user by email and administrator by user UUID
// 	case0:If user does not exists and administrator does not exists --> return error user does not exist
// 	case1:If user exists and administrator does not exists --> return error administrator does not exists
// 	case2:If user exists and administrator exists
// 			If password matches --> success, return administrator and no error
// 			If password does not match --> return error administrator incorrect credentials
func (uc *Usecase) Login(a *Administrator) (res *Administrator, e error) {
	if a.Email == nil || a.Password == nil {
		return nil, ErrAdministratorIncompleteDetails
	}

	byEmail := true
	u := user.User{Email: a.Email, Password: a.Password}

	userRead, _ := uc.Services.User.Retrieve(&u, byEmail)

	// case0
	if userRead == nil {
		return nil, user.ErrUserNotFound
	}

	a.UUID = userRead.UUID
	administratorRead, _ := uc.Services.Administrator.Read(a)
	administratorRead.Email = a.Email
	administratorRead.Password = a.Password

	// case1
	if userRead != nil && administratorRead == nil {
		return nil, ErrAdministratorNotFound
	}
	// case2
	if userRead != nil && administratorRead != nil {
		return uc.loginUserFadminF(userRead, administratorRead)
	}

	return nil, ErrAdministratorLogin
}
func (uc *Usecase) loginUserFadminF(u *user.User, a *Administrator) (*Administrator, error) {
	err := bcrypt.CompareHashAndPassword([]byte(*u.Password), []byte(*a.Password))
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		return nil, ErrAdministratorIncorrectCredentials
	} else if err != nil {
		return nil, err
	}

	return a, nil
}

// NewSubscription for an administrator. A new subscription for an
// administrator can be created so that the administrator can gain access to the
// church_fund_managing service after reaching the free usage limit
// (100 automated donation entries and 400 manual donation entries).
func (uc *Usecase) NewSubscription(s *Subscription, a *Administrator) (e error) {
	a.Subscription = s
	uc.Services.Administrator.Update(a)
	return nil
}

// ReadAdministrators by UUID belonging to a church using the following businiess logic:
// If church does not belong to administrator return error church does not belong to administrator
// Else return read administrators for church
func (uc *Usecase) ReadAdministrators(a *Administrator, c *church.Church) (res map[string]*Administrator, e error) {

	aRead, err := uc.Services.Administrator.Read(a)
	if err != nil {
		return nil, err
	}
	if aRead == nil {
		return nil, ErrAdministratorNotFound
	}
	if aRead.Churches == nil {
		return nil, ErrAdministratorNoChurches
	}

	_, ok := (*aRead.Churches)[*c.ID]
	if !ok {
		return nil, ErrAdministratorDoesNotBelongToChurch
	}
	byEmail := false
	cRead, err := uc.Services.Church.Read(c, byEmail)
	if err != nil {
		return nil, err
	}

	aa, err := uc.Services.Administrator.ReadMultiple(cRead.Administrators.Keys())
	if err != nil {
		return nil, err
	}

	res = make(map[string]*Administrator, len(aa))
	for _, elem := range aa {
		// Administrator model does not store email or password in database.
		// However, email is found on administrator's subscription.
		elem.Email = elem.Subscription.Customeremail
		res[*elem.UUID] = elem
	}

	return res, e
}

// EditAdministrator using the following business logic:
// If Churches or Subscription were to be updated ---> return field not editable error
// If the email or password is updated validate correct email and password length else ---> return validation error
func (uc *Usecase) EditAdministrator(a *Administrator) (res *Administrator, e error) {
	if a.Churches != nil || a.Subscription != nil {
		return nil, ErrAdministratorFieldNotEditable
	}
	aRead, err := uc.Services.Administrator.Read(a)
	if err != nil {
		return nil, err
	}
	if aRead == nil {
		return nil, ErrAdministratorNotFound
	}

	if a.Password != nil && validation.ValidatePassword(*a.Password) != nil {
		return nil, validation.ErrValidationPassword
	} else if a.Password != nil {
		hash, err := bcrypt.GenerateFromPassword([]byte(*a.Password), bcrypt.DefaultCost)
		if err != nil {
			return nil, err
		}
		hashstr := string(hash)
		a.Password = &hashstr
	}

	if a.Email != nil && validation.ValidateUserEmail(*a.Email) != nil {
		return nil, validation.ErrValidationUserEmail
	} else if a.Email != nil && aRead.Subscription != nil {
		aSubscription := aRead.Subscription
		aSubscription.Customeremail = a.Email
		a.Subscription = aRead.Subscription
	} else if a.Email != nil && aRead.Subscription == nil {
		freeUsageLimitCount := int64(0)
		subscriptionType := FreePlan
		customerEmail := a.Email

		subscription := &Subscription{
			Freeusagelimitcount: &freeUsageLimitCount,
			Customeremail:       customerEmail,
			Type:                &subscriptionType,
		}
		a.Subscription = subscription
	}

	if a.Password != nil || a.Email != nil {
		userUpdate := &user.User{UUID: aRead.UUID, Email: a.Email, Password: a.Password}
		byEmail := false
		if err := uc.Services.User.Update(userUpdate, byEmail); err != nil {
			return nil, err
		}

	}

	if err := uc.Services.Administrator.Update(a); err != nil {
		return nil, err
	}

	a.Churches = nil
	a.Subscription = nil

	return a, e
}

// EditAdministrators allows an administrator to edit another administrator within a church they are both part of.
// Edit an administrator using the following business logic:
// If Administrator has restricted access --> return error administrator not authorized to edit other administrators
// If the administrator to be deleted has non-restricted access and is the only non-restricted administrator of the church
// 		--> church will be marked for deletion after 14 days.
// If newChurchAdministrator is not nil
// 		If field to be edit is role ---> return error field not editable
//		Else update churchAdministrators with new churchAdministrator
// If newChurchAdministrator is nil
// 		If a deleted administrator has non-restricted access but the church has multiple administrators with at least one other administrator
// 		with non-restricted access:
// 			OR
// 		If a deleted administrator has restricted access:
// 				The administrator can simply be deleted
func (uc *Usecase) EditAdministrators(
	c *church.Church,
	a *Administrator,
	aToEdit *Administrator,
	newChurchAdministrator *church.Administrator,
) (res *church.Administrator, e error) {
	byEmail := false
	cRead, err := uc.Services.Church.Read(c, byEmail)
	if err != nil {
		return nil, err
	}
	aRead, err := uc.Services.Administrator.Read(a)
	if aRead == nil {
		return nil, ErrAdministratorNotFound
	} else if err != nil {
		return nil, err
	}
	if cRead.Administrators == nil {
		return nil, church.ErrChurchAdministrators
	}
	churchAdmin, ok := (*cRead.Administrators)[*a.UUID]
	if !ok {
		return nil, ErrAdministratorDoesNotBelongToChurch
	}

	if *churchAdmin.Access == domain.Restricted {
		return nil, ErrAdministratorNotAuthorized
	}

	if newChurchAdministrator != nil {
		return uc.updateChurchAdministrators(cRead, newChurchAdministrator, aToEdit)
	}

	return newChurchAdministrator, uc.deleteChurchAdministrators(cRead, aToEdit)
}
func (uc *Usecase) updateChurchAdministrators(c *church.Church, churchAdministrator *church.Administrator, aToEdit *Administrator) (*church.Administrator, error) {
	_, ok := (*c.Administrators)[*aToEdit.UUID]
	if !ok {
		return nil, ErrAdministratorDoesNotBelongToChurch
	}
	if *(*c.Administrators)[*aToEdit.UUID].Role == domain.Creator && *churchAdministrator.Access == domain.Restricted {
		return nil, ErrAdministratorFieldNotEditable
	}

	newChurchAdministrators := *c.Administrators
	newChurchAdministrators[*aToEdit.UUID] = churchAdministrator
	byEmail := false
	if err := uc.Services.Church.Update(&church.Church{ID: c.ID, Administrators: &newChurchAdministrators}, byEmail); err != nil {
		return nil, err
	}

	return churchAdministrator, nil
}
func (uc *Usecase) deleteChurchAdministrators(c *church.Church, aToDelete *Administrator) error {
	_, ok := (*c.Administrators)[*aToDelete.UUID]
	if !ok {
		return ErrAdministratorDoesNotBelongToChurch
	}

	if numberOfNonRestrictedMembers((*c.Administrators)) == 1 {
		markChurchForDeletion(c)
	}

	newChurchAdministrators := *c.Administrators
	delete(newChurchAdministrators, *aToDelete.UUID)
	byEmail := false
	if err := uc.Services.Church.Update(&church.Church{ID: c.ID, Administrators: &newChurchAdministrators}, byEmail); err != nil {
		return err
	}

	return nil
}
func numberOfNonRestrictedMembers(churchAdministrators church.Administrators) int64 {
	var count int64 = 0
	for _, churchAdministrator := range churchAdministrators {
		if *churchAdministrator.Access == domain.NonRestricted {
			count++
		}
	}
	return count
}
func markChurchForDeletion(c *church.Church) {
	//TODO
}

// AddChurch for the administrator using the following business logic:
// Find administrator by uuid and church by email
// If church or administrator not found --> return error
// Find church for administrator
// If administrator has church --> return error admin already has church
// Else --> for church found by email
//  case0: If password incorrect --> return error incorrect credentials for church
//	case1: Else --> add church with Restricted access, and Support role
func (uc *Usecase) AddChurch(a *Administrator, c *church.Church) (res *church.Church, e error) {
	if c.Email == nil || c.Password == nil {
		return nil, ErrAdministratorChurchIncompleteDetails
	}

	byEmail := true
	cRead, err := uc.Services.Church.Read(c, byEmail)
	if err != nil {
		return nil, err
	}
	aRead, err := uc.Services.Administrator.Read(a)
	if aRead == nil {
		return nil, ErrAdministratorNotFound
	} else if err != nil {
		return nil, err
	}

	var cAdministrators church.Administrators
	var access domain.Access
	var role domain.Role
	if cRead.Administrators == nil {
		access = domain.NonRestricted
		role = domain.Support
		cAdministrators = make(church.Administrators)
	} else {
		access = domain.Restricted
		role = domain.Support
		cAdministrators = *cRead.Administrators
	}
	if _, ok := cAdministrators[*aRead.UUID]; ok {
		return nil, ErrAdministratorExists
	}

	var aChurches Churches
	if aRead.Churches == nil {
		aChurches = make(Churches)
	} else {
		aChurches = *aRead.Churches
	}
	if _, ok := aChurches[*cRead.ID]; ok {
		return nil, ErrAdministratorChurchExists
	}

	// case0
	if ok := addChurchPasswordValidated(*cRead.Password, *c.Password); !ok {
		return nil, ErrAdministratorChurchIncorrectCredentilas
	}

	// case1
	cAdministrators[*aRead.UUID] = &church.Administrator{Access: &access, Role: &role}
	aChurches[*cRead.ID] = &Church{Access: &access, Role: &role}

	cUpdate := church.Church{Email: cRead.Email, Administrators: &cAdministrators}
	aUpdate := Administrator{UUID: aRead.UUID, Churches: &aChurches}

	if err := uc.Services.Church.Update(&cUpdate, byEmail); err != nil {
		return nil, err
	}
	if err := uc.Services.Administrator.Update(&aUpdate); err != nil {
		return nil, err
	}

	return cRead, nil
}
func addChurchPasswordValidated(hash string, password string) bool {
	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)); err != nil {
		return false
	}
	return true
}

// CreateChurch for the administrator using the following business logic:
// Find administrator by uuid and church by email
// If administrator not found --> return error
// If church found--> return error
// Else --> add church with Non-Restricted access, and Creator role
func (uc *Usecase) CreateChurch(a *Administrator, c *church.Church) (res *church.Church, e error) {
	if c.Email == nil || c.Password == nil {
		return nil, ErrAdministratorChurchIncompleteDetails
	}
	if err := validation.ValidateUserEmail(*c.Email); err != nil {
		return nil, validation.ErrValidationUserEmail
	} else if err := validation.ValidatePassword(*c.Password); err != nil {
		return nil, validation.ErrValidationPassword
	}

	byEmail := true
	cRead, err := uc.Services.Church.Read(c, byEmail)
	if err != nil {
		return nil, err
	}
	aRead, err := uc.Services.Administrator.Read(a)
	if err != nil {
		return nil, err
	}

	if aRead == nil {
		return nil, ErrAdministratorNotFound
	}
	if cRead != nil {
		return nil, church.ErrChurchExists
	}

	access := domain.NonRestricted
	role := domain.Creator

	cAdministrators := make(church.Administrators)
	cAdministrators[*aRead.UUID] = &church.Administrator{Access: &access, Role: &role}
	c.Administrators = &cAdministrators

	hash, err := bcrypt.GenerateFromPassword([]byte(*c.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	hashstr := string(hash)
	c.Password = &hashstr

	if err := uc.Services.Church.Create(c); err != nil {
		return nil, err
	}

	var aChurches Churches
	if aRead.Churches == nil {
		aChurches = make(Churches)
	} else {
		aChurches = *aRead.Churches
	}

	if _, ok := aChurches[*c.ID]; ok {
		return nil, ErrAdministratorChurchExists
	}

	aChurches[*c.ID] = &Church{Access: &access, Role: &role}
	aUpdate := Administrator{UUID: aRead.UUID, Churches: &aChurches}

	uc.Services.Administrator.Update(&aUpdate)

	return c, nil
}

// ReadChurches for an administrator using the following business logic:
// If administrator does not have churches return error administrator does
// not have churches.
// Else read and return churches
func (uc *Usecase) ReadChurches(a *Administrator) (res map[int64]*church.Church, e error) {
	aRead, err := uc.Services.Administrator.Read(a)
	if err != nil {
		return nil, err
	}
	if aRead == nil {
		return nil, ErrAdministratorNotFound
	}
	if aRead.Churches == nil {
		return nil, ErrAdministratorNoChurches
	}

	cc, err := uc.Services.Church.ReadMultiple(aRead.Churches.Keys())
	if err != nil {
		return nil, err
	}

	res = make(map[int64]*church.Church, len(cc))
	for _, elem := range cc {
		elem.Password = nil
		res[*elem.ID] = elem
	}

	return res, nil
}

// EditChurch for an administrator using the following business logic:
// If church is nil ---> use delete church function
// If Administrators,Donators,Accountstatement were to be updated ---> return field not editable error
// If the church to edit does not belong to the administrator ---> return does not belong error
// If the email or password for the church is updated and administrator is not the creator ---> return not authorized error
// If the email or password is updated validate correct email and password length else ---> return validation error
//	If password is updated ---> remove church from every administrator which was previously part of church except the creator
func (uc *Usecase) EditChurch(a *Administrator, c *church.Church, cID int64) (res *church.Church, e error) {
	if c != nil && (c.Administrators != nil || c.Donators != nil || c.Accountstatement != nil) {
		return nil, church.ErrChurchFieldNotEditable
	}
	aRead, err := uc.Services.Administrator.Read(a)
	if err != nil {
		return nil, err
	}
	if aRead == nil {
		return nil, ErrAdministratorNotFound
	}
	if aRead.Churches == nil {
		return nil, ErrAdministratorNoChurches
	}

	_, ok := (*aRead.Churches)[cID]
	if !ok {
		return nil, ErrAdministratorDoesNotBelongToChurch
	}

	byEmail := false
	cRead, err := uc.Services.Church.Read(&church.Church{ID: &cID}, byEmail)
	if err != nil {
		return nil, err
	}
	cAdministrators := cRead.Administrators
	if cAdministrators == nil {
		return nil, church.ErrChurchAdministrators
	}
	if c == nil {
		return nil, uc.deleteAdministratorChurch(aRead, cRead)
	}

	creatorUUID, _ := uc.findCreatorAdministrator(*cAdministrators)

	if (c.Password != nil || c.Email != nil) && creatorUUID != *aRead.UUID {
		return nil, church.ErrChurchFieldEditUnPriveledged
	}
	if c.Password != nil && validation.ValidatePassword(*c.Password) != nil {
		return nil, validation.ErrValidationPassword
	} else if c.Password != nil {
		// delete church from all administrators except creator
		uc.removeChurchFromSupportAdministrators(*cRead.ID, cRead.Administrators)
		uc.removeSupportAdministratorsFromChurch(cRead)
	}
	if c.Email != nil && validation.ValidateUserEmail(*c.Email) != nil {
		return nil, validation.ErrValidationUserEmail
	}

	if err := uc.Services.Church.Update(c, byEmail); err != nil {
		return nil, err
	}

	return c, e
}
func (uc *Usecase) removeSupportAdministratorsFromChurch(c *church.Church) error {
	// Iterating over keys of churchAdministrators. They keys are the administrator's UUID
	cAdministrators := *c.Administrators
	for uuid, a := range *c.Administrators {
		if a.Role == nil {
			continue
		}
		if a.Role != nil && *a.Role == domain.Creator {
			continue
		}
		delete(cAdministrators, uuid)

	}

	c.Administrators = &cAdministrators
	byEmail := false
	if err := uc.Services.Church.Update(c, byEmail); err != nil {
		return err
	}

	return nil
}
func (uc *Usecase) removeChurchFromSupportAdministrators(cID int64, churchAdministrators *church.Administrators) error {
	// Iterating over keys of churchAdministrators. They keys are the administrator's UUID
	for aUUID, a := range *churchAdministrators {
		if a.Role == nil {
			continue
		}
		if a.Role != nil && *a.Role == domain.Creator {
			continue
		}

		aRead, err := uc.Services.Administrator.Read(&Administrator{UUID: &aUUID})
		if err != nil {
			return err
		}
		if aReadChurches := aRead.Churches; aReadChurches != nil {
			delete((*aReadChurches), cID)
			aRead.Churches = aReadChurches
			if err := uc.Services.Administrator.Update(aRead); err != nil {
				return err
			}
		}
	}

	return nil
}
func (uc *Usecase) deleteAdministratorChurch(a *Administrator, c *church.Church) error {
	aChurches := (*a.Churches)
	cAdministrators := (*c.Administrators)
	aRole := *aChurches[*c.ID].Access

	delete(aChurches, *c.ID)
	delete(cAdministrators, *a.UUID)
	a.Churches = &aChurches
	c.Administrators = &cAdministrators

	if err := uc.Services.Administrator.Update(a); err != nil {
		return err
	}
	byEmail := false
	if err := uc.Services.Church.Update(c, byEmail); err != nil {
		return err
	}

	if aRole == domain.NonRestricted {
		if err := uc.removeChurchFromSupportAdministrators(*c.ID, c.Administrators); err != nil {
			return err
		}
		if err := uc.removeSupportAdministratorsFromChurch(c); err != nil {
			return err
		}

		return uc.deleteChurch(c)
	}

	return nil
}

// Preconditions: Administrators can delete a church only if they have non-restricted access. However, they can remove themselves from churches where they have non-restricted access.
// Steps:
// When a church is deleted all administrators,donators, and donations for that church are deleted. All administrators and donators will be notified.
// The administrator who deletes the church will be given a choice to download a backup of data and a choice to recover the church account within a 15 day window.
// When a church is deleted donations can no longer be made to that church and administrators will no longer have access to get any data from that church.
func (uc *Usecase) deleteChurch(c *church.Church) error {
	for uuid := range *c.Administrators {
		uc.notifyAdministratorChurchDeletion(uuid)
	}
	for id := range *c.Donators {
		uc.notifyDonatorChurchDeletion(id)
	}
	markChurchForDeletion(c)

	return nil
}
func (uc *Usecase) notifyAdministratorChurchDeletion(aUUID string) error {
	return nil
}
func (uc *Usecase) notifyDonatorChurchDeletion(dID int64) error {
	return nil
}

// CreateDonator for the a church using the following business logic:
// Find administrator by uuid and church by id
// If administrator or church not found -->
//		return error
// If new donator does not include firstname,lastname and [phone and address] or email -->
//  	return err
// If new donator has email --> case0
//  Find donator by email
//  If donator found and does not belong to church --> add donator to church donators
//  (any donations made in past are already linked by uuid)
//  Else If donator found and belongs to church --> return error
// 	Else If donator not found --> add new donator and add donator to church donators
// Else If new donator does not have email --> case1
//  Find donator by first name,last name,phone,and address
//  If donator found and belongs to church--> return error
//	Else if donator found and does not belong to church --> add donator to church donators
//  (any donations made in past are already linked by uuid)
//  Else if donator not found --> add new donator and add donator to church donators
func (uc *Usecase) CreateDonator(a *Administrator, c *church.Church, d *donator.Donator) (res *donator.Donator, e error) {
	byEmail := false
	cRead, err := uc.Services.Church.Read(c, byEmail)
	if err != nil {
		return nil, err
	}
	aRead, err := uc.Services.Administrator.Read(a)
	if err != nil {
		return nil, err
	}

	if aRead == nil {
		return nil, ErrAdministratorNotFound
	}
	if cRead == nil {
		return nil, church.ErrChurchNotFound
	}
	if (d.Firstname == nil && d.Lastname == nil) && (d.Email == nil || (d.Phone == nil && d.Address == nil)) {
		return nil, ErrAdministratorDonatorIncompleteDetails
	}

	if d.Email != nil {
		//case0
		byEmail := true
		res, e = uc.createDonator(d, cRead, byEmail)
	} else {
		//case1
		byEmail := false
		res, e = uc.createDonator(d, cRead, byEmail)
	}

	if e != nil {
		return nil, e
	}

	return res, nil

}
func (uc *Usecase) createDonator(d *donator.Donator, c *church.Church, byEmail bool) (res *donator.Donator, e error) {
	var dRead *donator.Donator
	if !byEmail {
		dRead, e = uc.Services.Donator.ReadWithFilter(
			d,
			&donator.Donator{
				Firstname: d.Firstname,
				Lastname:  d.Lastname,
				Phone:     d.Phone,
				Address:   d.Address,
			},
		)
	} else {
		dRead, e = uc.Services.Donator.Read(d, byEmail)
	}

	if e != nil {
		return nil, e
	}

	var cDonators church.Donators
	if c.Donators == nil {
		cDonators = make(church.Donators)
	} else {
		cDonators = *c.Donators
	}
	// Donator found
	if dRead != nil {
		//belongs to church--> return error
		if _, ok := cDonators[*dRead.ID]; ok {
			return nil, church.ErrChurchDonatorExists
		}
		// does not belong to church --> add donator to church donators
		// (any donations made in past are already linked by uuid)
		cDonators[*dRead.ID] = &church.Donator{UUID: dRead.UUID}
		cUpdate := church.Church{Email: c.Email, Donators: &cDonators}
		byEmail := true
		if err := uc.Services.Church.Update(&cUpdate, byEmail); err != nil {
			return nil, err
		}

		donationCount := int64(0)
		(*dRead.Churches)[*c.ID] = &donator.Church{
			Donationcount: &donationCount,
			Firstdonation: nil,
		}
		byEmail = false
		if err := uc.Services.Donator.Update(dRead, byEmail); err != nil {
			return nil, err
		}

		return dRead, nil
	}

	// donator not found --> add new donator and add donator to church donators
	donationCount := int64(0)
	d.Churches = &donator.Churches{
		*c.ID: &donator.Church{
			Donationcount: &donationCount,
			Firstdonation: nil,
		},
	}
	if err := uc.Services.Donator.Create(d); err != nil {
		return nil, err
	}

	cDonators[*d.ID] = &church.Donator{UUID: d.UUID}
	cUpdate := church.Church{ID: c.ID, Email: c.Email, Donators: &cDonators}
	if err := uc.Services.Church.Update(&cUpdate, byEmail); err != nil {
		return nil, err
	}

	return d, nil
}

// ReadDonators for an administrator's church using the following business logic:
// If administrator does not have belong to the church return error .
// Else read donators for church and return donators
func (uc *Usecase) ReadDonators(a *Administrator, c *church.Church) (res map[int64]*donator.Donator, e error) {
	aRead, err := uc.Services.Administrator.Read(a)
	if err != nil {
		return nil, err
	}
	if aRead == nil {
		return nil, ErrAdministratorNotFound
	}
	if aRead.Churches == nil {
		return nil, ErrAdministratorNoChurches
	}

	_, ok := (*aRead.Churches)[*c.ID]
	if !ok {
		return nil, ErrAdministratorDoesNotBelongToChurch
	}
	byEmail := false
	cRead, err := uc.Services.Church.Read(c, byEmail)
	if err != nil {
		return nil, err
	}

	dd, err := uc.Services.Donator.ReadMultiple(cRead.Donators.Keys())
	if err != nil {
		return nil, err
	}

	res = make(map[int64]*donator.Donator, len(dd))
	for _, elem := range dd {
		res[*elem.ID] = elem
	}

	return res, e
}

// EditDonator using the following business logic:
// If Churches,Accountstatement were to be updated ---> return field not editable error
// Validate if donator belongs to the church that administrator is part of ---> return err on validation fail
// If updating a donator creates a non-unique donator ---> Return err
// If email updated, send message to bridge
func (uc *Usecase) EditDonator(a *Administrator, c *church.Church, d *donator.Donator, dID int64) (res *donator.Donator, e error) {
	if d != nil && (d.Churches != nil || d.Accountstatement != nil) {
		return nil, donator.ErrDonatorFieldNotEditable
	}

	aRead, err := uc.Services.Administrator.Read(a)
	if err != nil {
		return nil, err
	}
	if aRead == nil {
		return nil, ErrAdministratorNotFound
	}
	if aRead.Churches == nil {
		return nil, ErrAdministratorNoChurches
	}

	_, ok := (*aRead.Churches)[*c.ID]
	if !ok {
		return nil, ErrAdministratorDoesNotBelongToChurch
	}
	byEmail := false
	cRead, err := uc.Services.Church.Read(c, byEmail)
	if err != nil {
		return nil, err
	}

	if cRead.Donators == nil {
		return nil, church.ErrChurchDonatorDoesNotExists
	}
	if _, ok := (*cRead.Donators)[dID]; !ok {
		return nil, church.ErrChurchDonatorDoesNotExists
	}
	if d == nil {
		return nil, uc.deleteChurchDonator(cRead, dID)
	}

	if dRead, err := uc.Services.Donator.ReadWithFilter(
		d,
		&donator.Donator{
			Firstname: d.Firstname,
			Lastname:  d.Lastname,
			Phone:     d.Phone,
			Address:   d.Address,
			Email:     d.Email,
		},
	); err != nil {
		return nil, err
	} else if dRead != nil {
		return nil, donator.ErrDonatorExists
	}

	if d.Email != nil {
		// Send message to bridge
	}

	if err := uc.Services.Donator.Update(d, byEmail); err != nil {
		return nil, err
	}

	d.Accountstatement = nil
	d.Churches = nil

	return d, e
}
func (uc *Usecase) deleteChurchDonator(c *church.Church, dID int64) error {
	cDonators := (*c.Donators)
	delete(cDonators, dID)
	c.Donators = &cDonators

	byEmail := false
	if err := uc.Services.Church.Update(c, byEmail); err != nil {
		return err
	}

	dRead, err := uc.Services.Donator.Read(&donator.Donator{ID: &dID}, byEmail)
	if err != nil {
		return err
	}

	if dRead.Churches == nil {
		return donator.ErrDonatorInternal
	}
	if dRead.Email == nil && len((*dRead.Churches).Keys()) == 1 {
		return uc.Services.Donator.Delete(&donator.Donator{ID: dRead.ID}, byEmail)
	}

	return nil
}

// CreateDonation allows administrator to add donation for a donator to a church
// using the following business logic:
// Find administrator who created church
// If administrator (free usage limit > 500) and (subscription plan type = free)
// 	--> return err to prevent creating a donation and notify that creator
// 		administrator to upgrade to standard plan.
// Else --> create two transactions (credit toward donator account and debit towards church account)
func (uc *Usecase) CreateDonation(a *Administrator, donation *transaction.Donation) (res *transaction.Donation, e error) {
	// c *church.Church, d *donator.Donator,
	byEmail := false
	cRead, _ := uc.Services.Church.Read(&church.Church{ID: donation.ChurchID}, byEmail)
	dRead, _ := uc.Services.Donator.Read(&donator.Donator{ID: donation.DonatorID}, byEmail)

	cAdministrators := cRead.Administrators
	if cAdministrators == nil {
		return nil, church.ErrChurchAdministrators
	}

	if !uc.administratorIsPartOfChurch(a, *cAdministrators) {
		return nil, church.ErrChurchAdministratorDoesNotBelong
	}

	creatorUUID, err := uc.findCreatorAdministrator(*cAdministrators)
	if err != nil {
		return nil, err
	}

	creator, err := uc.Services.Administrator.Read(&Administrator{UUID: &creatorUUID})
	if err != nil {
		return nil, err
	} else if creator == nil {
		return nil, ErrAdministratorNotFound
	}

	if *donation.Amount < 0 {
		return uc.refundDonation(creator, cRead, dRead, donation)
	}

	return uc.createDonation(creator, cRead, dRead, donation)

}
func (uc *Usecase) administratorIsPartOfChurch(a *Administrator, cAdministrators church.Administrators) bool {
	aUUID := a.UUID
	if _, ok := cAdministrators[*aUUID]; ok {
		return true
	}
	return false
}

// reuturns the creator administrator uuid
func (uc *Usecase) findCreatorAdministrator(ca church.Administrators) (string, error) {
	for k, v := range ca {
		if *v.Role == domain.Creator {
			return k, nil
		}
	}
	return "", church.ErrChurchCreatorDoesNotExists
}
func (uc *Usecase) updateChurchAccountStatement(c *church.Church, donation *transaction.Donation) error {
	closingBalance := *c.Accountstatement.Closingbalance
	closingBalance += *donation.Amount
	now := time.Now()
	c.Accountstatement.Date = &now

	byEmail := false
	return uc.Services.Church.Update(c, byEmail)
}
func (uc *Usecase) updateDonatorAccountStatement(c *church.Church, d *donator.Donator, donation *transaction.Donation) error {
	closingBalance := float64(0)
	if d.Accountstatement != nil && d.Accountstatement.Closingbalance != nil {
		closingBalance = *d.Accountstatement.Closingbalance
	}

	closingBalance += *donation.Amount
	now := time.Now()
	accountStatement := &domain.AccountStatement{Closingbalance: &closingBalance, Date: &now}
	d.Accountstatement = accountStatement

	donationCount := (*d.Churches)[*c.ID].Donationcount
	if donationCount == nil || *donationCount == 0 {
		dc := int64(1)
		donationCount = &dc
		firstDonation := time.Now()
		(*d.Churches)[*c.ID].Firstdonation = &firstDonation
	} else {
		*donationCount++
	}

	(*d.Churches)[*c.ID].Donationcount = donationCount

	byEmail := false
	return uc.Services.Donator.Update(d, byEmail)
}
func (uc *Usecase) createDonation(creator *Administrator, c *church.Church, d *donator.Donator, donation *transaction.Donation) (*transaction.Donation, error) {
	creatorSubscriptionPlanType := *creator.Subscription.Type
	freeUsageLimitCount := *creator.Subscription.Freeusagelimitcount

	if creatorSubscriptionPlanType == FreePlan && freeUsageLimitCount > 500 {
		return nil, ErrAdministratorFreeUsageLimitReached
	}

	creditType := transaction.Credit
	debitType := transaction.Debit

	createdAt, error := time.Parse("2006-01-02T15:04:05.999Z", *donation.Date)
	if error != nil {
		return nil, error
	}
	updatedat := time.Now().UTC()

	donatorAccountType := transaction.Donator
	churchAccountType := transaction.Church

	donation.Account = &donatorAccountType
	donationTowardsChurch := &transaction.Donation{
		DonatorID: d.ID,
		ChurchID:  c.ID,
		Amount:    donation.Amount,
		Type:      donation.Type,
		Currency:  donation.Currency,
		Account:   &churchAccountType,
		Category:  donation.Category,
		Details:   donation.Details,
		Date:      donation.Date,
	}

	creditTxn := &transaction.Transaction{
		DonatorID: d.ID,
		ChurchID:  c.ID,
		Amount:    donation.Amount,
		Type:      &creditType,
		Donation:  donation,
		CreatedAt: &createdAt,
		Updatedat: &updatedat,
	}
	debitTxn := &transaction.Transaction{
		DonatorID: d.ID,
		ChurchID:  c.ID,
		Amount:    donation.Amount,
		Type:      &debitType,
		Donation:  donationTowardsChurch,
		CreatedAt: &createdAt,
		Updatedat: &updatedat,
	}

	if err := uc.Services.Transaction.Create(creditTxn); err != nil {
		return nil, err
	} else if err := uc.Services.Transaction.Create(debitTxn); err != nil {
		return nil, err
	}

	freeUsageLimitCount++
	creator.Subscription.Freeusagelimitcount = &freeUsageLimitCount
	if err := uc.Services.Administrator.Update(creator); err != nil {
		return nil, err
	}

	if err := uc.updateChurchAccountStatement(c, donation); err != nil {
		return nil, err
	} else if err := uc.updateDonatorAccountStatement(c, d, donation); err != nil {
		return nil, err
	}

	return donation, nil
}
func (uc *Usecase) refundDonation(creator *Administrator, c *church.Church, d *donator.Donator, donation *transaction.Donation) (*transaction.Donation, error) {
	amount := *donation.Amount * -1
	donation.Amount = &amount

	creatorSubscriptionPlanType := *creator.Subscription.Type
	freeUsageLimitCount := *creator.Subscription.Freeusagelimitcount

	if creatorSubscriptionPlanType == FreePlan && freeUsageLimitCount > 500 {
		return nil, ErrAdministratorFreeUsageLimitReached
	}

	creditType := transaction.Credit
	debitType := transaction.Debit

	createdAt, error := time.Parse("2006-01-02T15:04:05.999Z", *donation.Date)
	if error != nil {
		return nil, error
	}
	updatedat := time.Now().UTC()

	donatorAccountType := transaction.Donator
	churchAccountType := transaction.Church

	donation.Account = &donatorAccountType
	donationTowardsChurch := &transaction.Donation{
		ChurchID:  c.ID,
		DonatorID: d.ID,
		Amount:    &amount,
		Type:      donation.Type,
		Currency:  donation.Currency,
		Account:   &churchAccountType,
		Category:  donation.Category,
		Details:   donation.Details,
		Date:      donation.Date,
	}

	creditTxn := &transaction.Transaction{
		DonatorID: d.ID,
		ChurchID:  c.ID,
		Amount:    &amount,
		Type:      &creditType,
		Donation:  donationTowardsChurch,
		CreatedAt: &createdAt,
		Updatedat: &updatedat,
	}
	debitTxn := &transaction.Transaction{
		DonatorID: d.ID,
		ChurchID:  c.ID,
		Amount:    &amount,
		Type:      &debitType,
		Donation:  donation,
		CreatedAt: &createdAt,
		Updatedat: &updatedat,
	}

	if err := uc.Services.Transaction.Create(creditTxn); err != nil {
		return nil, err
	} else if err := uc.Services.Transaction.Create(debitTxn); err != nil {
		return nil, err
	}

	if err := uc.updateChurchAccountStatement(c, donation); err != nil {
		return nil, err
	} else if err := uc.updateDonatorAccountStatement(c, d, donation); err != nil {
		return nil, err
	}

	return donation, nil
}

// ReadDonations made to a specific church
func (uc *Usecase) ReadDonations(a *Administrator, c *church.Church) (res map[int64][]*transaction.Donation, e error) {
	aRead, err := uc.Services.Administrator.Read(a)
	if err != nil {
		return nil, err
	}
	if aRead == nil {
		return nil, ErrAdministratorNotFound
	}
	if aRead.Churches == nil {
		return nil, ErrAdministratorNoChurches
	}

	_, ok := (*aRead.Churches)[*c.ID]
	if !ok {
		return nil, ErrAdministratorDoesNotBelongToChurch
	}
	byEmail := false
	cRead, err := uc.Services.Church.Read(c, byEmail)
	if err != nil {
		return nil, err
	}

	tt, err := uc.Services.Transaction.ReadMultiple(cRead.Donators.Keys(), *cRead.ID, nil, nil)
	if err != nil {
		return nil, err
	}
	res = extractDonations(tt)

	return res, e
}
func extractDonations(tt []*transaction.Transaction) (res map[int64][]*transaction.Donation) {
	donatorAccountTT := extractDonatorAccountTransactions(tt)

	res = make(map[int64][]*transaction.Donation, len(donatorAccountTT))
	for donatorID, tt := range donatorAccountTT {
		donatorDonationsMap := mapDonatorDonationsByCreatedAtTypeCategory(tt)

		values := extractValues(donatorDonationsMap)

		res[donatorID] = values
	}

	return res

}

// EditDonation made to a specific chruch by a specific donator using the following business logic:
// Find the donation transactions where donation date is equal to transaction updatedat
// If donation amount is negative and absolute amount is greater than donator acount debit + credit transaction amounts
//		-- -> return error updated amount results in negative amount donation
// Else Create two new transactions are created (a credit towards donator account and a debit
// towards the church account ). For these two new transactions, the updatedat
// field reflects the time at which update was applied and the createdat field
// matches the original donation.
func (uc *Usecase) EditDonation(a *Administrator, donation *transaction.Donation) (res *transaction.Donation, e error) {
	aRead, err := uc.Services.Administrator.Read(a)
	if err != nil {
		return nil, err
	}
	if aRead == nil {
		return nil, ErrAdministratorNotFound
	}
	if aRead.Churches == nil {
		return nil, ErrAdministratorNoChurches
	}

	_, ok := (*aRead.Churches)[*donation.ChurchID]
	if !ok {
		return nil, ErrAdministratorDoesNotBelongToChurch
	}

	donationTime, err := time.Parse("2006-01-02T15:04:05.999Z", *donation.Date)
	if err != nil {
		return nil, err
	}

	tt, err := uc.Services.Transaction.Read(*donation.DonatorID, *donation.ChurchID, &transaction.TimeRange{Lower: &donationTime, Upper: &donationTime}, nil)
	if err != nil {
		return nil, err
	}

	if *donation.Amount < 0 {
		var absoluteAmount float64 = -1 * *donation.Amount
		if (sumOfCreditsMinusDebitsForDonatorAccountTransaction(tt) - absoluteAmount) < 0 {
			return nil, transaction.ErrTransactionInternal
		}
	}

	return uc.CreateDonation(aRead, donation)
}

func sumOfCreditsMinusDebitsForDonatorAccountTransaction(tt []*transaction.Transaction) float64 {
	var sum float64 = 0
	for _, txn := range tt {
		if *txn.Donation.Account == transaction.Donator && *txn.Type == transaction.Credit {
			sum += *txn.Amount
		} else if *txn.Donation.Account == transaction.Donator && *txn.Type == transaction.Debit {
			sum -= *txn.Amount
		}
	}
	return sum
}

// DonationStatementsReport for one or all donators of a specific church is created using the following business logic:
func (uc *Usecase) DonationStatementsReport(a *Administrator, dr *transaction.DonationReport) (res *transaction.DonationReport, e error) {
	byEmail := false
	cRead, _ := uc.Services.Church.Read(&church.Church{ID: dr.ChurchID}, byEmail)
	if cRead.Donators == nil {
		return nil, church.ErrChurchDonatorDoesNotExists
	}
	donatorids := cRead.Donators.Keys()

	comparatorOperator := "="
	fieldName := "donation.category"
	values := make([]interface{}, len(dr.DonationCategories))
	for i, category := range dr.DonationCategories {
		values[i] = category
	}
	donationCategoriesFilter := createSingleFieldFilter(fieldName, values, comparatorOperator)

	transactions, err := uc.Services.Transaction.ReadMultipleWithFilter(donatorids, *dr.ChurchID, dr.TimeRange, nil, &donationCategoriesFilter)
	if err != nil {
		return nil, err
	}

	donatorAccountTT := extractDonatorAccountTransactions(transactions)

	donationsMap := make(map[int64]map[string]*transaction.Donation)
	for donatorID, tt := range donatorAccountTT {
		donatorDonationsMap := mapDonatorDonationsByCreatedAtTypeCategory(tt)
		donationsMap[donatorID] = donatorDonationsMap
	}

	aggregatedDonations := aggregateDonationsByFilter(donationsMap, dr.SumFilter)
	dr.Donations = aggregatedDonations

	summedByCategories := summedByCategories(aggregatedDonations)
	dr.DonationsSum = summedByCategories

	res = dr
	return res, e
}

// ChurchDonationReport for all donators of a specific church is created by using the following business logic:
func (uc *Usecase) ChurchDonationReport(a *Administrator, dr *transaction.DonationReport) (res *transaction.DonationReport, e error) {
	byEmail := false
	cRead, _ := uc.Services.Church.Read(&church.Church{ID: dr.ChurchID}, byEmail)
	if cRead.Donators == nil {
		return nil, church.ErrChurchDonatorDoesNotExists
	}
	donatorids := cRead.Donators.Keys()

	var donationCategoriesFilter *domain.Filter
	if dr.DonationCategories != nil {
		comparatorOperator := "="
		fieldName := "donation.category"
		values := make([]interface{}, len(dr.DonationCategories))
		for i, category := range dr.DonationCategories {
			values[i] = category
		}
		f := createSingleFieldFilter(fieldName, values, comparatorOperator)
		donationCategoriesFilter = &f
	}

	transactions, err := uc.Services.Transaction.ReadMultipleWithFilter(donatorids, *dr.ChurchID, dr.TimeRange, nil, donationCategoriesFilter)
	if err != nil {
		return nil, err
	}

	donatorAccountTT := extractDonatorAccountTransactions(transactions)

	donationsMap := make(map[int64]map[string]*transaction.Donation)
	for donatorID, tt := range donatorAccountTT {
		donatorDonationsMap := mapDonatorDonationsByCreatedAtTypeCategory(tt)
		donationsMap[donatorID] = donatorDonationsMap
	}

	aggregatedDonations := aggregateDonationsByFilter(donationsMap, dr.SumFilter)
	dr.Donations = aggregatedDonations

	summedByCategories := summedByCategories(aggregatedDonations)
	dr.DonationsSum = summedByCategories

	res = dr
	return res, e
}

// DonationStatementReport for one or many donators of a specific church is created by using the following business logic:
func (uc *Usecase) DonationStatementReport(a *Administrator, dr *transaction.DonationReport) (res *transaction.DonationReport, e error) {

	byEmail := false
	cRead, _ := uc.Services.Church.Read(&church.Church{ID: dr.ChurchID}, byEmail)
	if cRead.Donators == nil {
		return nil, church.ErrChurchDonatorDoesNotExists
	}
	donatorids := dr.DonatorIDs

	var donationCategoriesFilter *domain.Filter
	if dr.DonationCategories != nil {
		comparatorOperator := "="
		fieldName := "donation.category"
		values := make([]interface{}, len(dr.DonationCategories))
		for i, category := range dr.DonationCategories {
			values[i] = category
		}
		f := createSingleFieldFilter(fieldName, values, comparatorOperator)
		donationCategoriesFilter = &f
	}

	transactions, err := uc.Services.Transaction.ReadMultipleWithFilter(donatorids, *dr.ChurchID, dr.TimeRange, nil, donationCategoriesFilter)
	if err != nil {
		return nil, err
	}

	donatorAccountTT := extractDonatorAccountTransactions(transactions)

	donationsMap := make(map[int64]map[string]*transaction.Donation)
	for donatorID, tt := range donatorAccountTT {
		donatorDonationsMap := mapDonatorDonationsByCreatedAtTypeCategory(tt)
		donationsMap[donatorID] = donatorDonationsMap
	}

	aggregatedDonations := aggregateDonationsByFilter(donationsMap, dr.SumFilter)
	dr.Donations = aggregatedDonations

	//TODO summed by categories logic is not correct
	summedByCategories := summedByCategoriesForEachDonator(aggregatedDonations)
	dr.DonationsSum = summedByCategories

	res = dr
	return res, e
}

// Extract Donator Account transactions from an array of transactions.
func extractDonatorAccountTransactions(tt []*transaction.Transaction) (res map[int64][]*transaction.Transaction) {
	donatorAccountTT := make(map[int64][]*transaction.Transaction, len(tt))
	for _, txn := range tt {
		if *txn.Donation.Account == transaction.Donator {
			donatorAccountTT[*txn.DonatorID] = append(donatorAccountTT[*txn.DonatorID], txn)
		}

	}

	res = donatorAccountTT
	return res
}
func extractValues(d map[string]*transaction.Donation) (res []*transaction.Donation) {
	for _, value := range d {
		if *value.Amount != 0 {
			res = append(res, value)
		}
	}

	return res
}

// Extract donations from an array of donator account transactions and map
// donations by the time they were created at,the type(cash,online,check,etc), and category.
func mapDonatorDonationsByCreatedAtTypeCategory(tt []*transaction.Transaction) (res map[string]*transaction.Donation) {

	donatorDonationsMap := make(map[string]*transaction.Donation, len(tt))
	for _, txn := range tt {
		key := txn.CreatedAt.String() + "_" + txn.Donation.Type.String() + "_" + *txn.Donation.Category

		if _, donationsExist := donatorDonationsMap[key]; !donationsExist {
			if *txn.Type == transaction.Debit {
				newAmount := *txn.Donation.Amount * -1
				txn.Donation.Amount = &newAmount
			}
			donatorDonationsMap[key] = txn.Donation
			continue
		}

		newAmount := *(donatorDonationsMap[key].Amount)
		if *txn.Type == transaction.Debit {
			newAmount -= *txn.Amount
		} else {
			newAmount += *txn.Amount
		}

		donatorDonationsMap[key].Amount = &newAmount

	}
	res = donatorDonationsMap
	return res
}
func aggregateDonationsByFilter(dd map[int64]map[string]*transaction.Donation, sumFilter *transaction.SumFilter) (res transaction.DonationsMap) {

	// For each donator's donations aggregate by sumFilter.
	res = make(transaction.DonationsMap)
	for donatorID, donations := range dd {
		donatorDonationsMap := aggregateBySumFilter(donations, sumFilter)
		res[donatorID] = donatorDonationsMap
	}

	return res
}
func aggregateBySumFilter(d map[string]*transaction.Donation, sumFilter *transaction.SumFilter) (res transaction.DonatorDonationsMap) {

	a := make(transaction.DonatorDonationsMap)

	if sumFilter == nil || sumFilter.TimePeriod == nil {
		for groupingKey, donation := range d {
			a[groupingKey] = append([]*transaction.Donation{}, donation)
		}
		res = a
		return res
	}

	switch *sumFilter.TimePeriod {
	case transaction.Day:
		a = aggregateByDayAndCategory(d)
	case transaction.Week:
		a = aggregateByWeekAndCategory(d)
	case transaction.Month:
		a = aggregateByMonthAndCategory(d)
	}

	if sumFilter.Multiplier != nil {
		res = applyMultiplierToAggregated(a, sumFilter.Multiplier)
	} else {
		res = a
	}

	return res
}
func aggregateByDayAndCategory(d map[string]*transaction.Donation) (res transaction.DonatorDonationsMap) {
	donatorDonationMap := make(map[string]*transaction.Donation)
	for _, donation := range extractValues(d) {
		createdAt, _ := time.Parse("2006-01-02T15:04:05.999Z", *donation.Date)

		key := strconv.Itoa(createdAt.Year()) + "_" + strconv.Itoa(createdAt.YearDay()) + "_" + *donation.Category

		if _, donationsExist := donatorDonationMap[key]; !donationsExist {
			donatorDonationMap[key] = donation
			continue
		}

		newAmount := *(donatorDonationMap[key].Amount)
		newAmount += *donation.Amount

		donatorDonationMap[key].Amount = &newAmount

	}

	res = make(transaction.DonatorDonationsMap)
	for _, donation := range donatorDonationMap {
		createdAt, _ := time.Parse("2006-01-02T15:04:05.999Z", *donation.Date)
		key := strconv.Itoa(createdAt.Year()) + "_" + strconv.Itoa(createdAt.YearDay())

		if _, donationExist := res[key]; !donationExist {
			res[key] = make([]*transaction.Donation, 0)
		}
		res[key] = append(res[key], donation)
	}

	return res
}
func aggregateByWeekAndCategory(d map[string]*transaction.Donation) (res transaction.DonatorDonationsMap) {
	donatorDonationMap := make(map[string]*transaction.Donation)
	for _, donation := range extractValues(d) {
		createdAt, _ := time.Parse("2006-01-02T15:04:05.999Z", *donation.Date)
		year, week := createdAt.ISOWeek()
		key := strconv.Itoa(year) + "_" + strconv.Itoa(week) + "_" + *donation.Category

		if _, donationsExist := donatorDonationMap[key]; !donationsExist {
			donatorDonationMap[key] = donation
			continue
		}

		newAmount := *(donatorDonationMap[key].Amount)
		newAmount += *donation.Amount

		donatorDonationMap[key].Amount = &newAmount

	}

	res = make(transaction.DonatorDonationsMap)
	for _, donation := range donatorDonationMap {
		createdAt, _ := time.Parse("2006-01-02T15:04:05.999Z", *donation.Date)
		year, week := createdAt.ISOWeek()
		key := strconv.Itoa(year) + "_" + strconv.Itoa(week) + "_" + *donation.Category

		if _, donationExist := res[key]; !donationExist {
			res[key] = make([]*transaction.Donation, 0)
		}
		res[key] = append(res[key], donation)
	}

	return res
}
func aggregateByMonthAndCategory(d map[string]*transaction.Donation) (res transaction.DonatorDonationsMap) {
	donatorDonationMap := make(map[string]*transaction.Donation)
	for _, donation := range extractValues(d) {
		createdAt, _ := time.Parse("2006-01-02T15:04:05.999Z", *donation.Date)
		key := strconv.Itoa(createdAt.Year()) + "_" + createdAt.Month().String() + "_" + *donation.Category

		if _, donationsExist := donatorDonationMap[key]; !donationsExist {
			donatorDonationMap[key] = donation
			continue
		}

		newAmount := *(donatorDonationMap[key].Amount)
		newAmount += *donation.Amount

		donatorDonationMap[key].Amount = &newAmount

	}

	res = make(transaction.DonatorDonationsMap)
	for _, donation := range donatorDonationMap {
		createdAt, _ := time.Parse("2006-01-02T15:04:05.999Z", *donation.Date)
		key := strconv.Itoa(createdAt.Year()) + "_" + createdAt.Month().String() + "_" + *donation.Category

		if _, donationExist := res[key]; !donationExist {
			res[key] = make([]*transaction.Donation, 0)
		}
		res[key] = append(res[key], donation)
	}

	return res
}
func applyMultiplierToAggregated(d transaction.DonatorDonationsMap, multiplier *int64) (res transaction.DonatorDonationsMap) {
	// TODO: implement logic here to apply multiplier to aggregated
	res = d
	return res
}
func createSingleFieldFilter(fieldName string, values []interface{}, comparatorOperator string) (res domain.Filter) {
	if values != nil {
		fv := domain.FilterValue{ComparatorOperator: comparatorOperator, Value: values}
		res = domain.Filter{fieldName: []domain.FilterValue{fv}}
	}
	return res
}
func summedByCategories(d transaction.DonationsMap) (res transaction.DonationsSumMap) {
	res = make(transaction.DonationsSumMap)
	for _, donatorDonationsMap := range d {
		for groupingKey, donations := range donatorDonationsMap {
			for _, donation := range donations {

				if _, exists := res[*donation.Category]; !exists {
					res[*donation.Category] = make(transaction.DonationsSum)
					res[*donation.Category][groupingKey] = *donation.Amount
					continue
				}
				if _, amountExists := res[*donation.Category][groupingKey]; !amountExists {
					res[*donation.Category][groupingKey] = *donation.Amount
					continue
				}

				res[*donation.Category][groupingKey] += *donation.Amount
			}
		}
	}

	totalsMap := make(transaction.DonationsSumMap)
	for category, donationsSum := range res {

		for _, amount := range donationsSum {
			if _, exists := totalsMap[category]; !exists {
				totalsMap[category] = make(transaction.DonationsSum)
				totalsMap[category]["total"] = amount
				continue
			}
			if _, amountExists := totalsMap[category]["total"]; !amountExists {
				totalsMap[category]["total"] = amount
				continue
			}

			totalsMap[category]["total"] += amount
		}

	}
	for category, totalDonationsSum := range totalsMap {
		res[category]["total"] = totalDonationsSum["total"]
	}

	return res
}

func summedByCategoriesForEachDonator(d transaction.DonationsMap) (res transaction.DonationsSumMap) {
	tmpMap := make(transaction.DonationsSumMap)
	for donatorID, donatorDonationsMap := range d {
		for originalGroupingKey, donations := range donatorDonationsMap {
			groupingKey := strconv.FormatInt(donatorID, 10) + "_" + originalGroupingKey
			for _, donation := range donations {

				if _, exists := tmpMap[*donation.Category]; !exists {
					tmpMap[*donation.Category] = make(transaction.DonationsSum)
					tmpMap[*donation.Category][groupingKey] = *donation.Amount
					continue
				}
				if _, amountExists := tmpMap[*donation.Category][groupingKey]; !amountExists {
					tmpMap[*donation.Category][groupingKey] = *donation.Amount
					continue
				}

				tmpMap[*donation.Category][groupingKey] += *donation.Amount
			}
		}
	}

	res = make(transaction.DonationsSumMap)
	for category, donationsSum := range tmpMap {

		for groupingKey, amount := range donationsSum {
			donatorID := strings.Split(groupingKey, "_")[0]
			if _, exists := res[category]; !exists {
				res[category] = make(transaction.DonationsSum)
				res[category][donatorID+"_total"] = amount
				continue
			}
			if _, amountExists := res[category][donatorID+"_total"]; !amountExists {
				res[category][donatorID+"_total"] = amount
				continue
			}

			res[category][donatorID+"_total"] += amount
		}

	}

	return res
}
