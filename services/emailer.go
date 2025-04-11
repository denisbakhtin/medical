package services

import (
	"strconv"

	"github.com/denisbakhtin/medical/config"
	"gopkg.in/gomail.v2"
)

type Emailer interface {
	NotifyAdmin(replyTo, subject, body string)
}

type gmailer struct {
	config *config.Config
	logger Logger
}

func NewGmailer(config *config.Config, logger Logger) *gmailer {
	return &gmailer{config: config, logger: logger}
}

func (g *gmailer) NotifyAdmin(replyTo, subject, body string) {
	go func() {
		smtp := g.config.SMTP
		msg := gomail.NewMessage()
		msg.SetHeader("From", smtp.From)
		msg.SetHeader("To", smtp.To)
		if len(smtp.Cc) > 0 {
			msg.SetHeader("Cc", smtp.Cc)
		}
		msg.SetHeader("Subject", subject)
		msg.SetBody(
			"text/html",
			body,
		)

		port, _ := strconv.Atoi(smtp.Port)
		dialer := gomail.NewDialer(smtp.SMTP, port, smtp.User, smtp.Password)
		sender, err := dialer.Dial()
		if err != nil {
			g.logger.Error(err)
			return
		}
		if err := gomail.Send(sender, msg); err != nil {
			g.logger.Error(err)
			return
		}
	}()
}
