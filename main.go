package main

import (
	"fmt"
	"net/http"

	"./controllers"
	"./views"

	"github.com/gorilla/mux"
)

var homeView *views.View
var contactView *views.View

// var signupView *views.View
var faqView *views.View

func home(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	if err := homeView.Render(w, nil); err != nil {
		panic(err)
	}
}

func contact(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	if err := contactView.Render(w, nil); err != nil {
		panic(err)
	}
}

// func signup(w http.ResponseWriter, r *http.Request) {
// 	w.Header().Set("Content-Type", "text/html")
// 	if err := signupView.Render(w, nil); err != nil {
// 		panic(err)
// 	}
// }

func faq(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	if err := faqView.Render(w, nil); err != nil {
		panic(err)
	}
}

func NotFound(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusNotFound)
	fmt.Fprintf(w, "<h2>404 Not Found</h2><p>Sorry, we could not find "+r.RequestURI+"</p>")
}

func main() {
	homeView = views.NewView("bootstrap", "views/home.gohtml")
	contactView = views.NewView("bootstrap", "views/contact.gohtml")
	// signupView = views.NewView("bootstrap", "views/signup.gohtml")
	faqView = views.NewView("bootstrap", "views/faq.gohtml")

	usersC := controllers.NewUsers()

	r := mux.NewRouter()
	r.HandleFunc("/", home).Methods("GET")
	r.HandleFunc("/contact", contact).Methods("GET")
	r.HandleFunc("/faq", faq).Methods("GET")

	r.HandleFunc("/signup", usersC.New).Methods("GET")
	r.HandleFunc("/signup", usersC.Create).Methods("POST")

	r.NotFoundHandler = http.HandlerFunc(NotFound)

	http.ListenAndServe(":3000", r)
}
