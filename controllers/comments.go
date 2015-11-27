package controllers

import (
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"regexp"

	"github.com/denisbakhtin/medical/helpers"
	"github.com/denisbakhtin/medical/models"
)

//CommentIndex handles GET /admin/comments route
func CommentIndex(w http.ResponseWriter, r *http.Request) {
	tmpl := helpers.Template(r)
	data := helpers.DefaultData(r)
	T := helpers.T(r)
	if r.Method == "GET" {

		list, err := models.GetComments()
		if err != nil {
			w.WriteHeader(500)
			tmpl.Lookup("errors/500").Execute(w, helpers.ErrorData(err))
			return
		}
		data["Title"] = T("comments")
		data["Active"] = "comments"
		data["List"] = list
		tmpl.Lookup("comments/index").Execute(w, data)

	} else {
		err := fmt.Errorf("Method %q not allowed", r.Method)
		log.Printf("ERROR: %s\n", err)
		w.WriteHeader(405)
		tmpl.Lookup("errors/405").Execute(w, helpers.ErrorData(err))
	}
}

//CommentShow handles GET /comments/:id-slug route
func CommentShow(w http.ResponseWriter, r *http.Request) {
	tmpl := helpers.Template(r)
	data := helpers.DefaultData(r)
	T := helpers.T(r)
	if r.Method == "GET" {

		re := regexp.MustCompile("^[0-9]+")
		id := re.FindString(r.URL.Path[len("/comments/"):])
		comment, err := models.GetComment(id)
		if err != nil || !comment.Published {
			w.WriteHeader(404)
			tmpl.Lookup("errors/404").Execute(w, nil)
			return
		}
		//redirect to canonical url
		if r.URL.Path != comment.URL() {
			http.Redirect(w, r, comment.URL(), http.StatusSeeOther)
			return
		}
		data["Comment"] = comment
		data["SimilarComments"], _ = comment.GetSimilar()
		data["Article"], _ = models.GetArticle(comment.ArticleID)
		data["Title"] = comment.AuthorName + " " + T("asks_kineziologist")
		data["Active"] = ""
		tmpl.Lookup("comments/show").Execute(w, data)

	} else {
		err := fmt.Errorf("Method %q not allowed", r.Method)
		log.Printf("ERROR: %s\n", err)
		w.WriteHeader(405)
		tmpl.Lookup("errors/405").Execute(w, helpers.ErrorData(err))
	}
}

//CommentCreate handles /new_comment route
func CommentCreate(w http.ResponseWriter, r *http.Request) {
	session := helpers.Session(r)
	tmpl := helpers.Template(r)
	T := helpers.T(r)
	if r.Method == "POST" {

		r.ParseForm()
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

		comment := &models.Comment{
			ArticleID:   helpers.Atoi64(r.PostFormValue("article_id")),
			AuthorName:  r.PostFormValue("author_name"),
			AuthorEmail: r.PostFormValue("author_email"),
			Content:     r.PostFormValue("content"),
			Published:   false, //comments are published by admin via dashboard
		}

		if err := comment.Insert(); err != nil {
			log.Printf("ERROR: %s\n", err)
			w.WriteHeader(400)
			tmpl.Lookup("errors/400").Execute(w, helpers.ErrorData(err))
			return
		}
		session.AddFlash(T("thank_you_for_posting_question"), "comments")
		session.Save(r, w)
		http.Redirect(w, r, fmt.Sprintf("/articles/%d#comments", comment.ArticleID), 303)

	} else {
		err := fmt.Errorf("Method %q not allowed", r.Method)
		log.Printf("ERROR: %s\n", err)
		w.WriteHeader(405)
		tmpl.Lookup("errors/405").Execute(w, helpers.ErrorData(err))
	}
}

//CommentUpdate handles /admin/edit_comment/:id route
func CommentUpdate(w http.ResponseWriter, r *http.Request) {
	tmpl := helpers.Template(r)
	session := helpers.Session(r)
	data := helpers.DefaultData(r)
	T := helpers.T(r)
	if r.Method == "GET" {

		id := r.URL.Path[len("/admin/edit_comment/"):]
		comment, err := models.GetComment(id)
		if err != nil {
			w.WriteHeader(404)
			tmpl.Lookup("errors/404").Execute(w, helpers.ErrorData(err))
			return
		}

		data["Title"] = T("edit_comment")
		data["Active"] = "comments"
		data["Comment"] = comment
		data["Flash"] = session.Flashes("comments")
		session.Save(r, w)
		tmpl.Lookup("comments/form").Execute(w, data)

	} else if r.Method == "POST" {

		r.ParseForm()
		comment := &models.Comment{
			ID:        helpers.Atoi64(r.PostFormValue("id")),
			Content:   r.PostFormValue("content"),
			Answer:    r.PostFormValue("answer"),
			Published: helpers.Atob(r.PostFormValue("published")),
		}

		if err := comment.Update(); err != nil {
			session.AddFlash(err.Error(), "comments")
			session.Save(r, w)
			http.Redirect(w, r, r.RequestURI, 303)
			return
		}
		http.Redirect(w, r, "/admin/comments", 303)

	} else {
		err := fmt.Errorf("Method %q not allowed", r.Method)
		log.Printf("ERROR: %s\n", err)
		w.WriteHeader(405)
		tmpl.Lookup("errors/405").Execute(w, helpers.ErrorData(err))
	}
}

//CommentDelete handles /admin/delete_comment route
func CommentDelete(w http.ResponseWriter, r *http.Request) {
	tmpl := helpers.Template(r)

	if r.Method == "POST" {

		comment, err := models.GetComment(r.PostFormValue("id"))
		if err != nil {
			log.Printf("ERROR: %s\n", err)
			w.WriteHeader(404)
			tmpl.Lookup("errors/404").Execute(w, helpers.ErrorData(err))
		}

		if err := comment.Delete(); err != nil {
			log.Printf("ERROR: %s\n", err)
			w.WriteHeader(500)
			tmpl.Lookup("errors/500").Execute(w, helpers.ErrorData(err))
			return
		}
		http.Redirect(w, r, "/admin/comments", 303)

	} else {
		err := fmt.Errorf("Method %q not allowed", r.Method)
		log.Printf("ERROR: %s\n", err)
		w.WriteHeader(405)
		tmpl.Lookup("errors/405").Execute(w, helpers.ErrorData(err))
	}
}
