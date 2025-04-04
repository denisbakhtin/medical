package repos

import (
	"github.com/denisbakhtin/medical/models"
	"github.com/jinzhu/gorm"
)

type Pages interface {
	Get(id uint) (*models.Page, error)
	GetPublished(id uint) (*models.Page, error)
	GetAll() ([]models.Page, error)
	GetAllPublished() ([]models.Page, error)
	Create(*models.Page) error
	Update(*models.Page) error
	Delete(id uint) error
}

type PagesRepo struct {
	db *gorm.DB
}

func NewPagesRepo(db *gorm.DB) Pages {
	return &PagesRepo{db: db}
}

func (r *PagesRepo) Get(id uint) (*models.Page, error) {
	page := &models.Page{}
	err := r.db.First(page, id).Error
	return page, err
}

func (r *PagesRepo) GetPublished(id uint) (*models.Page, error) {
	page := &models.Page{}
	err := r.db.Where("published = ?", true).First(page, id).Error
	return page, err
}

func (r *PagesRepo) GetAll() ([]models.Page, error) {
	var list []models.Page
	err := r.db.Order("published desc, id desc").Find(&list).Error
	return list, err
}

func (r *PagesRepo) GetAllPublished() ([]models.Page, error) {
	var list []models.Page
	err := r.db.Where("published = ?", true).Order("id desc").Find(&list).Error
	return list, err
}

func (r *PagesRepo) Create(page *models.Page) error {
	return r.db.Create(page).Error
}

func (r *PagesRepo) Update(page *models.Page) error {
	return r.db.Save(page).Error
}

func (r *PagesRepo) Delete(id uint) error {
	page := &models.Page{}
	err := r.db.First(page, id).Error
	if err != nil {
		return err
	}
	return r.db.Delete(page).Error
}
