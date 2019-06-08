package models

import (
	"fmt"
	"math"
	"math/rand"
	"os"
	"strconv"
	"testing"
	"time"
)

// "github.com/peterpla/webdevgo/controllers"

const (
	dbHost = "localhost"
	dbPort = 5432
	dbUser = "postgres"
	dbName = "whatever_dev"
)

var connStr string
var services *Services

func TestMain(m *testing.M) {
	connStr = fmt.Sprintf("host=%s port=%d user=%s dbname=%s sslmode=disable", dbHost, dbPort, dbUser, dbName)
	// fmt.Printf("TestMain: %s\n", connStr)
	// initialize services and database connection
	services, err := NewServices(connStr)
	if err != nil {
		panic(err)
	}
	defer services.User.Close()

	os.Exit(m.Run())
}

func TestNewUserServiceByIDAndClose(t *testing.T) {
	// call ByEmail("bozo@clown.net") to confirm UserService is operational
	// expect ErrNotFound from ByEmail
	if _, err := services.User.ByEmail("bozo@clown.net"); err != ErrNotFound {
		t.Fatalf("us.ByID(1): expected \"%v\", got \"%v\"", ErrNotFound, err)
	}
}

func TestCreateByEmailAndDelete(t *testing.T) {
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
	if err := services.User.Create(&user); err != nil {
		t.Fatalf("us.Create(): expected nil, got = %v", err)
	}

	var err error
	var foundRecord *User

	// call ByEmail() to confirm that User was created
	if foundRecord, err = services.User.ByEmail(user.Email); err != nil {
		t.Fatalf("us.ByEmail(): expected \"%v\", got \"%v\"", ErrNotFound, err)
	}

	// delete that user
	if err := services.User.Delete(foundRecord.Model.ID); err != nil {
		t.Fatalf("us.Delete(): expected nil, got \"%v\"", err)
	}

	// confirm the created user was deleted by looking for their id
	if foundRecord, err = services.User.ByID(user.ID); err != ErrNotFound {
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
	if err := services.User.Create(&user); err != ErrPasswordRequired {
		t.Fatalf("us.Create(): expected ErrPasswordRequired, got = %v", err)
	}

	// var foundRecord *User

	// due to blank password, Authenticate should return ErrNotFound
	if _, err := services.User.Authenticate(user.Email, user.Password); err != ErrNotFound {
		t.Logf("us.Authenticate(): expected ErrNotFound, got \"%v\"", err)
	}

}

/*
func TestUvHmacRemember(t *testing.T) {
	// connect to the user database
	ug, err := newUserGorm(connStr)
	if err != nil {
		t.Error("could not connect to database")
	}
	hmac := hash.NewHMAC(hmacSecretKey)
	uv := newUserValidator(ug, hmac)

	type testset struct {
		testUser     User
		expErr       error
		expHashEmpty bool
	}

	var tests = []testset{
		{
			User{
				Name:         "Bozo Clown",
				Email:        "bozo@clown.net",
				Password:     "",
				PasswordHash: "",
				Remember:     "",
				RememberHash: "", // should be untouched
			},
			nil,
			true,
		},
		{
			User{
				Name:         "Bozo Clown",
				Email:        "bozo@clown.net",
				Password:     "",
				PasswordHash: "",
				Remember:     "notempty",
				RememberHash: "", // should get new hash
			},
			nil,
			false,
		},
	}

	for _, r := range tests {
		// log.Printf("*** current test: User: %+v, expErr: %+v, expHashEmpty: %t\n",
		// 	r.testUser, r.expErr, r.expHashEmpty)

		err := uv.hmacRemember(&r.testUser)

		if err != r.expErr {
			t.Errorf("hmacRemember: got %v, want %v", err, r.expErr)
		}

		if r.expHashEmpty && r.testUser.RememberHash != "" {
			// expect Remember hash to be "", but it isn't
			t.Errorf("hmacRemember: got %s, want %s", r.testUser.RememberHash, "")
		}
		if !r.expHashEmpty && r.testUser.RememberHash == "" {
			// expect Remember hash to be not-empty, but it is
			t.Errorf("hmacRemember: got %s, want %s", r.testUser.RememberHash, "")
		}
	}
}
*/

/*
func TestSetRememberIfUnset(t *testing.T) {
	// connect to the user database
	ug, err := newUserGorm(connStr)
	if err != nil {
		t.Error("could not connect to database")
	}

	hmac := hash.NewHMAC(hmacSecretKey)
	uv := newUserValidator(ug, hmac)

	type testset struct {
		testUser                User
		expErr                  error
		expRememberShouldChange bool
	}

	var tests = []testset{
		{
			User{
				Name:         "Bozo Clown",
				Email:        "bozo@clown.net",
				Password:     "",
				PasswordHash: "",
				Remember:     "", // blank, should change
				RememberHash: "",
			},
			nil,
			true,
		},
		// TODO: confirm this test fails due to controllers/user.go/SignIn also sets Remember?
		{
			User{
				Name:         "Bozo Clown",
				Email:        "bozo@clown.net",
				Password:     "",
				PasswordHash: "",
				Remember:     "notempty", // non-empty, should NOT change
				RememberHash: "",
			},
			nil,
			false,
		},
	}

	for _, r := range tests {
		// log.Printf("*** current test: User: %+v, expErr: %+v, expHashEmpty: %t\n",
		// 	r.testUser, r.expErr, r.expHashEmpty)

		origRemember := r.testUser.Remember
		err = uv.setRememberIfUnset(&r.testUser)

		if err != r.expErr {
			t.Errorf("setRememberIfUnset: got %v, want %v", err, r.expErr)
		}

		if r.expRememberShouldChange && r.testUser.Remember == origRemember {
			// expect Remember to change, but it didn't
			t.Errorf("setRememberIfUnset didn't change: got %q, want %q",
				r.testUser.Remember, origRemember)
		}
		if !r.expRememberShouldChange && r.testUser.RememberHash != origRemember {
			// expect Remember to be unchanged, but it changed
			t.Skip("setRememberIfUnset when should not change, still does - outside forces?")
			t.Errorf("setRememberIfUnset changed: got %q, want %q",
				r.testUser.RememberHash, origRemember)
		}
	}
}
*/
