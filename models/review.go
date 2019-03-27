package models

import (
	"fmt"
	"html/template"
	"time"

	"github.com/microcosm-cc/bluemonday"
)

//Review type contains client review
type Review struct {
	ID              uint   `json:"id" form:"id"`
	ArticleID       *uint  `json:"article_id"`
	AuthorName      string `json:"author_name" form:"author_name"`
	AuthorEmail     string `json:"author_email" form:"author_email"`
	Content         string `json:"content" form:"content"`
	Image           string `json:"image"`
	Video           string `json:"video" form:"video"`
	MetaKeywords    string `form:"meta_keywords"`
	MetaDescription string `form:"meta_description"`
	Published       bool   `json:"published" form:"published"`
	Captcha         string `form:"captcha" gorm:"-" db:"-"`
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

//Excerpt returns review excerpt, 100 char long
func (review *Review) Excerpt() string {
	policy := bluemonday.StrictPolicy()
	return truncate(policy.Sanitize(review.Content), 300)
}

//URL returns review url
func (review *Review) URL() string {
	return fmt.Sprintf("/reviews/%d", review.ID)
}

//HTMLContent returns parsed html content
func (review *Review) HTMLContent() template.HTML {
	return template.HTML(review.Content)
}
