package controllers

import (
	"github.com/gin-gonic/gin"
)

//Error404 shows fancy 404 error page
func Error404(c *gin.Context) {
	c.HTML(404, "errors/404", nil)
}
