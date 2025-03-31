package controllers

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/denisbakhtin/medical/config"
	"github.com/denisbakhtin/medical/helpers"
	"github.com/denisbakhtin/medical/models"
	"github.com/denisbakhtin/medical/views"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/microcosm-cc/bluemonday"
	"gopkg.in/gomail.v2"
)

// ReviewShow handles /reviews/:id route
func ReviewShow(c *gin.Context) {
	db := models.GetDB()
	session := sessions.Default(c)

	idslug := c.Param("idslug")
	id := helpers.Atouint(strings.Split(idslug, "-")[0])
	review := &models.Review{}
	db.First(review, id)
	if review.ID == 0 || !review.Published {
		c.HTML(404, "errors/404", nil)
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
		"Authenticated":   (session.Get("user_id") != nil),
	})
}

// ReviewsIndex handles GET /reviews route
func ReviewsIndex(c *gin.Context) {
	session := sessions.Default(c)
	db := models.GetDB()
	flashes := session.Flashes()
	_ = session.Save()

	var list []models.Review
	db.Where("published = ?", true).Order("id desc").Find(&list)
	c.HTML(200, "reviews/index", gin.H{
		"Title":           "Кинезиология - отзывы",
		"Active":          c.Request.RequestURI,
		"List":            list,
		"Flash":           flashes,
		"MetaDescription": "Отзывы пациентов о работе врача кинезиолога Ростовцева Е.В...",
		"MetaKeywords":    "кинезиология отзывы, прикладная кинезиология отзывы, отзывы пациентов",
		"Authenticated":   (session.Get("user_id") != nil),
	})
}

// ReviewsAdminIndex handles GET /admin/reviews route
func ReviewsAdminIndex(c *gin.Context) {
	db := models.GetDB()

	var list []models.Review
	db.Order("id desc").Find(&list)
	c.HTML(200, "reviews/admin/index", gin.H{
		"Title":  "Отзывы",
		"Active": "reviews",
		"List":   list,
	})
}

// ReviewCreateGet handles /new_review get request
func ReviewCreateGet(c *gin.Context) {
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
func ReviewCreatePost(c *gin.Context) {
	db := models.GetDB()
	session := sessions.Default(c)

	if err := c.Request.ParseMultipartForm(32 << 20); err != nil {
		c.HTML(500, "errors/500", helpers.ErrorData(err))
		return
	}
	review := &models.Review{}
	if c.Bind(review) == nil {
		policy := bluemonday.StrictPolicy()
		review.Content = fmt.Sprintf("<p>%s</p>", policy.Sanitize(review.Content))
		// simple captcha check
		captcha, err := base64.StdEncoding.DecodeString(review.Captcha)
		if err != nil {
			c.HTML(500, "errors/500", helpers.ErrorData(err))
			return
		}
		if string(captcha) != "100.00" {
			c.HTML(400, "errors/400", nil)
			return
		}
		review.Published = false

		if mpartFile, mpartHeader, err := c.Request.FormFile("image"); err == nil {
			defer mpartFile.Close()
			review.Image, err = saveFile(mpartHeader, mpartFile)
			if err != nil {
				c.HTML(500, "errors/500", helpers.ErrorData(err))
				return
			}
		}

		if err := db.Create(review).Error; err != nil {
			c.HTML(400, "errors/400", helpers.ErrorData(err))
			return
		}
		notifyAdminOfReview(review)
		session.AddFlash("Спасибо! Ваш отзыв будет опубликован после проверки.")
	} else {
		session.AddFlash("Ошибка! Внимательно проверьте заполнение всех полей!")
	}
	_ = session.Save()
	c.Redirect(303, "/reviews")
}

// ReviewAdminCreateGet handles /admin/new_review get request
func ReviewAdminCreateGet(c *gin.Context) {
	session := sessions.Default(c)
	flashes := session.Flashes()
	_ = session.Save()
	db := models.GetDB()

	var articles []models.Article
	db.Where("published = ?", true).Find(&articles)
	c.HTML(200, "reviews/admin/form", gin.H{
		"Title":    "Новый отзыв",
		"Active":   "reviews",
		"Articles": articles,
		"Flash":    flashes,
	})
}

// ReviewAdminCreatePost handles /admin/new_review post request
func ReviewAdminCreatePost(c *gin.Context) {
	session := sessions.Default(c)
	db := models.GetDB()

	if err := c.Request.ParseMultipartForm(32 << 20); err != nil {
		c.HTML(500, "errors/500", helpers.ErrorData(err))
		return
	}
	review := &models.Review{}
	if c.Bind(review) == nil {
		review.ArticleID = helpers.Atouintr(c.Request.FormValue("article_id"))
		if mpartFile, mpartHeader, err := c.Request.FormFile("image"); err == nil {
			defer mpartFile.Close()
			review.Image, err = saveFile(mpartHeader, mpartFile)
			if err != nil {
				c.HTML(500, "errors/500", helpers.ErrorData(err))
				return
			}
		}

		if err := db.Create(review).Error; err != nil {
			c.HTML(400, "errors/400", helpers.ErrorData(err))
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
func ReviewAdminUpdateGet(c *gin.Context) {
	session := sessions.Default(c)
	flashes := session.Flashes()
	_ = session.Save()
	db := models.GetDB()

	id := helpers.Atouint(c.Param("id"))
	review := &models.Review{}
	db.First(review, id)
	if review.ID == 0 {
		c.HTML(404, "errors/404", nil)
		return
	}

	var articles []models.Article
	db.Where("published = ?", true).Find(&articles)
	c.HTML(200, "reviews/admin/form", gin.H{
		"Title":    "Редактировать отзыв",
		"Active":   "reviews",
		"Review":   review,
		"Articles": articles,
		"Flash":    flashes,
	})
}

// ReviewAdminUpdatePost handles /admin/edit_review/:id post request
func ReviewAdminUpdatePost(c *gin.Context) {
	session := sessions.Default(c)
	db := models.GetDB()

	if err := c.Request.ParseMultipartForm(32 << 20); err != nil {
		c.HTML(500, "errors/500", helpers.ErrorData(err))
		return
	}
	review := &models.Review{}
	if c.Bind(review) == nil {
		review.ArticleID = helpers.Atouintr(c.Request.FormValue("article_id"))
		if mpartFile, mpartHeader, err := c.Request.FormFile("image"); err == nil {
			defer mpartFile.Close()
			review.Image, err = saveFile(mpartHeader, mpartFile)
			if err != nil {
				c.HTML(500, "errors/500", helpers.ErrorData(err))
				return
			}
		}

		if err := db.Model(&models.Review{}).Updates(review).Error; err != nil {
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
func ReviewUpdateGet(c *gin.Context) {
	session := sessions.Default(c)
	flashes := session.Flashes()
	_ = session.Save()
	db := models.GetDB()

	id := getIDFromToken(c.Request.FormValue("token"))
	review := &models.Review{}
	db.First(review, id)
	if review.ID == 0 || review.Published {
		err := errors.New("отзыв не найден или уже был опубликован и не подлежит редактированию")
		c.HTML(404, "errors/404", helpers.ErrorData(err))
		return
	}

	var articles []models.Article
	db.Where("published = ?", true).Find(&articles)
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
func ReviewUpdatePost(c *gin.Context) {
	session := sessions.Default(c)
	db := models.GetDB()

	if err := c.Request.ParseMultipartForm(32 << 20); err != nil {
		c.HTML(500, "errors/500", helpers.ErrorData(err))
		return
	}
	review := &models.Review{}
	if err := c.Bind(review); err == nil {

		if mpartFile, mpartHeader, err := c.Request.FormFile("image"); err == nil {
			defer mpartFile.Close()
			review.Image, err = saveFile(mpartHeader, mpartFile)
			if err != nil {
				c.HTML(500, "errors/500", helpers.ErrorData(err))
				return
			}
		}

		if err := db.Model(&models.Review{}).Updates(review).Error; err != nil {
			session.AddFlash(err.Error())
			_ = session.Save()
			c.Redirect(303, c.Request.RequestURI)
			return
		}
		session.AddFlash("Отзыв был успешно сохранен")
	} else {
		log.Println(err)
		session.AddFlash("Ошибка! Внимательно проверьте заполнение всех полей")
	}
	_ = session.Save()
	c.Redirect(303, "/reviews")
}

// ReviewAdminDelete handles /admin/delete_review route
func ReviewAdminDelete(c *gin.Context) {
	db := models.GetDB()

	id := helpers.Atouint(c.Request.PostFormValue("id"))
	review := &models.Review{}
	db.First(review, id)
	if review.ID == 0 {
		c.HTML(404, "errors/404", nil)
	}

	if err := db.Delete(review).Error; err != nil {
		c.HTML(500, "errors/500", helpers.ErrorData(err))
		return
	}
	c.Redirect(303, "/admin/reviews")
}

func notifyAdminOfReview(review *models.Review) {
	tmpl := views.GetTemplates()
	go func() {
		data := map[string]interface{}{
			"Review": review,
			"Token":  createTokenFromID(review.ID),
		}
		var b bytes.Buffer
		if err := tmpl.Lookup("emails/review").Execute(&b, data); err != nil {
			log.Printf("ERROR: %s\n", err)
			return
		}

		smtp := config.GetConfig().SMTP
		msg := gomail.NewMessage()
		msg.SetHeader("From", smtp.From)
		msg.SetHeader("To", smtp.To)
		if len(smtp.Cc) > 0 {
			msg.SetHeader("Cc", smtp.Cc)
		}
		msg.SetHeader("Subject", fmt.Sprintf("Новый отзыв на сайте %s: %s", config.GetConfig().Domain, review.AuthorName))
		msg.SetBody(
			"text/html",
			b.String(),
		)

		port, _ := strconv.Atoi(smtp.Port)
		dialer := gomail.NewDialer(smtp.SMTP, port, smtp.User, smtp.Password)
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
