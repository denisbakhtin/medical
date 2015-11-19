package controllers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/denisbakhtin/medical/helpers"
	"github.com/denisbakhtin/medical/models"
)

//SignIn handles /signin route
func SignIn(w http.ResponseWriter, r *http.Request) {
	tmpl := helpers.Template(r)
	session := helpers.Session(r)
	data := helpers.DefaultData(r)
	T := helpers.T(r)

	if r.Method == "GET" {

		data["Title"] = T("sign_in")
		data["Active"] = "signin"
		data["Flash"] = session.Flashes()
		session.Save(r, w)
		tmpl.Lookup("auth/signin").Execute(w, data)

	} else if r.Method == "POST" {

		user := &models.User{
			Email:    r.PostFormValue("email"),
			Password: r.PostFormValue("password"),
		}
		//check existence
		userDB, _ := models.GetUserByEmail(user.Email)
		if userDB.ID == 0 {
			log.Printf("ERROR: Login failed, IP: %s, Email: %s\n", r.RemoteAddr, user.Email)
			session.AddFlash("Email or password incorrect")
			session.Save(r, w)
			http.Redirect(w, r, "/signin", 303)
			return
		}
		//create user
		if err := userDB.ComparePassword(user.Password); err != nil {
			log.Printf("ERROR: Login failed, IP: %s, Email: %s\n", r.RemoteAddr, user.Email)
			session.AddFlash("Email or password incorrect")
			session.Save(r, w)
			http.Redirect(w, r, "/signin", 303)
			return
		}

		session.Values["user_id"] = userDB.ID
		session.Save(r, w)
		http.Redirect(w, r, "/", 303)

	} else {
		err := fmt.Errorf("Method %q not allowed", r.Method)
		log.Printf("ERROR: %s\n", err)
		tmpl.Lookup("errors/405").Execute(w, helpers.ErrorData(err))
	}
}

//SignUp handles /signup route
func SignUp(w http.ResponseWriter, r *http.Request) {
	tmpl := helpers.Template(r)
	session := helpers.Session(r)
	data := helpers.DefaultData(r)
	T := helpers.T(r)

	if r.Method == "GET" {

		data["Title"] = T("sign_up")
		data["Active"] = "signup"
		data["Flash"] = session.Flashes()
		session.Save(r, w)
		tmpl.Lookup("auth/signup").Execute(w, data)

	} else if r.Method == "POST" {

		user := &models.User{
			Email:    r.PostFormValue("email"),
			Password: r.PostFormValue("password"),
		}
		//check existence
		userDB, _ := models.GetUserByEmail(user.Email)
		if userDB.ID != 0 {
			session.AddFlash("User exists")
			session.Save(r, w)
			http.Redirect(w, r, "/signup", 303)
			return
		}
		//create user
		err := user.HashPassword()
		if err != nil {
			session.AddFlash("Error whilst registering user.")
			session.Save(r, w)
			log.Printf("ERROR: can't register user: %v", err)
			http.Redirect(w, r, "/signup", 303)
			return
		}
		if err := user.Insert(); err != nil {
			session.AddFlash("Error whilst registering user.")
			session.Save(r, w)
			log.Printf("ERROR: can't register user: %v", err)
			http.Redirect(w, r, "/signup", 303)
			return
		}
		session.Values["user_id"] = user.ID
		session.Save(r, w)
		http.Redirect(w, r, "/", 303)

	} else {
		err := fmt.Errorf("Method %q not allowed", r.Method)
		log.Printf("ERROR: %s\n", err)
		tmpl.Lookup("errors/405").Execute(w, helpers.ErrorData(err))
	}
}

//Logout handles /logout route
func Logout(w http.ResponseWriter, r *http.Request) {
	//any method will do :3
	session := helpers.Session(r)
	delete(session.Values, "user_id")
	session.Save(r, w)
	http.Redirect(w, r, "/", 303)
}
