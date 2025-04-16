package repos

import (
	"github.com/denisbakhtin/medical/models"
	"github.com/jinzhu/gorm"
)

// Pages is an interface for "static" (contacs, about, etc.) pages repository
type Pages interface {
	Get(id uint) (*models.Page, error)
	GetPublished(id uint) (*models.Page, error)
	GetAll() ([]models.Page, error)
	GetAllPublished() ([]models.Page, error)
	Create(*models.Page) error
	Update(*models.Page) error
	Delete(id uint) error
}

// PagesRepo implements Pages repository interface
type PagesRepo struct {
	db *gorm.DB
}

// NewPagesRepo creates an instance of PagesRepo
func NewPagesRepo(db *gorm.DB) *PagesRepo {
	return &PagesRepo{db: db}
}

// Get returns a page by ID
func (r *PagesRepo) Get(id uint) (*models.Page, error) {
	page := &models.Page{}
	err := r.db.First(page, id).Error
	return page, err
}

// GetPublished returns a published page by ID
func (r *PagesRepo) GetPublished(id uint) (*models.Page, error) {
	page := &models.Page{}
	err := r.db.Where("published = ?", true).First(page, id).Error
	return page, err
}

// GetAll returns all pages
func (r *PagesRepo) GetAll() ([]models.Page, error) {
	var list []models.Page
	err := r.db.Order("published desc, id desc").Find(&list).Error
	return list, err
}

// GetAllPublished returns all published pages
func (r *PagesRepo) GetAllPublished() ([]models.Page, error) {
	var list []models.Page
	err := r.db.Where("published = ?", true).Order("id desc").Find(&list).Error
	return list, err
}

// Create creates a new page in db
func (r *PagesRepo) Create(page *models.Page) error {
	return r.db.Create(page).Error
}

// Update updates a page in db, ID must be non-zero
func (r *PagesRepo) Update(page *models.Page) error {
	return r.db.Save(page).Error
}

// Delete removes a page from db by ID
func (r *PagesRepo) Delete(id uint) error {
	page := &models.Page{}
	err := r.db.First(page, id).Error
	if err != nil {
		return err
	}
	return r.db.Delete(page).Error
}
