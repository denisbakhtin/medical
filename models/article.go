package models

import (
	"fmt"
	"html/template"
	"regexp"
	"time"
)

//Article type contains article info
type Article struct {
	ID             int64 `db:"id"`
	Name           string
	Slug           string
	Excerpt        string
	Content        string
	SellingPreface string `db:"selling_preface"`
	Published      bool
	CreatedAt      time.Time `db:"created_at"`
	UpdatedAt      time.Time `db:"updated_at"`
	//calculated fields
	TopComments []Comment `db:"-"`
	Comments    []Comment `db:"-"`
}

//Insert stores Article record in db
func (article *Article) Insert() error {
	if len(article.Slug) == 0 {
		article.Slug = createSlug(article.Name)
	}
	err := db.QueryRow(
		`INSERT INTO articles(name, slug, excerpt, content, published, selling_preface, created_at, updated_at) 
		VALUES($1,$2,$3,$4,$5,$6,$7,$7) RETURNING id`,
		article.Name,
		article.Slug,
		article.Excerpt,
		article.Content,
		article.Published,
		article.SellingPreface,
		time.Now(),
	).Scan(&article.ID)
	return err
}

//Update updates Article record in db
func (article *Article) Update() error {
	if len(article.Slug) == 0 {
		article.Slug = createSlug(article.Name)
	}
	_, err := db.Exec(
		"UPDATE articles SET name=$2, slug=$3, excerpt=$4, content=$5, published=$6, selling_preface=$7, updated_at=$8 WHERE id=$1",
		article.ID,
		article.Name,
		article.Slug,
		article.Excerpt,
		article.Content,
		article.Published,
		article.SellingPreface,
		time.Now(),
	)
	return err
}

//Delete removes Article record from db
func (article *Article) Delete() error {
	_, err := db.Exec("DELETE FROM articles WHERE id=$1", article.ID)
	return err
}

//HTMLContent returns parsed html content
func (article *Article) HTMLContent() template.HTML {
	return template.HTML(article.Content)
}

//GetCommentCount returns comment count
func (article *Article) GetCommentCount() int {
	count := 0
	db.Get(&count, "SELECT count(id) FROM comments WHERE published=$1 AND article_id=$2", true, article.ID)
	return count
}

//GetImage returns extracts first image url from article content
func (article *Article) GetImage() string {
	re := regexp.MustCompile(`<img[^<>]+src="([^"]+)"[^<>]*>`)
	res := re.FindStringSubmatch(article.Content)
	if len(res) == 2 {
		return res[1]
	}
	return ""
}

//URL returns article url
func (article *Article) URL() string {
	return fmt.Sprintf("/articles/%d-%s", article.ID, article.Slug)
}

//GetArticle loads Article record by its ID
func GetArticle(id interface{}) (*Article, error) {
	article := &Article{}
	err := db.Get(article, "SELECT * FROM articles WHERE id=$1", id)
	if err != nil {
		return article, err
	}
	article.TopComments, err = GetTopCommentsByArticleID(article.ID)
	if err != nil {
		return article, err
	}
	article.Comments, err = GetCommentsByArticleID(article.ID)
	return article, err
}

//GetArticles returns a slice f articles
func GetArticles() ([]Article, error) {
	var list []Article
	err := db.Select(&list, "SELECT * FROM articles ORDER BY articles.id DESC")
	return list, err
}

//GetPublishedArticles returns a slice of published articles
func GetPublishedArticles() ([]Article, error) {
	var list []Article
	err := db.Select(&list, "SELECT * FROM articles WHERE published=$1 ORDER BY articles.id DESC", true)
	return list, err
}

//GetRecentArticles returns a slice of last 7 published articles
func GetRecentArticles() ([]Article, error) {
	var list []Article
	err := db.Select(&list, "SELECT id, name FROM articles WHERE published=$1 ORDER BY id DESC LIMIT 8", true)
	return list, err
}

//SearchArticles returns a slice of articles, matching query
func SearchArticles(query string) ([]Article, error) {
	var list []Article
	err := db.Select(
		&list,
		`SELECT * FROM articles 
		WHERE to_tsvector('russian', name || ' ' || content) @@ to_tsquery('russian', $1) AND 
		published=$2 
		ORDER BY articles.id DESC`,
		query,
		true,
	)
	return list, err
}
