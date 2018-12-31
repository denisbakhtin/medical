package controllers

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"log"
	"strconv"

	"github.com/denisbakhtin/medical/helpers"
	"github.com/denisbakhtin/medical/models"
	"github.com/denisbakhtin/medical/system"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"gopkg.in/gomail.v2"
)

//CommentsAdminIndex handles GET /admin/comments route
func CommentsAdminIndex(c *gin.Context) {
	db := models.GetDB()

	var list []models.Comment
	db.Order("id desc").Find(&list)
	c.HTML(200, "comments/admin/index", gin.H{
		"Title":  "Вопросы посетителей",
		"Active": "comments",
		"List":   list,
	})
}

//CommentCreatePost handles /new_comment route
func CommentCreatePost(c *gin.Context) {
	session := sessions.Default(c)
	db := models.GetDB()

	comment := &models.Comment{}
	if c.Bind(comment) == nil {
		//simple captcha check
		captcha, err := base64.StdEncoding.DecodeString(comment.Captcha)
		if err != nil {
			c.HTML(500, "errors/500", helpers.ErrorData(err))
			return
		}
		if string(captcha) != "100.00" {
			c.HTML(400, "errors/400", nil)
			return
		}
		comment.Published = false //leave unpublished
		if err := db.Create(comment).Error; err != nil {
			c.HTML(400, "errors/400", helpers.ErrorData(err))
			return
		}
		notifyAdminOfComment(comment)
		session.AddFlash("Спасибо! Ваш вопрос будет опубликован после проверки.")
		session.Save()
		c.Redirect(303, fmt.Sprintf("/articles/%d#comments", comment.ArticleID))
	} else {
		session.AddFlash("Ошибка! Внимательно проверьте заполнение всех полей!")
		session.Save()
		c.Redirect(303, "/")
	}
}

//CommentAdminUpdateGet handles /admin/edit_comment/:id get request
func CommentAdminUpdateGet(c *gin.Context) {
	session := sessions.Default(c)
	flashes := session.Flashes()
	session.Save()
	db := models.GetDB()

	id := c.Param("id")
	comment := &models.Comment{}
	db.First(comment, id)
	if comment.ID == 0 {
		c.HTML(404, "errors/404", nil)
		return
	}

	c.HTML(200, "comments/admin/form", gin.H{
		"Title":   "Редактировать вопрос",
		"Active":  "comments",
		"Comment": comment,
		"Flash":   flashes,
	})
}

//CommentAdminUpdatePost handles /admin/edit_comment/:id post request
func CommentAdminUpdatePost(c *gin.Context) {
	session := sessions.Default(c)
	db := models.GetDB()

	comment := &models.Comment{}
	if c.Bind(comment) == nil {
		if err := db.Save(comment).Error; err != nil {
			session.AddFlash(err.Error())
			session.Save()
			c.Redirect(303, c.Request.RequestURI)
			return
		}
		c.Redirect(303, "/admin/comments")
	} else {
		session.AddFlash("Ошибка! Внимательно проверьте заполнение полей!")
		session.Save()
		c.Redirect(303, c.Request.RequestURI)
	}
}

//CommentPublicUpdate handles /edit_comment?token=:secure_token route
/*
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
			notifyClientOfComment(comment)
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
*/

//CommentAdminDelete handles /admin/delete_comment route
func CommentAdminDelete(c *gin.Context) {
	db := models.GetDB()

	comment := &models.Comment{}
	db.First(comment, c.Request.PostFormValue("id"))
	if comment.ID == 0 {
		c.HTML(404, "errors/404", nil)
	}

	if err := db.Delete(comment).Error; err != nil {
		c.HTML(500, "errors/500", helpers.ErrorData(err))
		return
	}
	c.Redirect(303, "/admin/comments")

}

func notifyAdminOfComment(comment *models.Comment) {
	//closure is needed here, as r may be released by the time func finishes
	tmpl := system.GetTemplates()
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
func notifyClientOfComment(comment *models.Comment) {
	if len(comment.AuthorEmail) == 0 {
		return
	}
	tmpl := system.GetTemplates()
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

//CommentsIndex handles GET /comments/:id route, where :id is the article id
func CommentsIndex(c *gin.Context) {
	db := models.GetDB()
	id := c.Param("id")

	article := models.Article{}
	db.First(&article, id)

	var list []models.Comment
	db.Where("article_id = ?", id).Order("answer desc, id desc").Find(&list)
	c.HTML(200, "comments/index", gin.H{
		"Title":   "Вопросы посетителей",
		"Active":  "comments",
		"Article": &article,
		"List":    list,
	})
}
