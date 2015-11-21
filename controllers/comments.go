package controllers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/denisbakhtin/medical/helpers"
	"github.com/denisbakhtin/medical/models"
	"github.com/gorilla/context"
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

//CommentCreate handles /new_comment route
func CommentCreate(w http.ResponseWriter, r *http.Request) {
	session := helpers.Session(r)
	tmpl := helpers.Template(r)
	T := helpers.T(r)
	if r.Method == "POST" {

		if _, ok := session.Values["oauth_name"]; !ok {
			err := fmt.Errorf("You are not authorized to post comments.")
			log.Printf("ERROR: %s\n", err)
			w.WriteHeader(405)
			tmpl.Lookup("errors/405").Execute(w, helpers.ErrorData(err))
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

//TODO: adopt to update via secure link
//CommentReply handles ...
func CommentReply(w http.ResponseWriter, r *http.Request) {
	tmpl := helpers.Template(r)
	session := helpers.Session(r)
	data := helpers.DefaultData(r)
	T := helpers.T(r)
	if r.Method == "GET" {

		user := context.Get(r, "user").(*models.User)
		//TODO: get comment by ID
		comment := &models.Comment{
			ArticleID:  0,
			AuthorName: user.Name,
		}

		data["Title"] = T("reply")
		data["Active"] = "comments"
		data["Comment"] = comment
		data["Flash"] = session.Flashes("comments")
		session.Save(r, w)
		tmpl.Lookup("comments/form").Execute(w, data)

	} else if r.Method == "POST" {

		//TODO: get comment by ID, then update
		comment := &models.Comment{
			ArticleID:   helpers.Atoi64(r.PostFormValue("article_id")),
			AuthorName:  r.PostFormValue("author_name"),
			AuthorEmail: r.PostFormValue("author_email"),
			Content:     r.PostFormValue("content"),
			Answer:      r.PostFormValue("answer"),
			Published:   helpers.Atob(r.PostFormValue("published")),
		}

		if err := comment.Insert(); err != nil {
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
