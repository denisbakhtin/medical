package controllers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/denisbakhtin/medical/helpers"
	"github.com/denisbakhtin/medical/models"
)

//ReviewShow handles /reviews/:id route
func ReviewShow(w http.ResponseWriter, r *http.Request) {
	tmpl := helpers.Template(r)
	data := helpers.DefaultData(r)
	if r.Method == "GET" {

		id := r.URL.Path[len("/reviews/"):]
		review, err := models.GetReview(id)
		if err != nil || !review.Published {
			w.WriteHeader(404)
			tmpl.Lookup("errors/404").Execute(w, nil)
			return
		}
		data["Review"] = review
		data["Title"] = review.Excerpt()
		data["Active"] = fmt.Sprintf("reviews/%s", id)
		tmpl.Lookup("reviews/show").Execute(w, data)

	} else {
		err := fmt.Errorf("Method %q not allowed", r.Method)
		log.Printf("ERROR: %s\n", err)
		w.WriteHeader(405)
		tmpl.Lookup("errors/405").Execute(w, helpers.ErrorData(err))
	}
}

//ReviewIndex handles GET /admin/reviews route
func ReviewIndex(w http.ResponseWriter, r *http.Request) {
	tmpl := helpers.Template(r)
	data := helpers.DefaultData(r)
	T := helpers.T(r)
	if r.Method == "GET" {

		list, err := models.GetReviews()
		if err != nil {
			w.WriteHeader(500)
			tmpl.Lookup("errors/500").Execute(w, helpers.ErrorData(err))
			return
		}
		data["Title"] = T("reviews")
		data["Active"] = "reviews"
		data["List"] = list
		tmpl.Lookup("reviews/index").Execute(w, data)

	} else {
		err := fmt.Errorf("Method %q not allowed", r.Method)
		log.Printf("ERROR: %s\n", err)
		w.WriteHeader(405)
		tmpl.Lookup("errors/405").Execute(w, helpers.ErrorData(err))
	}
}

//ReviewCreate handles /new_review route
func ReviewCreate(w http.ResponseWriter, r *http.Request) {
	session := helpers.Session(r)
	tmpl := helpers.Template(r)
	data := helpers.DefaultData(r)
	T := helpers.T(r)
	if r.Method == "GET" {

		data["Title"] = T("new_review")
		data["Active"] = "reviews"
		data["Flash"] = session.Flashes("reviews")
		session.Save(r, w)
		tmpl.Lookup("reviews/form").Execute(w, data)

	} else if r.Method == "POST" {

		if _, ok := session.Values["oauth_name"]; !ok {
			err := fmt.Errorf("You are not authorized to post reviews.")
			log.Printf("ERROR: %s\n", err)
			w.WriteHeader(405)
			tmpl.Lookup("errors/405").Execute(w, helpers.ErrorData(err))
			return
		}

		review := &models.Review{
			AuthorName:  r.PostFormValue("author_name"),
			AuthorEmail: r.PostFormValue("author_email"),
			Content:     r.PostFormValue("content"),
			Image:       "",    //TODO: upload image
			Published:   false, //reviews are published by admin via dashboard
		}

		if err := review.Insert(); err != nil {
			log.Printf("ERROR: %s\n", err)
			w.WriteHeader(400)
			tmpl.Lookup("errors/400").Execute(w, helpers.ErrorData(err))
			return
		}
		session.AddFlash(T("thank_you_for_posting_review"), "reviews")
		session.Save(r, w)
		http.Redirect(w, r, "/reviews", 303)

	} else {
		err := fmt.Errorf("Method %q not allowed", r.Method)
		log.Printf("ERROR: %s\n", err)
		w.WriteHeader(405)
		tmpl.Lookup("errors/405").Execute(w, helpers.ErrorData(err))
	}
}

//ReviewUpdate handles /admin/edit_review/:id route
func ReviewUpdate(w http.ResponseWriter, r *http.Request) {
	tmpl := helpers.Template(r)
	session := helpers.Session(r)
	data := helpers.DefaultData(r)
	T := helpers.T(r)
	if r.Method == "GET" {

		id := r.URL.Path[len("/admin/edit_review/"):]
		review, err := models.GetReview(id)
		if err != nil {
			w.WriteHeader(404)
			tmpl.Lookup("errors/404").Execute(w, helpers.ErrorData(err))
			return
		}

		data["Title"] = T("edit_review")
		data["Active"] = "reviews"
		data["Review"] = review
		data["Flash"] = session.Flashes("reviews")
		session.Save(r, w)
		tmpl.Lookup("reviews/form").Execute(w, data)

	} else if r.Method == "POST" {

		r.ParseForm()
		review := &models.Review{
			ID:          helpers.Atoi64(r.PostFormValue("id")),
			AuthorName:  r.PostFormValue("author_name"),
			AuthorEmail: r.PostFormValue("author_email"),
			Content:     r.PostFormValue("content"),
			Image:       "", //TODO: upload image if changed
			Published:   helpers.Atob(r.PostFormValue("published")),
		}

		if err := review.Update(); err != nil {
			session.AddFlash(err.Error(), "reviews")
			session.Save(r, w)
			http.Redirect(w, r, r.RequestURI, 303)
			return
		}
		http.Redirect(w, r, "/admin/reviews", 303)

	} else {
		err := fmt.Errorf("Method %q not allowed", r.Method)
		log.Printf("ERROR: %s\n", err)
		w.WriteHeader(405)
		tmpl.Lookup("errors/405").Execute(w, helpers.ErrorData(err))
	}
}

//ReviewDelete handles /admin/delete_review route
func ReviewDelete(w http.ResponseWriter, r *http.Request) {
	tmpl := helpers.Template(r)

	if r.Method == "POST" {

		review, err := models.GetReview(r.PostFormValue("id"))
		if err != nil {
			log.Printf("ERROR: %s\n", err)
			w.WriteHeader(404)
			tmpl.Lookup("errors/404").Execute(w, helpers.ErrorData(err))
		}

		if err := review.Delete(); err != nil {
			log.Printf("ERROR: %s\n", err)
			w.WriteHeader(500)
			tmpl.Lookup("errors/500").Execute(w, helpers.ErrorData(err))
			return
		}
		http.Redirect(w, r, "/admin/reviews", 303)

	} else {
		err := fmt.Errorf("Method %q not allowed", r.Method)
		log.Printf("ERROR: %s\n", err)
		w.WriteHeader(405)
		tmpl.Lookup("errors/405").Execute(w, helpers.ErrorData(err))
	}
}
