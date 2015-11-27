package models

import (
	"fmt"
	"time"
)

//Review type contains client review
type Review struct {
	ID          int64     `json:"id" db:"id"`
	AuthorName  string    `json:"author_name" db:"author_name"`
	AuthorEmail string    `json:"author_email" db:"author_email"`
	Content     string    `json:"content"`
	Image       string    `json:"image"`
	Published   bool      `json:"published"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

//Insert stores Review in db
func (review *Review) Insert() error {
	err := db.QueryRow(
		`INSERT INTO reviews(author_name, author_email, content, image, published, created_at, updated_at) 
		VALUES($1,$2,$3,$4,$5,$6,$6) RETURNING id`,
		review.AuthorName,
		review.AuthorEmail,
		review.Content,
		review.Image,
		review.Published,
		time.Now(),
	).Scan(&review.ID)
	return err
}

//Update updates Review record in db
func (review *Review) Update() error {
	_, err := db.Exec(
		`UPDATE reviews 
		SET author_name=$2, author_email=$3, content=$4, image=$5, published=$6, updated_at=$7 
		WHERE id=$1`,
		review.ID,
		review.AuthorName,
		review.AuthorEmail,
		review.Content,
		review.Image,
		review.Published,
		time.Now(),
	)
	return err
}

//Delete removes Review from db.
func (review *Review) Delete() error {
	_, err := db.Exec("DELETE FROM reviews WHERE id=$1", review.ID)
	return err
}

//Excerpt returns review excerpt, 100 char long
func (review *Review) Excerpt() string {
	return truncate(review.Content, 300)
}

//URL returns review url
func (review *Review) URL() string {
	return fmt.Sprintf("/reviews/%d", review.ID)
}

//GetReview returns Review record by its ID.
func GetReview(id interface{}) (*Review, error) {
	review := &Review{}
	err := db.Get(review, "SELECT * FROM reviews WHERE id=$1", id)
	return review, err
}

//GetReviews returns a slice of reviews
func GetReviews() ([]Review, error) {
	var list []Review
	err := db.Select(&list, "SELECT * FROM reviews ORDER BY id DESC")
	return list, err
}

//GetPublishedReviews returns a slice published of reviews
func GetPublishedReviews() ([]Review, error) {
	var list []Review
	err := db.Select(&list, "SELECT * FROM reviews WHERE published=$1 ORDER BY id DESC", true)
	return list, err
}

//GetRecentReviews returns a slice of last 7 published reviews
func GetRecentReviews() ([]Review, error) {
	var list []Review
	err := db.Select(&list, "SELECT * FROM reviews WHERE published=$1 ORDER BY id DESC LIMIT 7", true)
	return list, err
}
