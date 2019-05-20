package models

import (
	"fmt"
	"os"
	"testing"
)

const (
	dbHost = "localhost"
	dbPort = 5432
	dbUser = "postgres"
	dbName = "whatever_dev"
)

var connStr string

func TestMain(m *testing.M) {
	connStr = fmt.Sprintf("host=%s port=%d user=%s dbname=%s sslmode=disable", dbHost, dbPort, dbUser, dbName)
	os.Exit(m.Run())
}

func TestNewUserServiceByIDAndClose(t *testing.T) {
	// connect to the user database
	us, err := NewUserService(connStr)
	if err != nil {
		t.Fatalf("NewUser(psqlInfo): err = %v; want nil", err)
	}
	defer us.Close()

	// call ByID(1) to confirm UserService is operational
	// expect ErrNotFound from ByID(1)
	if _, err = us.ByID(1); err != ErrNotFound {
		t.Fatalf("us.ByID(1): expected \"%v\", got \"%v\"", ErrNotFound, err)
	}
}

func TestCreateByEmailAndDelete(t *testing.T) {
	// connect to the user database
	us, err := NewUserService(connStr)
	if err != nil {
		t.Fatalf("NewUser(psqlInfo): expected nil, got = %v", err)
	}
	defer us.Close()

	// Create a user
	user := User{
		Name:  "Test1 User",
		Email: "test1@test.com",
		Age:   18,
	}
	if err := us.Create(&user); err != nil {
		t.Fatalf("us.Create(): expected nil, got = %v", err)
	}

	var foundRecord *User

	// call ByEmail() to confirm that User was created
	if foundRecord, err = us.ByEmail(user.Email); err != nil {
		t.Fatalf("us.ByEmail(): expected \"%v\", got \"%v\"", ErrNotFound, err)
	}

	// delete that user
	if err = us.Delete(foundRecord.Model.ID); err != nil {
		t.Fatalf("us.Delete(): expected nil, got \"%v\"", err)
	}

	// confirm that user deleted by calling ByID()
	if foundRecord, err = us.ByID(user.ID); err != ErrNotFound {
		t.Fatalf("us.ByID(): expected \"%v\", got \"%v\"", ErrNotFound, err)
	}

}
