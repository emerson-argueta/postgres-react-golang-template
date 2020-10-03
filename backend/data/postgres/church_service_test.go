package postgres_test

import (
	"reflect"
	"testing"

	"trustdonations.org/m/v2/domain"
	"trustdonations.org/m/v2/domain/church"
)

// Ensure an church can be created and retrieved.
func TestChurchService_CreateChurch(t *testing.T) {
	c := MustOpenClient()
	defer c.Close()
	s := c.ChurchService()

	name_c := "TEST"
	email_c := "test@test.com"
	password_c := "test1234"
	access := domain.NonRestricted
	role := domain.Creator
	donatoruuid := "TESTDONATOR"
	closingbalance := float64(100.23)
	curTime := c.Now()
	ch := church.Church{
		Name:             &name_c,
		Email:            &email_c,
		Password:         &password_c,
		Administrators:   &church.Administrators{"TESTADMINISTRATOR_UUID": {Access: &access, Role: &role}},
		Donators:         &church.Donators{1: {UUID: &donatoruuid}},
		Accountstatement: &domain.AccountStatement{Closingbalance: &closingbalance, Date: &curTime},
	}

	if err := s.Create(&ch); err != nil {
		t.Fatal(err)
	}

	// Retrieve admin compare.
	byEmail := true
	if other, err := s.Read(&ch, byEmail); err != nil {
		t.Fatal(err)
	} else if !reflect.DeepEqual(&ch, other) {
		t.Fatalf("unexpected admin: %#v", other)
	}

	if err := s.Delete(&ch, byEmail); err != nil {
		t.Fatal(err)
	}
}

func TestChurchService_CreateAdmin_ErrChurchExists(t *testing.T) {
	c := MustOpenClient()
	defer c.Close()
	s := c.ChurchService()

	name_c := "TEST"
	email_c := "test@test.com"
	password_c := "test1234"
	ch := church.Church{
		Name:     &name_c,
		Email:    &email_c,
		Password: &password_c,
	}

	if err := s.Create(&ch); err != nil {
		t.Fatal(err)
	}
	if err := s.Create(&ch); err != church.ErrChurchExists {
		t.Fatal(err)
	}
	// Clean up database
	byEmail := false
	if err := s.Delete(&ch, byEmail); err != nil {
		t.Fatal(err)
	}
}

func TestChurchService_UpdateChurch(t *testing.T) {
	c := MustOpenClient()
	defer c.Close()
	s := c.ChurchService()

	name_c1 := "TEST0"
	email_c1 := "test_0@test.com"
	password_c1 := "test1234"
	church1 := church.Church{
		Name:     &name_c1,
		Email:    &email_c1,
		Password: &password_c1,
	}

	name_c2 := "TEST1"
	email_c2 := "test_1@test.com"
	password_c2 := "test1234"
	church2 := church.Church{
		Name:     &name_c2,
		Email:    &email_c2,
		Password: &password_c2,
	}
	// Create new church
	if err := s.Create(&church1); err != nil {
		t.Fatal(err)
	} else if err := s.Create(&church2); err != nil {
		t.Fatal(err)
	}

	name_c1_u := "TEST0_UPDATE"
	email_c1_u := "test_0_update@test.com"
	password_c1_u := "test1234_update"
	church1Update := church.Church{
		ID:       church1.ID,
		Name:     &name_c1_u,
		Email:    &email_c1_u,
		Password: &password_c1_u,
	}

	name_c2_u := "TEST1_UPDATE"
	email_c2_u := "test_1_update@test.com"
	password_c2_u := "test1234_update"
	admin2Update := church.Church{
		ID:       church2.ID,
		Name:     &name_c2_u,
		Email:    &email_c2_u,
		Password: &password_c2_u,
	}
	// Update admins.
	byEmail := false
	if err := s.Update(&church1Update, byEmail); err != nil {
		t.Fatal(err)
	} else if err := s.Update(&admin2Update, byEmail); err != nil {
		t.Fatal(err)
	}

	// Verify church1 updated searching church1.
	if d, err := s.Read(&church1, byEmail); err != nil {
		t.Fatal(err)
	} else if !reflect.DeepEqual(&church1Update, d) {
		t.Fatalf("unexpected admin: %#v", d)
	}

	// Verify admin2 updated searching admin2.
	if d, err := s.Read(&church2, byEmail); err != nil {
		t.Fatal(err)
	} else if !reflect.DeepEqual(&admin2Update, d) {
		t.Fatalf("unexpected admin: %#v", d)
	}

	// Clean up database searching church by id.
	if err := s.Delete(&church1, byEmail); err != nil {
		t.Fatal(err)
	} else if err := s.Delete(&church2, byEmail); err != nil {
		t.Fatal(err)
	}

}

func TestChurchService_DeleteChurch_ErrChurchNotFound(t *testing.T) {
	c := MustOpenClient()
	defer c.Close()
	s := c.ChurchService()

	name_c := "TEST"
	email_c := "test@test.com"
	password_c := "test1234"
	ch := church.Church{
		Name:     &name_c,
		Email:    &email_c,
		Password: &password_c,
	}

	byEmail := true
	if err := s.Delete(&ch, byEmail); err != church.ErrChurchNotFound {
		t.Fatal(err)
	}
}
