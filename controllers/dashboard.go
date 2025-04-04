package controllers

import "github.com/gin-gonic/gin"

// Dashboard handles GET /admin route
func (ap *Application) Dashboard(c *gin.Context) {
	c.HTML(200, "dashboard/show", gin.H{
		"Title":         "Панель управления",
		"Authenticated": true,
	})
}
