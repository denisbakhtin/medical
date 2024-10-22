package models

import (
	"fmt"
	"html/template"
	"regexp"
	"strings"
	"time"
)

// Info type represents info article
type Info struct {
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

// HTMLContent returns parsed html content
func (info *Info) HTMLContent() template.HTML {
	return template.HTML(info.Content)
}

// GetImage returns extracts first image url from article content
func (info *Info) GetImage() string {
	re := regexp.MustCompile(`<img[^<>]+src="([^"]+)"[^<>]*>`)
	res := re.FindStringSubmatch(info.Content)
	if len(res) == 2 {
		return res[1]
	}
	return ""
}

// URL returns article url
func (info *Info) URL() string {
	return fmt.Sprintf("/info/%d-%s", info.ID, info.Slug)
}

// BeforeCreate gorm hook
func (info *Info) BeforeCreate() (err error) {
	if strings.TrimSpace(info.Slug) == "" {
		info.Slug = createSlug(info.Name)
	}
	return
}

// BeforeSave gorm hook
func (info *Info) BeforeSave() (err error) {
	if strings.TrimSpace(info.Slug) == "" {
		info.Slug = createSlug(info.Name)
	}
	return
}
