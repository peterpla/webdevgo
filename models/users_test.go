package models

import (
	"fmt"
	"math"
	"math/rand"
	"os"
	"strconv"
	"testing"
	"time"

	"golang.org/x/crypto/bcrypt"
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
	fmt.Printf("TestMain: %s\n", connStr)
	os.Exit(m.Run())
}

func TestNewUserServiceByIDAndClose(t *testing.T) {
	// connect to the user database
	us, err := NewUserService(connStr)
	if err != nil {
		t.Fatalf("NewUser(psqlInfo): err = %v; want nil", err)
	}
	defer us.Close()

	// call ByEmail("bozo@clown.net") to confirm UserService is operational
	// expect ErrNotFound from ByEmail
	if _, err = us.ByEmail("bozo@clown.net"); err != ErrNotFound {
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

	/* ********** ********** ********** ********** ********** */
	// TEST 1: Create, ByEmail, Delete, ByID

	// random 16-bit integer to add to test username
	rand.Seed(time.Now().UnixNano())
	rnd := rand.Intn(math.MaxUint16)
	r := strconv.Itoa(rnd)

	name := fmt.Sprintf("Test%s User", r)
	email := fmt.Sprintf("test%s@test.com", r)
	pwd := fmt.Sprintf("test%sPASS", r)

	// Create a test user
	user := User{
		Name:     name,
		Email:    email,
		Password: pwd,
	}
	// fmt.Printf("User: %+v", user)
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

	// confirm the created user was deleted by looking for their id
	if foundRecord, err = us.ByID(user.ID); err != ErrNotFound {
		t.Fatalf("us.ByID(): expected \"%v\", got \"%v\"", ErrNotFound, err)
	}

	/* ********** ********** ********** ********** ********** */
	// TEST 2: Create (no password), Authenticate, Delete, ByID

	rand.Seed(time.Now().UnixNano())
	rnd = rand.Intn(math.MaxUint16)
	r = strconv.Itoa(rnd)

	// Create a test user
	user = User{
		Name:     "",
		Email:    "",
		Password: "",
	}

	user.Name = fmt.Sprintf("Test%s User", r)
	user.Email = fmt.Sprintf("test%s@test.com", r)
	// keep Password as empty string

	fmt.Printf("User: %+v\n", user)
	if err := us.Create(&user); err != nil {
		t.Fatalf("us.Create(): expected nil, got = %v", err)
	}

	// var foundRecord *User

	// call ByEmail() to confirm that User was created
	if foundRecord, err = us.ByEmail(user.Email); err != nil {
		t.Fatalf("us.ByEmail(): expected \"%v\", got \"%v\"", ErrNotFound, err)
	}

	// due to blank password, Authenticate should return an error from bcrypt package
	if _, err := us.Authenticate(user.Email, user.Password); err != bcrypt.ErrHashTooShort {
		t.Logf("us.Authenticate(): expected ErrHashTooShort, got \"%v\"", err)
	}

	// delete that user
	if err = us.Delete(foundRecord.Model.ID); err != nil {
		t.Fatalf("us.Delete(): expected nil, got \"%v\"", err)
	}

	// confirm the created user was deleted by looking for their id
	if foundRecord, err = us.ByID(user.ID); err != ErrNotFound {
		t.Fatalf("us.ByID(): expected \"%v\", got \"%v\"", ErrNotFound, err)
	}

}
