package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

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

	db.AutoMigrate(&User{})

	name, email := getInfo()

	u := &User{
		Name:  name,
		Email: email,
	}
	if err = db.Create(u).Error; err != nil {
		panic(err)
	}
	fmt.Printf("Created record: %+v\n", u)
}

func getInfo() (name, email string) {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("What is your name?")
	name, _ = reader.ReadString('\n')
	name = strings.TrimSpace(name)

	fmt.Println("What is your email?")
	email, _ = reader.ReadString('\n')
	email = strings.TrimSpace(email)

	return name, email
}
