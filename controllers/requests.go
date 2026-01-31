package controllers

import (
	"bytes"
	"fmt"

	"github.com/denisbakhtin/medical/helpers"
	"github.com/denisbakhtin/medical/models"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// RequestCreatePost handles /new_request route
func (app *Application) RequestCreatePost(c *gin.Context) {
	session := sessions.Default(c)

	request := &models.Request{}
	if c.Bind(request) == nil {
		if err := request.Valid(); err != nil {
			c.HTML(400, "errors/500", helpers.ErrorData(err))
			return
		}
		app.notifyAdminOfRequest(request)
		session.AddFlash("Спасибо, что оставили заявку на приём. В ближайшее время наш специалист свяжется с Вами по указанному телефону и согласует детали")
	} else {
		session.AddFlash("Ошибка! Проверьте внимательно заполнение всех полей!")
	}
	_ = session.Save()
	c.Redirect(303, "/")
}

// notifyAdminOfRequest sends a notification email to admin
func (app *Application) notifyAdminOfRequest(request *models.Request) {
	data := map[string]any{
		"Request": request,
	}
	var b bytes.Buffer
	if err := app.Template.Lookup("requests/request").Execute(&b, data); err != nil {
		app.Logger.Error(err)
		return
	}

	subject := fmt.Sprintf("Заявка на приём %s: %s", app.Config.Domain, request.Name)
	app.Emailer.NotifyAdmin("", subject, b.String())
}
