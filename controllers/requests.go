package controllers

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/denisbakhtin/medical/helpers"
	"github.com/denisbakhtin/medical/models"
	"github.com/denisbakhtin/medical/system"
	"gopkg.in/gomail.v2"
)

//RequestCreate handles /new_request route
func RequestCreate(w http.ResponseWriter, r *http.Request) {
	session := helpers.Session(r)
	tmpl := helpers.Template(r)
	if r.Method == "POST" {

		r.ParseForm()

		request := &models.Request{
			Name:    r.PostFormValue("name"),
			Phone:   r.PostFormValue("phone"),
			Comment: r.PostFormValue("comment"),
		}

		notifyAdminOfRequest(r, request)
		session.AddFlash("Спасибо, что оставили заявку на приём. В ближайшее время наш специалист свяжется с Вами по указанному телефону и согласует детали")
		session.Save(r, w)
		http.Redirect(w, r, "/", 303)

	} else {
		err := fmt.Errorf("Method %q not allowed", r.Method)
		log.Printf("ERROR: %s\n", err)
		w.WriteHeader(405)
		tmpl.Lookup("errors/405").Execute(w, helpers.ErrorData(err))
	}
}

func notifyAdminOfRequest(r *http.Request, request *models.Request) {
	//closure is needed here, as r may be released by the time func finishes
	tmpl := helpers.Template(r)
	go func() {
		data := map[string]interface{}{
			"Request": request,
		}
		var b bytes.Buffer
		if err := tmpl.Lookup("requests/request").Execute(&b, data); err != nil {
			log.Printf("ERROR: %s\n", err)
			return
		}

		smtp := system.GetConfig().SMTP
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
