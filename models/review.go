package models

import (
	"fmt"
	"time"
)

//Review type contains client review
type Review struct {
	ID              uint   `json:"id"`
	ArticleID       *uint  `json:"article_id"`
	AuthorName      string `json:"author_name"`
	AuthorEmail     string `json:"author_email"`
	Content         string `json:"content"`
	Image           string `json:"image"`
	Video           string `json:"video"`
	MetaKeywords    string
	MetaDescription string
	Published       bool      `json:"published"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
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
