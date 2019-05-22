package controllers

import (
	"fmt"
	"net/http"

	"../models"
	"../rand"
	"../views"
)

type Users struct {
	NewView   *views.View
	LoginView *views.View
	us        *models.UserService
}

func NewUsers(us *models.UserService) *Users {
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
	if err := u.NewView.Render(w, nil); err != nil {
		panic(err)
	}
}

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
	var form SignupForm
	if err := parseForm(r, &form); err != nil {
		panic(err)
	}
	user := models.User{
		Name:     form.Name,
		Email:    form.Email,
		Password: form.Password,
	}

	if err := u.us.Create(&user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	/*
		fmt.Fprintf(w, "User is: %+v\n", user) // echo to web page
		fmt.Printf("User is: %+v\n", user)     // echo to stdout
	*/

	// sign in the newly-created user
	err := u.signIn(w, &user)
	if err != nil {
		// TODO: replace the temporary debugging message below
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// redirect to the cookie test page to test the cookie
	http.Redirect(w, r, "/cookietest", http.StatusFound)
}

type LoginForm struct {
	Email    string `schema:"email"`
	Password string `schema:"password"`
}

// Login is used to process the login form when a user
// tries to log in as an existing user (via email & pw)
//
// POST /login
func (u *Users) Login(w http.ResponseWriter, r *http.Request) {
	form := LoginForm{}
	if err := parseForm(r, &form); err != nil {
		panic(err)
	}

	user, err := u.us.Authenticate(form.Email, form.Password)

	if err != nil {
		switch err {
		case models.ErrNotFound:
			fmt.Fprintln(w, "Invalid email address.")
		case models.ErrInvalidPassword:
			fmt.Fprintln(w, "Invalid password.")
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// successful Authenticate, sign in the user
	err = u.signIn(w, user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/cookietest", http.StatusFound)
}

// signIn is used to sign in the given user via cookies
func (u *Users) signIn(w http.ResponseWriter, user *models.User) error {
	// the user argument should have been saved in the database earlier

	// if a remember token was not provided, create one
	if user.Remember == "" {
		token, err := rand.RememberToken()
		if err != nil {
			return err
		}
		user.Remember = token

		// update the user's database record with the new remember token
		// so we can look it up later
		err = u.us.Update(user)
		if err != nil {
			return err
		}
	}

	// add the remember token to a cookie
	cookie := http.Cookie{
		Name:  "remember_token",
		Value: user.Remember,
	}
	http.SetCookie(w, &cookie)
	return nil
}

// CookieTest is used to display cookies set on the current user
func (u *Users) CookieTest(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("email")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprintln(w, "Email is:", cookie.Value)
}
