package main

import (
	"fmt"
	"net/http"
)

func main() {

	fmt.Printf("index function type: %T\n", index)
	// r := httprouter.New()

	// r.HandlerFunc(("GET", "/", index)

}

func index(w http.ResponseWriter, r *http.Request) {
	w.Header()
}
