package models

import (
	"time"
)

//Comment type contains article comments
type Comment struct {
	ID          int64     `json:"id" db:"id"`
	ArticleID   int64     `json:"article_id" db:"article_id"`
	AuthorName  string    `json:"author_name" db:"author_name"`
	AuthorEmail string    `json:"author_email" db:"author_email"`
	Content     string    `json:"content"`
	Answer      string    `json:"answer"`
	Published   bool      `json:"published"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

//Insert stores Comment  in db
func (comment *Comment) Insert() error {
	err := db.QueryRow(
		`INSERT INTO comments(article_id, author_name, author_email, content, answer, published, created_at, updated_at) 
		VALUES($1,$2,$3,$4,$5,$6,$7,$7) RETURNING id`,
		comment.ArticleID,
		comment.AuthorName,
		comment.AuthorEmail,
		comment.Content,
		comment.Answer,
		comment.Published,
		time.Now(),
	).Scan(&comment.ID)
	return err
}

//Update updates Comment record in db
func (comment *Comment) Update() error {
	_, err := db.Exec(
		`UPDATE comments 
		SET content=$2, answer=$3, published=$4, updated_at=$5 
		WHERE id=$1`,
		comment.ID,
		comment.Content,
		comment.Answer,
		comment.Published,
		time.Now(),
	)
	return err
}

//Delete removes Comment from db.
func (comment *Comment) Delete() error {
	_, err := db.Exec("DELETE FROM comments WHERE id=$1", comment.ID)
	return err
}

//Excerpt returns comment excerpt, 100 char long
func (comment *Comment) Excerpt() string {
	return truncate(comment.Content, 20)
}

//GetComment returns Comment record by its ID.
func GetComment(id interface{}) (*Comment, error) {
	comment := &Comment{}
	err := db.Get(comment, "SELECT * FROM comments WHERE id=$1", id)
	return comment, err
}

//GetComments returns a slice of comments
func GetComments() ([]Comment, error) {
	var list []Comment
	err := db.Select(&list, "SELECT * FROM comments ORDER BY id DESC")
	return list, err
}

//GetPublishedComments returns a slice published of comments
func GetPublishedComments() ([]Comment, error) {
	var list []Comment
	err := db.Select(&list, "SELECT * FROM comments WHERE published=$1 ORDER BY id DESC", true)
	return list, err
}

//GetCommentsByArticleID returns a slice of published comments, associated with given article
func GetCommentsByArticleID(articleID int64) ([]Comment, error) {
	var list []Comment
	err := db.Select(
		&list,
		`SELECT * FROM comments 
		WHERE published=$1 AND article_id=$2 
		ORDER BY id DESC`,
		true,
		articleID,
	)
	return list, err
}
