package controllers

import (
	"fmt"
	"log"
	"net/http"

	"encoding/base64"
	"github.com/denisbakhtin/medical/helpers"
	"github.com/denisbakhtin/medical/models"
)

//ReviewShow handles /reviews/:id route
func ReviewShow(w http.ResponseWriter, r *http.Request) {
	tmpl := helpers.Template(r)
	data := helpers.DefaultData(r)
	T := helpers.T(r)
	if r.Method == "GET" {

		id := r.URL.Path[len("/reviews/"):]
		review, err := models.GetReview(id)
		if err != nil || !review.Published {
			w.WriteHeader(404)
			tmpl.Lookup("errors/404").Execute(w, nil)
			return
		}
		data["Review"] = review
		data["Title"] = T("testimonial_title") + ". " + review.AuthorName
		data["Active"] = "/reviews"
		tmpl.Lookup("reviews/show").Execute(w, data)

	} else {
		err := fmt.Errorf("Method %q not allowed", r.Method)
		log.Printf("ERROR: %s\n", err)
		w.WriteHeader(405)
		tmpl.Lookup("errors/405").Execute(w, helpers.ErrorData(err))
	}
}

//ReviewPublicIndex handles GET /reviews route
func ReviewPublicIndex(w http.ResponseWriter, r *http.Request) {
	tmpl := helpers.Template(r)
	data := helpers.DefaultData(r)
	T := helpers.T(r)
	if r.Method == "GET" {

		list, err := models.GetPublishedReviews()
		if err != nil {
			log.Printf("ERROR: %s\n", err)
			w.WriteHeader(500)
			tmpl.Lookup("errors/500").Execute(w, helpers.ErrorData(err))
			return
		}
		data["Title"] = T("testimonials_title")
		data["Active"] = r.RequestURI
		data["List"] = list
		tmpl.Lookup("reviews/public-index").Execute(w, data)

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

//ReviewPublicCreate handles /new_review route
func ReviewPublicCreate(w http.ResponseWriter, r *http.Request) {
	session := helpers.Session(r)
	tmpl := helpers.Template(r)
	data := helpers.DefaultData(r)
	T := helpers.T(r)
	if r.Method == "GET" {

		data["Title"] = T("new_review")
		data["Active"] = "reviews"
		data["Flash"] = session.Flashes("reviews")
		session.Save(r, w)
		tmpl.Lookup("reviews/public-form").Execute(w, data)

	} else if r.Method == "POST" {

		r.ParseMultipartForm(32 << 20)
		//simple captcha check
		captcha, err := base64.StdEncoding.DecodeString(r.FormValue("captcha"))
		if err != nil {
			w.WriteHeader(500)
			tmpl.Lookup("errors/500").Execute(w, helpers.ErrorData(err))
			return
		}
		if string(captcha) != "100.00" {
			w.WriteHeader(400)
			tmpl.Lookup("errors/400").Execute(w, nil)
			return
		}

		review := &models.Review{
			AuthorName:  r.FormValue("author_name"),
			AuthorEmail: r.FormValue("author_email"),
			Content:     r.FormValue("content"),
			Published:   false, //reviews are published by admin via dashboard
		}

		if mpartFile, mpartHeader, err := r.FormFile("image"); err == nil {
			defer mpartFile.Close()
			review.Image, err = saveFile(mpartHeader, mpartFile)
			if err != nil {
				log.Println(err.Error())
				w.WriteHeader(500)
				tmpl.Lookup("errors/500").Execute(w, helpers.ErrorData(err))
				return
			}
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

//ReviewCreate handles /admin/new_review route
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

		r.ParseMultipartForm(32 << 20)
		review := &models.Review{
			AuthorName:  r.FormValue("author_name"),
			AuthorEmail: r.FormValue("author_email"),
			Content:     r.FormValue("content"),
			Published:   helpers.Atob(r.FormValue("published")),
		}

		if mpartFile, mpartHeader, err := r.FormFile("image"); err == nil {
			defer mpartFile.Close()
			review.Image, err = saveFile(mpartHeader, mpartFile)
			if err != nil {
				log.Println(err.Error())
				w.WriteHeader(500)
				tmpl.Lookup("errors/500").Execute(w, helpers.ErrorData(err))
				return
			}
		}

		if err := review.Insert(); err != nil {
			log.Printf("ERROR: %s\n", err)
			w.WriteHeader(400)
			tmpl.Lookup("errors/400").Execute(w, helpers.ErrorData(err))
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

		r.ParseMultipartForm(32 << 20)
		rev, _ := models.GetReview(r.FormValue("id"))
		if rev.ID == 0 {
			w.WriteHeader(400)
			tmpl.Lookup("errors/400").Execute(w, nil)
			return
		}

		review := &models.Review{
			ID:          rev.ID,
			AuthorName:  r.FormValue("author_name"),
			AuthorEmail: r.FormValue("author_email"),
			Content:     r.FormValue("content"),
			Image:       rev.Image,
			Published:   helpers.Atob(r.FormValue("published")),
		}
		if mpartFile, mpartHeader, err := r.FormFile("image"); err == nil {
			defer mpartFile.Close()
			review.Image, err = saveFile(mpartHeader, mpartFile)
			if err != nil {
				log.Println(err.Error())
				w.WriteHeader(500)
				tmpl.Lookup("errors/500").Execute(w, helpers.ErrorData(err))
				return
			}
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
