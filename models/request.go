package models

import (
	"encoding/base64"
	"errors"
	"strings"
	"unicode"
)

// Request appointment view model
type Request struct {
	Name    string `form:"name" binding:"required"`
	Phone   string `form:"phone" binding:"required"`
	Comment string `form:"comment" binding:"required"`
	Captcha string `json:"captcha" form:"captcha" gorm:"-" db:"-"`
}

// Valid checks that user form is valid
func (r *Request) Valid() error {
	captcha, err := base64.StdEncoding.DecodeString(r.Captcha)
	if err != nil {
		return err
	}
	if string(captcha) != "100.00" {
		return errors.New("Неверная captcha")
	}
	// simple spam protection
	if strings.Contains(strings.ToLower(r.Comment), "href") {
		return errors.New("Ссылки недопустимы")
	}
	return nil
}

func (r *Request) Tel() string {
	digits := func(d rune) rune {
		if unicode.IsDigit(d) {
			return d
		}
		return -1 //remove
	}
	phone := strings.Map(digits, r.Phone)
	if len(phone) > 0 && phone[0] == '8' {
		phone = "+7" + phone[1:]
	}
	// emails block <a href="tel:..."></a> links or encrypt them, so use plain text
	return phone
}
