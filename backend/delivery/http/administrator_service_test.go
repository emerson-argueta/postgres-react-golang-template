package http_test

import (
	"bytes"
	"errors"
	"log"
	"reflect"
	"strings"
	"testing"

	"emersonargueta/m/v1/delivery/http"
	"emersonargueta/m/v1/delivery/middleware"
	"emersonargueta/m/v1/domain/administrator"
	mockadministrator "emersonargueta/m/v1/domain/administrator/mock"
	"emersonargueta/m/v1/domain/church"
	mockchurch "emersonargueta/m/v1/domain/church/mock"
	"emersonargueta/m/v1/domain/donator"
	mockdonator "emersonargueta/m/v1/domain/donator/mock"
	"emersonargueta/m/v1/domain/transaction"
	mocktransaction "emersonargueta/m/v1/domain/transaction/mock"
	"emersonargueta/m/v1/user"
	mockuser "emersonargueta/m/v1/user/mock"

	"golang.org/x/crypto/bcrypt"
)

// AdministratorHandler represents a test wrapper for http.AdminHandler.
type AdministratorHandler struct {
	*http.AdministratorHandler
	AdministratorService mockadministrator.AdministratorService
	UserService          mockuser.UserService
	ChurchService        mockchurch.ChurchService
	DonatorService       mockdonator.DonatorService
	TransactionService   mocktransaction.TransactionService

	LogOutput bytes.Buffer
}

// NewAdministratorHandler returns a new instance of AdministratorHandler.
func NewAdministratorHandler() *AdministratorHandler {
	h := &AdministratorHandler{AdministratorHandler: http.NewAdministratorHandler()}

	h.AdministratorHandler.Usecase.Services.Administrator = &h.AdministratorService
	h.AdministratorHandler.Usecase.Services.User = &h.UserService
	h.AdministratorHandler.Usecase.Services.Church = &h.ChurchService
	h.AdministratorHandler.Usecase.Services.Donator = &h.DonatorService

	h.AdministratorHandler.Authorization.Administrator.Usecase.Services.Administrator = &h.AdministratorService

	h.Logger = log.New(VerboseWriter(&h.LogOutput), "", log.LstdFlags)
	return h
}

func TestAdminService_Register(t *testing.T) {
	t.Run("OK", testAdministratorService_Register)
	t.Run("ErrAdminExists", testAdministratorService_Register_ErrAdministratorExists)
}

// Ensure service can register a new administrator.
func testAdministratorService_Register(t *testing.T) {
	s, c := MustOpenServerClient()
	defer s.Close()

	// Mock services used by register business logic.
	s.Handler.AdministratorHandler.UserService.RegisterFn = func(u *user.User) error {
		if *u.Email != "test@test.com" {
			t.Fatalf("unexpected user email: %v", u.Email)
		} else if err := bcrypt.CompareHashAndPassword([]byte(*u.Password), []byte("test1234")); err != nil {
			t.Fatalf("unexpected user password: %v", u.Password)
		}
		uuid := "TEST"
		u.UUID = &uuid
		return nil
	}
	s.Handler.AdministratorHandler.AdministratorService.CreateFn = func(a *administrator.Administrator) error {
		if *a.UUID != "TEST" {
			t.Fatalf("unexpected administrator uuid: %v", a.UUID)
		} else if *a.Firstname != "testFirstname" {
			t.Fatalf("unexpected administrator firstname: %v", a.Firstname)
		} else if *a.Lastname != "testLastname" {
			t.Fatalf("unexpected administrator lastname: %v", a.Lastname)
		} else if *a.Email != "test@test.com" {
			t.Fatalf("unexpected administrator email: %v", a.Email)
		} else if err := bcrypt.CompareHashAndPassword([]byte(*a.Password), []byte("test1234")); err != nil {
			t.Fatalf("unexpected administrator password: %v", a.Password)
		}
		return nil
	}

	s.Handler.AdministratorHandler.AdministratorService.CreateManagementSessionFn = func() error { return nil }

	s.Handler.AdministratorHandler.AdministratorService.EndManagementSessionFn = func() error { return nil }

	s.Handler.AdministratorHandler.UserService.RetrieveFn = func(u *user.User, byEmail bool) (*user.User, error) {
		return nil, nil
	}
	s.Handler.AdministratorHandler.AdministratorService.ReadFn = func(a *administrator.Administrator) (*administrator.Administrator, error) {
		return nil, nil
	}

	firstname_a := "testFirstname"
	lastname_a := "testLastname"
	email_a := "test@test.com"
	password_a := "test1234"
	a := administrator.Administrator{
		Firstname: &firstname_a,
		Lastname:  &lastname_a,
		Email:     &email_a,
		Password:  &password_a,
	}

	// Create admin.
	if err := c.AdministratorService().Register(&a); err != nil {
		t.Fatal(err)
	}
	if *a.Firstname != "testFirstname" {
		t.Fatalf("unexpected administrator firstname: %v", a.Firstname)
	} else if *a.Lastname != "testLastname" {
		t.Fatalf("unexpected administrator lastname: %v", a.Lastname)
	} else if *a.Email != "test@test.com" {
		t.Fatalf("unexpected administrator email: %v", a.Email)
	}
}

// Ensure service returns an error if admin already exists.
func testAdministratorService_Register_ErrAdministratorExists(t *testing.T) {
	s, c := MustOpenServerClient()
	defer s.Close()

	s.Handler.AdministratorHandler.AdministratorService.CreateManagementSessionFn = func() error { return nil }

	s.Handler.AdministratorHandler.AdministratorService.EndManagementSessionFn = func() error { return nil }

	s.Handler.AdministratorHandler.UserService.RetrieveFn = func(u *user.User, byEmail bool) (*user.User, error) {
		uuid := "TEST"
		return &user.User{UUID: &uuid}, nil
	}
	s.Handler.AdministratorHandler.AdministratorService.ReadFn = func(a *administrator.Administrator) (*administrator.Administrator, error) {
		uuid := "TEST"
		return &administrator.Administrator{UUID: &uuid}, nil
	}

	email_a := "test@test.com"
	password_a := "test1234"
	if err := c.AdministratorService().Register(
		&administrator.Administrator{Email: &email_a, Password: &password_a},
	); err != administrator.ErrAdministratorExists {
		t.Fatal(err)
	}
}

func TestAdministratorService_Login(t *testing.T) {
	t.Run("OK", testAdministratorService_Login)
	t.Run("IncorrectPassword", testAdministratorService_Login_IncorrectPassword)
	t.Run("NotFound", testAdministratorService_Login_NotFound)
}

func testAdministratorService_Login(t *testing.T) {
	s, c := MustOpenServerClient()
	defer s.Close()

	// Mock service used by login business logic.
	s.Handler.AdministratorHandler.UserService.RetrieveFn = func(u *user.User, byEmail bool) (*user.User, error) {
		if *u.Email != "test@test.com" {
			t.Fatalf("unexpected email: %v", u.Email)
		} else if *u.Password != "test1234" {
			t.Fatalf("unexpected password: %v", u.Password)
		}
		uuid := "TEST"
		u.UUID = &uuid
		return u, nil
	}
	s.Handler.AdministratorHandler.AdministratorService.ReadFn = func(a *administrator.Administrator) (*administrator.Administrator, error) {
		if *a.Email != "test@test.com" {
			t.Fatalf("unexpected email: %v", a.Email)
		} else if *a.Password != "test1234" {
			t.Fatalf("unexpected password: %v", a.Password)
		}
		return a, nil
	}

	s.Handler.AdministratorHandler.AdministratorService.CreateManagementSessionFn = func() error { return nil }

	s.Handler.AdministratorHandler.AdministratorService.EndManagementSessionFn = func() error { return nil }

	// Retrieve admin.
	email_a := "test@test.com"
	password_a := "test1234"
	_, err := c.AdministratorService().Login(&administrator.Administrator{Email: &email_a, Password: &password_a})
	if err != nil {
		t.Fatal(err)
	}
}

func testAdministratorService_Login_IncorrectPassword(t *testing.T) {
	s, c := MustOpenServerClient()
	defer s.Close()

	// Mock service used by login usecase.
	s.Handler.AdministratorHandler.UserService.RetrieveFn = func(u *user.User, byEmail bool) (*user.User, error) {
		uuid := "TEST"
		u.UUID = &uuid
		password, _ := bcrypt.GenerateFromPassword([]byte("XXXXXXXX"), bcrypt.DefaultCost)
		passwordstr := string(password)
		u.Password = &passwordstr
		return u, nil
	}
	s.Handler.AdministratorHandler.AdministratorService.ReadFn = func(a *administrator.Administrator) (*administrator.Administrator, error) {
		return a, nil
	}

	s.Handler.AdministratorHandler.AdministratorService.CreateManagementSessionFn = func() error { return nil }

	s.Handler.AdministratorHandler.AdministratorService.EndManagementSessionFn = func() error { return nil }

	email_a := "test@test.com"
	password_a := "test1234"
	if _, err := c.AdministratorService().Login(&administrator.Administrator{Email: &email_a, Password: &password_a}); err != administrator.ErrAdministratorIncorrectCredentials {
		t.Fatal(err)
	}
}

func testAdministratorService_Login_NotFound(t *testing.T) {
	s, c := MustOpenServerClient()
	defer s.Close()

	// Mock service used by login usecase.
	s.Handler.AdministratorHandler.UserService.RetrieveFn = func(u *user.User, byEmail bool) (*user.User, error) {
		return nil, nil
	}
	s.Handler.AdministratorHandler.AdministratorService.ReadFn = func(a *administrator.Administrator) (*administrator.Administrator, error) {
		return nil, nil
	}

	s.Handler.AdministratorHandler.AdministratorService.CreateManagementSessionFn = func() error { return nil }

	s.Handler.AdministratorHandler.AdministratorService.EndManagementSessionFn = func() error { return nil }

	email_a := "test@test.com"
	password_a := "test1234"
	if _, err := c.AdministratorService().Login(&administrator.Administrator{Email: &email_a, Password: &password_a}); err != administrator.ErrAdministratorNotFound {
		t.Fatal(err)
	}
}

func TestAdministratorService_Authorize(t *testing.T) {
	t.Run("OK", testAdministratorService_Authorize)
}

func testAdministratorService_Authorize(t *testing.T) {
	s, c := MustOpenServerClient()
	defer s.Close()

	// Mock service used by refreshtoken route.
	s.Handler.AdministratorHandler.AdministratorService.ReadFn = func(a *administrator.Administrator) (*administrator.Administrator, error) {
		return a, nil
	}
	uuid_a := "TEST"
	token_a, _ := middleware.GenerateTokenPair(uuid_a, middleware.AccestokenLimit, middleware.RefreshtokenLimit)

	if _, err := c.AdministratorService().Authorize(token_a); err != nil {
		t.Fatal(err)
	}
}

func TestAdministratorService_Read(t *testing.T) {
	t.Run("OK", testAdministratorService_Read)
	t.Run("NotFound", testAdminService_Read_NotFound)
	t.Run("ErrInternal", testAdminService_Read_ErrInternal)
}

// Ensure service can return an admin.
func testAdministratorService_Read(t *testing.T) {
	token, _ := middleware.GenerateTokenPair("TEST", middleware.AccestokenLimit, middleware.RefreshtokenLimit)
	s, c := MustOpenServerClient()
	defer s.Close()

	uuid := "TEST"
	firstname := "FIRST_NAME"
	lastname := "LAST_NAME"
	email := "test@test.com"
	password := "XXXXXXXXX"
	// Mock service.
	s.Handler.AdministratorHandler.AdministratorService.ReadFn = func(a *administrator.Administrator) (*administrator.Administrator, error) {

		return &administrator.Administrator{
			UUID:      &uuid,
			Firstname: &firstname,
			Lastname:  &lastname,
			Email:     &email,
			Password:  &password,
		}, nil
	}
	s.Handler.AdministratorHandler.UserService.RetrieveFn = func(u *user.User, byEmail bool) (*user.User, error) {

		return &user.User{
			UUID:     &uuid,
			Email:    &email,
			Password: &password,
		}, nil
	}

	// Retrieve admin.
	a, err := c.Services.Administrator.Read(token)
	if err != nil {
		t.Fatal(err)
	} else if !reflect.DeepEqual(a, &administrator.Administrator{
		Firstname: &firstname,
		Lastname:  &lastname,
		Email:     &email,
	}) {
		t.Fatalf("unexpected admin: %#v", a)
	}
}

// Ensure service handles fetching a non-existent admin.
func testAdminService_Read_NotFound(t *testing.T) {
	token, _ := middleware.GenerateTokenPair("TEST", middleware.AccestokenLimit, middleware.RefreshtokenLimit)
	s, c := MustOpenServerClient()
	defer s.Close()

	// Mock service.
	s.Handler.AdministratorHandler.AdministratorService.ReadFn = func(a *administrator.Administrator) (*administrator.Administrator, error) {
		return nil, nil
	}

	// Retrieve admin.
	if a, err := c.Services.Administrator.Read(token); err != nil {
		t.Fatal(err)
	} else if a != nil {
		t.Fatal("expected nil admin")
	}
}

// Ensure service returns an internal error if an error occurs.
func testAdminService_Read_ErrInternal(t *testing.T) {
	token, _ := middleware.GenerateTokenPair("TEST", middleware.AccestokenLimit, middleware.RefreshtokenLimit)
	s, c := MustOpenServerClient()
	defer s.Close()

	// Mock service.
	s.Handler.AdministratorHandler.AdministratorService.ReadFn = func(a *administrator.Administrator) (*administrator.Administrator, error) {
		return nil, errors.New("marker")
	}

	// Retrieve admin.
	if _, err := c.Services.Administrator.Read(token); err != administrator.ErrAdministratorInternal {
		t.Fatal(err)
	} else if !strings.Contains(s.Handler.AdministratorHandler.LogOutput.String(), "marker") {
		t.Fatalf("expected log output")
	}
}

func TestAdministratorService_ChurchManagement(t *testing.T) {
	t.Run("OK", testAdministratorService_AddChurch)
	t.Run("OK", testAdministratorService_CreateChurch)
}

func testAdministratorService_AddChurch(t *testing.T) {
	token, _ := middleware.GenerateTokenPair("TEST", middleware.AccestokenLimit, middleware.RefreshtokenLimit)
	s, c := MustOpenServerClient()
	defer s.Close()

	// Mock service.
	s.Handler.AdministratorHandler.AdministratorService.ReadFn = func(a *administrator.Administrator) (*administrator.Administrator, error) {
		return a, nil
	}
	s.Handler.AdministratorHandler.AdministratorService.UpdateFn = func(a *administrator.Administrator) error {
		return nil
	}
	s.Handler.AdministratorHandler.ChurchService.ReadFn = func(c *church.Church, byEmail bool) (res *church.Church, e error) {
		password, _ := bcrypt.GenerateFromPassword([]byte("test1234"), bcrypt.DefaultCost)
		passwordstr := string(password)
		id := int64(1)
		res = &church.Church{}
		res.ID = &id
		res.Password = &passwordstr
		res.Email = c.Email
		return res, nil
	}

	s.Handler.AdministratorHandler.AdministratorService.CreateManagementSessionFn = func() error { return nil }
	s.Handler.AdministratorHandler.ChurchService.CreateManagementSessionFn = func() error { return nil }

	s.Handler.AdministratorHandler.AdministratorService.EndManagementSessionFn = func() error { return nil }
	s.Handler.AdministratorHandler.ChurchService.EndManagementSessionFn = func() error { return nil }

	uuid := "TEST"
	admin := &administrator.Administrator{UUID: &uuid}
	church_email := "testchurch@testchurch.com"
	church_password := "test1234"
	newchurch := &church.Church{Email: &church_email, Password: &church_password}

	if _, err := c.Services.Administrator.AddChurch(token, admin, newchurch); err != nil {
		t.Fatal(err)
	}
}

func testAdministratorService_CreateChurch(t *testing.T) {
	token, _ := middleware.GenerateTokenPair("TEST", middleware.AccestokenLimit, middleware.RefreshtokenLimit)
	s, c := MustOpenServerClient()
	defer s.Close()

	// Mock service.
	s.Handler.AdministratorHandler.AdministratorService.ReadFn = func(a *administrator.Administrator) (*administrator.Administrator, error) {
		return a, nil
	}
	s.Handler.AdministratorHandler.AdministratorService.UpdateFn = func(a *administrator.Administrator) error {
		return nil
	}
	s.Handler.AdministratorHandler.ChurchService.CreateFn = func(c *church.Church) error {
		id := int64(1)
		c.ID = &id
		return nil
	}
	s.Handler.AdministratorHandler.ChurchService.ReadFn = func(c *church.Church, byEmail bool) (res *church.Church, e error) {
		return res, nil
	}

	s.Handler.AdministratorHandler.AdministratorService.CreateManagementSessionFn = func() error { return nil }
	s.Handler.AdministratorHandler.ChurchService.CreateManagementSessionFn = func() error { return nil }

	s.Handler.AdministratorHandler.AdministratorService.EndManagementSessionFn = func() error { return nil }
	s.Handler.AdministratorHandler.ChurchService.EndManagementSessionFn = func() error { return nil }

	uuid := "TEST"
	admin := &administrator.Administrator{UUID: &uuid}
	church_email := "testchurch@testchurch.com"
	church_password := "test1234"
	newchurch := &church.Church{Email: &church_email, Password: &church_password}

	if _, err := c.Services.Administrator.CreateChurch(token, admin, newchurch); err != nil {
		t.Fatal(err)
	}

}

func TestAdministratorService_DonatorManagement(t *testing.T) {
	t.Run("OK", testAdministratorService_CreateDonator)
}

func testAdministratorService_CreateDonator(t *testing.T) {
	token, _ := middleware.GenerateTokenPair("TEST", middleware.AccestokenLimit, middleware.RefreshtokenLimit)
	s, c := MustOpenServerClient()
	defer s.Close()

	// Mock service.
	s.Handler.AdministratorHandler.AdministratorService.ReadFn = func(a *administrator.Administrator) (*administrator.Administrator, error) {
		uuid := "TEST"
		a.UUID = &uuid
		return a, nil
	}
	s.Handler.AdministratorHandler.ChurchService.ReadFn = func(c *church.Church, byEmail bool) (*church.Church, error) {
		id := int64(1)
		c.ID = &id
		return c, nil
	}
	s.Handler.AdministratorHandler.ChurchService.UpdateFn = func(c *church.Church, byEmail bool) error {
		return nil
	}
	s.Handler.AdministratorHandler.DonatorService.CreateFn = func(d *donator.Donator) error {
		id := int64(1)
		d.ID = &id
		return nil
	}
	s.Handler.AdministratorHandler.DonatorService.ReadFn = func(d *donator.Donator, byEmail bool) (res *donator.Donator, e error) {
		id := int64(1)
		d.ID = &id

		d.Churches = &donator.Churches{}

		res = d
		return res, nil
	}
	s.Handler.AdministratorHandler.DonatorService.ReadWithFilterFn = func(d *donator.Donator, fd *donator.Donator) (res *donator.Donator, e error) {
		id := int64(1)
		d.ID = &id

		d.Churches = &donator.Churches{}

		res = d
		return res, nil
	}
	s.Handler.AdministratorHandler.DonatorService.UpdateFn = func(d *donator.Donator, byEmail bool) error {
		return nil
	}

	s.Handler.AdministratorHandler.AdministratorService.CreateManagementSessionFn = func() error { return nil }
	s.Handler.AdministratorHandler.ChurchService.CreateManagementSessionFn = func() error { return nil }
	s.Handler.AdministratorHandler.DonatorService.CreateManagementSessionFn = func() error { return nil }

	s.Handler.AdministratorHandler.AdministratorService.EndManagementSessionFn = func() error { return nil }
	s.Handler.AdministratorHandler.ChurchService.EndManagementSessionFn = func() error { return nil }
	s.Handler.AdministratorHandler.DonatorService.EndManagementSessionFn = func() error { return nil }

	id := int64(1)
	ch := &church.Church{ID: &id}

	firstname := "test_firstname"
	lastname := "test_lastname"
	email := "test_email"
	address := "test_address"
	phone := "test_phone"

	if _, err := c.Services.Administrator.CreateDonator(
		token,
		ch, &donator.Donator{
			Firstname: &firstname,
			Lastname:  &lastname,
			Email:     &email,
		},
	); err != nil {
		t.Fatal(err)
	}

	if _, err := c.Services.Administrator.CreateDonator(
		token,
		ch, &donator.Donator{
			Firstname: &firstname,
			Lastname:  &lastname,
			Address:   &address,
			Phone:     &phone,
		},
	); err != nil {
		t.Fatal(err)
	}

	s.Handler.AdministratorHandler.DonatorService.ReadFn = func(d *donator.Donator, byEmail bool) (res *donator.Donator, e error) {
		return nil, nil
	}
	s.Handler.AdministratorHandler.DonatorService.ReadWithFilterFn = func(d *donator.Donator, fd *donator.Donator) (res *donator.Donator, e error) {
		return nil, nil
	}
	if _, err := c.Services.Administrator.CreateDonator(
		token,
		ch, &donator.Donator{
			Firstname: &firstname,
			Lastname:  &lastname,
			Email:     &email,
		},
	); err != nil {
		t.Fatal(err)
	}

	if _, err := c.Services.Administrator.CreateDonator(
		token,
		ch, &donator.Donator{
			Firstname: &firstname,
			Lastname:  &lastname,
			Address:   &address,
			Phone:     &phone,
		},
	); err != nil {
		t.Fatal(err)
	}
}

func TestAdministratorService_DonationManagement(t *testing.T) {
	t.Run("OK", testAdministratorService_CreateDonation)
}

func testAdministratorService_CreateDonation(t *testing.T) {
	token, _ := middleware.GenerateTokenPair("TEST", middleware.AccestokenLimit, middleware.RefreshtokenLimit)
	s, c := MustOpenServerClient()
	defer s.Close()

	// Mock service.
	s.Handler.AdministratorHandler.AdministratorService.ReadFn = func(a *administrator.Administrator) (*administrator.Administrator, error) {
		return a, nil
	}
	s.Handler.AdministratorHandler.AdministratorService.UpdateFn = func(a *administrator.Administrator) error {
		return nil
	}
	s.Handler.AdministratorHandler.TransactionService.CreateFn = func(txn *transaction.Transaction) error {
		return nil
	}
	s.Handler.AdministratorHandler.ChurchService.UpdateFn = func(c *church.Church, byEmail bool) error {
		return nil
	}
	s.Handler.AdministratorHandler.DonatorService.UpdateFn = func(d *donator.Donator, byEmail bool) error {
		return nil
	}

	s.Handler.AdministratorHandler.AdministratorService.CreateManagementSessionFn = func() error { return nil }
	s.Handler.AdministratorHandler.ChurchService.CreateManagementSessionFn = func() error { return nil }
	s.Handler.AdministratorHandler.DonatorService.CreateManagementSessionFn = func() error { return nil }
	s.Handler.AdministratorHandler.TransactionService.CreateManagementSessionFn = func() error { return nil }

	s.Handler.AdministratorHandler.AdministratorService.EndManagementSessionFn = func() error { return nil }
	s.Handler.AdministratorHandler.ChurchService.EndManagementSessionFn = func() error { return nil }
	s.Handler.AdministratorHandler.DonatorService.EndManagementSessionFn = func() error { return nil }
	s.Handler.AdministratorHandler.TransactionService.EndManagementSessionFn = func() error { return nil }

	id := int64(1)
	ch := &church.Church{ID: &id}
	d := &donator.Donator{ID: &id}

	amount := float64(100.11)
	donationType := transaction.Cash
	currency := "usd"
	accountType := transaction.Donator
	category := "general"
	details := ""
	donation := &transaction.Donation{
		Amount:   &amount,
		Type:     &donationType,
		Currency: &currency,
		Account:  &accountType,
		Category: &category,
		Details:  &details,
	}

	c.Services.Administrator.CreateDonation(token, ch, d, donation)

}

// func TestAdminService_Update(t *testing.T) {
// 	t.Run("OK", testAdminService_Update)
// 	t.Run("ErrAdminNotFound", testAdminService_Update_ErrAdminNotFound)
// 	t.Run("ErrInternal", testAdminService_Update_ErrInternal)
// }

// // Ensure service can set the level of an existing admin.
// func testAdminService_Update(t *testing.T) {
// 	token, _ := middleware.GenerateTokenPair(1, middleware.AccestokenLimit, middleware.RefreshtokenLimit)
// 	s, c := MustOpenServerClient()
// 	defer s.Close()

// 	// Mock service.
// 	s.Handler.AdminHandler.AdminService.UpdateFn = func(a *admin.Admin, byEmail bool) error {
// 		if a.ID != 1 {
// 			t.Fatalf("unexpected admin id: %d", a.ID)
// 		} else if a.Firstname != "FIRST_NAME" {
// 			t.Fatalf("unexpected Firstname: %v", a.Firstname)
// 		} else if a.Lastname != "LAST_NAME" {
// 			t.Fatalf("unexpected Lastname: %v", a.Lastname)
// 		} else if a.Email != "test@test.com" {
// 			t.Fatalf("unexpected email: %v", a.Email)
// 		} else if err := bcrypt.CompareHashAndPassword([]byte(a.Password), []byte("XXXXXXXXX")); err != nil {
// 			t.Fatalf("unexpected password updatefn: %v", a.Password)
// 		}
// 		return nil
// 	}

// 	a := admin.Admin{ID: 1, Firstname: "FIRST_NAME", Lastname: "LAST_NAME", Email: "test@test.com", Password: "XXXXXXXXX", Token: token}
// 	// Set admin level.
// 	err := c.AdminService().Update(&a, true)
// 	if err != nil {
// 		t.Fatal(err)
// 	} else if a.ID != 1 {
// 		t.Fatalf("unexpected admin id: %d", a.ID)
// 	} else if a.Firstname != "FIRST_NAME" {
// 		t.Fatalf("unexpected Firstname: %v", a.Firstname)
// 	} else if a.Lastname != "LAST_NAME" {
// 		t.Fatalf("unexpected Lastname: %v", a.Lastname)
// 	} else if a.Email != "test@test.com" {
// 		t.Fatalf("unexpected email: %v", a.Email)
// 	} else if a.Password != "XXXXXXXXX" {
// 		t.Fatalf("unexpected password: %v", a.Password)
// 	}
// }

// // Ensure service returns an error if the admin doesn't exist.
// func testAdminService_Update_ErrAdminNotFound(t *testing.T) {
// 	token, _ := middleware.GenerateTokenPair(1, middleware.AccestokenLimit, middleware.RefreshtokenLimit)
// 	s, c := MustOpenServerClient()
// 	defer s.Close()

// 	// Mock service.
// 	s.Handler.AdminHandler.AdminService.UpdateFn = func(a *admin.Admin, byEmail bool) error {
// 		return admin.ErrAdminNotFound
// 	}

// 	a := admin.Admin{ID: 1, Firstname: "FIRST_NAME", Lastname: "LAST_NAME", Email: "test@test.com", Password: "XXXXXXXXX", Token: token}
// 	if err := c.AdminService().Update(&a, true); err != admin.ErrAdminNotFound {
// 		t.Fatal(err)
// 	}
// }

// // Ensure service returns an error if an internal error occurs.
// func testAdminService_Update_ErrInternal(t *testing.T) {
// 	token, _ := middleware.GenerateTokenPair(1, middleware.AccestokenLimit, middleware.RefreshtokenLimit)
// 	s, c := MustOpenServerClient()
// 	defer s.Close()
// 	s.Handler.AdminHandler.AdminService.UpdateFn = func(a *admin.Admin, byEmail bool) error {
// 		return errors.New("marker")
// 	}

// 	a := admin.Admin{ID: 1, Firstname: "FIRST_NAME", Lastname: "LAST_NAME", Email: "test@test.com", Password: "XXXXXXXXX", Token: token}
// 	if err := c.AdminService().Update(&a, true); err != admin.ErrAdminInternal {
// 		t.Fatal(err)
// 	} else if !strings.Contains(s.Handler.AdminHandler.LogOutput.String(), "marker") {
// 		t.Fatalf("expected log output")
// 	}
// }

// func TestAdminService_Delete(t *testing.T) {
// 	t.Run("OK", testAdminService_Delete)
// 	t.Run("ErrAdminNotFound", testAdminService_Delete_ErrAdminNotFound)
// 	t.Run("ErrInternal", testAdminService_Delete_ErrInternal)
// }

// func testAdminService_Delete(t *testing.T) {
// 	token, _ := middleware.GenerateTokenPair(1, middleware.AccestokenLimit, middleware.RefreshtokenLimit)
// 	s, c := MustOpenServerClient()
// 	defer s.Close()

// 	// Mock service.
// 	s.Handler.AdminHandler.AdminService.DeleteFn = func(a *admin.Admin, byEmail bool) error {
// 		if a.ID != 1 {
// 			t.Fatalf("unexpected admin id: %d", a.ID)
// 		}
// 		return nil
// 	}

// 	a := admin.Admin{ID: 1, Token: token}
// 	if err := c.AdminService().Delete(&a, true); err != nil {
// 		t.Fatal(err)
// 	}
// }

// func testAdminService_Delete_ErrAdminNotFound(t *testing.T) {
// 	token, _ := middleware.GenerateTokenPair(1, middleware.AccestokenLimit, middleware.RefreshtokenLimit)
// 	s, c := MustOpenServerClient()
// 	defer s.Close()

// 	// Mock service.
// 	s.Handler.AdminHandler.AdminService.DeleteFn = func(a *admin.Admin, byEmail bool) error {
// 		return admin.ErrAdminNotFound
// 	}

// 	a := admin.Admin{ID: 1, Token: token}
// 	if err := c.AdminService().Delete(&a, true); err != admin.ErrAdminNotFound {
// 		t.Fatal(err)
// 	}
// }
// func testAdminService_Delete_ErrInternal(t *testing.T) {
// 	token, _ := middleware.GenerateTokenPair(1, middleware.AccestokenLimit, middleware.RefreshtokenLimit)
// 	s, c := MustOpenServerClient()
// 	defer s.Close()
// 	s.Handler.AdminHandler.AdminService.DeleteFn = func(a *admin.Admin, byEmail bool) error {
// 		return admin.Error("marker")
// 	}

// 	a := admin.Admin{ID: 1, Token: token}
// 	if err := c.AdminService().Delete(&a, true); err != admin.ErrAdminInternal {
// 		t.Fatal(err)
// 	} else if !strings.Contains(s.Handler.AdminHandler.LogOutput.String(), "marker") {
// 		t.Fatalf("expected log output")
// 	}
// }
