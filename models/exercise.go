package models

import (
	"fmt"
	"html/template"
	"regexp"
	"strings"
	"time"
)

//Exercise type contains exercise info
type Exercise struct {
	ID              uint   `form:"id"`
	Name            string `form:"name"`
	Slug            string `form:"slug"`
	Content         string `form:"content"`
	Image           string `form:"image"`
	Video           string `form:"video"`
	MetaKeywords    string `form:"meta_keywords"`
	MetaDescription string `form:"meta_description"`
	Published       bool   `form:"published"`
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

//HTMLContent returns parsed html content
func (e *Exercise) HTMLContent() template.HTML {
	return template.HTML(e.Content)
}

//GetImage returns extracts first image url from article content
func (e *Exercise) GetImage() string {
	re := regexp.MustCompile(`<img[^<>]+src="([^"]+)"[^<>]*>`)
	res := re.FindStringSubmatch(e.Content)
	if len(res) == 2 {
		return res[1]
	}
	return ""
}

//URL returns article url
func (e *Exercise) URL() string {
	return fmt.Sprintf("/exercises/%d-%s", e.ID, e.Slug)
}

//BeforeCreate gorm hook
func (e *Exercise) BeforeCreate() (err error) {
	if strings.TrimSpace(e.Slug) == "" {
		e.Slug = createSlug(e.Name)
	}
	return
}

//BeforeSave gorm hook
func (e *Exercise) BeforeSave() (err error) {
	if strings.TrimSpace(e.Slug) == "" {
		e.Slug = createSlug(e.Name)
	}
	return
}
