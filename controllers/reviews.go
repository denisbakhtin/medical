package controllers

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"github.com/denisbakhtin/medical/helpers"
	"github.com/denisbakhtin/medical/models"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/microcosm-cc/bluemonday"
)

// ReviewShow handles /reviews/:id route
func (app *Application) ReviewShow(c *gin.Context) {
	session := sessions.Default(c)

	idslug := c.Param("idslug")
	id := helpers.Atouint(strings.Split(idslug, "-")[0])
	review, err := app.ReviewsRepo.GetPublished(id)
	if err != nil {
		app.Error(c, err)
		return
	}
	// redirect to canonical url
	if c.Request.URL.Path != review.URL() {
		c.Redirect(301, review.URL())
		return
	}
	c.HTML(200, "reviews/show", gin.H{
		"Review":          review,
		"Title":           "Отзыв о работе кинезиолога: " + review.AuthorName,
		"Active":          "/reviews",
		"MetaDescription": review.MetaDescription,
		"MetaKeywords":    review.MetaKeywords,
		"Authenticated":   app.authenticated(session),
	})
}

// ReviewsIndex handles GET /reviews route
func (app *Application) ReviewsIndex(c *gin.Context) {
	session := sessions.Default(c)
	flashes := session.Flashes()
	_ = session.Save()

	list, err := app.ReviewsRepo.GetAllPublished()
	if err != nil {
		app.Error(c, err)
		return
	}

	c.HTML(200, "reviews/index", gin.H{
		"Title":           "Кинезиология - отзывы",
		"Active":          c.Request.RequestURI,
		"List":            list,
		"Flash":           flashes,
		"MetaDescription": "Отзывы пациентов о работе врача кинезиолога Ростовцева Е.В...",
		"MetaKeywords":    "кинезиология отзывы, прикладная кинезиология отзывы, отзывы пациентов",
		"Authenticated":   app.authenticated(session),
	})
}

// ReviewsAdminIndex handles GET /admin/reviews route
func (app *Application) ReviewsAdminIndex(c *gin.Context) {
	list, err := app.ReviewsRepo.GetAll()
	if err != nil {
		app.Error(c, err)
		return
	}

	c.HTML(200, "reviews/admin/index", gin.H{
		"Title":  "Отзывы",
		"Active": "reviews",
		"List":   list,
	})
}

// ReviewCreateGet handles /new_review get request
func (app *Application) ReviewCreateGet(c *gin.Context) {
	session := sessions.Default(c)
	flashes := session.Flashes()
	_ = session.Save()

	c.HTML(200, "reviews/form", gin.H{
		"Title":  "Новый отзыв",
		"Active": "reviews",
		"Flash":  flashes,
	})
}

// ReviewCreatePost handles /new_review post request
func (app *Application) ReviewCreatePost(c *gin.Context) {
	session := sessions.Default(c)

	if err := c.Request.ParseMultipartForm(32 << 20); err != nil {
		app.Error(c, err)
		return
	}
	review := &models.Review{}
	if c.Bind(review) == nil {
		policy := bluemonday.StrictPolicy()
		review.Content = fmt.Sprintf("<p>%s</p>", policy.Sanitize(review.Content))
		// simple captcha check
		captcha, err := base64.StdEncoding.DecodeString(review.Captcha)
		if err != nil {
			app.Error(c, err)
			return
		}
		if string(captcha) != "100.00" {
			c.HTML(400, "errors/400", nil)
			return
		}
		review.Published = false

		if mpartFile, mpartHeader, err := c.Request.FormFile("image"); err == nil {
			defer mpartFile.Close()
			review.Image, err = app.saveFile(mpartHeader, mpartFile)
			if err != nil {
				app.Error(c, err)
				return
			}
		}

		if err := app.ReviewsRepo.Create(review); err != nil {
			app.Error(c, err)
			return
		}
		app.notifyAdminOfReview(review)
		session.AddFlash("Спасибо! Ваш отзыв будет опубликован после проверки.")
	} else {
		session.AddFlash("Ошибка! Внимательно проверьте заполнение всех полей!")
	}
	_ = session.Save()
	c.Redirect(303, "/reviews")
}

// ReviewAdminCreateGet handles /admin/new_review get request
func (app *Application) ReviewAdminCreateGet(c *gin.Context) {
	session := sessions.Default(c)
	flashes := session.Flashes()
	_ = session.Save()

	articles, err := app.ArticlesRepo.GetAllPublished()
	if err != nil {
		app.Error(c, err)
		return
	}
	c.HTML(200, "reviews/admin/form", gin.H{
		"Title":    "Новый отзыв",
		"Active":   "reviews",
		"Articles": articles,
		"Flash":    flashes,
	})
}

// ReviewAdminCreatePost handles /admin/new_review post request
func (app *Application) ReviewAdminCreatePost(c *gin.Context) {
	session := sessions.Default(c)

	if err := c.Request.ParseMultipartForm(32 << 20); err != nil {
		app.Error(c, err)
		return
	}
	review := &models.Review{}
	if c.Bind(review) == nil {
		review.ArticleID = helpers.Atouintr(c.Request.FormValue("article_id"))
		if mpartFile, mpartHeader, err := c.Request.FormFile("image"); err == nil {
			defer mpartFile.Close()
			review.Image, err = app.saveFile(mpartHeader, mpartFile)
			if err != nil {
				app.Error(c, err)
				return
			}
		}

		if err := app.ReviewsRepo.Create(review); err != nil {
			app.Error(c, err)
			return
		}
		c.Redirect(303, "/admin/reviews")
	} else {
		session.AddFlash("Ошибка! Проверьте внимательно заполнение всех полей!")
		_ = session.Save()
		c.Redirect(303, "/admin/new_review")
	}
}

// ReviewAdminUpdateGet handles /admin/edit_review/:id get request
func (app *Application) ReviewAdminUpdateGet(c *gin.Context) {
	session := sessions.Default(c)
	flashes := session.Flashes()
	_ = session.Save()

	id := helpers.Atouint(c.Param("id"))
	review, err := app.ReviewsRepo.Get(id)
	if err != nil {
		app.Error(c, err)
		return
	}

	articles, err := app.ArticlesRepo.GetAllPublished()
	if err != nil {
		app.Error(c, err)
		return
	}

	c.HTML(200, "reviews/admin/form", gin.H{
		"Title":    "Редактировать отзыв",
		"Active":   "reviews",
		"Review":   review,
		"Articles": articles,
		"Flash":    flashes,
	})
}

// ReviewAdminUpdatePost handles /admin/edit_review/:id post request
func (app *Application) ReviewAdminUpdatePost(c *gin.Context) {
	session := sessions.Default(c)

	if err := c.Request.ParseMultipartForm(32 << 20); err != nil {
		app.Error(c, err)
		return
	}
	review := &models.Review{}
	if c.Bind(review) == nil {
		review.ArticleID = helpers.Atouintr(c.Request.FormValue("article_id"))
		if mpartFile, mpartHeader, err := c.Request.FormFile("image"); err == nil {
			defer mpartFile.Close()
			review.Image, err = app.saveFile(mpartHeader, mpartFile)
			if err != nil {
				app.Error(c, err)
				return
			}
		}

		if err := app.ReviewsRepo.Update(review); err != nil {
			session.AddFlash(err.Error())
			_ = session.Save()
			c.Redirect(303, c.Request.RequestURI)
			return
		}
		c.Redirect(303, "/admin/reviews")
	} else {
		session.AddFlash("Ошибка! Проверьте внимательно заполнение всех полей!")
		_ = session.Save()
		c.Redirect(303, c.Request.RequestURI)
	}
}

// ReviewUpdateGet handles /edit_review?token=:secure_token get request
func (app *Application) ReviewUpdateGet(c *gin.Context) {
	session := sessions.Default(c)
	flashes := session.Flashes()
	_ = session.Save()

	idStr := app.getIDFromToken(c.Request.FormValue("token"))
	id := helpers.Atouint(idStr)
	review, err := app.ReviewsRepo.Get(id)
	if err != nil || review.Published {
		err = errors.New("отзыв не найден или уже был опубликован и не подлежит редактированию")
		app.Error(c, err)
		return
	}

	articles, err := app.ArticlesRepo.GetAllPublished()
	if err != nil {
		app.Error(c, err)
		return
	}

	review.Published = true // set default to true
	c.HTML(200, "reviews/form", gin.H{
		"Title":      "Редактировать отзыв",
		"Articles":   articles,
		"Active":     "reviews",
		"Review":     review,
		"SecureEdit": true,
		"Flash":      flashes,
	})
}

// ReviewUpdatePost handles /edit_review post request
func (app *Application) ReviewUpdatePost(c *gin.Context) {
	session := sessions.Default(c)

	if err := c.Request.ParseMultipartForm(32 << 20); err != nil {
		app.Error(c, err)
		return
	}
	review := &models.Review{}
	if err := c.Bind(review); err == nil {
		if mpartFile, mpartHeader, err := c.Request.FormFile("image"); err == nil {
			defer mpartFile.Close()
			review.Image, err = app.saveFile(mpartHeader, mpartFile)
			if err != nil {
				app.Error(c, err)
				return
			}
		}

		if err := app.ReviewsRepo.Update(review); err != nil {
			session.AddFlash(err.Error())
			_ = session.Save()
			c.Redirect(303, c.Request.RequestURI)
			return
		}
		session.AddFlash("Отзыв был успешно сохранен")
	} else {
		app.Logger.Error(err)
		session.AddFlash("Ошибка! Внимательно проверьте заполнение всех полей")
	}
	_ = session.Save()
	c.Redirect(303, "/reviews")
}

// ReviewAdminDelete handles /admin/delete_review route
func (app *Application) ReviewAdminDelete(c *gin.Context) {
	id := helpers.Atouint(c.Request.PostFormValue("id"))

	if err := app.ReviewsRepo.Delete(id); err != nil {
		app.Error(c, err)
		return
	}
	c.Redirect(303, "/admin/reviews")
}

func (app *Application) notifyAdminOfReview(review *models.Review) {
	data := map[string]any{
		"Review": review,
		"Token":  app.createTokenFromID(review.ID),
	}
	var b bytes.Buffer
	if err := app.Template.Lookup("emails/review").Execute(&b, data); err != nil {
		app.Logger.Error(err)
		return
	}

	subject := fmt.Sprintf("Новый отзыв на сайте %s: %s", app.Config.Domain, review.AuthorName)
	app.Emailer.NotifyAdmin("", subject, b.String())
}
