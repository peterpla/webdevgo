package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"./views"

	"github.com/gorilla/mux"
)

var homeView *views.View
var contactView *views.View

func home(w http.ResponseWriter, r *http.Request) {
	log.Println("entered home()")
	// log.Printf("homeView: %T, %+v", homeView, *homeView)
	// log.Printf("homeView.Template: %T, %+v", homeView.Template, *homeView.Template)

	w.Header().Set("Content-Type", "text/html")

	if err := homeView.Template.ExecuteTemplate(w, homeView.Layout, nil); err != nil {
		log.Printf("homeView.Execute() returned error: %v", err)
		os.Exit(1)
	}
	log.Println("exiting home()")
}

func contact(w http.ResponseWriter, r *http.Request) {
	log.Println("entered contact()")
	// log.Printf("contactView: %T, %+v", contactView, *contactView)
	// log.Printf("contactView.Template: %T, %+v", contactView.Template, *contactView.Template)

	w.Header().Set("Content-Type", "text/html")

	if err := contactView.Template.ExecuteTemplate(w, contactView.Layout, nil); err != nil {
		log.Printf("contactView.Execute() returned error: %v", err)
		os.Exit(1)
	}
	log.Println("exiting contact()")
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
	homeView = views.NewView("bootstrap", "views/home.gohtml")
	log.Printf("homeView: %T, %+v", homeView, *homeView)
	log.Printf("homeView.Template: %T, %+v", homeView.Template, *homeView.Template)

	contactView = views.NewView("bootstrap", "views/contact.gohtml")
	log.Printf("contactView: %T, %+v", contactView, *contactView)
	log.Printf("contactView.Template: %T, %+v", contactView.Template, *contactView.Template)

	r := mux.NewRouter()
	r.HandleFunc("/", home)
	r.HandleFunc("/contact", contact)
	r.HandleFunc("/faq", Faq)
	r.NotFoundHandler = http.HandlerFunc(NotFound)

	http.ListenAndServe(":3000", r)
}
