package controllers

import (
	"errors"

	"github.com/denisbakhtin/medical/helpers"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

// Error404 shows fancy 404 error page, used directly mainly by gin.Engine.NoRoute default 404 handler
// In all other cases prefer Error(...) handler
func (app *Application) Error404(c *gin.Context) {
	c.HTML(404, "errors/404", nil)
}

// Error is an http error handler, it checks error type and returns an according http error page with http code
func (app *Application) Error(c *gin.Context, err error) {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		c.HTML(404, "errors/404", nil)
	} else {
		app.Logger.Errorf("Http code: 500, Reason: %v\n", err)
		c.HTML(500, "errors/500", helpers.ErrorData(err))
	}
}
