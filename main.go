package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

var homeTemplate *template.Template

func handlerFunc(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "<h1>Welcome to my awesome site!</h1>")
}

func Home(w http.ResponseWriter, r *http.Request) {
	var err error

	homeTemplate, err = template.ParseFiles("views/home.gohtml")
	if err != nil {
		panic(err)
	}

	w.Header().Set("Content-Type", "text/html")
	if err := homeTemplate.Execute(w, nil); err != nil {
		log.Printf("homeTemplate.Execute() returned error: %v", err)
		os.Exit(1)
	}
}

func Contact(w http.ResponseWriter, r *http.Request) {
	log.Printf("entered Contact")
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprintf(w, "To get in touch, please send an email "+
		"to <a href=\"mailto:support@lenslocked.com\">"+
		"support@lenslocked.com</a>.")
}

func Faq(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprintf(w, "<h1>FAQ</h1><p>This is my awesome FAQ page!</p>")
}

func NotFound(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusNotFound)
	fmt.Fprintf(w, "<h2>404 Not Found</h2><p>Sorry, we could not find "+r.RequestURI+"</p>")
}

func main() {
	// homeTemplate, err = template.ParseFiles("views/home.gohtml")
	// if err != nil {
	// 	panic(err)
	// }

	r := mux.NewRouter()
	r.HandleFunc("/", Home)
	r.HandleFunc("/contact", Contact)
	r.HandleFunc("/faq", Faq)
	r.NotFoundHandler = http.HandlerFunc(NotFound)

	http.ListenAndServe(":3000", r)
}
