package repos

import (
	"github.com/denisbakhtin/medical/models"
	"github.com/jinzhu/gorm"
)

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

type InfosRepo struct {
	db *gorm.DB
}

func NewInfosRepo(db *gorm.DB) Infos {
	return &InfosRepo{db: db}
}
func (r *InfosRepo) GetPublished(id uint) (*models.Info, error) {
	info := &models.Info{}
	err := r.db.Where("published = ?", true).First(info, id).Error
	return info, err
}

func (r *InfosRepo) GetPublishedCount() (int, error) {
	var totalInfos int
	err := r.db.Model(models.Info{}).Where("published = ?", true).Count(&totalInfos).Error
	return totalInfos, err
}

func (r *InfosRepo) GetPublishedPage(page, perPage int) ([]models.Info, error) {
	var infos []models.Info
	err := r.db.Where("published = ?", true).Order("id desc").Limit(perPage).Offset((page - 1) * perPage).Find(&infos).Error
	return infos, err
}

func (r *InfosRepo) Get(id uint) (*models.Info, error) {
	info := &models.Info{}
	err := r.db.First(&info, id).Error
	return info, err
}

func (r *InfosRepo) GetAll() ([]models.Info, error) {
	var infos []models.Info
	err := r.db.Order("published desc, id desc").Find(&infos).Error
	return infos, err
}

func (r *InfosRepo) GetAllPublished() ([]models.Info, error) {
	var infos []models.Info
	err := r.db.Where("published = ?").Order("id desc").Find(&infos).Error
	return infos, err
}

func (r *InfosRepo) Create(info *models.Info) error {
	return r.db.Create(info).Error
}

func (r *InfosRepo) Update(info *models.Info) error {
	return r.db.Save(info).Error
}

func (r *InfosRepo) Delete(id uint) error {
	info := &models.Info{}
	err := r.db.First(&info, id).Error
	if err != nil {
		return err
	}
	return r.db.Delete(info).Error
}
