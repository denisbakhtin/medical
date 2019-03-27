package helpers

import (
	"fmt"
	"html/template"
	"math"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/denisbakhtin/medical/models"
	"github.com/gin-gonic/gin"
)

//MenuItem represents main menu item
type MenuItem struct {
	URL      string
	Title    string
	CSSClass string
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

//YearNow returns current year
func YearNow() string {
	return fmt.Sprintf("%d", time.Now().Year())
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

//IsFirstInTheRow checks if index element is the first one in the row of inrow elements
func IsFirstInTheRow(index int, inrow int) bool {
	if inrow == 0 || index%inrow == 0 {
		return true
	}
	return false
}

//IsLastInTheRow checks if index element is the last one in the row of inrow elements, or the last in sequence
func IsLastInTheRow(index int, inrow int, length int) bool {
	if inrow == 0 || index%inrow == inrow-1 || index == length-1 {
		return true
	}
	return false
}

//IsLast checks if index element is the last in sequence
func IsLast(index int, length int) bool {
	return index == length-1
}

//MainMenu returns the list of main menu items
func MainMenu() []MenuItem {
	db := models.GetDB()
	about := &models.Page{}
	db.First(about, 4)
	contacts := &models.Page{}
	db.First(contacts, 7)
	seans := &models.Page{}
	db.First(seans, 10)
	menu := []MenuItem{
		MenuItem{
			URL:   "/reviews",
			Title: "Отзывы",
		},
		MenuItem{
			URL:   seans.URL(),
			Title: "Приём",
		},
		MenuItem{
			URL:   "/articles",
			Title: "Лечение",
		},
		MenuItem{
			URL:   about.URL(),
			Title: "Врач кинезиолог",
		},
		MenuItem{
			URL:   contacts.URL(),
			Title: "Контакты",
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
			Title: "О враче",
		},
		MenuItem{
			URL:   "#withoutpain",
			Title: "Этапы лечения",
		},
		MenuItem{
			URL:   "/reviews",
			Title: "Отзывы",
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

//SellingPreface is the beginning of the selling block partial
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

//ReplacePromoTill replaces {{promoTill}} placeholder with PromoTill() function result
func ReplacePromoTill(source template.HTML) template.HTML {
	return template.HTML(strings.Replace(string(source), "{{promoTill}}", PromoTill(), -1))
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

//AllReviews returns a slice of all published reviews
func AllReviews() (reviews []models.Review) {
	models.GetDB().Where("published = ?", true).Order("id desc").Find(&reviews)
	return
}

//CSSVersion - a closure returning css version function
func CSSVersion(path string) func() string {
	return func() string {
		file, err := os.Stat(path)
		if err != nil {
			return timeToString(time.Now())
		}
		modified := file.ModTime()
		return timeToString(modified)
	}
}

func timeToString(t time.Time) string {
	return fmt.Sprintf("%04d%02d%02d-%02d%02d", t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute())
}

//ArticleIDComments retrieves article id from a list of comments
func ArticleIDComments(comments []models.Comment) uint {
	if len(comments) == 0 {
		return 0
	}
	return comments[0].ArticleID
}

//Min returns int minimum of a & b
func Min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

//Max returns int maximum of a & b
func Max(a, b int) int {
	if a < b {
		return b
	}
	return a
}

//Pagination stores pagination element
type Pagination struct {
	Class string
	URL   string
	Title string
	Rel   string
}

func pageQuery(query *url.Values, page int) string {
	if page > 1 {
		query.Set("page", fmt.Sprintf("%d", page))
	} else {
		query.Del("page")
	}
	return query.Encode()
}

//CurrentPage retrieves page number query parameter
func CurrentPage(c *gin.Context) int {
	currentPage := 1
	if pageStr := c.Query("page"); pageStr != "" {
		currentPage, _ = strconv.Atoi(pageStr)
	}
	currentPage = int(math.Max(float64(1), float64(currentPage)))
	return currentPage
}

//Paginator creates paginator
func Paginator(currentPage, totalPages int, curURL *url.URL) []Pagination {
	currentPage = int(Max(1, Min(int(currentPage), int(totalPages))))
	queryValues, err := url.ParseQuery(curURL.RawQuery)
	if err != nil {
		queryValues = url.Values{}
	}
	nextID := 5
	lastID := 6
	if totalPages < 3 {
		nextID = 4
		lastID = 5
	}
	//first + last + prev + next + 3 adjusent == 7
	pagination := make([]Pagination, 7)

	if totalPages > 1 {
		//prev links
		if currentPage > 1 {
			newURL := *curURL
			newURL.RawQuery = pageQuery(&queryValues, 1)
			pagination[0] = Pagination{Class: "first_page", URL: newURL.RequestURI(), Title: "Первая"}
			newURL.RawQuery = pageQuery(&queryValues, currentPage-1)
			pagination[1] = Pagination{Class: "previous_page", URL: newURL.RequestURI(), Title: "Пред.", Rel: "prev"}
		} else {
			pagination[0] = Pagination{Class: "first_page disabled", URL: "", Title: "Первая"}
			pagination[1] = Pagination{Class: "previous_page disabled", URL: "", Title: "Пред."}
		}

		//page numbers
		switch currentPage {
		case 1:
			pagination[2] = Pagination{Class: "active", URL: "", Title: "1"}
			if 2 <= totalPages {
				newURL := *curURL
				newURL.RawQuery = pageQuery(&queryValues, 2)
				pagination[3] = Pagination{Class: "", URL: newURL.RequestURI(), Title: "2"}
			}
			if 3 <= totalPages {
				newURL := *curURL
				newURL.RawQuery = pageQuery(&queryValues, 3)
				pagination[4] = Pagination{Class: "", URL: newURL.RequestURI(), Title: "3"}
			}
		case totalPages:
			if 3 <= totalPages {
				pagination[4] = Pagination{Class: "active", URL: "", Title: fmt.Sprintf("%d", totalPages)}
				newURL := *curURL
				newURL.RawQuery = pageQuery(&queryValues, totalPages-1)
				pagination[3] = Pagination{Class: "", URL: newURL.RequestURI(), Title: fmt.Sprintf("%d", totalPages-1)}

				newURL.RawQuery = pageQuery(&queryValues, totalPages-2)
				pagination[2] = Pagination{Class: "", URL: newURL.RequestURI(), Title: fmt.Sprintf("%d", totalPages-2)}
			}
			if 2 == totalPages {
				pagination[3] = Pagination{Class: "active", URL: "", Title: fmt.Sprintf("%d", totalPages)}

				newURL := *curURL
				newURL.RawQuery = pageQuery(&queryValues, totalPages-1)
				pagination[2] = Pagination{Class: "", URL: newURL.RequestURI(), Title: fmt.Sprintf("%d", totalPages-1)}
			}

		default:
			newURL := *curURL
			newURL.RawQuery = pageQuery(&queryValues, currentPage+1)
			pagination[4] = Pagination{Class: "", URL: newURL.RequestURI(), Title: fmt.Sprintf("%d", currentPage+1)}
			pagination[3] = Pagination{Class: "active", URL: "", Title: fmt.Sprintf("%d", currentPage)}
			newURL.RawQuery = pageQuery(&queryValues, currentPage-1)
			pagination[2] = Pagination{Class: "", URL: newURL.RequestURI(), Title: fmt.Sprintf("%d", currentPage-1)}
		}

		//next links
		if currentPage < totalPages {
			newURL := *curURL
			newURL.RawQuery = pageQuery(&queryValues, currentPage+1)
			pagination[nextID] = Pagination{Class: "next_page", URL: newURL.RequestURI(), Title: "След.", Rel: "next"}

			newURL.RawQuery = pageQuery(&queryValues, totalPages)
			pagination[lastID] = Pagination{Class: "last_page", URL: newURL.RequestURI(), Title: "Последняя"}

		} else {
			pagination[lastID] = Pagination{Class: "last_page disabled", URL: "", Title: "Последняя"}
			pagination[nextID] = Pagination{Class: "next_page disabled", URL: "", Title: "След."}
		}
		return pagination
	}
	return nil
}
