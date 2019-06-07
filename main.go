// Package main ... [TODO: add documentation]
package main

import (
	"fmt"
	"net/http"

	"github.com/peterpla/webdevgo/controllers"
	"github.com/peterpla/webdevgo/models"
	"github.com/peterpla/webdevgo/views"

	"github.com/gorilla/mux"
)

const (
	dbHost = "localhost"
	dbPort = 5432
	dbUser = "postgres"
	// password = "" // DO NOT use empty-string password when NO password is set!
	dbName = "whatever_dev"
)

var homeView *views.View
var contactView *views.View
var faqView *views.View

// NotFound produces 404 Not Found responses for not-found URLs
func NotFound(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusNotFound)
	fmt.Fprintf(w, "<h2>404 Not Found</h2><p>Sorry, we could not find "+r.RequestURI+"</p>")
}

func main() {
	// initialize views
	homeView = views.NewView("bootstrap", "static/home")
	contactView = views.NewView("bootstrap", "static/contact")
	faqView = views.NewView("bootstrap", "static/faq")

	// initialize controllers
	staticC := controllers.NewStatic()
	galleriesC := controllers.NewGalleries()

	// the Users controller requires a UserService, make that happen
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbName)
	us, err := models.NewUserService(psqlInfo)
	if err != nil {
		panic(err)
	}
	defer us.Close()
	// us.DestructiveReset()

	us.AutoMigrate()

	usersC := controllers.NewUsers(us)

	// define routing
	r := mux.NewRouter()
	r.Handle("/", staticC.Home).Methods("GET")
	r.Handle("/contact", staticC.Contact).Methods("GET")
	r.Handle("/faq", staticC.Faq).Methods("GET")
	r.Handle("/gallery", galleriesC.Gallery).Methods("GET")

	r.HandleFunc("/signup", usersC.New).Methods("GET")
	r.HandleFunc("/signup", usersC.Create).Methods("POST")

	r.Handle("/login", usersC.LoginView).Methods("GET")
	r.HandleFunc("/login", usersC.Login).Methods("POST")

	r.HandleFunc("/cookietest", usersC.CookieTest).Methods("GET")

	r.NotFoundHandler = http.HandlerFunc(NotFound)

	http.ListenAndServe(":3000", r)
}
