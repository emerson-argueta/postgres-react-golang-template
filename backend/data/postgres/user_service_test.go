package postgres_test

import (
	"emersonargueta/m/v1/identity"
	"emersonargueta/m/v1/identity/user"
	"reflect"
	"testing"
)

// Ensure an user can be created and retrieved.
func TestUserService_CreateUser(t *testing.T) {
	c := MustOpenClient()
	defer c.Close()
	s := c.UserService()

	uuid_u := "TEST"
	email_u := "test@test.com"
	password_u := "test1234"

	u := user.User{
		UUID:     &uuid_u,
		Email:    &email_u,
		Password: &password_u,
	}

	if err := s.CreateUser(&u); err != nil {
		t.Fatal(err)
	}

	// Retrieve user and compare.
	if other, err := s.RetrieveUser(*u.UUID); err != nil {
		t.Fatal(err)
	} else if !reflect.DeepEqual(&u, other) {
		t.Fatalf("unexpected user: %#v", other)
	}

	if err := s.DeleteUser(*u.UUID); err != nil {
		t.Fatal(err)
	}
}

func TestUserService_CreateUser_ErrUserExists(t *testing.T) {
	c := MustOpenClient()
	defer c.Close()
	s := c.UserService()

	email_u := "test@test.com"
	password_u := "test1234"

	u := user.User{
		Email:    &email_u,
		Password: &password_u,
	}

	if err := s.CreateUser(&u); err != nil {
		t.Fatal(err)
	}
	if err := s.CreateUser(&u); err != identity.ErrUserExists {
		t.Fatal(err)
	}
	// Clean up database
	if err := s.DeleteUser(*u.UUID); err != nil {
		t.Fatal(err)
	}
}

func TestUserService_UpdateUser(t *testing.T) {
	c := MustOpenClient()
	defer c.Close()
	s := c.UserService()

	email_u1 := "hello@test.com"
	password_u1 := "world1234"
	user1 := user.User{
		Email:    &email_u1,
		Password: &password_u1,
	}

	email_u2 := "dog@test.com"
	password_u2 := "people1234"
	user2 := user.User{
		Email:    &email_u2,
		Password: &password_u2,
	}

	if err := s.CreateUser(&user1); err != nil {
		t.Fatal(err)
	} else if err := s.CreateUser(&user2); err != nil {
		t.Fatal(err)
	}
	email_u1_u := "hello_update@test.com"
	password_u1_u := "world_update1234"
	user1Update := user.User{
		UUID:     user1.UUID,
		Email:    &email_u1_u,
		Password: &password_u1_u,
	}

	email_u2_u := "dog_update@test.com"
	password_u2_u := "people_update1234"
	user2Update := user.User{
		UUID:     user2.UUID,
		Email:    &email_u2_u,
		Password: &password_u2_u,
	}

	if err := s.UpdateUser(&user1Update); err != nil {
		t.Fatal(err)
	} else if err := s.UpdateUser(&user2Update); err != nil {
		t.Fatal(err)
	}

	// Verify user1 updated searching user1.
	if d, err := s.RetrieveUser(*user1.UUID); err != nil {
		t.Fatal(err)
	} else if !reflect.DeepEqual(&user1Update, d) {
		t.Fatalf("unexpected user: %#v", d)
	}

	// Verify user2 updated searching user2.
	if d, err := s.RetrieveUser(*user2.UUID); err != nil {
		t.Fatal(err)
	} else if !reflect.DeepEqual(&user2Update, d) {
		t.Fatalf("unexpected user: %#v", d)
	}

	// Clean up database.
	if err := s.DeleteUser(*user1Update.UUID); err != nil {
		t.Fatal(err)
	} else if err := s.DeleteUser(*user2Update.UUID); err != nil {
		t.Fatal(err)
	}

}

func TestUserService_DeleteUser_ErrUserNotFound(t *testing.T) {
	c := MustOpenClient()
	defer c.Close()
	s := c.UserService()

	uuid_u := "TEST"
	email_u := "hello@world.com"
	password_u := "XXX"
	u := user.User{
		UUID:     &uuid_u,
		Email:    &email_u,
		Password: &password_u,
	}

	if err := s.DeleteUser(*u.UUID); err != identity.ErrUserNotFound {
		t.Fatal(err)
	}
}
