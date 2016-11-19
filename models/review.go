package models

import (
	"fmt"
	"time"
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
	return truncate(review.Content, 300)
}

//URL returns review url
func (review *Review) URL() string {
	return fmt.Sprintf("/reviews/%d", review.ID)
}

//GetRecentReviewsByArticle returns a slice of last published reviews by article
/*
func GetRecentReviewsByArticle(aID int64) ([]Review, error) {
	var list []Review
	err := db.Select(&list, "SELECT * FROM reviews WHERE published=$1 AND article_id=$2 ORDER BY video DESC, id DESC LIMIT 7", true, aID)
	return list, err
}
*/
