package controllers

import (
	"html/template"
	"path"

	"github.com/denisbakhtin/medical/helpers"
	"github.com/denisbakhtin/medical/views"
)

// loadTemplate parses html templates from html directory
// Some functions require database access and public path from config
func (app *Application) loadTemplate() *template.Template {
	tmpl := template.New("").Funcs(template.FuncMap{
		"isActive":          helpers.IsActive,
		"stringInSlice":     helpers.StringInSlice,
		"dateTime":          helpers.DateTime,
		"date":              helpers.Date,
		"yearNow":           helpers.YearNow,
		"mainMenu":          app.MenusRepo.Main,
		"scrollMenu":        app.MenusRepo.Scrolled,
		"oddEvenClass":      helpers.OddEvenClass,
		"truncate":          helpers.Truncate,
		"sellingPreface":    helpers.SellingPreface,
		"promoTill":         helpers.PromoTill,
		"replacePromoTill":  helpers.ReplacePromoTill,
		"cityList":          helpers.CityList,
		"eqRU":              helpers.EqRU,
		"allReviews":        app.ReviewsRepo.GetLastPublished,
		"isFirstInTheRow":   helpers.IsFirstInTheRow,
		"isLastInTheRow":    helpers.IsLastInTheRow,
		"isLast":            helpers.IsLast,
		"cssVersion":        helpers.FileVersion(path.Join(app.Config.Public, "css", "application.css")),
		"jsVersion":         helpers.FileVersion(path.Join(app.Config.Public, "js", "application.js")),
		"articleIdComments": helpers.ArticleIDComments,
		"domain":            func() string { return app.Config.Domain },
		"fullDomain":        func() template.HTML { return template.HTML(app.Config.FullDomain) },
	})

	var err error
	// parser can't walk subdirectories, so just list all patterns
	tmpl, err = tmpl.ParseFS(views.TemplateFiles,
		"html/*.gohtml",
		"html/*/*.gohtml",
		"html/*/*/*.gohtml")
	if err != nil {
		app.Logger.Fatal(err)
	}
	return tmpl
}
