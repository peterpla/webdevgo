package main

import (
	"fmt"

	"github.com/peterpla/webdevgo/models"
)

const (
	dbHost = "localhost"
	dbPort = 5432
	dbUser = "postgres"
	// password = "" // DO NOT use empty-string password when NO password is set!
	dbName = "whatever_dev"
)

func main() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbName)
	us, err := models.NewUserService(psqlInfo)
	if err != nil {
		panic(err)
	}
	defer us.Close()
	us.DestructiveReset()

	user := models.User{
		Name:     "Michael Scott",
		Email:    "michael@dundermifflin.com",
		Password: "bestboss",
	}

	err = us.Create(&user)
	if err != nil {
		panic(err)
	}

	// Verify that the user object has its Remember and RememberHash
	// set during Create()
	fmt.Printf("User: %+v\n", user)
	if user.Remember == "" {
		panic("Invalid remember token")
	}

	// Verify we can lookup a user with that remember token
	user2, err := us.ByRemember(user.Remember)
	if err != nil {
		panic(err)
	}
	fmt.Printf("User2: %+v\n", user2)

}
