package models

import (
	"fmt"
	"testing"
	"time"
)

func testingUserService() (*UserService, error) {
	const (
		host     = "localhost"
		port     = 5432
		user     = "postgres"
		password = "root"
		dbname   = "lenslocked_test"
	)

	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	us, err := NewUserService(dsn)
	if err != nil {
		return nil, err
	}

	us.ResetDB()

	return us, nil
}

func TestCreateUser(t *testing.T) {
	us, err := testingUserService()
	if err != nil {
		t.Fatal(err)
	}
	user := User{
		Name:  "Ridley Scott",
		Email: "ridleyscott@gmail.com",
	}
	start := time.Now()
	err = us.Create(&user)
	if err != nil {
		t.Fatal(err)
	}
	if user.ID == 0 {
		t.Errorf("Expected ID > 0, Received %d", user.ID)
	}

	elapsed := time.Now().Sub(start)

	if time.Since(user.CreatedAt) > elapsed*5*time.Second {
		t.Errorf("Expected CreatedAt to be recent. Received %s", user.CreatedAt)
	}
	if time.Since(user.UpdatedAt) > elapsed*5*time.Second {
		t.Errorf("Expected UpdatedAt to be recent. Received %s", user.UpdatedAt)
	}
}
