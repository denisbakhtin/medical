package system

import (
	"github.com/denisbakhtin/medical/models"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
)

//Authenticated is authentication middleware, enabled by router for protected routes
func Authenticated() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		user := models.User{}
		db := models.GetDB()
		if session.Get("user_id") != nil {
			db.First(&user, session.Get("user_id"))
		}
		if user.ID == 0 {
			c.AbortWithStatus(403)
		}
		c.Next()
	}
}
