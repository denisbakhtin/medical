package system

import (
	"html/template"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/denisbakhtin/medical/helpers"
)

var tmpl *template.Template

func loadTemplates() {
	tmpl = template.New("").Funcs(template.FuncMap{
		"isActive":          helpers.IsActive,
		"stringInSlice":     helpers.StringInSlice,
		"dateTime":          helpers.DateTime,
		"date":              helpers.Date,
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
		"cssVersion":        helpers.CSSVersion(path.Join(GetConfig().Public, "css", "application.css")),
		"articleIdComments": helpers.ArticleIDComments,
	})

	fn := func(path string, f os.FileInfo, err error) error {
		if f.IsDir() != true && strings.HasSuffix(f.Name(), ".gohtml") {
			var err error
			tmpl, err = tmpl.ParseFiles(path)
			if err != nil {
				return err
			}
		}
		return nil
	}

	err := filepath.Walk("views", fn)
	if err != nil {
		log.Panic(err)
	}
}

//GetTemplates exports loaded templates
func GetTemplates() *template.Template {
	return tmpl
}
