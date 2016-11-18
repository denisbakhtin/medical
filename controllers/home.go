package controllers

import (
	"net/http"

	"github.com/denisbakhtin/medical/helpers"
	"github.com/denisbakhtin/medical/models"
)

//Home handles GET / route
func Home(w http.ResponseWriter, r *http.Request) {
	tmpl, data := helpers.Template(r), helpers.DefaultData(r)
	session := helpers.Session(r)
	T := helpers.T(r)
	db := models.GetDB()
	if r.RequestURI != "/" {
		w.WriteHeader(404)
		tmpl.Lookup("errors/404").Execute(w, nil)
		return
	}
	page := &models.Page{}
	db.First(page, 1)
	data["Title"] = T("site_name_full")
	data["Page"] = page
	data["Active"] = "/"
	data["Flash"] = session.Flashes()
	data["TitleSuffix"] = ""
	data["MetaDescription"] = "Прикладная кинезиология МиоБаланс - восстановление баланса обмена веществ, опорно-двигательного аппарата и нервной системы..."
	session.Save(r, w)
	tmpl.Lookup("home/show").Execute(w, data)
}
