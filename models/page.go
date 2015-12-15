package models

import (
	"fmt"
	"html/template"
	"time"
)

//Page type contains page info
type Page struct {
	ID        int64     `json:"id" db:"id"`
	Name      string    `json:"name"`
	Slug      string    `json:"slug"`
	Content   string    `json:"content"`
	Published bool      `json:"published"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

//Insert stores Page struct in db
func (page *Page) Insert() error {
	if len(page.Slug) == 0 {
		page.Slug = createSlug(page.Name)
	}
	err := db.QueryRow(
		`INSERT INTO pages(name, slug, content, published, created_at, updated_at) 
		VALUES($1,$2,$3,$4,$5,$5) RETURNING id`,
		page.Name,
		page.Slug,
		page.Content,
		page.Published,
		time.Now(),
	).Scan(&page.ID)
	return err
}

//Update updates Page record in db
func (page *Page) Update() error {
	if len(page.Slug) == 0 {
		page.Slug = createSlug(page.Name)
	}
	_, err := db.Exec(
		"UPDATE pages SET name=$2, slug=$3, content=$4, published=$5, updated_at=$6 WHERE id=$1",
		page.ID,
		page.Name,
		page.Slug,
		page.Content,
		page.Published,
		time.Now(),
	)
	return err
}

//Delete removes Page record from db
func (page *Page) Delete() error {
	_, err := db.Exec("DELETE FROM pages WHERE id=$1", page.ID)
	return err
}

//HTMLContent returns parsed html content
func (page *Page) HTMLContent() template.HTML {
	return template.HTML(page.Content)
}

//URL returns page url
func (page *Page) URL() string {
	return fmt.Sprintf("/pages/%d-%s", page.ID, page.Slug)
}

//GetPage loads page record by its id
func GetPage(id interface{}) (*Page, error) {
	page := &Page{}
	err := db.Get(page, "SELECT * FROM pages WHERE id=$1", id)
	return page, err
}

//GetPages returns a slice of all pages
func GetPages() ([]Page, error) {
	var list []Page
	err := db.Select(&list, "SELECT * FROM pages ORDER BY id")
	return list, err
}

//GetPublishedPages returns a slice of published pages
func GetPublishedPages() ([]Page, error) {
	var list []Page
	err := db.Select(&list, "SELECT * FROM pages WHERE published=$1 ORDER BY id", true)
	return list, err
}
