package postgres_test

import (
	"reflect"
	"testing"

	"emersonargueta/m/v1/domain"
	"emersonargueta/m/v1/domain/donator"
)

// Ensure an donator can be created and retrieved.
func TestDonatorService_CreateDonator(t *testing.T) {
	c := MustOpenClient()
	defer c.Close()
	s := c.DonatorService()

	firstname_d := "TEST"
	email_d := "test@test.com"
	donationcount := int64(100)
	closingbalance := float64(100.23)
	curTime := c.Now()
	d := donator.Donator{
		Firstname:        &firstname_d,
		Email:            &email_d,
		Churches:         &donator.Churches{1: {Donationcount: &donationcount, Firstdonation: &curTime}},
		Accountstatement: &domain.AccountStatement{Closingbalance: &closingbalance, Date: &curTime},
	}

	if err := s.Create(&d); err != nil {
		t.Fatal(err)
	}

	// Retrieve admin compare.
	byEmail := true
	if other, err := s.Read(&d, byEmail); err != nil {
		t.Fatal(err)
	} else if !reflect.DeepEqual(&d, other) {
		t.Fatalf("unexpected admin: %#v", other)
	}

	if err := s.Delete(&d, byEmail); err != nil {
		t.Fatal(err)
	}
}

func TestDonatorService_ReadWithFilterDonator(t *testing.T) {
	c := MustOpenClient()
	defer c.Close()
	s := c.DonatorService()

	firstname_d := "test_firstname"
	lastname_d := "test_lastname"
	email_d := "test@test.com"
	donationcount := int64(100)
	closingbalance := float64(100.23)
	curTime := c.Now()
	d := donator.Donator{
		Firstname:        &firstname_d,
		Lastname:         &lastname_d,
		Email:            &email_d,
		Churches:         &donator.Churches{1: {Donationcount: &donationcount, Firstdonation: &curTime}},
		Accountstatement: &domain.AccountStatement{Closingbalance: &closingbalance, Date: &curTime},
	}

	if err := s.Create(&d); err != nil {
		t.Fatal(err)
	}

	// Retrieve admin compare.
	byEmail := true
	if other, err := s.ReadWithFilter(&d, &donator.Donator{Firstname: &firstname_d, Lastname: &lastname_d, Email: &email_d}); err != nil {
		t.Fatal(err)
	} else if !reflect.DeepEqual(&d, other) {
		t.Fatalf("unexpected admin: %#v", other)
	}

	if err := s.Delete(&d, byEmail); err != nil {
		t.Fatal(err)
	}
}

func TestDonatorService_CreateDonator_ErrDonatorExists(t *testing.T) {
	c := MustOpenClient()
	defer c.Close()
	s := c.DonatorService()

	firstname_d := "TEST"
	email_d := "test@test.com"
	d := donator.Donator{
		Firstname: &firstname_d,
		Email:     &email_d,
	}

	if err := s.Create(&d); err != nil {
		t.Fatal(err)
	}
	if err := s.Create(&d); err != donator.ErrDonatorExists {
		t.Fatal(err)
	}
	// Clean up database
	byEmail := false
	if err := s.Delete(&d, byEmail); err != nil {
		t.Fatal(err)
	}
}

func TestDonatorService_UpdateDonator(t *testing.T) {
	c := MustOpenClient()
	defer c.Close()
	s := c.DonatorService()

	firstname_d1 := "TEST0"
	email_d1 := "test_0@test.com"
	donator1 := donator.Donator{
		Firstname: &firstname_d1,
		Email:     &email_d1,
	}

	firstname_d2 := "TEST1"
	email_d2 := "test_1@test.com"
	donator2 := donator.Donator{
		Firstname: &firstname_d2,
		Email:     &email_d2,
	}
	// Create new donator
	if err := s.Create(&donator1); err != nil {
		t.Fatal(err)
	} else if err := s.Create(&donator2); err != nil {
		t.Fatal(err)
	}

	firstname_d1_u := "TEST0_UPDATE"
	email_d1_u := "test_0_update@test.com"
	donator1Update := donator.Donator{
		ID:        donator1.ID,
		Firstname: &firstname_d1_u,
		Email:     &email_d1_u,
	}

	firstname_d2_u := "TEST1_UPDATE"
	email_d2_u := "test_1_update@test.com"
	donator2Update := donator.Donator{
		ID:        donator2.ID,
		Firstname: &firstname_d2_u,
		Email:     &email_d2_u,
	}
	// Update admins.
	byEmail := false
	if err := s.Update(&donator1Update, byEmail); err != nil {
		t.Fatal(err)
	} else if err := s.Update(&donator2Update, byEmail); err != nil {
		t.Fatal(err)
	}

	// Verify donator1 updated searching donator1.
	if d, err := s.Read(&donator1, byEmail); err != nil {
		t.Fatal(err)
	} else if !reflect.DeepEqual(&donator1Update, d) {
		t.Fatalf("unexpected admin: %#v", d)
	}

	// Verify admin2 updated searching admin2.
	if d, err := s.Read(&donator2, byEmail); err != nil {
		t.Fatal(err)
	} else if !reflect.DeepEqual(&donator2Update, d) {
		t.Fatalf("unexpected admin: %#v", d)
	}

	// Clean up database searching donator by id.
	if err := s.Delete(&donator1, byEmail); err != nil {
		t.Fatal(err)
	} else if err := s.Delete(&donator2, byEmail); err != nil {
		t.Fatal(err)
	}

}

func TestDonatorService_DeleteDonator_ErrDonatorNotFound(t *testing.T) {
	c := MustOpenClient()
	defer c.Close()
	s := c.DonatorService()

	firstname_d := "TEST"
	email_d := "test@test.com"
	d := donator.Donator{
		Firstname: &firstname_d,
		Email:     &email_d,
	}

	byEmail := true
	if err := s.Delete(&d, byEmail); err != donator.ErrDonatorNotFound {
		t.Fatal(err)
	}
}
