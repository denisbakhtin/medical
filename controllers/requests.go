package controllers

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/denisbakhtin/medical/helpers"
	"github.com/denisbakhtin/medical/models"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// RequestCreatePost handles /new_request route
func (ap *Application) RequestCreatePost(c *gin.Context) {
	session := sessions.Default(c)

	request := &models.Request{}
	if c.Bind(request) == nil {
		captcha, err := base64.StdEncoding.DecodeString(request.Captcha)
		if err != nil {
			c.HTML(500, "errors/500", helpers.ErrorData(err))
			return
		}
		if string(captcha) != "100.00" {
			c.HTML(400, "errors/400", nil)
			return
		}
		if !strings.Contains(strings.ToLower(request.Comment), "href") {
			ap.notifyAdminOfRequest(request)
		}
		session.AddFlash("Спасибо, что оставили заявку на приём. В ближайшее время наш специалист свяжется с Вами по указанному телефону и согласует детали")
	} else {
		session.AddFlash("Ошибка! Проверьте внимательно заполнение всех полей!")
	}
	_ = session.Save()
	c.Redirect(303, "/")
}

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
