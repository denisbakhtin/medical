package models

import (
	"fmt"
	"time"
)

//Comment type contains article comments
type Comment struct {
	ID          int64     `json:"id" db:"id"`
	ArticleID   int64     `json:"article_id" db:"article_id"`
	AuthorCity  string    `json:"author_city" db:"author_city"`
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
		`INSERT INTO comments(article_id, author_city, author_name, author_email, content, answer, published, created_at, updated_at) 
		VALUES($1,$2,$3,$4,$5,$6,$7,$8,$8) RETURNING id`,
		comment.ArticleID,
		comment.AuthorCity,
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
		SET content=$2, answer=$3, published=$4, author_city=$5, updated_at=$6 
		WHERE id=$1`,
		comment.ID,
		comment.Content,
		comment.Answer,
		comment.Published,
		comment.AuthorCity,
		time.Now(),
	)
	return err
}

//Delete removes Comment from db.
func (comment *Comment) Delete() error {
	_, err := db.Exec("DELETE FROM comments WHERE id=$1", comment.ID)
	return err
}

//Title returns comment excerpt, 100 char long
func (comment *Comment) Title() string {
	return truncate(comment.Content, 50) + "..."
}

//Excerpt returns comment excerpt, 100 char long
func (comment *Comment) Excerpt() string {
	return truncate(comment.Content, 20)
}

//URL returns comment url
func (comment *Comment) URL() string {
	return fmt.Sprintf(
		"/comments/%d-%s",
		comment.ID,
		createSlug(truncate(comment.Content, 70)),
	)
}

//GetSimilar returns a slice of similar (adjacent) comments
func (comment *Comment) GetSimilar() ([]Comment, error) {
	var list []Comment
	err := db.Select(
		&list,
		`(SELECT * FROM comments WHERE id < $1 AND article_id=$2 AND published=$3 ORDER BY id DESC LIMIT $4)
		UNION
		(SELECT * FROM comments WHERE id > $1 AND article_id=$2 AND published=$3 ORDER BY id ASC LIMIT $4)
		ORDER BY id ASC`,
		comment.ID,
		comment.ArticleID,
		true,
		5,
	)
	if err != nil {
		return nil, err
	}
	if len(list) < 10 {
		var list2 []Comment
		err := db.Select(
			&list2,
			`(SELECT * FROM comments WHERE id < $1 AND article_id!=$2 AND published=$3 ORDER BY id DESC LIMIT $4)
			UNION
			(SELECT * FROM comments WHERE id > $1 AND article_id!=$2 AND published=$3 ORDER BY id ASC LIMIT $4)
			ORDER BY id ASC LIMIT $4`,
			comment.ID,
			comment.ArticleID,
			true,
			5,
		)
		if err != nil {
			return nil, err
		}
		list = append(list, list2...)
	}
	return list, err
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

//GetCommentsByArticleID returns a slice of published comments, associated with the given article
func GetCommentsByArticleID(articleID int64) ([]Comment, error) {
	var list []Comment
	var list2 []Comment
	err := db.Select(
		&list,
		`SELECT * FROM comments 
		WHERE published=false AND answer!='' AND article_id=$1 AND author_city=$2
		ORDER BY id DESC`,
		articleID,
		"Москва",
	)
	if err != nil {
		return nil, err
	}
	err = db.Select(
		&list2,
		`SELECT * FROM comments 
		WHERE published=false AND answer!='' AND article_id=$1 AND author_city!=$2
		ORDER BY id DESC`,
		articleID,
		"Москва",
	)
	if err != nil {
		return nil, err
	}
	list = append(list, list2...)
	return list, err
}

//GetTopCommentsByArticleID returns a slice of top (latest, or rated) published comments, associated with given article
func GetTopCommentsByArticleID(articleID int64) ([]Comment, error) {
	var list []Comment
	limit := 10
	err := db.Select(
		&list,
		`SELECT * FROM comments 
		WHERE published=$1 AND article_id=$2 AND author_city=$3
		ORDER BY id DESC LIMIT $4`,
		true,
		articleID,
		"Москва",
		limit,
	)
	return list, err
}
