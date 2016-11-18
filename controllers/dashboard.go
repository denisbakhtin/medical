package controllers

import (
	"net/http"

	"github.com/denisbakhtin/medical/helpers"
)

//Dashboard handles GET /admin route
func Dashboard(w http.ResponseWriter, r *http.Request) {
	tmpl := helpers.Template(r)
	data := helpers.DefaultData(r)
	data["Title"] = "Панель управления"
	tmpl.Lookup("dashboard/show").Execute(w, data)
}
