package controllers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/denisbakhtin/medical/helpers"
	"github.com/denisbakhtin/medical/models"
)

//TagShow handles GET /tags/:name route
func TagShow(w http.ResponseWriter, r *http.Request) {
	tmpl := helpers.Template(r)
	data := helpers.DefaultData(r)
	if r.Method == "GET" {

		name := r.URL.Path[len("/tags/"):]
		tag, err := models.GetTag(name)
		if err != nil {
			w.WriteHeader(404)
			tmpl.Lookup("errors/404").Execute(w, nil)
			return
		}
		data["Tag"] = tag
		data["Title"] = tag.Name
		data["Active"] = fmt.Sprintf("tags/%s", name)
		tmpl.Lookup("tags/show").Execute(w, data)

	} else {
		err := fmt.Errorf("Method %q not allowed", r.Method)
		log.Printf("ERROR: %s\n", err)
		w.WriteHeader(405)
		tmpl.Lookup("errors/405").Execute(w, helpers.ErrorData(err))
	}
}

//TagIndex handles GET /admin/tags route
func TagIndex(w http.ResponseWriter, r *http.Request) {
	tmpl := helpers.Template(r)
	data := helpers.DefaultData(r)
	T := helpers.T(r)
	if r.Method == "GET" {

		list, err := models.GetTags()
		if err != nil {
			w.WriteHeader(500)
			tmpl.Lookup("errors/500").Execute(w, helpers.ErrorData(err))
			return
		}
		data["Title"] = T("tags")
		data["Active"] = "tags"
		data["List"] = list
		tmpl.Lookup("tags/index").Execute(w, data)

	} else {
		err := fmt.Errorf("Method %q not allowed", r.Method)
		log.Printf("ERROR: %s\n", err)
		w.WriteHeader(405)
		tmpl.Lookup("errors/405").Execute(w, helpers.ErrorData(err))
	}
}

//TagCreate handles /admin/new_tag route
func TagCreate(w http.ResponseWriter, r *http.Request) {
	tmpl := helpers.Template(r)
	session := helpers.Session(r)
	data := helpers.DefaultData(r)
	T := helpers.T(r)
	if r.Method == "GET" {

		data["Title"] = T("new_tag")
		data["Active"] = "tags"
		data["Flash"] = session.Flashes()
		session.Save(r, w)
		tmpl.Lookup("tags/form").Execute(w, data)

	} else if r.Method == "POST" {

		tag := &models.Tag{
			Name: r.PostFormValue("name"),
		}

		if err := tag.Insert(); err != nil {
			session.AddFlash(err.Error())
			session.Save(r, w)
			http.Redirect(w, r, "/admin/new_tag", 303)
			return
		}
		http.Redirect(w, r, "/admin/tags", 303)

	} else {
		err := fmt.Errorf("Method %q not allowed", r.Method)
		log.Printf("ERROR: %s\n", err)
		w.WriteHeader(405)
		tmpl.Lookup("errors/405").Execute(w, helpers.ErrorData(err))
	}
}

//TagDelete handles /admin/delete_tag route
func TagDelete(w http.ResponseWriter, r *http.Request) {
	tmpl := helpers.Template(r)

	if r.Method == "POST" {

		tag, err := models.GetTag(r.PostFormValue("name"))
		if err != nil {
			log.Printf("ERROR: %s\n", err)
			w.WriteHeader(404)
			tmpl.Lookup("errors/404").Execute(w, helpers.ErrorData(err))
		}

		if err := tag.Delete(); err != nil {
			log.Printf("ERROR: %s\n", err)
			w.WriteHeader(500)
			tmpl.Lookup("errors/500").Execute(w, helpers.ErrorData(err))
			return
		}
		http.Redirect(w, r, "/admin/tags", 303)

	} else {
		err := fmt.Errorf("Method %q not allowed", r.Method)
		log.Printf("ERROR: %s\n", err)
		w.WriteHeader(405)
		tmpl.Lookup("errors/405").Execute(w, helpers.ErrorData(err))
	}
}
