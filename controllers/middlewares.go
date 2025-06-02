package controllers

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

const userIDkey = "user_id"
const userKey = "user"

// FilterAuthenticated is authentication middleware, that blocks unauthenticated requests to protected urls
func (app *Application) FilterAuthenticated() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		if session.Get(userIDkey) != nil {
			if id, ok := session.Get(userIDkey).(uint); ok {
				user, err := app.UsersRepo.Get(id)
				if err == nil {
					c.Set(userKey, *user)
					c.Next()
					return
				}
			}
		}

		c.AbortWithStatus(403)
	}
}
