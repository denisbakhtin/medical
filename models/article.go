package models

import (
	"fmt"
	"html/template"
	"regexp"
	"time"

	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday"
)

//Article type contains article info
type Article struct {
	ID        int64     `json:"id" db:"id"`
	Name      string    `json:"name"`
	Slug      string    `json:"slug"`
	Content   string    `json:"content"`
	Published bool      `json:"published"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
	//calculated fields
	Comments []Comment `json:"comments" db:"comments"`
}

//Insert stores Article record in db
func (article *Article) Insert() error {
	if len(article.Slug) == 0 {
		article.Slug = createSlug(article.Name)
	}
	err := db.QueryRow(
		`INSERT INTO articles(name, slug, content, published, created_at, updated_at) 
		VALUES($1,$2,$3,$4,$5,$5) RETURNING id`,
		article.Name,
		article.Slug,
		article.Content,
		article.Published,
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
		"UPDATE articles SET name=$2, slug=$3, content=$4, published=$5, updated_at=$6 WHERE id=$1",
		article.ID,
		article.Name,
		article.Slug,
		article.Content,
		article.Published,
		time.Now(),
	)
	return err
}

//Delete removes Article record from db
func (article *Article) Delete() error {
	_, err := db.Exec("DELETE FROM articles WHERE id=$1", article.ID)
	return err
}

//Excerpt returns article excerpt, 300 char long. Html tags are stripped.
func (article *Article) Excerpt() template.HTML {
	//you can sanitize, cut it down, add images, etc
	policy := bluemonday.StrictPolicy() //remove all html tags
	sanitized := policy.Sanitize(string(blackfriday.MarkdownCommon([]byte(article.Content))))
	excerpt := template.HTML(truncate(sanitized, 300) + "...")
	return excerpt
}

//HTMLContent returns parsed html content
func (article *Article) HTMLContent() template.HTML {
	return template.HTML(string(blackfriday.MarkdownCommon([]byte(article.Content))))
}

//GetCommentCount returns comment count
func (article *Article) GetCommentCount() int {
	count := 0
	db.Get(&count, "SELECT count(id) FROM comments WHERE published=$1 AND article_id=$2", true, article.ID)
	return count
}

//GetImage returns extracts first image url from article content
func (article *Article) GetImage() string {
	content := string(blackfriday.MarkdownCommon([]byte(article.Content)))
	re := regexp.MustCompile(`<img[^<>]+src="([^"]+)"[^<>]*>`)
	res := re.FindStringSubmatch(content)
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
	err := db.Select(&list, "SELECT id, name FROM articles WHERE published=$1 ORDER BY id DESC LIMIT 7", true)
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
