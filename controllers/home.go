package controllers

import (
	"github.com/denisbakhtin/medical/models"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

const homePageId = 1

// Home handles GET / route
func Home(c *gin.Context) {
	db := models.GetDB()
	page := &models.Page{}
	db.First(page, homePageId)
	session := sessions.Default(c)
	flashes := session.Flashes()
	_ = session.Save()

	c.HTML(200, "home/show", gin.H{
		"Title":           "Кинезиология Миобаланс",
		"Page":            page,
		"Active":          "/",
		"Flash":           flashes,
		"TitleSuffix":     "| Доктор Ростовцев Е.В.",
		"MetaDescription": "Прикладная кинезиология МиоБаланс - восстановление баланса обмена веществ, опорно-двигательного аппарата и нервной системы...",
		"Authenticated":   (session.Get("user_id") != nil),
	})
}
