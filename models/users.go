package models

import (
	"errors"
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"golang.org/x/crypto/bcrypt"

	"../hash"
	"../rand"
)

type User struct {
	gorm.Model
	Name         string
	Email        string `gorm:"not null;unique_index"`
	Password     string `gorm:"-"`
	PasswordHash string `gorm:"not null"`
	Remember     string `gorm:"-"`
	RememberHash string `gorm:"not null;unique_index"`
}

// userValidator is our validation layer that validates
// and normalizes data before passing it on to the next
// UserDB in our interface chain
type userValidator struct {
	UserDB
}

// UserDB is used to interact with the users database
//
// For pretty much all single user queries:
// If the user is found, return the user and nil
// If the user is not found, return nil and ErrNotFound
// If another error occurs, return the error we receive, which
// may not be an error generated by the models package.
//
// For single user queries, any error but ErrNotFound should
// probably result in a 500 ErrInternalServerError until we make
// "public" facing errors.
type UserDB interface {
	// Methods for querying for a single user
	ByID(id uint) (*User, error)
	ByEmail(email string) (*User, error)
	ByRemember(token string) (*User, error)

	// Methods for altering a single user
	Create(user *User) error
	UpdateWithRememberHash(user *User) error
	Delete(id uint) error

	// Close a database connection
	Close() error

	// Migration helpers
	AutoMigrate() error
	DestructiveReset() error
}

// a compile-time error below indicates the userService type no longer matches
// the UserService interface. They should match.
var _ UserService = &userService{}

// UserService interface methods are used to work with the user model
type UserService interface {
	// Authenticate verifies the provided email address and
	// password are correct.
	// If correct, return the corresponding user and nil.
	// Otherwise, return ErrNotFound, ErrInvalidPassword, or
	// pass along an error received from deeper in the stack.
	Authenticate(email string, password string) (*User, error)
	UserDB
}

type userService struct {
	UserDB
}

// a compile-time error below indicates the UserDB type no longer matches
// the userGorm interface. They should match.
var _ UserDB = &userGorm{}

// userGorm represents our database interaction layer
// and implements the UserDB interface fully
type userGorm struct {
	db   *gorm.DB
	hmac hash.HMAC
}

const hmacSecretKey = "secret-hmac-key"
const userPwPepper = "secret-random-string"

// newUserGorm returns a pointer to a new userGorm instance,
// effectively a connection to the user database.
// Only newUserGorm knows or cares which SQL-style database we're using.
func newUserGorm(connectionInfo string) (*userGorm, error) {
	db, err := gorm.Open("postgres", connectionInfo)
	if err != nil {
		return nil, err
	}
	db.LogMode(true)
	hmac := hash.NewHMAC(hmacSecretKey)
	return &userGorm{
		db:   db,
		hmac: hmac,
	}, nil
}

// NewUserService returns a UserService INTERFACE that other
// packages will use to access the user database.
//
// To change to a NoSQL database, replace the newUserGorm call with
// a comparable call to open a different database.
func NewUserService(connectionInfo string) (UserService, error) {
	// log.Printf("enter NewUserService, connectionInfo: %s", connectionInfo)
	ug, err := newUserGorm(connectionInfo)
	if err != nil {
		return nil, err
	}

	// MAGIC: somehow returning the address of a userService is equivalent to
	// returning a UserService interface ???
	return &userService{
		UserDB: userValidator{
			UserDB: ug,
		},
	}, nil
}

// Close the UserService database connection
func (ug *userGorm) Close() error {
	// log.Printf("enter UserService.Close")
	return ug.db.Close()
}

var (
	// ErrNotFound is returned when the query executes successfully
	// but returned zero rows. I.e., the resource cannot be found
	// in the database.
	ErrNotFound = errors.New("models: resource not found")

	// ErrInvalidID is returned when an invalid ID is provided
	// to a method like Delete.
	ErrInvalidID = errors.New("models: ID provided was invalid")

	// ErrInvalidPassword is returned when an invalid password
	// is dtected when attempting to authenticate a user.
	ErrInvalidPassword = errors.New("models: incorrect password provided")
)

// Create expects the Name, Email and Password fields to be populated, and
// will populate the remaining fields before creating the  User record in the database.
// GORM will populate the gorm.Model data including the ID, CreatedAt, and
// UpdatedAt fields.
func (ug *userGorm) Create(user *User) error {
	pwBytes := []byte(user.Password + userPwPepper)
	hashedBytes, err := bcrypt.GenerateFromPassword(pwBytes, bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user.PasswordHash = string(hashedBytes) // store the PasswordHash in the database
	user.Password = ""                      // ... but overwrite the Password immediately (does not reach DB)

	if user.Remember == "" { // created/populated at login, so expected to be empty
		token, err := rand.RememberToken()
		if err != nil {
			return err
		}
		user.Remember = token // store in the DB
	} else {
		fmt.Printf("user.Remember unexpectedly not empty: \"%s\"\n", user.Remember)
	}
	user.RememberHash = ug.hmac.Hash(user.Remember) // store in the DB, will be

	return ug.db.Create(user).Error
}

// Authenticate will authenticate a user using the
// provided email address and password.
// If the email address provided is invalid, return
//   nil, ErrNotFound
// If the password provided is invalid, return
//   nil, ErrInvalidPassword
// If both the email and password are valid (success), return
//   user, nil
// If there is another error, return
//   nil, error
func (us *userService) Authenticate(email string, password string) (*User, error) {
	foundUser, err := us.ByEmail(email)
	if err != nil {
		return nil, err // pass on ByEmail's error return, email not found in the database
	}

	// test the provided password against the stored PasswordHash
	err = bcrypt.CompareHashAndPassword(
		[]byte(foundUser.PasswordHash),
		[]byte(password+userPwPepper))

	switch err {
	case nil:
		return foundUser, nil // SUCCESS, return user populated with fields from DB
	case bcrypt.ErrMismatchedHashAndPassword:
		return nil, ErrInvalidPassword // password did not produce matching hash
	default:
		return nil, err // some other error
	}
}

// ByID will look up a user with the provided ID.
// If the user is found, return a nil error
// If the user is not found, return ErrNotFound
// If there is another error, return and error with
// more information about what went wrong. This
// may not be an error generated by the models package.
//
// As a general rule, any error but ErrNotFound should
// probably result in a 500 error.
func (ug *userGorm) ByID(id uint) (*User, error) {
	var user User
	db := ug.db.Where("id = ?", id)
	err := first(db, &user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// ByEmail looks up a user with the given email address and
// returns that user.
// If the user is found, we will return a nil error.
// If the user is not found, we will return ErrNotFound
// If there is another error, we will return an error with
// more information about what went wrong. This may not be
// an error genereated by the models package.
func (ug *userGorm) ByEmail(email string) (*User, error) {
	var user User
	db := ug.db.Where("email = ?", email)
	err := first(db, &user)
	return &user, err
}

// ByRemember looks up a user with the given remember token
// and returns that user. This method will handle hasning
// the token for us.
// Errors are the same as ByEmail above
func (ug *userGorm) ByRemember(token string) (*User, error) {
	var user User
	rememberHash := ug.hmac.Hash(token)
	err := first(ug.db.Where("remember_hash = ?", rememberHash), &user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// UpdateWithRememberHash will update the user's DB record with the
// RememberHash value from the Remember field of the provided User object
func (ug *userGorm) UpdateWithRememberHash(user *User) error {
	if user.Remember != "" {
		user.RememberHash = ug.hmac.Hash(user.Remember)
	} else {
		fmt.Println("user.Remember unexpectedly empty")
	}
	// save the updated user record (with RememberHash) to the DB
	return ug.db.Save(user).Error
}

// Delete will delete the user with the provided ID
func (ug *userGorm) Delete(id uint) error {
	if id == 0 {
		return ErrInvalidID
	}
	user := User{Model: gorm.Model{ID: id}}
	return ug.db.Delete(&user).Error
}

// DestructiveReset drops the user table and rebuilds it
func (ug *userGorm) DestructiveReset() error {
	err := ug.db.DropTableIfExists(&User{}).Error
	if err != nil {
		return err
	}
	return ug.AutoMigrate()
}

// AutoMigrate will attempt to automaticaly migrate
// the Users table
func (ug *userGorm) AutoMigrate() error {
	if err := ug.db.AutoMigrate(&User{}).Error; err != nil {
		return err
	}
	return nil
}

// first will query using the provided gorm.DB pointer,
// get the first item returned, and store it into dst. If
// the query returns nothing, return ErrNotFound.
func first(db *gorm.DB, dst interface{}) error {
	err := db.First(dst).Error
	if err == gorm.ErrRecordNotFound {
		return ErrNotFound
	}
	return err
}
