package main

import (
	"fmt"

	"../models"

	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
)

const (
	dbHost = "localhost"
	dbPort = 5432
	dbUser = "postgres"
	// password = "" // DO NOT use empty-string password when NO password is set!
	dbName = "whatever_dev"
)

type User struct { // database table "users"
	gorm.Model
	Name   string
	Email  string `gorm:"not null;unique_index"`
	Orders []Order
}

type Order struct {
	gorm.Model
	UserID      uint
	Amount      int
	Description string
}

func main() {

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbName)

	us, err := models.NewUserService(psqlInfo)
	if err != nil {
		panic(err)
	}
	defer us.Close()

	// This will reset the database on every run, but is fine
	// for testing things out.
	us.DestructiveReset()

	// Create a user
	user := models.User{
		Name:  "Michael Scott",
		Email: "michael@dundermifflin.com",
	}
	if err := us.Create(&user); err != nil {
		panic(err)
	}

	foundUser, err := us.ByEmail("michael@dundermifflin.com")
	if err != nil {
		panic(err)
	}
	fmt.Println(foundUser)
}

// func createOrder(db *gorm.DB, user User, amount int, desc string) {
// 	db.Create(&Order{
// 		UserID:      user.ID,
// 		Amount:      amount,
// 		Description: desc,
// 	})
// 	if db.Error != nil {
// 		panic(db.Error)
// 	}
// }
