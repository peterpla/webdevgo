package models

import (
	"errors"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type User struct {
	gorm.Model
	Name  string
	Email string
	Age   uint // new.gohtml enforces min=18, max=120
}

type UserService struct {
	db *gorm.DB
}

// NewUserService returns a connection to the database holding User objects
func NewUserService(connectionInfo string) (*UserService, error) {
	// log.Printf("enter NewUserService, connectionInfo: %s", connectionInfo)
	db, err := gorm.Open("postgres", connectionInfo)
	if err != nil {
		return nil, err
	}
	db.LogMode(true)
	// postgresql://[user[:password]@][netloc][:port][/dbname]
	// fmt.Println("Successfully connected to database!")

	return &UserService{
		db: db,
	}, nil
}

// Close the UserService database connection
func (us *UserService) Close() error {
	// log.Printf("enter UserService.Close")
	return us.db.Close()
}

var (
	// ErrNotFound is returned when the query executes successfully
	// but returned zero rows. I.e., the resource cannot be found
	// in the database.
	ErrNotFound = errors.New("models: resource not found")

	// ErrInvalidID is returned when an invalid ID is provided
	// to a method like Delete.
	ErrInvalidID = errors.New("models: ID provided was invalid")
)

// Create will create the provided User record in the database,
// and backfill gorm.Model data including the ID, CreatedAt, and
// UpdatedAt fields.
func (us *UserService) Create(user *User) error {
	return us.db.Create(user).Error
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
func (us *UserService) ByID(id uint) (*User, error) {
	var user User
	db := us.db.Where("id = ?", id)
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
func (us *UserService) ByEmail(email string) (*User, error) {
	var user User
	db := us.db.Where("email = ?", email)
	err := first(db, &user)
	return &user, err
}

// ByAge will look up the first user with the provided Age.
// If the user is found, return a nil error
// If the user is not found, return ErrNotFound
// If there is another error, return and error with
// more information about what went wrong. This
// may not be an error generated by the models package.
func (us *UserService) ByAge(age uint) (*User, error) {
	var user User
	db := us.db.Where("age = ?", age)
	err := first(db, &user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// Update will update the provided user with all of the data
// in the provided User object
func (us *UserService) Update(user *User) error {
	return us.db.Save(user).Error
}

// Delete will delete the user with the provided ID
func (us *UserService) Delete(id uint) error {
	if id == 0 {
		return ErrInvalidID
	}
	user := User{Model: gorm.Model{ID: id}}
	return us.db.Delete(&user).Error
}

// DestructiveReset drops the user table and rebuilds it
func (us *UserService) DestructiveReset() error {
	err := us.db.DropTableIfExists(&User{}).Error
	if err != nil {
		return err
	}
	return us.AutoMigrate()
}

// AutoMigrate will attempt to automaticaly migrate
// the Users table
func (us *UserService) AutoMigrate() error {
	if err := us.db.AutoMigrate(&User{}).Error; err != nil {
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
