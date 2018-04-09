package controllers

import (
	"fmt"
	"strings"

	"github.com/denisbakhtin/medical/helpers"
	"github.com/denisbakhtin/medical/models"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
)

//ArticleShow handles GET /articles/:id-slug route
func ArticleShow(c *gin.Context) {
	db := models.GetDB()
	session := sessions.Default(c)

	idslug := c.Param("idslug")
	id := helpers.Atouint(strings.Split(idslug, "-")[0])
	article := &models.Article{}
	db.First(article, id)
	if article.ID == 0 || !article.Published {
		c.HTML(404, "errors/404", nil)
		return
	}
	//redirect to canonical url
	if c.Request.URL.Path != article.URL() {
		c.Redirect(303, article.URL())
		return
	}
	var testimonials []models.Review
	db.Where("published = ? and article_id = ?", true, article.ID).Order("created_at desc").Find(&testimonials)
	topComments := models.GetTopComments(article.ID)
	comments := models.GetComments(article.ID)
	article.Comments = append(topComments, comments...)
	imageurl := ""
	if img := article.GetImage(); len(img) > 0 {
		imageurl = fmt.Sprintf("http://%s%s", c.Request.Host, img)
	}
	flashes := session.Flashes()
	session.Save()
	c.HTML(200, "articles/show", gin.H{
		"Article":         article,
		"Testimonials":    testimonials,
		"Title":           article.Name,
		"Active":          "/articles",
		"MetaDescription": article.MetaDescription,
		"MetaKeywords":    article.MetaKeywords,
		"Ogheadprefix":    "og: http://ogp.me/ns# fb: http://ogp.me/ns/fb# article: http://ogp.me/ns/article#",
		"Ogtitle":         article.Name,
		"Ogurl":           fmt.Sprintf("http://%s%s", c.Request.Host, article.URL()),
		"Ogtype":          "article",
		"Ogdescription":   article.Excerpt,
		"Ogimage":         imageurl,
		"Flash":           flashes,
		"Authenticated":   (session.Get("user_id") != nil),
	})
}

//ArticlesIndex handles GET /articles route
func ArticlesIndex(c *gin.Context) {
	db := models.GetDB()
	session := sessions.Default(c)

	var list []models.Article
	if err := db.Where("published = ?", true).Order("id desc").Find(&list).Error; err != nil {
		c.HTML(500, "errors/500", helpers.ErrorData(err))
		return
	}

	var infos []models.Info
	if err := db.Where("published = ?", true).Order("id desc").Find(&infos).Error; err != nil {
		c.HTML(500, "errors/500", helpers.ErrorData(err))
		return
	}
	c.HTML(200, "articles/index", gin.H{
		"Title":           "Кинезиология во врачебной практике",
		"Active":          c.Request.RequestURI,
		"List":            list,
		"Infos":           infos,
		"MetaDescription": "Статьи о кинезиологической практике лечения заболеваний опорно-двигательного аппарата...",
		"MetaKeywords":    "кинезиология, статьи, лечение болей, прикладная кинезиология",
		"Authenticated":   (session.Get("user_id") != nil),
	})
}

//ArticlesAdminIndex handles GET /admin/articles route
func ArticlesAdminIndex(c *gin.Context) {
	db := models.GetDB()

	var list []models.Article
	if err := db.Order("published desc, id desc").Find(&list).Error; err != nil {
		c.HTML(500, "errors/500", helpers.ErrorData(err))
		return
	}
	c.HTML(200, "articles/admin/index", gin.H{
		"Title":  "Статьи",
		"Active": "articles",
		"List":   list,
	})
}

//ArticleAdminCreateGet handles /admin/new_article route
func ArticleAdminCreateGet(c *gin.Context) {
	session := sessions.Default(c)
	flashes := session.Flashes()
	session.Save()

	c.HTML(200, "articles/admin/form", gin.H{
		"Title":  "Новая статья",
		"Active": "articles",
		"Flash":  flashes,
	})
}

//ArticleAdminCreatePost handles /admin/new_article post request
func ArticleAdminCreatePost(c *gin.Context) {
	session := sessions.Default(c)
	db := models.GetDB()

	article := &models.Article{}
	if c.Bind(article) == nil {
		if err := db.Create(article).Error; err != nil {
			session.AddFlash(err.Error())
			session.Save()
			c.Redirect(303, "/admin/new_article")
			return
		}
		c.Redirect(303, "/admin/articles")
	} else {
		session.AddFlash("Ошибка! Проверьте внимательно заполнение всех полей!")
		session.Save()
		c.Redirect(303, "/admin/new_article")
	}
}

//ArticleAdminUpdateGet handles /admin/edit_article/:id get request
func ArticleAdminUpdateGet(c *gin.Context) {
	session := sessions.Default(c)
	flashes := session.Flashes()
	session.Save()
	db := models.GetDB()

	id := c.Param("id")
	article := &models.Article{}
	db.First(article, id)
	if article.ID == 0 {
		c.HTML(404, "errors/404", nil)
		return
	}

	c.HTML(200, "articles/admin/form", gin.H{
		"Title":   "Редактировать статью",
		"Active":  "articles",
		"Article": article,
		"Flash":   flashes,
	})
}

//ArticleAdminUpdatePost handles /admin/edit_article/:id post request
func ArticleAdminUpdatePost(c *gin.Context) {
	session := sessions.Default(c)
	db := models.GetDB()

	article := &models.Article{}
	if c.Bind(article) == nil {
		if err := db.Save(article).Error; err != nil {
			session.AddFlash(err.Error())
			session.Save()
			c.Redirect(303, c.Request.RequestURI)
			return
		}
		c.Redirect(303, "/admin/articles")
	} else {
		session.AddFlash("Ошибка! Проверьте внимательно заполнение всех полей!")
		session.Save()
		c.Redirect(303, c.Request.RequestURI)
	}
}

//ArticleAdminDelete handles /admin/delete_article route
func ArticleAdminDelete(c *gin.Context) {
	db := models.GetDB()

	article := &models.Article{}
	db.First(article, c.Request.PostFormValue("id"))
	if article.ID == 0 {
		c.HTML(404, "errors/404", nil)
		return
	}

	if err := db.Delete(article).Error; err != nil {
		c.HTML(500, "errors/500", helpers.ErrorData(err))
		return
	}
	c.Redirect(303, "/admin/articles")
}
