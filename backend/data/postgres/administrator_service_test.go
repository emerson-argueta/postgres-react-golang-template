package postgres_test

import (
	"reflect"
	"testing"

	"emersonargueta/m/v1/domain"
	"emersonargueta/m/v1/domain/administrator"
)

// Ensure an administrator can be created and retrieved.
func TestAdministratorService_CreateAdministrator(t *testing.T) {
	c := MustOpenClient()
	defer c.Close()
	s := c.AdministratorService()

	uuid_a := "TEST"
	firstname_a := "Hello"
	lastname_a := "World"

	subscriptionType := administrator.FreePlan
	role := domain.Creator
	access := domain.NonRestricted

	freeusagelimitcount := int64(100)
	customeremail := "test@test.com"
	a := administrator.Administrator{
		UUID:         &uuid_a,
		Firstname:    &firstname_a,
		Lastname:     &lastname_a,
		Churches:     &administrator.Churches{1: &administrator.Church{Role: &role, Access: &access}},
		Subscription: &administrator.Subscription{Freeusagelimitcount: &freeusagelimitcount, Customeremail: &customeremail, Type: &subscriptionType},
	}

	if err := s.Create(&a); err != nil {
		t.Fatal(err)
	}

	// Retrieve admin compare.
	if other, err := s.Read(&a); err != nil {
		t.Fatal(err)
	} else if !reflect.DeepEqual(&a, other) {
		t.Fatalf("unexpected admin: %#v", other)
	}

	if err := s.Delete(&a); err != nil {
		t.Fatal(err)
	}
}

func TestAdministratorService_CreateAdministrator_ErrAdministratorExists(t *testing.T) {
	c := MustOpenClient()
	defer c.Close()
	s := c.AdministratorService()

	uuid_a := "TEST"
	firstname_a := "Hello"
	lastname_a := "World"
	a := administrator.Administrator{
		UUID:      &uuid_a,
		Firstname: &firstname_a,
		Lastname:  &lastname_a,
	}

	if err := s.Create(&a); err != nil {
		t.Fatal(err)
	}
	if err := s.Create(&a); err != administrator.ErrAdministratorExists {
		t.Fatal(err)
	}
	// Clean up database
	if err := s.Delete(&a); err != nil {
		t.Fatal(err)
	}
}

func TestAdministratorService_UpdateAdministrator(t *testing.T) {
	c := MustOpenClient()
	defer c.Close()
	s := c.AdministratorService()

	subscriptionType := administrator.FreePlan
	uuid_admin1 := "TEST0"
	firstname_admin1 := "Hello"
	lastname_admin1 := "World"
	address_admin1 := "TESTADDRESS"

	freeusagelimitcount_admin1 := int64(100)
	customeremail_admin1 := "test@test.com"
	admin1 := administrator.Administrator{
		UUID:         &uuid_admin1,
		Firstname:    &firstname_admin1,
		Lastname:     &lastname_admin1,
		Address:      &address_admin1,
		Subscription: &administrator.Subscription{Freeusagelimitcount: &freeusagelimitcount_admin1, Customeremail: &customeremail_admin1, Type: &subscriptionType},
	}
	uuid_admin2 := "TEST1"
	firstname_admin2 := "Dog"
	lastname_admin2 := "People"
	admin2 := administrator.Administrator{
		UUID:      &uuid_admin2,
		Firstname: &firstname_admin2,
		Lastname:  &lastname_admin2,
	}
	// Create new administrators
	if err := s.Create(&admin1); err != nil {
		t.Fatal(err)
	} else if err := s.Create(&admin2); err != nil {
		t.Fatal(err)
	}

	firstname_admin1_u := "Hello_update"
	lastname_admin1_u := "World_update"
	admin1Update := administrator.Administrator{

		UUID:      admin1.UUID,
		Firstname: &firstname_admin1_u,
		Lastname:  &lastname_admin1_u,
	}

	firstname_admin2_u := "Dog_update"
	lastname_admin2_u := "People_update"
	admin2Update := administrator.Administrator{

		UUID:      admin2.UUID,
		Firstname: &firstname_admin2_u,
		Lastname:  &lastname_admin2_u,
	}
	// Update administrators.
	if err := s.Update(&admin1Update); err != nil {
		t.Fatal(err)
	} else if err := s.Update(&admin2Update); err != nil {
		t.Fatal(err)
	}

	// Verify admin1 updated searching admin1.
	if d, err := s.Read(&admin1); err != nil {
		t.Fatal(err)
	} else if !reflect.DeepEqual(&admin1Update, d) {
		t.Fatalf("unexpected admin: %#v", d)
	}

	// Verify admin2 updated searching admin2.
	if d, err := s.Read(&admin2); err != nil {
		t.Fatal(err)
	} else if !reflect.DeepEqual(&admin2Update, d) {
		t.Fatalf("unexpected admin: %#v", d)
	}

	// Clean up database searching admins by id.
	if err := s.Delete(&admin1); err != nil {
		t.Fatal(err)
	} else if err := s.Delete(&admin2); err != nil {
		t.Fatal(err)
	}

}

func TestAdministratorService_DeleteAdministrator_ErrAdministratorNotFound(t *testing.T) {
	c := MustOpenClient()
	defer c.Close()
	s := c.AdministratorService()

	firstname_a := "Hello"
	lastname_a := "World"
	email_a := "hello@world.com"
	password_a := "XXX"
	a := administrator.Administrator{
		Firstname: &firstname_a,
		Lastname:  &lastname_a,
		Email:     &email_a,
		Password:  &password_a,
	}

	if err := s.Delete(&a); err != administrator.ErrAdministratorNotFound {
		t.Fatal(err)
	}
}
