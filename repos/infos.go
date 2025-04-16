package repos

import (
	"github.com/denisbakhtin/medical/models"
	"github.com/jinzhu/gorm"
)

// Infos is an interface for Info pages repository
type Infos interface {
	GetPublished(id uint) (*models.Info, error)
	GetPublishedCount() (int, error)
	GetPublishedPage(page, perPage int) ([]models.Info, error)
	Get(id uint) (*models.Info, error)
	GetAll() ([]models.Info, error)
	GetAllPublished() ([]models.Info, error)
	Create(*models.Info) error
	Update(*models.Info) error
	Delete(id uint) error
}

// InfosRepo implements Infos repository interface
type InfosRepo struct {
	db *gorm.DB
}

// NewInfosRepo creates an instance of InfosRepo
func NewInfosRepo(db *gorm.DB) *InfosRepo {
	return &InfosRepo{db: db}
}

// GetPublished returns a published info by ID
func (r *InfosRepo) GetPublished(id uint) (*models.Info, error) {
	info := &models.Info{}
	err := r.db.Where("published = ?", true).First(info, id).Error
	return info, err
}

// GetPublishedCount returns a total count of published infos
func (r *InfosRepo) GetPublishedCount() (int, error) {
	var totalInfos int
	err := r.db.Model(models.Info{}).Where("published = ?", true).Count(&totalInfos).Error
	return totalInfos, err
}

// GetPublishedPage returns a paginated list of published infos
func (r *InfosRepo) GetPublishedPage(page, perPage int) ([]models.Info, error) {
	var infos []models.Info
	err := r.db.Where("published = ?", true).Order("id desc").Limit(perPage).Offset((page - 1) * perPage).Find(&infos).Error
	return infos, err
}

// Get returns an info by ID
func (r *InfosRepo) Get(id uint) (*models.Info, error) {
	info := &models.Info{}
	err := r.db.First(&info, id).Error
	return info, err
}

// GetAll returns a list of all infos
func (r *InfosRepo) GetAll() ([]models.Info, error) {
	var infos []models.Info
	err := r.db.Order("published desc, id desc").Find(&infos).Error
	return infos, err
}

// GetAllPublished returns a list of published infos
func (r *InfosRepo) GetAllPublished() ([]models.Info, error) {
	var infos []models.Info
	err := r.db.Where("published = ?").Order("id desc").Find(&infos).Error
	return infos, err
}

// Create creates a new info record in db
func (r *InfosRepo) Create(info *models.Info) error {
	return r.db.Create(info).Error
}

// Update updates an info record in db, ID must be non-zero
func (r *InfosRepo) Update(info *models.Info) error {
	return r.db.Save(info).Error
}

// Delete removes an info from db by ID
func (r *InfosRepo) Delete(id uint) error {
	info := &models.Info{}
	err := r.db.First(&info, id).Error
	if err != nil {
		return err
	}
	return r.db.Delete(info).Error
}
