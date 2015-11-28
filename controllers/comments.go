package controllers

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"github.com/denisbakhtin/medical/helpers"
	"github.com/denisbakhtin/medical/models"
	"github.com/denisbakhtin/medical/system"
	"gopkg.in/gomail.v2"
	"log"
	"net/http"
	"regexp"
	"strconv"
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

//CommentPublicCreate handles /new_comment route
func CommentPublicCreate(w http.ResponseWriter, r *http.Request) {
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
		notifyAdminOfComment(r, comment)
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

//CommentPublicUpdate handles /edit_comment?token=:secure_token route
func CommentPublicUpdate(w http.ResponseWriter, r *http.Request) {
	tmpl := helpers.Template(r)
	session := helpers.Session(r)
	data := helpers.DefaultData(r)
	T := helpers.T(r)
	if r.Method == "GET" {

		id := getIDFromToken(r.FormValue("token"))
		comment, err := models.GetComment(id)
		if err != nil || comment.Published {
			err := fmt.Errorf(T("comment_not_found_or_already_published"))
			w.WriteHeader(404)
			tmpl.Lookup("errors/404").Execute(w, helpers.ErrorData(err))
			return
		}

		comment.Published = true //set default to true
		data["Title"] = T("edit_comment")
		data["Active"] = "comments"
		data["Comment"] = comment
		data["SecureEdit"] = true
		data["Flash"] = session.Flashes("comments")
		session.Save(r, w)
		tmpl.Lookup("comments/public-edit-form").Execute(w, data)

	} else if r.Method == "POST" {

		r.ParseForm()
		comment := &models.Comment{
			ID:          helpers.Atoi64(r.PostFormValue("id")),
			AuthorName:  r.PostFormValue("author_name"),
			AuthorEmail: r.PostFormValue("author_email"),
			Content:     r.PostFormValue("content"),
			Answer:      r.PostFormValue("answer"),
			Published:   helpers.Atob(r.PostFormValue("published")),
		}

		if err := comment.Update(); err != nil {
			session.AddFlash(err.Error(), "comments")
			session.Save(r, w)
			http.Redirect(w, r, r.RequestURI, 303)
			return
		}
		if comment.Published {
			notifyClientOfComment(r, comment)
			postCommentOnSocialWalls(r, comment)
		}
		session.AddFlash(T("comment_has_been_successfully_updated"))
		session.Save(r, w)
		http.Redirect(w, r, "/", 303)

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

func notifyAdminOfComment(r *http.Request, comment *models.Comment) {
	//closure is needed here, as r may be released by the time func finishes
	tmpl := helpers.Template(r)
	T := helpers.T(r)
	go func() {
		data := map[string]interface{}{
			"Comment": comment,
			"Token":   createTokenFromID(comment.ID),
		}
		var b bytes.Buffer
		if err := tmpl.Lookup("emails/question").Execute(&b, data); err != nil {
			log.Printf("ERROR: %s\n", err)
			return
		}

		smtp := system.GetConfig().SMTP
		msg := gomail.NewMessage()
		msg.SetHeader("From", smtp.From)
		msg.SetHeader("To", smtp.To)
		if len(smtp.Cc) > 0 {
			msg.SetHeader("Cc", smtp.Cc)
		}
		msg.SetHeader("Subject", T("new_comment_has_been_created", map[string]interface{}{"Name": comment.AuthorName}))
		msg.SetBody(
			"text/html",
			b.String(),
		)

		port, _ := strconv.Atoi(smtp.Port)
		dialer := gomail.NewPlainDialer(smtp.SMTP, port, smtp.User, smtp.Password)
		sender, err := dialer.Dial()
		if err != nil {
			log.Printf("ERROR: %s\n", err)
			return
		}
		if err := gomail.Send(sender, msg); err != nil {
			log.Printf("ERROR: %s\n", err)
			return
		}
	}()
}

//notifyClientOfComment sends notification email to comment(question) author
func notifyClientOfComment(r *http.Request, comment *models.Comment) {
	if len(comment.AuthorEmail) == 0 {
		return
	}
	tmpl := helpers.Template(r)
	T := helpers.T(r)
	go func() {
		data := map[string]interface{}{
			"Comment": comment,
		}
		var b bytes.Buffer
		if err := tmpl.Lookup("emails/answer").Execute(&b, data); err != nil {
			log.Printf("ERROR: %s\n", err)
			return
		}

		smtp := system.GetConfig().SMTP
		msg := gomail.NewMessage()
		msg.SetHeader("From", smtp.From)
		msg.SetHeader("To", comment.AuthorEmail)
		msg.SetHeader("Subject", T("your_question_has_been_answered"))
		msg.SetBody(
			"text/html",
			b.String(),
		)

		port, _ := strconv.Atoi(smtp.Port)
		dialer := gomail.NewPlainDialer(smtp.SMTP, port, smtp.User, smtp.Password)
		sender, err := dialer.Dial()
		if err != nil {
			log.Printf("ERROR: %s\n", err)
			return
		}
		if err := gomail.Send(sender, msg); err != nil {
			log.Printf("ERROR: %s\n", err)
			return
		}
	}()
}

func postCommentOnSocialWalls(r *http.Request, comment *models.Comment) {
	//TODO
}
