package main

import (
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
)

const (
	host = "localhost"
	port = 5432
	user = "postgres"
	// password = "" // DO NOT use empty-string password when NO password is set!
	dbname = "whatever_dev"
)

type User struct { // database table "users"
	gorm.Model
	Name  string
	Email string `gorm:"not null;unique_index"`
}

func main() {

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"dbname=%s sslmode=disable",
		host, port, user, dbname)

	db, err := gorm.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	db.LogMode(true)

	// postgresql://[user[:password]@][netloc][:port][/dbname]
	fmt.Printf("Successfully connected! postgresql://%s:\"%s\"@%s:%d/%s\n", user, "", host, port, dbname)

	var users []User
	db.Find(&users)
	if db.Error != nil {
		panic(db.Error)
	}
	fmt.Println("Retrieved", len(users), "users.")
	fmt.Println(users)
}
