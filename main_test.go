package main

import (
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestViewHandlers(t *testing.T) {

	type testset struct {
		method   string
		url      string
		handler  http.HandlerFunc
		expected int
	}

	var tests = []testset{
		{"GET", "/contact", Contact, http.StatusOK},
		{"GET", "/faq", Faq, http.StatusOK},
		{"GET", "/blah", NotFound, http.StatusNotFound},
		{"GET", "/", Home, http.StatusOK},
	}

	for _, r := range tests {
		log.Printf("*** current test: method: %s, url: %s, handler: %T / %p, expected: %d", r.method, r.url, r.handler, r.handler, r.expected)

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
		log.Printf("response Code: %d", rr.Code)

		// check the status code is what we expect.
		if status := rr.Code; status != r.expected {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, r.expected)
		} else {
			log.Println("*** test PASSED")
		}
	}

}
