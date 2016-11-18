package helpers

import (
	"fmt"
	"time"

	"github.com/denisbakhtin/medical/models"
)

//MenuItem represents main menu item
type MenuItem struct {
	URL      string
	Title    string //will be passed to T i18n function
	CssClass string
	IsActive bool
}

//IsActive checks uri against currently active (uri, or nil) and returns "active" if they are equal
func IsActive(active interface{}, uri string) string {
	if s, ok := active.(string); ok {
		if s == uri {
			return "active"
		}
	}
	return ""
}

//DateTime prints timestamp in human format
func DateTime(t time.Time) string {
	return fmt.Sprintf("%d-%02d-%02d %02d:%02d:%02d", t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second())
}

//Date prints date part of timestamp
func Date(t time.Time) string {
	return fmt.Sprintf("%d-%02d-%02d", t.Year(), t.Month(), t.Day())
}

//StringInSlice returns true if value is in list slice
func StringInSlice(value string, list []string) bool {
	for i := range list {
		if value == list[i] {
			return true
		}
	}
	return false
}

//OddEvenClass returns odd or even class depending on the index
func OddEvenClass(index int) string {
	//range indexes start with zero %)
	if (index+1)%2 == 1 {
		return "odd"
	}
	return "even"
}

//MainMenu returns the list of main menu items
func MainMenu() []MenuItem {
	db := models.GetDB()
	about := &models.Page{}
	db.First(about, 4)
	contacts := &models.Page{}
	db.First(contacts, 7)
	menu := []MenuItem{
		MenuItem{
			URL:   "/reviews",
			Title: "reviews",
		},
		MenuItem{
			URL:   "/articles",
			Title: "articles",
		},
		MenuItem{
			URL:      about.URL(),
			Title:    "about_doctor",
			CssClass: "small",
		},
		MenuItem{
			URL:      contacts.URL(),
			Title:    "contacts",
			CssClass: "small",
		},
	}
	return menu
}

//ScrollMenu returns the list of scroll menu items
func ScrollMenu() []MenuItem {
	db := models.GetDB()
	about := &models.Page{}
	db.First(about, 4)
	menu := []MenuItem{
		MenuItem{
			URL:   about.URL(),
			Title: "about_doctor",
		},
		MenuItem{
			URL:   "#withoutpain",
			Title: "treatment_stages",
		},
		MenuItem{
			URL:   "/reviews",
			Title: "reviews",
		},
	}
	return menu
}

//Truncate truncates string to n chars
func Truncate(s string, n int) string {
	runes := []rune(s)
	if len(runes) > n {
		return string(runes[:n]) + "..."
	}
	return s
}

func SellingPreface() string {
	return "Выяснить причины возникновения жалоб и пройти кинезиологическое тестирование можно во время:"
}

//PromoTill returns promotion text
func PromoTill() string {
	now := time.Now()
	wday := now.Weekday()
	endofweek := now.Add(time.Duration(6-int(wday)) * 24 * time.Hour)
	return fmt.Sprintf("до %d %s", endofweek.Day(), mon(endofweek.Month()))
}

//CityList returns the list of cities for comments form
func CityList() []string {
	return []string{
		"Москва",
		"Санкт-Петербург",
		"Московская обл.",
		"Новосибирск",
		"Екатеринбург",
		"Нижний Новгород",
		"Казань",
		"Самара",
		"Челябинск",
		"Омск",
		"Ростов-на-Дону",
		"Волгоград",
		"Воронеж",
		"Не указан в списке",
	}
}

//EqRU compares *uint to uint
func EqRU(r *uint, i uint) bool {
	if r == nil {
		if i == 0 {
			return true
		}
		return false
	}
	return (*r == i)
}

func mon(m time.Month) string {
	switch m {
	case 1:
		return "января"
	case 2:
		return "февраля"
	case 3:
		return "марта"
	case 4:
		return "апреля"
	case 5:
		return "мая"
	case 6:
		return "июня"
	case 7:
		return "июля"
	case 8:
		return "августа"
	case 9:
		return "сентября"
	case 10:
		return "октября"
	case 11:
		return "ноября"
	case 12:
		return "декабря"
	default:
		return ""
	}
}
