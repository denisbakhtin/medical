package controllers

import (
	"log"
	"net/http"

	"fmt"

	"github.com/denisbakhtin/medical/helpers"
	"github.com/denisbakhtin/medical/models"
)

//UserIndex handles GET /admin/users route
func UserIndex(w http.ResponseWriter, r *http.Request) {
	tmpl := helpers.Template(r)
	data := helpers.DefaultData(r)
	T := helpers.T(r)
	db := models.GetDB()
	if r.Method == "GET" {

		var list []models.User
		db.Find(&list)
		data["Title"] = T("users")
		data["Active"] = "users"
		data["List"] = list
		tmpl.Lookup("users/index").Execute(w, data)

	} else {
		err := fmt.Errorf("Method %q not allowed", r.Method)
		log.Printf("ERROR: %s\n", err)
		w.WriteHeader(405)
		tmpl.Lookup("errors/405").Execute(w, helpers.ErrorData(err))
	}
}

//UserCreate handles /admin/new_user route
func UserCreate(w http.ResponseWriter, r *http.Request) {
	tmpl := helpers.Template(r)
	session := helpers.Session(r)
	data := helpers.DefaultData(r)
	T := helpers.T(r)
	db := models.GetDB()
	if r.Method == "GET" {

		data["Title"] = T("new_user")
		data["Active"] = "users"
		data["Flash"] = session.Flashes()
		session.Save(r, w)
		tmpl.Lookup("users/form").Execute(w, data)

	} else if r.Method == "POST" {

		user := &models.User{
			Name:     r.PostFormValue("name"),
			Email:    r.PostFormValue("email"),
			Password: r.PostFormValue("password"),
		}

		if err := db.Create(user).Error; err != nil {
			session.AddFlash(err.Error())
			session.Save(r, w)
			http.Redirect(w, r, "/admin/new_user", 303)
			return
		}
		http.Redirect(w, r, "/admin/users", 303)

	} else {
		err := fmt.Errorf("Method %q not allowed", r.Method)
		log.Printf("ERROR: %s\n", err)
		w.WriteHeader(405)
		tmpl.Lookup("errors/405").Execute(w, helpers.ErrorData(err))
	}
}

//UserUpdate handles /admin/edit_user/:id route
func UserUpdate(w http.ResponseWriter, r *http.Request) {
	tmpl := helpers.Template(r)
	session := helpers.Session(r)
	data := helpers.DefaultData(r)
	T := helpers.T(r)
	db := models.GetDB()
	if r.Method == "GET" {

		id := r.URL.Path[len("/admin/edit_user/"):]
		user := &models.User{}
		db.First(user, id)
		if user.ID == 0 {
			w.WriteHeader(404)
			tmpl.Lookup("errors/404").Execute(w, nil)
			return
		}

		data["Title"] = T("edit_user")
		data["Active"] = "users"
		data["User"] = user
		data["Flash"] = session.Flashes()
		session.Save(r, w)
		tmpl.Lookup("users/form").Execute(w, data)

	} else if r.Method == "POST" {

		user := &models.User{
			ID:       helpers.Atouint(r.PostFormValue("id")),
			Name:     r.PostFormValue("name"),
			Email:    r.PostFormValue("email"),
			Password: r.PostFormValue("password"),
		}

		if err := db.Save(user).Error; err != nil {
			session.AddFlash(err.Error())
			session.Save(r, w)
			http.Redirect(w, r, r.RequestURI, 303)
			return
		}
		http.Redirect(w, r, "/admin/users", 303)

	} else {
		err := fmt.Errorf("Method %q not allowed", r.Method)
		log.Printf("ERROR: %s\n", err)
		w.WriteHeader(405)
		tmpl.Lookup("errors/405").Execute(w, helpers.ErrorData(err))
	}
}

//UserDelete handles /admin/delete_user route
func UserDelete(w http.ResponseWriter, r *http.Request) {
	tmpl := helpers.Template(r)
	db := models.GetDB()

	if r.Method == "POST" {

		user := &models.User{}
		db.First(user, r.PostFormValue("id"))
		if user.ID == 0 {
			w.WriteHeader(404)
			tmpl.Lookup("errors/404").Execute(w, nil)
		}

		if err := db.Delete(user).Error; err != nil {
			log.Printf("ERROR: %s\n", err)
			w.WriteHeader(500)
			tmpl.Lookup("errors/500").Execute(w, helpers.ErrorData(err))
			return
		}
		http.Redirect(w, r, "/admin/users", 303)

	} else {
		err := fmt.Errorf("Method %q not allowed", r.Method)
		log.Printf("ERROR: %s\n", err)
		w.WriteHeader(405)
		tmpl.Lookup("errors/405").Execute(w, helpers.ErrorData(err))
	}
}
