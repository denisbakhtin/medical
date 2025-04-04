package controllers

import (
	"fmt"
	"math"
	"strings"

	"github.com/denisbakhtin/medical/helpers"
	"github.com/denisbakhtin/medical/models"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// ArticleShow handles GET /articles/:id-slug route
func (app *Application) ArticleShow(c *gin.Context) {
	session := sessions.Default(c)
	idslug := c.Param("idslug")
	id := helpers.Atouint(strings.Split(idslug, "-")[0])
	article, err := app.ArticlesRepo.GetPublished(id)
	if err != nil {
		app.Error(c, err)
		return
	}
	// redirect to canonical url
	if c.Request.URL.Path != article.URL() {
		c.Redirect(301, article.URL())
		return
	}
	reviews, err := app.ReviewsRepo.GetPublishedByArticle(article.ID)
	if err != nil {
		app.Error(c, err)
		return
	}
	topcomments, err := app.CommentsRepo.GetTopByArticle(article.ID)
	if err != nil {
		app.Error(c, err)
		return
	}

	if len(topcomments) > 0 {
		article.Comments = topcomments
	} else {
		article.Comments, err = app.CommentsRepo.GetUnpublishedByArticle(article.ID)
		if err != nil {
			app.Error(c, err)
			return
		}
	}

	imageurl := ""
	if img := article.GetImage(); len(img) > 0 {
		imageurl = fmt.Sprintf("%s%s", app.Config.FullDomain, img)
	}
	flashes := session.Flashes()
	_ = session.Save()
	c.HTML(200, "articles/show", gin.H{
		"Article":         article,
		"Testimonials":    reviews,
		"Title":           article.Name,
		"Active":          "/articles",
		"MetaDescription": article.MetaDescription,
		"MetaKeywords":    article.MetaKeywords,
		"Ogheadprefix":    "og: http://ogp.me/ns# fb: http://ogp.me/ns/fb# article: http://ogp.me/ns/article#",
		"Ogtitle":         article.Name,
		"Ogurl":           app.fullURL(article.URL()),
		"Ogtype":          "article",
		"Ogdescription":   article.Excerpt,
		"Ogimage":         imageurl,
		"Flash":           flashes,
		"Authenticated":   app.authenticated(session),
	})
}

// ArticlesIndex handles GET /articles route
func (app *Application) ArticlesIndex(c *gin.Context) {
	session := sessions.Default(c)

	list, err := app.ArticlesRepo.GetAllPublished()
	if err != nil {
		app.Error(c, err)
		return
	}

	totalInfos, err := app.InfosRepo.GetPublishedCount()
	if err != nil {
		app.Error(c, err)
		return
	}

	infosPerPage := 8
	totalPages := int(math.Ceil(float64(totalInfos) / float64(infosPerPage)))
	currentPage := helpers.CurrentPage(c)
	infos, err := app.InfosRepo.GetPublishedPage(currentPage, infosPerPage)
	if err != nil {
		app.Error(c, err)
		return
	}

	c.HTML(200, "articles/index", gin.H{
		"Title":           "Кинезиология во врачебной практике",
		"Active":          c.Request.RequestURI,
		"List":            list,
		"Infos":           infos,
		"MetaDescription": "Статьи о кинезиологической практике лечения заболеваний опорно-двигательного аппарата...",
		"MetaKeywords":    "кинезиология, статьи, лечение болей, прикладная кинезиология",
		"Authenticated":   app.authenticated(session),
		"Pagination":      helpers.Paginator(currentPage, totalPages, c.Request.URL),
		"CurrentPage":     currentPage,
	})
}

// ArticlesAdminIndex handles GET /admin/articles route
func (app *Application) ArticlesAdminIndex(c *gin.Context) {
	list, err := app.ArticlesRepo.GetAll()
	if err != nil {
		app.Error(c, err)
		return
	}
	c.HTML(200, "articles/admin/index", gin.H{
		"Title":  "Статьи",
		"Active": "articles",
		"List":   list,
	})
}

// ArticleAdminCreateGet handles /admin/new_article route
func (app *Application) ArticleAdminCreateGet(c *gin.Context) {
	session := sessions.Default(c)
	flashes := session.Flashes()
	_ = session.Save()

	c.HTML(200, "articles/admin/form", gin.H{
		"Title":  "Новая статья",
		"Active": "articles",
		"Flash":  flashes,
	})
}

// ArticleAdminCreatePost handles /admin/new_article post request
func (app *Application) ArticleAdminCreatePost(c *gin.Context) {
	session := sessions.Default(c)

	article := &models.Article{}
	if c.Bind(article) == nil {
		if err := app.ArticlesRepo.Create(article); err != nil {
			session.AddFlash(err.Error())
			_ = session.Save()
			c.Redirect(303, "/admin/new_article")
			return
		}
		c.Redirect(303, "/admin/articles")
	} else {
		session.AddFlash("Ошибка! Проверьте внимательно заполнение всех полей!")
		_ = session.Save()
		c.Redirect(303, "/admin/new_article")
	}
}

// ArticleAdminUpdateGet handles /admin/edit_article/:id get request
func (app *Application) ArticleAdminUpdateGet(c *gin.Context) {
	session := sessions.Default(c)
	flashes := session.Flashes()
	_ = session.Save()

	id := helpers.Atouint(c.Param("id"))
	article, err := app.ArticlesRepo.Get(id)
	if err != nil {
		app.Error(c, err)
		return
	}

	c.HTML(200, "articles/admin/form", gin.H{
		"Title":   "Редактировать статью",
		"Active":  "articles",
		"Article": article,
		"Flash":   flashes,
	})
}

// ArticleAdminUpdatePost handles /admin/edit_article/:id post request
func (app *Application) ArticleAdminUpdatePost(c *gin.Context) {
	session := sessions.Default(c)

	article := &models.Article{}
	if c.Bind(article) == nil {
		if err := app.ArticlesRepo.Update(article); err != nil {
			session.AddFlash(err.Error())
			_ = session.Save()
			c.Redirect(303, c.Request.RequestURI)
			return
		}
		c.Redirect(303, "/admin/articles")
	} else {
		session.AddFlash("Ошибка! Проверьте внимательно заполнение всех полей!")
		_ = session.Save()
		c.Redirect(303, c.Request.RequestURI)
	}
}

// ArticleAdminDelete handles /admin/delete_article route
func (app *Application) ArticleAdminDelete(c *gin.Context) {
	id := helpers.Atouint(c.Request.PostFormValue("id"))
	if err := app.ArticlesRepo.Delete(id); err != nil {
		app.Error(c, err)
		return
	}
	c.Redirect(303, "/admin/articles")
}
