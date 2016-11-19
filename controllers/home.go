package controllers

import (
	"github.com/denisbakhtin/medical/models"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
)

//Home handles GET / route
func Home(c *gin.Context) {
	db := models.GetDB()
	page := &models.Page{}
	db.First(page, 1)
	session := sessions.Default(c)
	flashes := session.Flashes()
	session.Save()

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
