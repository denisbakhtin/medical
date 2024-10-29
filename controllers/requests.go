package controllers

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/denisbakhtin/medical/config"
	"github.com/denisbakhtin/medical/helpers"
	"github.com/denisbakhtin/medical/models"
	"github.com/denisbakhtin/medical/views"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"gopkg.in/gomail.v2"
)

// RequestCreatePost handles /new_request route
func RequestCreatePost(c *gin.Context) {
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
			notifyAdminOfRequest(request)
		}
		session.AddFlash("Спасибо, что оставили заявку на приём. В ближайшее время наш специалист свяжется с Вами по указанному телефону и согласует детали")
	} else {
		session.AddFlash("Ошибка! Проверьте внимательно заполнение всех полей!")
	}
	session.Save()
	c.Redirect(303, "/")
}

func notifyAdminOfRequest(request *models.Request) {
	// closure is needed here, as r may be released by the time func finishes
	tmpl := views.GetTemplates()
	go func() {
		data := map[string]interface{}{
			"Request": request,
		}
		var b bytes.Buffer
		if err := tmpl.Lookup("requests/request").Execute(&b, data); err != nil {
			log.Printf("ERROR: %s\n", err)
			return
		}

		smtp := config.GetConfig().SMTP
		msg := gomail.NewMessage()
		msg.SetHeader("From", smtp.From)
		msg.SetHeader("To", smtp.To)
		if len(smtp.Cc) > 0 {
			msg.SetHeader("Cc", smtp.Cc)
		}
		msg.SetHeader("Subject", fmt.Sprintf("Заявка на приём www.miobalans.ru: %s", request.Name))
		msg.SetBody(
			"text/html",
			b.String(),
		)

		port, _ := strconv.Atoi(smtp.Port)
		dialer := gomail.NewPlainDialer(smtp.SMTP, port, smtp.User, smtp.Password)
		sender, err := dialer.Dial()
		if err != nil {
			log.Printf("ERROR: %s\n", err)
			return
		}
		if err := gomail.Send(sender, msg); err != nil {
			log.Printf("ERROR: %s\n", err)
			return
		}
	}()
}
