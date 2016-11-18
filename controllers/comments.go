package controllers

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/denisbakhtin/medical/helpers"
	"github.com/denisbakhtin/medical/models"
	"github.com/denisbakhtin/medical/system"
	"gopkg.in/gomail.v2"
)

//CommentIndex handles GET /admin/comments route
func CommentIndex(w http.ResponseWriter, r *http.Request) {
	tmpl := helpers.Template(r)
	data := helpers.DefaultData(r)
	db := models.GetDB()
	if r.Method == "GET" {

		var list []models.Comment
		db.Order("id desc").Find(&list)
		data["Title"] = "Вопросы посетителей"
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

//CommentPublicCreate handles /new_comment route
func CommentPublicCreate(w http.ResponseWriter, r *http.Request) {
	session := helpers.Session(r)
	tmpl := helpers.Template(r)
	db := models.GetDB()
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
			ArticleID:   helpers.Atouint(r.PostFormValue("article_id")),
			AuthorCity:  r.PostFormValue("author_city"),
			AuthorName:  r.PostFormValue("author_name"),
			AuthorEmail: r.PostFormValue("author_email"),
			Content:     r.PostFormValue("content"),
			Published:   false, //comments are published by admin via dashboard
		}

		if err := db.Create(comment).Error; err != nil {
			log.Printf("ERROR: %s\n", err)
			w.WriteHeader(400)
			tmpl.Lookup("errors/400").Execute(w, helpers.ErrorData(err))
			return
		}
		notifyAdminOfComment(r, comment)
		session.AddFlash("Спасибо! Ваш вопрос будет опубликован после проверки.", "comments")
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
	db := models.GetDB()
	if r.Method == "GET" {

		id := r.URL.Path[len("/admin/edit_comment/"):]
		comment := &models.Comment{}
		db.First(comment, id)
		if comment.ID == 0 {
			w.WriteHeader(404)
			tmpl.Lookup("errors/404").Execute(w, nil)
			return
		}

		data["Title"] = "Редактировать вопрос"
		data["Active"] = "comments"
		data["Comment"] = comment
		data["Flash"] = session.Flashes("comments")
		session.Save(r, w)
		tmpl.Lookup("comments/form").Execute(w, data)

	} else if r.Method == "POST" {

		r.ParseForm()
		comment := &models.Comment{
			ID:          helpers.Atouint(r.PostFormValue("id")),
			ArticleID:   helpers.Atouint(r.PostFormValue("article_id")),
			AuthorCity:  r.PostFormValue("author_city"),
			AuthorName:  r.PostFormValue("author_name"),
			AuthorEmail: r.PostFormValue("author_email"),
			Content:     r.PostFormValue("content"),
			Answer:      r.PostFormValue("answer"),
			Published:   helpers.Atob(r.PostFormValue("published")),
		}

		if err := db.Save(comment).Error; err != nil {
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
	db := models.GetDB()
	if r.Method == "GET" {

		id := getIDFromToken(r.FormValue("token"))
		comment := &models.Comment{}
		db.First(comment, id)
		if comment.ID == 0 || comment.Published {
			err := fmt.Errorf("Вопрос не существует или уже был опубликован ранее")
			w.WriteHeader(404)
			tmpl.Lookup("errors/404").Execute(w, helpers.ErrorData(err))
			return
		}

		comment.Published = true //set default to true
		data["Title"] = "Редактировать вопрос"
		data["Active"] = "comments"
		data["Comment"] = comment
		data["SecureEdit"] = true
		data["Flash"] = session.Flashes("comments")
		session.Save(r, w)
		tmpl.Lookup("comments/public-edit-form").Execute(w, data)

	} else if r.Method == "POST" {

		r.ParseForm()
		comment := &models.Comment{
			ID:          helpers.Atouint(r.PostFormValue("id")),
			AuthorName:  r.PostFormValue("author_name"),
			AuthorEmail: r.PostFormValue("author_email"),
			AuthorCity:  r.PostFormValue("author_city"),
			Content:     r.PostFormValue("content"),
			Answer:      r.PostFormValue("answer"),
			Published:   helpers.Atob(r.PostFormValue("published")),
		}

		if err := db.Save(comment).Error; err != nil {
			session.AddFlash(err.Error(), "comments")
			session.Save(r, w)
			http.Redirect(w, r, r.RequestURI, 303)
			return
		}
		if comment.Published {
			notifyClientOfComment(r, comment)
		}
		session.AddFlash("Вопрос был успешно сохранен")
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
	db := models.GetDB()

	if r.Method == "POST" {

		comment := &models.Comment{}
		db.First(comment, r.PostFormValue("id"))
		if comment.ID == 0 {
			w.WriteHeader(404)
			tmpl.Lookup("errors/404").Execute(w, nil)
		}

		if err := db.Delete(comment).Error; err != nil {
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
		if len(comment.AuthorEmail) > 0 {
			msg.SetHeader("Reply-To", comment.AuthorEmail)
		}
		if len(smtp.Cc) > 0 {
			msg.SetHeader("Cc", smtp.Cc)
		}
		msg.SetHeader("Subject", fmt.Sprintf("Новый вопрос на сайте www.miobalans.ru: %s", comment.AuthorName))
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
		msg.SetHeader("Subject", "Врач ответил на ваш вопрос на сайте www.miobalans.ru")
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
