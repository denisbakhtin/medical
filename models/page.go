package models

import (
	"fmt"
	"html/template"
	"strings"
	"time"
)

//Page type contains page info
type Page struct {
	ID              uint   `form:"id"`
	Name            string `form:"name"`
	Slug            string `form:"slug"`
	Content         string `form:"content"`
	MetaKeywords    string `form:"meta_keywords"`
	MetaDescription string `form:"meta_description"`
	Published       bool   `form:"published"`
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

//HTMLContent returns parsed html content
func (page *Page) HTMLContent() template.HTML {
	return template.HTML(page.Content)
}

//URL returns page url
func (page *Page) URL() string {
	return fmt.Sprintf("/pages/%d-%s", page.ID, page.Slug)
}

func (page *Page) BeforeCreate() (err error) {
	if strings.TrimSpace(page.Slug) == "" {
		page.Slug = createSlug(page.Name)
	}
	return
}

func (page *Page) BeforeSave() (err error) {
	if strings.TrimSpace(page.Slug) == "" {
		page.Slug = createSlug(page.Name)
	}
	return
}
