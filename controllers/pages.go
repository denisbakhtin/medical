package controllers

import (
	"fmt"
	"log"
	"net/http"
	"regexp"

	"github.com/denisbakhtin/medical/helpers"
	"github.com/denisbakhtin/medical/models"
)

//PageShow handles /pages/:id route
func PageShow(w http.ResponseWriter, r *http.Request) {
	tmpl := helpers.Template(r)
	data := helpers.DefaultData(r)
	if r.Method == "GET" {

		re := regexp.MustCompile("^[0-9]+")
		id := re.FindString(r.URL.Path[len("/pages/"):])
		page, err := models.GetPage(id)
		if err != nil || !page.Published {
			w.WriteHeader(404)
			tmpl.Lookup("errors/404").Execute(w, nil)
			return
		}
		//redirect to canonical url
		if r.URL.Path != page.Url() {
			http.Redirect(w, r, page.Url(), http.StatusSeeOther)
			return
		}
		data["Page"] = page
		data["Title"] = page.Name
		data["Active"] = page.Url()
		tmpl.Lookup("pages/show").Execute(w, data)

	} else {
		err := fmt.Errorf("Method %q not allowed", r.Method)
		log.Printf("ERROR: %s\n", err)
		w.WriteHeader(405)
		tmpl.Lookup("errors/405").Execute(w, helpers.ErrorData(err))
	}
}

//PageIndex handles GET /admin/pages route
func PageIndex(w http.ResponseWriter, r *http.Request) {
	tmpl := helpers.Template(r)
	data := helpers.DefaultData(r)
	T := helpers.T(r)
	if r.Method == "GET" {

		list, err := models.GetPages()
		if err != nil {
			w.WriteHeader(500)
			tmpl.Lookup("errors/500").Execute(w, helpers.ErrorData(err))
			return
		}
		data["Title"] = T("pages")
		data["Active"] = "pages"
		data["List"] = list
		tmpl.Lookup("pages/index").Execute(w, data)

	} else {
		err := fmt.Errorf("Method %q not allowed", r.Method)
		log.Printf("ERROR: %s\n", err)
		w.WriteHeader(405)
		tmpl.Lookup("errors/405").Execute(w, helpers.ErrorData(err))
	}
}

//PageCreate handles /admin/new_page route
func PageCreate(w http.ResponseWriter, r *http.Request) {
	tmpl := helpers.Template(r)
	session := helpers.Session(r)
	data := helpers.DefaultData(r)
	T := helpers.T(r)
	if r.Method == "GET" {

		data["Title"] = T("new_page")
		data["Active"] = "pages"
		data["Flash"] = session.Flashes()
		session.Save(r, w)
		tmpl.Lookup("pages/form").Execute(w, data)

	} else if r.Method == "POST" {

		page := &models.Page{
			Name:      r.PostFormValue("name"),
			Slug:      r.PostFormValue("slug"),
			Content:   r.PostFormValue("content"),
			Published: helpers.Atob(r.PostFormValue("published")),
		}

		if err := page.Insert(); err != nil {
			session.AddFlash(err.Error())
			session.Save(r, w)
			http.Redirect(w, r, "/admin/new_page", 303)
			return
		}
		http.Redirect(w, r, "/admin/pages", 303)

	} else {
		err := fmt.Errorf("Method %q not allowed", r.Method)
		log.Printf("ERROR: %s\n", err)
		w.WriteHeader(405)
		tmpl.Lookup("errors/405").Execute(w, helpers.ErrorData(err))
	}
}

//PageUpdate handles /admin/edit_page/:id route
func PageUpdate(w http.ResponseWriter, r *http.Request) {
	tmpl := helpers.Template(r)
	session := helpers.Session(r)
	data := helpers.DefaultData(r)
	T := helpers.T(r)
	if r.Method == "GET" {

		id := r.URL.Path[len("/admin/edit_page/"):]
		page, err := models.GetPage(id)
		if err != nil {
			w.WriteHeader(400)
			tmpl.Lookup("errors/400").Execute(w, helpers.ErrorData(err))
			return
		}

		data["Title"] = T("edit_page")
		data["Active"] = "pages"
		data["Page"] = page
		data["Flash"] = session.Flashes()
		session.Save(r, w)
		tmpl.Lookup("pages/form").Execute(w, data)

	} else if r.Method == "POST" {

		page := &models.Page{
			ID:        helpers.Atoi64(r.PostFormValue("id")),
			Name:      r.PostFormValue("name"),
			Slug:      r.PostFormValue("slug"),
			Content:   r.PostFormValue("content"),
			Published: helpers.Atob(r.PostFormValue("published")),
		}

		if err := page.Update(); err != nil {
			session.AddFlash(err.Error())
			session.Save(r, w)
			http.Redirect(w, r, r.RequestURI, 303)
			return
		}
		http.Redirect(w, r, "/admin/pages", 303)

	} else {
		err := fmt.Errorf("Method %q not allowed", r.Method)
		log.Printf("ERROR: %s\n", err)
		w.WriteHeader(405)
		tmpl.Lookup("errors/405").Execute(w, helpers.ErrorData(err))
	}
}

//PageDelete handles /admin/delete_page route
func PageDelete(w http.ResponseWriter, r *http.Request) {
	tmpl := helpers.Template(r)

	if r.Method == "POST" {

		page, err := models.GetPage(r.PostFormValue("id"))
		if err != nil {
			log.Printf("ERROR: %s\n", err)
			w.WriteHeader(404)
			tmpl.Lookup("errors/404").Execute(w, helpers.ErrorData(err))
		}

		if err := page.Delete(); err != nil {
			log.Printf("ERROR: %s\n", err)
			w.WriteHeader(500)
			tmpl.Lookup("errors/500").Execute(w, helpers.ErrorData(err))
			return
		}
		http.Redirect(w, r, "/admin/pages", 303)

	} else {
		err := fmt.Errorf("Method %q not allowed", r.Method)
		log.Printf("ERROR: %s\n", err)
		w.WriteHeader(405)
		tmpl.Lookup("errors/405").Execute(w, helpers.ErrorData(err))
	}
}
