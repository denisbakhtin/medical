package models

import (
	"fmt"
	"html/template"
	"regexp"
	"strings"
	"time"
)

//Article type contains article info
type Article struct {
	ID              uint
	Name            string
	Slug            string
	Excerpt         string
	Content         string
	SellingPreface  string
	MetaKeywords    string
	MetaDescription string
	Published       bool
	CreatedAt       time.Time
	UpdatedAt       time.Time
	//calculated
	Comments []Comment
}

//HTMLContent returns parsed html content
func (article *Article) HTMLContent() template.HTML {
	return template.HTML(article.Content)
}

//GetCommentCount returns comment count
func (article *Article) GetCommentCount() int {
	count := 0
	db.Where("article_id = ? AND answer <> ?", article.ID, "").Model(&Comment{}).Count(&count)
	return count
}

//GetImage returns extracts first image url from article content
func (article *Article) GetImage() string {
	re := regexp.MustCompile(`<img[^<>]+src="([^"]+)"[^<>]*>`)
	res := re.FindStringSubmatch(article.Content)
	if len(res) == 2 {
		return res[1]
	}
	return ""
}

//URL returns article url
func (article *Article) URL() string {
	return fmt.Sprintf("/articles/%d-%s", article.ID, article.Slug)
}

func (article *Article) BeforeCreate() (err error) {
	if strings.TrimSpace(article.Slug) == "" {
		article.Slug = createSlug(article.Name)
	}
	return
}

func (article *Article) BeforeSave() (err error) {
	if strings.TrimSpace(article.Slug) == "" {
		article.Slug = createSlug(article.Name)
	}
	return
}
