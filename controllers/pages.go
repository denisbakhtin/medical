package controllers

import (
	"fmt"
	"strings"

	"github.com/denisbakhtin/medical/helpers"
	"github.com/denisbakhtin/medical/models"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// PageShow handles /pages/:id route
func PageShow(c *gin.Context) {
	db := models.GetDB()
	session := sessions.Default(c)

	idslug := c.Param("idslug")
	id := helpers.Atouint(strings.Split(idslug, "-")[0])
	page := &models.Page{}
	db.First(page, id)
	if page.ID == 0 || !page.Published {
		c.HTML(404, "errors/404", nil)
		return
	}
	// redirect to canonical url
	if c.Request.URL.Path != page.URL() {
		c.Redirect(301, page.URL())
		return
	}
	c.HTML(200, "pages/show", gin.H{
		"Page":            page,
		"Title":           page.Name,
		"Active":          page.URL(),
		"MetaDescription": page.MetaDescription,
		"MetaKeywords":    page.MetaKeywords,
		"Ogtitle":         page.Name,
		"Ogurl":           fmt.Sprintf("http://%s%s", c.Request.Host, page.URL()),
		"Ogtype":          "article",
		"Ogdescription":   page.MetaDescription,
		"Authenticated":   (session.Get("user_id") != nil),
	})
}

// PagesAdminIndex handles GET /admin/pages route
func PagesAdminIndex(c *gin.Context) {
	db := models.GetDB()

	var list []models.Page
	db.Order("published desc, id desc").Find(&list)
	c.HTML(200, "pages/admin/index", gin.H{
		"Title":  "Страницы",
		"Active": "pages",
		"List":   list,
	})
}

// PageAdminCreateGet handles /admin/new_page get request
func PageAdminCreateGet(c *gin.Context) {
	session := sessions.Default(c)
	flashes := session.Flashes()
	_ = session.Save()

	c.HTML(200, "pages/admin/form", gin.H{
		"Title":  "Новая страница",
		"Active": "pages",
		"Flash":  flashes,
	})
}

// PageAdminCreatePost handles /admin/new_page post request
func PageAdminCreatePost(c *gin.Context) {
	session := sessions.Default(c)
	db := models.GetDB()

	page := &models.Page{}
	if c.Bind(page) == nil {
		if err := db.Create(page).Error; err != nil {
			session.AddFlash(err.Error())
			_ = session.Save()
			c.Redirect(303, "/admin/new_page")
			return
		}
		c.Redirect(303, "/admin/pages")
	} else {
		session.AddFlash("Ошибка! Проверьте заполнение всех полей!")
		_ = session.Save()
		c.Redirect(303, "/admin/new_page")
	}
}

// PageAdminUpdateGet handles /admin/edit_page/:id get request
func PageAdminUpdateGet(c *gin.Context) {
	session := sessions.Default(c)
	flashes := session.Flashes()
	_ = session.Save()
	db := models.GetDB()

	id := c.Param("id")
	page := &models.Page{}
	db.First(page, id)
	if page.ID == 0 {
		c.HTML(404, "errors/400", nil)
		return
	}

	c.HTML(200, "pages/admin/form", gin.H{
		"Title":  "Редактировать страницу",
		"Active": "pages",
		"Page":   page,
		"Flash":  flashes,
	})
}

// PageAdminUpdatePost handles /admin/edit_page/:id post request
func PageAdminUpdatePost(c *gin.Context) {
	session := sessions.Default(c)
	db := models.GetDB()

	page := &models.Page{}
	if c.Bind(page) == nil {

		if err := db.Save(page).Error; err != nil {
			session.AddFlash(err.Error())
			_ = session.Save()
			c.Redirect(303, c.Request.RequestURI)
			return
		}
		c.Redirect(303, "/admin/pages")
	} else {
		session.AddFlash("Ошибка! Проверьте внимательно заполнение всех полей!")
		_ = session.Save()
		c.Redirect(303, c.Request.RequestURI)
	}
}

// PageAdminDelete handles /admin/delete_page route
func PageAdminDelete(c *gin.Context) {
	db := models.GetDB()

	page := &models.Page{}
	db.First(page, c.Request.PostFormValue("id"))
	if page.ID == 0 {
		c.HTML(404, "errors/404", nil)
	}

	if err := db.Delete(page).Error; err != nil {
		c.HTML(500, "errors/500", helpers.ErrorData(err))
		return
	}
	c.Redirect(303, "/admin/pages")
}
