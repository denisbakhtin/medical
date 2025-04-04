package models

import "time"

// Comment type contains article comments
type Comment struct {
	ID          uint      `json:"id" form:"id"`
	ArticleID   uint      `json:"article_id" form:"article_id"`
	AuthorCity  string    `json:"author_city" form:"author_city"`
	AuthorName  string    `json:"author_name" form:"author_name"`
	AuthorEmail string    `json:"author_email" form:"author_email"`
	Content     string    `json:"content" form:"content"`
	Answer      string    `json:"answer" form:"answer"`
	Published   bool      `json:"published" form:"published"`
	Captcha     string    `json:"captcha" form:"captcha" gorm:"-" db:"-"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Title returns comment excerpt, 100 char long
func (comment *Comment) Title() string {
	return truncate(comment.Content, 50) + "..."
}

// Excerpt returns comment excerpt, 100 char long
func (comment *Comment) Excerpt() string {
	return truncate(comment.Content, 20)
}
