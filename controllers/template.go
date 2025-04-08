package controllers

import (
	"html/template"
	"path"

	"github.com/denisbakhtin/medical/config"
	"github.com/denisbakhtin/medical/helpers"
	"github.com/denisbakhtin/medical/repos"
	"github.com/denisbakhtin/medical/views"
)

// loadTemplate parses html templates from html directory
// Some functions require database access and public path from config
func loadTemplate(menus repos.Menus, reviews repos.Reviews, config *config.Config) *template.Template {
	tmpl := template.New("").Funcs(template.FuncMap{
		"isActive":          helpers.IsActive,
		"stringInSlice":     helpers.StringInSlice,
		"dateTime":          helpers.DateTime,
		"date":              helpers.Date,
		"yearNow":           helpers.YearNow,
		"mainMenu":          menus.Main,
		"scrollMenu":        menus.Scrolled,
		"oddEvenClass":      helpers.OddEvenClass,
		"truncate":          helpers.Truncate,
		"sellingPreface":    helpers.SellingPreface,
		"promoTill":         helpers.PromoTill,
		"replacePromoTill":  helpers.ReplacePromoTill,
		"cityList":          helpers.CityList,
		"eqRU":              helpers.EqRU,
		"allReviews":        reviews.GetLastPublished,
		"isFirstInTheRow":   helpers.IsFirstInTheRow,
		"isLastInTheRow":    helpers.IsLastInTheRow,
		"isLast":            helpers.IsLast,
		"cssVersion":        helpers.FileVersion(path.Join(config.Public, "css", "application.css")),
		"jsVersion":         helpers.FileVersion(path.Join(config.Public, "js", "application.js")),
		"articleIdComments": helpers.ArticleIDComments,
		"domain":            func() string { return config.Domain },
		"fullDomain":        func() template.HTML { return template.HTML(config.FullDomain) },
	})

	var err error
	// parser can't walk subdirectories, so just list all patterns
	tmpl, err = tmpl.ParseFS(views.TemplateFiles,
		"html/*.gohtml",
		"html/*/*.gohtml",
		"html/*/*/*.gohtml")
	if err != nil {
		panic(err)
	}
	return tmpl
}
