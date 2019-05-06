package main

import (
	"fmt"

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

	db, err := gorm.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	db.LogMode(true)
	db.AutoMigrate(&User{}, &Order{})

	// postgresql://[user[:password]@][netloc][:port][/dbname]
	fmt.Printf("Successfully connected! postgresql://%s:\"%s\"@%s:%d/%s\n", dbUser, "", dbHost, dbPort, dbName)

	// retrieve orders by user
	var newU User
	newU.Name = "Peter Plamondon"
	fmt.Printf("Before: newU: %+v", newU)

	// db.Preload("Orders").Find(&newU)
	db.Preload("Orders").First(&newU)
	if db.Error != nil {
		panic(db.Error)
	}
	fmt.Printf("After: newU: %+v", newU)

	fmt.Printf("Email: %s\n", newU.Email)
	fmt.Printf("Number of orders: %d\n", len(newU.Orders))
	fmt.Printf("Orders: %+v\n", newU.Orders)
}

func createOrder(db *gorm.DB, user User, amount int, desc string) {
	db.Create(&Order{
		UserID:      user.ID,
		Amount:      amount,
		Description: desc,
	})
	if db.Error != nil {
		panic(db.Error)
	}
}
