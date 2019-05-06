package main

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

const (
	host   = "localhost"
	port   = 5432
	user   = "postgres"
	dbname = "whatever_dev"
)

func main() {

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"dbname=%s sslmode=disable",
		host, port, user, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		panic(err)
	}

	fmt.Println("Successfully connected!")

	// var name = "Jon Calhoun"
	// var email = "jon@calhoun.io"

	// _, err = db.Exec(`
	//   INSERT INTO users(name, email)
	//   VALUES ($1, $2)`,
	// 	name, email)
	// if err != nil {
	// 	panic(err)
	// }

	// fmt.Printf("Successfully inserted %s, %s!\n", name, email)

	// name = "Jon2 Calhoun2"
	// email = "jon2@calhoun2.io"

	// row := db.QueryRow(`
	//   INSERT INTO users(name, email)
	//   VALUES($1, $2) RETURNING id`,
	// 	name, email)

	// err = row.Scan(&id)
	// if err != nil {
	// 	panic(err)
	// }

	// fmt.Printf("Successfully inserted %s, %s!\n", name, email)

	var id int
	var name, email string

	rows, err := db.Query(`
  	  SELECT id, name, email
  	  FROM users
		WHERE email = $1
		OR ID > $2`,
		"jon@calhoun.io", 3)
	if err != nil {
		panic(err)
	}

	for rows.Next() {
		rows.Scan(&id, &name, &email)
		fmt.Println("ID:", id, "Name:", name, "Email:", email)
	}

	db.Close()
}
