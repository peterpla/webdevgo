package controllers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/peterpla/webdevgo/models"
	"github.com/peterpla/webdevgo/rand"
	"github.com/peterpla/webdevgo/views"
)

// Users ... [add documentation]
type Users struct {
	NewView   *views.View
	LoginView *views.View
	us        models.UserService
}

// NewUsers ... [add documentation]
func NewUsers(us models.UserService) *Users {
	return &Users{
		NewView:   views.NewView("bootstrap", "users/new"),
		LoginView: views.NewView("bootstrap", "users/login"),
		us:        us,
	}
}

// New is used to render the form where a user can
// create a new user account
//
// GET /signup
func (u *Users) New(w http.ResponseWriter, r *http.Request) {
	u.NewView.Render(w, nil)
}

// SignupForm ... [add documentation]
type SignupForm struct {
	Name     string `schema:"name"`
	Email    string `schema:"email"`
	Password string `schema:"password"`
}

// Create is used to process the signup form when a user
// tries to create a new user account
//
// POST /signup
func (u *Users) Create(w http.ResponseWriter, r *http.Request) {
	var vd views.Data
	var form SignupForm

	if err := parseForm(r, &form); err != nil {
		vd.SetAlert(err)
		u.NewView.Render(w, vd)
		return
	}
	user := models.User{
		Name:     form.Name,
		Email:    form.Email,
		Password: form.Password,
	}

	if err := u.us.Create(&user); err != nil {
		vd.SetAlert(err)
		u.NewView.Render(w, vd)
		return
	}

	// fmt.Fprintf(w, "User is: %+v\n", user) // echo to web page
	// fmt.Printf("User is: %+v\n", user)     // echo to stdout

	// sign in the newly-created user
	// remember token is NOT provided
	err := u.signIn(w, &user)
	if err != nil {
		// user resource was created, but we could not signin; likely
		// a transient problem, so redirect the user to login. Not optimal, but
		// we think this won't happen often.
		//
		// Log it so we can see if that's a valid assumption :)
		log.Println("RARE: after user creation, signin failed")
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	// redirect to to pull Remember token from the user's cookie
	// and confirm it matches stored RememberHash
	http.Redirect(w, r, "/cookietest", http.StatusFound)
}

// LoginForm ... [add documentation]
type LoginForm struct {
	Email    string `schema:"email"`
	Password string `schema:"password"`
}

// Login is used to process the login form when a user
// tries to log in as an existing user (via email & pw)
//
// POST /login
func (u *Users) Login(w http.ResponseWriter, r *http.Request) {
	var vd views.Data
	var form LoginForm

	if err := parseForm(r, &form); err != nil {
		vd.SetAlert(err)
		u.LoginView.Render(w, vd)
		return
	}

	user, err := u.us.Authenticate(form.Email, form.Password)

	if err != nil {
		switch err {
		case models.ErrNotFound:
			vd.AlertError("No user exists with that email address")
		default:
			vd.SetAlert(err)
		}
		u.LoginView.Render(w, vd)
		return
	}

	// SUCCESS from Authenticate, sign in the user
	// Remember token is NOT populated in User object
	err = u.signIn(w, user)
	if err != nil {
		vd.SetAlert(err)
		u.LoginView.Render(w, vd)
		return
	}

	// redirect to to pull Remember token from the user's cookie
	// and confirm it matches stored RememberHash
	http.Redirect(w, r, "/cookietest", http.StatusFound)
}

// signIn is used to sign in the given user, confirming the Remember token
// from the user's cookie hashes to the RememberHash value stored in the user's
// DB record
func (u *Users) signIn(w http.ResponseWriter, user *models.User) error {
	// at Signup/Create, the user was saved in the database

	if user.Remember == "" { //
		token, err := rand.RememberToken() // generate a new Remember token
		if err != nil {
			return err
		}
		user.Remember = token

		// update the user's record with the RememberHash (but NOT the Remember token!)
		// and write the user's record to the DB so we can look it up later
		err = u.us.Update(user)
		if err != nil {
			return err
		}
	} else {
		fmt.Printf("user.Remember unexpectedly not empty: \"%s\"\n", user.Remember)
	}

	// add the remember token to a cookie, the only place it is stored
	cookie := http.Cookie{
		Name:     "remember_token",
		Value:    user.Remember,
		HttpOnly: true,
	}
	http.SetCookie(w, &cookie)
	return nil
}

// CookieTest displays the cookie set on the current user
func (u *Users) CookieTest(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("remember_token")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	user, err := u.us.ByRemember(cookie.Value)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Println(w, user)
}
