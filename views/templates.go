package views

import (
	"embed"
	"html/template"
	"log"
	"path"

	_ "embed"

	"github.com/denisbakhtin/medical/config"
	"github.com/denisbakhtin/medical/helpers"
)

//go:embed html/*
var files embed.FS

var tmpl *template.Template

// Load parses html templates from html directory
func Load() {
	tmpl = template.New("").Funcs(template.FuncMap{
		"isActive":          helpers.IsActive,
		"stringInSlice":     helpers.StringInSlice,
		"dateTime":          helpers.DateTime,
		"date":              helpers.Date,
		"yearNow":           helpers.YearNow,
		"mainMenu":          helpers.MainMenu,
		"scrollMenu":        helpers.ScrollMenu,
		"oddEvenClass":      helpers.OddEvenClass,
		"truncate":          helpers.Truncate,
		"sellingPreface":    helpers.SellingPreface,
		"promoTill":         helpers.PromoTill,
		"replacePromoTill":  helpers.ReplacePromoTill,
		"cityList":          helpers.CityList,
		"eqRU":              helpers.EqRU,
		"allReviews":        helpers.AllReviews,
		"isFirstInTheRow":   helpers.IsFirstInTheRow,
		"isLastInTheRow":    helpers.IsLastInTheRow,
		"isLast":            helpers.IsLast,
		"cssVersion":        helpers.FileVersion(path.Join(config.GetConfig().Public, "css", "application.css")),
		"jsVersion":         helpers.FileVersion(path.Join(config.GetConfig().Public, "js", "application.js")),
		"articleIdComments": helpers.ArticleIDComments,
	})

	var err error
	// parser can't walk subdirectories, so just list all patterns
	tmpl, err = tmpl.ParseFS(files,
		"html/*.gohtml",
		"html/*/*.gohtml",
		"html/*/*/*.gohtml")
	if err != nil {
		log.Panic(err)
	}
}

// GetTemplates exports loaded templates
func GetTemplates() *template.Template {
	return tmpl
}
