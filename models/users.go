package models

import (
	"log"
	"regexp"
	"strings"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres" // load gorm's postgres driver
	"golang.org/x/crypto/bcrypt"

	"github.com/peterpla/webdevgo/hash"
	"github.com/peterpla/webdevgo/rand"
)

// User ... [TODO: add documentation]
type User struct {
	gorm.Model
	Name         string
	Email        string `gorm:"not null;unique_index"`
	Password     string `gorm:"-"`
	PasswordHash string `gorm:"not null"`
	Remember     string `gorm:"-"`
	RememberHash string `gorm:"not null;unique_index"`
}

// userValidator is our validation/normalization layer that
// validates and normalizes data before passing it along our
// interface chain
type userValidator struct {
	UserDB
	hmac       hash.HMAC
	emailRegex *regexp.Regexp
}

// a compile-time error below indicates the UserDB type no longer matches
// the userGorm interface. They should match.
var _ UserDB = &userGorm{}

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
	Update(user *User) error
	Delete(id uint) error
}

// The Public method on the modelError type effectively "white lists" these errors
// for display to end users, after some massaging. Otherwise the error message text
// would be repleaced with AlertMsgGeneric.
type modelError string

var (
	// ErrNotFound is returned when the query executes successfully
	// but returned zero rows. I.e., the resource cannot be found
	// in the database.
	ErrNotFound modelError = "models: resource not found"

	// ErrIDInvalid is returned when an invalid ID is provided
	// to a method like Delete.
	ErrIDInvalid modelError = "models: ID provided was invalid"

	// ErrPasswordRequired is returned when a password is empty (after
	// whitespace trimmed) or password hash is empty
	ErrPasswordRequired modelError = "models: password is required"

	// ErrPasswordTooShort is returned when a user specifies a password
	// shorter than 8 characters
	ErrPasswordTooShort modelError = "models: password must be at least 8 characters long"

	// ErrPasswordIncorrect is returned when an invalid password
	// is dtected when attempting to authenticate a user.
	ErrPasswordIncorrect modelError = "models: incorrect password provided"

	// ErrEmailRequired is returned when an email address is not
	// provided when creating a user
	ErrEmailRequired modelError = "models: email address is required"

	// ErrEmailInvalid is returned when an email address provided
	// fails our regular expression test
	ErrEmailInvalid modelError = "models: email address is not valid"

	// ErrEmailTaken is returned when an update or create is attempted
	// specifying an email address that is already in use (found in the database)
	ErrEmailTaken modelError = "models: email address is already taken"

	// ErrRememberRequired is returned when a create or update
	// is attempted without a user Remember token hash
	ErrRememberRequired modelError = "models: remember token is required"

	// ErrRememberTooShort is returned when a Remember token
	// is not at least 32 bytes
	ErrRememberTooShort modelError = "models: remember token must be at least 32 bytes"
)

// userGorm represents our database interaction layer
// and implements the UserDB interface fully
type userGorm struct {
	db *gorm.DB
}

// a compile-time error below indicates the userService type no longer matches
// the UserService interface. They should match.
var _ UserService = &userService{}

// UserService interface methods are used to work with the user model
type UserService interface {
	// Authenticate verifies the provided email address and
	// password are correct.
	// If correct, return the corresponding user and nil.
	// Otherwise, return ErrNotFound, ErrPasswordIncorrect, or
	// pass along an error received from deeper in the stack.
	Authenticate(email string, password string) (*User, error)
	UserDB
}

type userService struct {
	UserDB
}

const hmacSecretKey = "secret-hmac-key"
const userPwPepper = "secret-random-string"

// NewUserService returns a UserService INTERFACE that other
// packages will use to access the user database.
func NewUserService(db *gorm.DB) UserService {
	ug := &userGorm{db}
	log.Printf("enter NewUserService, ug: %+v", ug)

	hmac := hash.NewHMAC(hmacSecretKey)
	uv := newUserValidator(ug, hmac)

	u := &userService{
		UserDB: uv,
	}
	return u
}

// newUserValidator returns a pointer to a userValidator instance
func newUserValidator(udb UserDB, hmac hash.HMAC) *userValidator {
	return &userValidator{
		UserDB:     udb,
		hmac:       hmac,
		emailRegex: regexp.MustCompile(`[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,16}$`),
	}
}

/* ********** ********** ********** */
/*         userService methods      */

// Authenticate will authenticate a user using the
// provided email address and password.
// If the email address provided is invalid, return
//   nil, ErrNotFound
// If the password provided is invalid, return
//   nil, ErrPasswordIncorrect
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
		return nil, ErrPasswordIncorrect // password did not produce matching hash
	default:
		return nil, err // some other error
	}
}

/* ********** ********** ********** */
/*            userGorm methods      */

// Create expects the Name, Email and Password fields to validated and
// normalized, and will create the user database record, populating
// the gorm.Model data including the ID, CreatedAt, and UpdatedAt fields.
func (ug *userGorm) Create(user *User) error {
	// fmt.Printf("enter Create, user=%+v\n", user)
	return ug.db.Create(user).Error
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
	// fmt.Printf("enter ByID, id=%d\n", id)
	var user User
	db := ug.db.Where("id = ?", id)
	err := first(db, &user)
	if err != nil {
		return nil, err
	}
	// fmt.Printf("exit ByID, user:%+v\n", user)
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
	// fmt.Printf("enter ByEmail, email=%s\n", email)
	var user User
	db := ug.db.Where("email = ?", email)
	err := first(db, &user)
	// fmt.Printf("exit ByEmail, user:%+v\n", user)
	return &user, err
}

// ByRemember looks up a user by remember token hash and returns
// that user.
// Errors are the same as ByEmail above
func (ug *userGorm) ByRemember(rememberHash string) (*User, error) {
	// fmt.Printf("enter ByRemember, rememberHash=%s\n", rememberHash)
	var user User
	err := first(ug.db.Where("remember_hash = ?", rememberHash), &user)
	if err != nil {
		return nil, err
	}
	// fmt.Printf("exit ByRemember, user:%+v\n", user)
	return &user, nil
}

// Update expects the Name, Email and Password fields to be
// validated and normalized, and will update the user's DB record with the
// provided User object
func (ug *userGorm) Update(user *User) error {
	// fmt.Printf("enter Update, user=%+v\n", user)
	return ug.db.Save(user).Error
}

// Delete expects the user ID to be validated and normalized, and will
// delete the user with the provided ID
func (ug *userGorm) Delete(id uint) error {
	// fmt.Printf("enter Delete, id=%d\n", id)
	user := User{Model: gorm.Model{ID: id}}
	// fmt.Printf("calling ug.db.Delete passing user=%+v\n", user)
	return ug.db.Delete(&user).Error
}

/* ********** ********** ********** */
/*         modelError methods       */

// Error implements the expected Error() method
func (e modelError) Error() string {
	return string(e)
}

// Public returns the white-listed error string
func (e modelError) Public() string {
	// strip off initial "models: ", capitalize the first word of the Error string
	s := strings.Replace(string(e), "models: ", "", 1)
	split := strings.Split(s, " ")
	split[0] = strings.Title(split[0])
	return strings.Join(split, " ")
}

/* ********** ********** ********** */
/*       userValidator methods      */

// Create will validate arguments, create the password hash, overwrite the password
// value with an empty string, set the remember token and hash; then pass to the
// database layer to create the user record in the database
func (uv *userValidator) Create(user *User) error {
	/*
		if user.Password == "" {
			panic(ErrPasswordIncorrect)
		}
	*/
	err := runUserValFns(user,
		uv.passwordRequired,     // 1 - sequence matters!
		uv.passwordMinLength,    // 2 - sequence matters!
		uv.bcryptPassword,       // 3 - sequence matters!
		uv.passwordHashRequired, // 4 - sequence matters!
		uv.setRememberIfUnset,
		uv.rememberMinBytes, // after setRemember - sequence matters!
		uv.hmacRemember,
		uv.rememberHashRequired, // after hmacRemember - sequence matters!
		uv.normalizeEmail,
		uv.requireEmail,
		uv.emailFormat,
		uv.emailIsAvail)
	if err != nil {
		return err
	}
	return uv.UserDB.Create(user)
}

// Update will set (normalize) the remember hash, then pass to the database layer to
// update the user record in the database.
func (uv *userValidator) Update(user *User) error {
	err := runUserValFns(user,
		uv.passwordMinLength,    // 1 - sequence matters!
		uv.bcryptPassword,       // 2 - sequence matters!
		uv.passwordHashRequired, // 3 - sequence matters!
		uv.rememberMinBytes,
		uv.hmacRemember,
		uv.rememberHashRequired, // after hmacRemember - sequence matters!
		uv.normalizeEmail,
		uv.requireEmail,
		uv.emailFormat,
		uv.emailIsAvail)
	if err != nil {
		return err
	}
	return uv.UserDB.Update(user)
}

// Delete will validate the provided user ID, then pass to the database layer to
// delete the user record from the database.
func (uv *userValidator) Delete(id uint) error {
	var user User
	user.ID = id

	err := runUserValFns(&user, uv.idGreaterThan(0))
	if err != nil {
		return err
	}
	return uv.UserDB.Delete(id)
}

// ByEmail normalizes the email address before passing it to the database layer
// to perform the query
func (uv *userValidator) ByEmail(email string) (*User, error) {
	user := User{
		Email: email,
	}
	err := runUserValFns(&user, uv.normalizeEmail)
	if err != nil {
		return nil, err
	}
	return uv.UserDB.ByEmail(user.Email)
}

// ByRemember normalization: hash the remember token and then pass it
// to UserDB's ByRemember
func (uv *userValidator) ByRemember(token string) (*User, error) {
	user := User{
		Remember: token,
	}
	if err := runUserValFns(&user, uv.hmacRemember); err != nil {
		return nil, err
	}
	return uv.UserDB.ByRemember(user.RememberHash)
}

// hmacRemember calculates and stores in User the remember token hash
func (uv *userValidator) hmacRemember(user *User) error {
	if user.Remember == "" {
		return nil
	}
	user.RememberHash = uv.hmac.Hash(user.Remember)
	return nil
}

// setRememberIfUnset ensures User has a remember token
func (uv *userValidator) setRememberIfUnset(user *User) error {
	if user.Remember != "" {
		return nil
	}
	token, err := rand.RememberToken()
	if err != nil {
		return err
	}
	user.Remember = token
	return nil
}

// idGreaterThan ensures the ID is greater than the provided argument
func (uv *userValidator) idGreaterThan(n uint) userValFn {
	return userValFn(func(user *User) error {
		if user.ID <= n {
			return ErrIDInvalid
		}
		return nil
	})
}

// normalize email address by converting to lower case and trimming whitespace
func (uv *userValidator) normalizeEmail(user *User) error {
	user.Email = strings.ToLower(user.Email)
	user.Email = strings.TrimSpace(user.Email)
	return nil
}

// ensure email address is present
func (uv *userValidator) requireEmail(user *User) error {
	if user.Email == "" {
		return ErrEmailRequired
	}
	return nil
}

// ensure email address matches our regular expression test
func (uv *userValidator) emailFormat(user *User) error {
	if user.Email == "" {
		return nil
	}
	if !uv.emailRegex.MatchString(user.Email) {
		return ErrEmailInvalid
	}
	return nil
}

// ensure specified email is not already being used
// (i.e., is not already present in the database)
func (uv *userValidator) emailIsAvail(user *User) error {
	existing, err := uv.ByEmail(user.Email)
	if err == ErrNotFound {
		// Email address is available: we didn't find a user with this email
		return nil
	}
	if err != nil {
		// some other error occurred, return it
		return err
	}

	// if we get here, a user record in the database uses this email address
	// is it the same user trying to update their existing email address?
	if user.ID != existing.ID {
		// a different user than the one we're updating, so email is taken
		return ErrEmailTaken
	}
	return nil // user is updating their existing email address
}

// ensure password meets minimum length requirement
func (uv *userValidator) passwordMinLength(user *User) error {
	if user.Password == "" { // password required handled elsewhere
		return nil
	}
	if len(user.Password) < 8 {
		return ErrPasswordTooShort
	}
	return nil
}

// ensure password is not empty
func (uv *userValidator) passwordRequired(user *User) error {
	if user.Password == "" {
		return ErrPasswordRequired
	}
	return nil
}

// ensure password hash is not empty
func (uv *userValidator) passwordHashRequired(user *User) error {
	if user.PasswordHash == "" {
		return ErrPasswordRequired
	}
	return nil
}

// ensure Remember hash is >= 32 bytes
func (uv *userValidator) rememberMinBytes(user *User) error {
	if user.Remember == "" {
		// Remember tokens arent' always updated, so see if one is provided.
		// If not, trust other validations will catch any errors.
		return nil
	}
	n, err := rand.NBytes(user.Remember)
	if err != nil {
		return err
	}

	if n < 32 { // CAUTION: hard-coded constant, unlikely to change
		return ErrRememberTooShort
	}

	return nil
}

// ensure Remember hash is provided
func (uv *userValidator) rememberHashRequired(user *User) error {
	if user.RememberHash == "" {
		return ErrRememberRequired
	}
	return nil
}

/* ********** ********** ********** */
/*       userValidator helpers      */

// all user validation/normalization functions implement this signature
// to simplify runUserValFns
type userValFn func(*User) error

// iterate through the sequence of userValFn-conforming validation/normalization functions
func runUserValFns(user *User, fns ...userValFn) error {
	for _, fn := range fns {
		if err := fn(user); err != nil {
			return err
		}
	}
	return nil
}

// bcryptPassword will hash a user's password with an
// app-wide pepper and becrypt, which salts for us
func (uv *userValidator) bcryptPassword(user *User) error {
	if user.Password == "" {
		// Nothing to do if a new password wasn't provided.
		return nil
	}

	pwBytes := []byte(user.Password + userPwPepper)
	hashedBytes, err := bcrypt.GenerateFromPassword(pwBytes, bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user.PasswordHash = string(hashedBytes) // save the PasswordHash in the user object
	user.Password = ""                      // ... but overwrite the Password immediately (does not reach DB)
	return nil
}

/* ********** ********** ********** */
/*            helper methods        */

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
