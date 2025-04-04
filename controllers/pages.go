package controllers

import (
	"strings"

	"github.com/denisbakhtin/medical/helpers"
	"github.com/denisbakhtin/medical/models"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// PageShow handles /pages/:id route
func (app *Application) PageShow(c *gin.Context) {
	session := sessions.Default(c)

	idslug := c.Param("idslug")
	id := helpers.Atouint(strings.Split(idslug, "-")[0])
	page, err := app.PagesRepo.GetPublished(id)
	if err != nil {
		app.Error(c, err)
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
		"Ogurl":           app.fullURL(page.URL()),
		"Ogtype":          "article",
		"Ogdescription":   page.MetaDescription,
		"Authenticated":   app.authenticated(session),
	})
}

// PagesAdminIndex handles GET /admin/pages route
func (app *Application) PagesAdminIndex(c *gin.Context) {
	list, err := app.PagesRepo.GetAll()
	if err != nil {
		app.Error(c, err)
		return
	}

	c.HTML(200, "pages/admin/index", gin.H{
		"Title":  "Страницы",
		"Active": "pages",
		"List":   list,
	})
}

// PageAdminCreateGet handles /admin/new_page get request
func (app *Application) PageAdminCreateGet(c *gin.Context) {
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
func (app *Application) PageAdminCreatePost(c *gin.Context) {
	session := sessions.Default(c)

	page := &models.Page{}
	if c.Bind(page) == nil {
		if err := app.PagesRepo.Create(page); err != nil {
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
func (app *Application) PageAdminUpdateGet(c *gin.Context) {
	session := sessions.Default(c)
	flashes := session.Flashes()
	_ = session.Save()

	id := helpers.Atouint(c.Param("id"))
	page, err := app.PagesRepo.Get(id)
	if err != nil {
		app.Error(c, err)
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
func (app *Application) PageAdminUpdatePost(c *gin.Context) {
	session := sessions.Default(c)

	page := &models.Page{}
	if c.Bind(page) == nil {
		if err := app.PagesRepo.Update(page); err != nil {
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
func (app *Application) PageAdminDelete(c *gin.Context) {
	id := helpers.Atouint(c.Request.PostFormValue("id"))
	if err := app.PagesRepo.Delete(id); err != nil {
		app.Error(c, err)
		return
	}
	c.Redirect(303, "/admin/pages")
}
