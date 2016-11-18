package helpers

import (
	"html/template"
	"net/http"
	"strconv"

	"github.com/denisbakhtin/medical/models"
	"github.com/gorilla/context"
	"github.com/gorilla/csrf"
	"github.com/gorilla/sessions"
	"github.com/nicksnyder/go-i18n/i18n"
)

//DefaultData returns common to all pages template data
func DefaultData(r *http.Request) map[string]interface{} {
	var testimonials []models.Review
	models.GetDB().Where("published = ?", true).Order("id desc").Limit(7).Find(&testimonials)
	t := T(r)
	return map[string]interface{}{
		"ActiveUser":      context.Get(r, "user"), //signed in models.User
		"Active":          "",                     //active uri shortening for menu item highlight
		"Title":           "",                     //page title
		"TitleSuffix":     t("title_suffix"),
		"MetaDescription": "",
		"SignupEnabled":   context.Get(r, "signup_enabled"), //signup route is enabled (otherwise everyone can signup ;)
		"Testimonials":    testimonials,
		csrf.TemplateTag:  csrf.TemplateField(r),
	}
}

//ErrorData returns template data for error
func ErrorData(err error) map[string]interface{} {
	return map[string]interface{}{
		"Title": err.Error(),
		"Error": err.Error(),
	}
}

//T returns i18n.TranslateFunc for current locale
func T(r *http.Request) i18n.TranslateFunc {
	return context.Get(r, "T").(i18n.TranslateFunc)
}

//Template returns parsed *html/template.Template
func Template(r *http.Request) *template.Template {
	return context.Get(r, "template").(*template.Template)
}

//Session returns current session
func Session(r *http.Request) *sessions.Session {
	return context.Get(r, "session").(*sessions.Session)
}

//Atoi64 converts string to int64, returns 0 if error
func Atoi64(s string) int64 {
	i, _ := strconv.ParseInt(s, 10, 64)
	return i
}

//Atouint converts string to uint, returns 0 if error
func Atouint(s string) uint {
	i, _ := strconv.ParseUint(s, 10, 32)
	return uint(i)
}

//Atoi64r converts string to int64 reference
func Atoi64r(s string) *int64 {
	i, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return nil
	}
	return &i
}

//Atouintr converts string to uint reference
func Atouintr(s string) *uint {
	i, err := strconv.ParseUint(s, 10, 32)
	if err != nil {
		return nil
	}
	ui := uint(i)
	return &ui
}

//Atob converts string to bool
func Atob(s string) bool {
	b, _ := strconv.ParseBool(s)
	return b
}
