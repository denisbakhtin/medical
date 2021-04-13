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
	ID              uint   `form:"id"`
	Name            string `form:"name"`
	Slug            string `form:"slug"`
	Excerpt         string `form:"excerpt"`
	Content         string `form:"content"`
	SellingBlock    string `form:"selling_block"`
	MetaKeywords    string `form:"meta_keywords"`
	MetaDescription string `form:"meta_description"`
	Published       bool   `form:"published"`
	CreatedAt       time.Time
	UpdatedAt       time.Time
	//calculated
	Comments []Comment
}

//HTMLContent returns parsed html content
func (article *Article) HTMLContent() template.HTML {
	return template.HTML(article.Content)
}

//HTMLSellingBlock returns parsed html selling block
func (article *Article) HTMLSellingBlock() template.HTML {
	return template.HTML(article.SellingBlock)
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

//BeforeCreate gorm hook
func (article *Article) BeforeCreate() (err error) {
	if strings.TrimSpace(article.Slug) == "" {
		article.Slug = createSlug(article.Name)
	}
	return
}

//BeforeSave gorm hook
func (article *Article) BeforeSave() (err error) {
	if strings.TrimSpace(article.Slug) == "" {
		article.Slug = createSlug(article.Name)
	}
	return
}
