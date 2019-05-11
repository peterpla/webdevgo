package models

import (
	"fmt"
	"log"
	"testing"
)

func TestUsers(t *testing.T) {
	log.Printf("entered TestUsers")

	testCases := []struct {
		name           string
		dbHost         string // "localhost"
		dbPort         uint   // 5432
		dbUser         string // "postgres"
		dbSkipPassword bool   // true if skip password portion of connection string
		dbPassword     string // ignored if dbSkipPassword == true
		dbName         string // "whatever_dev"
	}{
		{"Basic pass", "localhost", 5432, "postgres", true, "", "ignore_test"},
		{"Fail with bad password", "localhost", 5432, "postgres", false, "badPassword", "ignore_test"},
	}

	var psqlInfo string
	for _, r := range testCases {
		t.Run(r.name, func(t *testing.T) {
			if r.dbSkipPassword {
				// nil password, skip password portion of connection string
				psqlInfo = fmt.Sprintf("host=%s port=%d user=%s "+
					"dbname=%s sslmode=disable",
					r.dbHost, r.dbPort, r.dbUser, r.dbName)
			} else {
				// include non-nil password in connection string
				psqlInfo = fmt.Sprintf("host=%s port=%d user=%s password=%s "+
					"dbname=%s sslmode=disable",
					r.dbHost, r.dbPort, r.dbUser, r.dbPassword, r.dbName)
			}

			// connect to the user database
			us, err := NewUserService(psqlInfo)
			if err != nil {
				t.Log(err)
				t.Fail()
			}

			// close the User Service connection created by NewUserService
			us.Close()
		})
	}
	t.Fatalf("incomplete implementation, results invalid")
}
