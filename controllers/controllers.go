package controllers

import (
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/denisbakhtin/medical/models"
	"github.com/denisbakhtin/medical/system"
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
		c.Set("user", user)
		c.Next()
	}
}

//createTokenFromId creates secure token for id
func createTokenFromID(ID uint) string {
	digest := sha1.New().Sum([]byte(fmt.Sprintf("%d-%s", ID, system.GetConfig().Salt)))
	return base64.URLEncoding.EncodeToString(
		[]byte(fmt.Sprintf("%d-%x", ID, digest)),
	)
}

//getIDFromToken deciphers token and returns review ID. Returns empty string if error
func getIDFromToken(token string) string {
	idDigest, err := base64.URLEncoding.DecodeString(token)
	if err != nil {
		return ""
	}
	if sl := strings.Split(string(idDigest), "-"); len(sl) == 2 {
		digest := sha1.New().Sum([]byte(fmt.Sprintf("%s-%s", sl[0], system.GetConfig().Salt)))
		if fmt.Sprintf("%x", digest) == sl[1] {
			return sl[0]
		}
	}
	return ""
}
