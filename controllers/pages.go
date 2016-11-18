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
	db := models.GetDB()
	if r.Method == "GET" {

		re := regexp.MustCompile("^[0-9]+")
		id := re.FindString(r.URL.Path[len("/pages/"):])
		page := &models.Page{}
		db.First(page, id)
		if page.ID == 0 || !page.Published {
			w.WriteHeader(404)
			tmpl.Lookup("errors/404").Execute(w, nil)
			return
		}
		//redirect to canonical url
		if r.URL.Path != page.URL() {
			http.Redirect(w, r, page.URL(), http.StatusSeeOther)
			return
		}
		data["Page"] = page
		data["Title"] = page.Name
		data["Active"] = page.URL()
		data["MetaDescription"] = page.MetaDescription
		data["MetaKeywords"] = page.MetaKeywords
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
	db := models.GetDB()
	if r.Method == "GET" {

		var list []models.Page
		db.Order("published desc, id desc").Find(&list)
		data["Title"] = "Страницы"
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
	db := models.GetDB()
	if r.Method == "GET" {

		data["Title"] = "Новая страница"
		data["Active"] = "pages"
		data["Flash"] = session.Flashes()
		session.Save(r, w)
		tmpl.Lookup("pages/form").Execute(w, data)

	} else if r.Method == "POST" {

		page := &models.Page{
			Name:            r.PostFormValue("name"),
			Slug:            r.PostFormValue("slug"),
			Content:         r.PostFormValue("content"),
			MetaDescription: r.PostFormValue("meta_description"),
			MetaKeywords:    r.PostFormValue("meta_keywords"),
			Published:       helpers.Atob(r.PostFormValue("published")),
		}

		if err := db.Create(page).Error; err != nil {
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
	db := models.GetDB()
	if r.Method == "GET" {

		id := r.URL.Path[len("/admin/edit_page/"):]
		page := &models.Page{}
		db.First(page, id)
		if page.ID == 0 {
			w.WriteHeader(400)
			tmpl.Lookup("errors/400").Execute(w, nil)
			return
		}

		data["Title"] = "Редактировать страницу"
		data["Active"] = "pages"
		data["Page"] = page
		data["Flash"] = session.Flashes()
		session.Save(r, w)
		tmpl.Lookup("pages/form").Execute(w, data)

	} else if r.Method == "POST" {

		page := &models.Page{
			ID:              helpers.Atouint(r.PostFormValue("id")),
			Name:            r.PostFormValue("name"),
			Slug:            r.PostFormValue("slug"),
			Content:         r.PostFormValue("content"),
			MetaDescription: r.PostFormValue("meta_description"),
			MetaKeywords:    r.PostFormValue("meta_keywords"),
			Published:       helpers.Atob(r.PostFormValue("published")),
		}

		if err := db.Save(page).Error; err != nil {
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
	db := models.GetDB()

	if r.Method == "POST" {

		page := &models.Page{}
		db.First(page, r.PostFormValue("id"))
		if page.ID == 0 {
			w.WriteHeader(404)
			tmpl.Lookup("errors/404").Execute(w, nil)
		}

		if err := db.Delete(page).Error; err != nil {
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
