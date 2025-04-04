package controllers

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// Authenticated is authentication middleware, enabled by router for protected routes
func (app *Application) FilterAuthenticated() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		if session.Get("user_id") != nil {
			if id, ok := session.Get("user_id").(uint); ok {
				user, err := app.UsersRepo.Get(id)
				if err == nil {
					c.Set("user", *user)
					c.Next()
					return
				}
			}
		}

		c.AbortWithStatus(403)
	}
}
