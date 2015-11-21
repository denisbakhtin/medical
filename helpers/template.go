package helpers

import (
	"fmt"
	"time"

	"github.com/denisbakhtin/medical/models"
)

type menuItem struct {
	Url      string
	Title    string //will be passed to T i18n function
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

//RecentArticles returns the list of recent articles
func RecentArticles() []models.Article {
	list, _ := models.GetRecentArticles()
	return list
}

//MainMenu returns the list of main menu items
func MainMenu() []menuItem {
	about, _ := models.GetPage(4)
	prices, _ := models.GetPage(5)
	cure, _ := models.GetPage(6)
	contacts, _ := models.GetPage(7)
	menu := []menuItem{
		menuItem{
			Url:   about.Url(),
			Title: "about_doctor",
		},
		menuItem{
			Url:   cure.Url(),
			Title: "what_we_cure",
		},
		menuItem{
			Url:   "/articles",
			Title: "articles",
		},
		menuItem{
			Url:   prices.Url(),
			Title: "prices",
		},
		menuItem{
			Url:   "/reviews",
			Title: "reviews",
		},
		menuItem{
			Url:   contacts.Url(),
			Title: "contacts",
		},
	}
	return menu
}
