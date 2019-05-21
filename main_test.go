package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"./controllers"
	"./models"
)

func TestViewHandlers(t *testing.T) {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbName)
	us, err := models.NewUserService(psqlInfo)
	if err != nil {
		panic(err)
	}
	defer us.Close()

	type testset struct {
		method   string
		url      string
		handler  http.HandlerFunc
		expected int
	}

	staticC := controllers.NewStatic()
	usersC := controllers.NewUsers(us)
	galleriesC := controllers.NewGalleries()

	var tests = []testset{
		{"GET", "/blah", NotFound, http.StatusNotFound},
		{"GET", "/contact", staticC.Contact.ServeHTTP, http.StatusOK},
		{"GET", "/faq", staticC.Faq.ServeHTTP, http.StatusOK},
		{"GET", "/", staticC.Home.ServeHTTP, http.StatusOK},
		{"GET", "/galleries/new", galleriesC.Gallery.ServeHTTP, http.StatusOK},
		{"GET", "/signup", usersC.New, http.StatusOK},
		// {"POST", "/signup", usersC.Create, http.StatusOK}, // need to populate form body
	}

	for _, r := range tests {
		log.Printf("*** current test: method: %s, url: %s, handler: %T, expected: %d", r.method, r.url, r.handler, r.expected)

		// create a request to pass to our handler
		// we don't have any query parameters for now, so we'll pass
		// 'nil' as the third parameter
		req, err := http.NewRequest(r.method, r.url, nil)
		if err != nil {
			t.Fatal(err)
		}

		// create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response
		rr := httptest.NewRecorder()

		handler := http.HandlerFunc(r.handler)

		// our handlers satisfy http.Handler, so we can call their ServeHTTP method
		// directly and pass in our Request and ResponseRecorder.
		handler.ServeHTTP(rr, req)

		// check the status code is what we expect.
		if status := rr.Code; status != r.expected {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, r.expected)
		}
	}

}
