package models

import (
	"fmt"
	"html/template"
	"strings"
	"time"
)

//Page type contains page info
type Page struct {
	ID              uint   `json:"id"`
	Name            string `json:"name"`
	Slug            string `json:"slug"`
	Content         string `json:"content"`
	MetaKeywords    string
	MetaDescription string
	Published       bool      `json:"published"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
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
