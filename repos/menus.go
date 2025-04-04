package repos

import (
	"github.com/denisbakhtin/medical/models"
	"github.com/jinzhu/gorm"
)

type Menus interface {
	Main() []models.MenuItem
	Scrolled() []models.MenuItem
}

type MenusRepo struct {
	db *gorm.DB
}

func NewMenusRepo(db *gorm.DB) Menus {
	return &MenusRepo{db: db}
}

// Main returns the list of main menu items
func (r *MenusRepo) Main() []models.MenuItem {
	about := &models.Page{}
	r.db.First(about, models.AboutPageId)
	contacts := &models.Page{}
	r.db.First(contacts, models.ContactsPageId)
	seans := &models.Page{}
	r.db.First(seans, models.SessionPageId)
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

// Scrolled returns the list of scroll menu items
func (r *MenusRepo) Scrolled() []models.MenuItem {
	about := &models.Page{}
	r.db.First(about, models.AboutPageId)
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
