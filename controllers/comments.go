package controllers

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"math"

	"github.com/denisbakhtin/medical/helpers"
	"github.com/denisbakhtin/medical/models"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// CommentsAdminIndex handles GET /admin/comments route
func (app *Application) CommentsAdminIndex(c *gin.Context) {
	list, err := app.CommentsRepo.GetAll()
	if err != nil {
		app.Error(c, err)
		return
	}
	c.HTML(200, "comments/admin/index", gin.H{
		"Title":  "Вопросы посетителей",
		"Active": "comments",
		"List":   list,
	})
}

// CommentCreatePost handles /new_comment route
func (app *Application) CommentCreatePost(c *gin.Context) {
	session := sessions.Default(c)

	comment := &models.Comment{}
	if c.Bind(comment) == nil {
		// simple captcha check
		captcha, err := base64.StdEncoding.DecodeString(comment.Captcha)
		if err != nil {
			app.Error(c, err)
			return
		}
		if string(captcha) != "100.00" {
			c.HTML(400, "errors/400", nil)
			return
		}
		comment.Published = false // set unpublished
		if err := app.CommentsRepo.Create(comment); err != nil {
			app.Error(c, err)
			return
		}
		app.notifyAdminOfComment(comment)
		session.AddFlash("Спасибо! Ваш вопрос будет опубликован после проверки.")
		_ = session.Save()
		c.Redirect(303, fmt.Sprintf("/articles/%d#comments", comment.ArticleID))
	} else {
		session.AddFlash("Ошибка! Внимательно проверьте заполнение всех полей!")
		_ = session.Save()
		c.Redirect(303, "/")
	}
}

// CommentAdminUpdateGet handles /admin/edit_comment/:id get request
func (app *Application) CommentAdminUpdateGet(c *gin.Context) {
	session := sessions.Default(c)
	flashes := session.Flashes()
	_ = session.Save()

	id := helpers.Atouint(c.Param("id"))
	comment, err := app.CommentsRepo.Get(id)
	if err != nil {
		app.Error(c, err)
		return
	}

	c.HTML(200, "comments/admin/form", gin.H{
		"Title":   "Редактировать вопрос",
		"Active":  "comments",
		"Comment": comment,
		"Flash":   flashes,
	})
}

// CommentAdminUpdatePost handles /admin/edit_comment/:id post request
func (app *Application) CommentAdminUpdatePost(c *gin.Context) {
	session := sessions.Default(c)

	comment := &models.Comment{}
	if c.Bind(comment) == nil {
		if err := app.CommentsRepo.Update(comment); err != nil {
			session.AddFlash(err.Error())
			_ = session.Save()
			c.Redirect(303, c.Request.RequestURI)
			return
		}
		c.Redirect(303, "/admin/comments")
	} else {
		session.AddFlash("Ошибка! Внимательно проверьте заполнение полей!")
		_ = session.Save()
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
		app.Logger.Error(err)
		w.WriteHeader(405)
		tmpl.Lookup("errors/405").Execute(w, helpers.ErrorData(err))
	}
}
*/

// CommentAdminDelete handles /admin/delete_comment route
func (app *Application) CommentAdminDelete(c *gin.Context) {
	id := helpers.Atouint(c.Request.PostFormValue("id"))
	if err := app.CommentsRepo.Delete(id); err != nil {
		app.Error(c, err)
		return
	}
	c.Redirect(303, "/admin/comments")
}

func (app *Application) notifyAdminOfComment(comment *models.Comment) {
	data := map[string]interface{}{
		"Comment": comment,
		"Token":   app.createTokenFromID(comment.ID),
	}
	var b bytes.Buffer
	if err := app.Template.Lookup("emails/question").Execute(&b, data); err != nil {
		app.Logger.Error(err)
		return
	}

	subject := fmt.Sprintf("Новый вопрос на сайте %s: %s", app.Config.Domain, comment.AuthorName)
	app.Emailer.NotifyAdmin(comment.AuthorEmail, subject, b.String())
}

// CommentsIndex handles GET /comments/:id route, where :id is the article id
func (app *Application) CommentsIndex(c *gin.Context) {
	id := helpers.Atouint(c.Param("id"))

	article, err := app.ArticlesRepo.Get(id)
	if err != nil {
		app.Error(c, err)
		return
	}

	totalCount, err := app.CommentsRepo.GetCountByArticle(id)
	if err != nil {
		app.Error(c, err)
		return
	}
	perPage := 15
	totalPages := int(math.Ceil(float64(totalCount) / float64(perPage)))
	currentPage := helpers.CurrentPage(c)

	list, err := app.CommentsRepo.GetByArticlePage(id, currentPage, perPage)
	if err != nil {
		app.Error(c, err)
		return
	}

	c.HTML(200, "comments/index", gin.H{
		"Title":           "Вопросы посетителей",
		"Active":          "comments",
		"Article":         &article,
		"MetaDescription": fmt.Sprintf("Вопросы и ответы к статье: %s", article.Name),
		"Pagination":      helpers.Paginator(currentPage, totalPages, c.Request.URL),
		"List":            list,
	})
}
