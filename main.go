package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

var err error
var homeTemplate, contactTemplate *template.Template

func home(w http.ResponseWriter, r *http.Request) {
	homeTemplate, err = template.ParseFiles("views/home.gohtml", "views/layouts/footer.gohtml")
	if err != nil {
		panic(err)
	}

	w.Header().Set("Content-Type", "text/html")
	// log.Printf("http.ResponseWriter w: %T, %+v", w, w)
	if err := homeTemplate.Execute(w, nil); err != nil {
		log.Printf("homeTemplate.Execute() returned error: %v", err)
		os.Exit(1)
	}
}

func contact(w http.ResponseWriter, r *http.Request) {
	contactTemplate, err = template.ParseFiles("views/contact.gohtml", "views/layouts/footer.gohtml")
	if err != nil {
		panic(err)
	}

	w.Header().Set("Content-Type", "text/html")
	if err := contactTemplate.Execute(w, nil); err != nil {
		log.Printf("contactTemplate.Execute() returned error: %v", err)
		os.Exit(1)
	}
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
	// contactTemplate, err = template.ParseFiles("views/contact.gohtml")
	// if err != nil {
	// 	panic(err)
	// }

	r := mux.NewRouter()
	r.HandleFunc("/", home)
	r.HandleFunc("/contact", contact)
	r.HandleFunc("/faq", Faq)
	r.NotFoundHandler = http.HandlerFunc(NotFound)

	http.ListenAndServe(":3000", r)
}
