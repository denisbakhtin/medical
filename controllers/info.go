package controllers

import (
	"fmt"
	"strings"

	"github.com/denisbakhtin/medical/helpers"
	"github.com/denisbakhtin/medical/models"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// InfoShow handles GET /info/:id-slug route
func InfoShow(c *gin.Context) {
	db := models.GetDB()
	session := sessions.Default(c)

	idslug := c.Param("idslug")
	id := helpers.Atouint(strings.Split(idslug, "-")[0])
	info := &models.Info{}
	db.First(info, id)
	if info.ID == 0 || !info.Published {
		c.HTML(404, "errors/404", nil)
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
		"Ogurl":           fmt.Sprintf("http://%s%s", c.Request.Host, info.URL()),
		"Ogtype":          "article",
		"Flash":           flashes,
		"Authenticated":   (session.Get("user_id") != nil),
	})
}

// InfoAdminIndex handles GET /admin/info route
func InfoAdminIndex(c *gin.Context) {
	db := models.GetDB()

	var list []models.Info
	if err := db.Order("published desc, id desc").Find(&list).Error; err != nil {
		c.HTML(500, "errors/500", helpers.ErrorData(err))
		return
	}
	c.HTML(200, "info/admin/index", gin.H{
		"Title":  "Информационные материалы",
		"Active": "info",
		"List":   list,
	})
}

// InfoAdminCreateGet handles /admin/new_info get request
func InfoAdminCreateGet(c *gin.Context) {
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
func InfoAdminCreatePost(c *gin.Context) {
	session := sessions.Default(c)
	db := models.GetDB()

	info := &models.Info{}
	if c.Bind(info) == nil {
		if err := db.Create(info).Error; err != nil {
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
func InfoAdminUpdateGet(c *gin.Context) {
	session := sessions.Default(c)
	flashes := session.Flashes()
	_ = session.Save()
	db := models.GetDB()

	id := helpers.Atouint(c.Param("id"))
	info := &models.Info{}
	db.First(info, id)
	if info.ID == 0 {
		c.HTML(404, "errors/404", nil)
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
func InfoAdminUpdatePost(c *gin.Context) {
	session := sessions.Default(c)
	db := models.GetDB()

	info := &models.Info{}
	if c.Bind(info) == nil {
		if err := db.Save(info).Error; err != nil {
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
func InfoAdminDelete(c *gin.Context) {
	db := models.GetDB()

	id := helpers.Atouint(c.Request.PostFormValue("id"))
	info := &models.Info{}
	db.First(info, id)
	if info.ID == 0 {
		c.HTML(404, "errors/404", nil)
		return
	}

	if err := db.Delete(info).Error; err != nil {
		c.HTML(500, "errors/500", helpers.ErrorData(err))
		return
	}
	c.Redirect(303, "/admin/info")
}
