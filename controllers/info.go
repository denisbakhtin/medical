package controllers

import (
	"strings"

	"github.com/denisbakhtin/medical/helpers"
	"github.com/denisbakhtin/medical/models"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// InfoShow handles GET /info/:id-slug route
func (app *Application) InfoShow(c *gin.Context) {
	session := sessions.Default(c)

	idslug := c.Param("idslug")
	id := helpers.Atouint(strings.Split(idslug, "-")[0])
	info, err := app.InfosRepo.GetPublished(id)
	if err != nil {
		app.Error(c, err)
		return
	}
	// redirect to canonical url
	if c.Request.URL.Path != info.URL() {
		c.Redirect(301, info.URL())
		return
	}
	flashes := session.Flashes()
	_ = session.Save()
	c.HTML(200, "info/show", gin.H{
		"Info":            info,
		"Title":           info.Name,
		"Active":          "/info",
		"MetaDescription": info.MetaDescription,
		"MetaKeywords":    info.MetaKeywords,
		"Ogheadprefix":    "og: http://ogp.me/ns# fb: http://ogp.me/ns/fb# article: http://ogp.me/ns/article#",
		"Ogtitle":         info.Name,
		"Ogurl":           app.fullURL(info.URL()),
		"Ogtype":          "article",
		"Flash":           flashes,
		"Authenticated":   app.authenticated(session),
	})
}

// InfoAdminIndex handles GET /admin/info route
func (app *Application) InfoAdminIndex(c *gin.Context) {
	list, err := app.InfosRepo.GetAll()
	if err != nil {
		app.Error(c, err)
		return
	}

	c.HTML(200, "info/admin/index", gin.H{
		"Title":    "Информационные материалы",
		"Active":   "info",
		"List":     list,
		"AllCount": len(list),
		"PublishedCount": helpers.CountFunc(list, func(e models.Info) bool {
			return e.Published
		}),
		"UnpublishedCount": helpers.CountFunc(list, func(e models.Info) bool {
			return !e.Published
		}),
	})
}

// InfoAdminCreateGet handles /admin/new_info get request
func (app *Application) InfoAdminCreateGet(c *gin.Context) {
	session := sessions.Default(c)
	flashes := session.Flashes()
	_ = session.Save()

	c.HTML(200, "info/admin/form", gin.H{
		"Title":  "Новый материал",
		"Active": "info",
		"Flash":  flashes,
	})
}

// InfoAdminCreatePost handles /admin/new_info post request
func (app *Application) InfoAdminCreatePost(c *gin.Context) {
	session := sessions.Default(c)

	info := &models.Info{}
	if c.Bind(info) == nil {
		if err := app.InfosRepo.Create(info); err != nil {
			session.AddFlash(err.Error())
			_ = session.Save()
			c.Redirect(303, "/admin/new_info")
			return
		}
		c.Redirect(303, "/admin/info")
	} else {
		session.AddFlash("Ошибка! Проверьте внимательно заполнение всех полей!")
		_ = session.Save()
		c.Redirect(303, "/admin/new_info")
	}
}

// InfoAdminUpdateGet handles /admin/edit_info/:id get request
func (app *Application) InfoAdminUpdateGet(c *gin.Context) {
	session := sessions.Default(c)
	flashes := session.Flashes()
	_ = session.Save()

	id := helpers.Atouint(c.Param("id"))
	info, err := app.InfosRepo.Get(id)
	if err != nil {
		app.Error(c, err)
		return
	}

	c.HTML(200, "info/admin/form", gin.H{
		"Title":  "Редактировать материал",
		"Active": "info",
		"Info":   info,
		"Flash":  flashes,
	})
}

// InfoAdminUpdatePost handles /admin/edit_info/:id post request
func (app *Application) InfoAdminUpdatePost(c *gin.Context) {
	session := sessions.Default(c)

	info := &models.Info{}
	if c.Bind(info) == nil {
		if err := app.InfosRepo.Update(info); err != nil {
			session.AddFlash(err.Error())
			_ = session.Save()
			c.Redirect(303, c.Request.RequestURI)
			return
		}
		c.Redirect(303, "/admin/info")
	} else {
		session.AddFlash("Ошибка! Проверьте внимательно заполнение всех полей!")
		_ = session.Save()
		c.Redirect(303, c.Request.RequestURI)
	}
}

// InfoAdminDelete handles /admin/delete_info route
func (app *Application) InfoAdminDelete(c *gin.Context) {
	id := helpers.Atouint(c.Request.PostFormValue("id"))

	if err := app.InfosRepo.Delete(id); err != nil {
		app.Error(c, err)
		return
	}
	c.Redirect(303, "/admin/info")
}
