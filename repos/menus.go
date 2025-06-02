package repos

import (
	"github.com/denisbakhtin/medical/models"
	"github.com/jinzhu/gorm"
)

// Menus is an interface for menu items repository
type Menus interface {
	Main() []models.MenuItem
	Scrolled() []models.MenuItem
}

// MenusRepo implements Menus repository interface
type MenusRepo struct {
	db *gorm.DB
}

// NewMenusRepo creates an instance of MenusRepo
func NewMenusRepo(db *gorm.DB) *MenusRepo {
	return &MenusRepo{db: db}
}

// Main returns the list of main menu items
func (r *MenusRepo) Main() []models.MenuItem {
	about := &models.Page{}
	r.db.First(about, models.AboutPageID)
	contacts := &models.Page{}
	r.db.First(contacts, models.ContactsPageID)
	seans := &models.Page{}
	r.db.First(seans, models.SessionPageID)
	var articles []models.Article
	r.db.Where("published = ?", true).Order("id desc").Find(&articles)
	submenu := make([]models.MenuItem, 0, 10)
	for i := range articles {
		submenu = append(submenu, models.MenuItem{URL: articles[i].URL(), Title: articles[i].Name})
	}
	menu := []models.MenuItem{
		{
			URL:   seans.URL(),
			Title: "Приём",
		},
		{
			URL:      "/articles",
			Title:    "Лечение",
			Children: submenu,
		},
		{
			URL:   about.URL(),
			Title: "Врач кинезиолог",
		},
		{
			URL:   "/reviews",
			Title: "Отзывы",
		},
		{
			URL:   contacts.URL(),
			Title: "Контакты",
		},
		{
			URL:   "/exercises",
			Title: "Упражнения",
		},
	}
	return menu
}

// Scrolled returns the list of visible menu items when a web-page is scrolled down
func (r *MenusRepo) Scrolled() []models.MenuItem {
	about := &models.Page{}
	r.db.First(about, models.AboutPageID)
	menu := []models.MenuItem{
		{
			URL:   about.URL(),
			Title: "О враче",
		},
		{
			URL:   "#withoutpain",
			Title: "Этапы лечения",
		},
		{
			URL:   "/reviews",
			Title: "Отзывы",
		},
	}
	return menu
}
