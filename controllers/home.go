package controllers

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

const homePageId = 1

// Home handles GET / route
func (app *Application) Home(c *gin.Context) {
	page, err := app.PagesRepo.Get(homePageId)
	if err != nil {
		app.Error(c, err)
		return
	}
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
		"Authenticated":   app.authenticated(session),
	})
}
