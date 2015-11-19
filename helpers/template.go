package helpers

import (
	"fmt"
	"time"

	"github.com/denisbakhtin/medical/models"
)

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

//StringInSlice returns true if value is in list slice
func StringInSlice(value string, list []string) bool {
	for i := range list {
		if value == list[i] {
			return true
		}
	}
	return false
}

//RecentPosts returns the list of recent blog posts
func RecentPosts() []models.Post {
	list, _ := models.GetRecentPosts()
	return list
}

//Tags returns the list of blog tags
func Tags() []models.Tag {
	list, _ := models.GetNotEmptyTags()
	return list
}

//Archives returns the list of blog archives
func Archives() []models.Post {
	list, _ := models.GetPostMonths()
	return list
}
