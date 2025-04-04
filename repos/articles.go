package repos

import (
	"github.com/denisbakhtin/medical/models"
	"github.com/jinzhu/gorm"
)

type Articles interface {
	Get(id uint) (*models.Article, error)
	GetPublished(id uint) (*models.Article, error)
	GetAll() ([]models.Article, error)
	GetAllPublished() ([]models.Article, error)
	Create(*models.Article) error
	Update(*models.Article) error
	Delete(id uint) error
}

type ArticlesRepo struct {
	db *gorm.DB
}

func NewArticlesRepo(db *gorm.DB) Articles {
	return &ArticlesRepo{db: db}
}

func (r *ArticlesRepo) Get(id uint) (*models.Article, error) {
	article := &models.Article{}
	err := r.db.First(article, id).Error
	return article, err
}

func (r *ArticlesRepo) GetPublished(id uint) (*models.Article, error) {
	article := &models.Article{}
	err := r.db.Where("published = ?", true).First(article, id).Error
	return article, err
}

func (r *ArticlesRepo) GetAll() ([]models.Article, error) {
	var list []models.Article
	err := r.db.Order("published desc, id desc").Find(&list).Error
	return list, err
}

func (r *ArticlesRepo) GetAllPublished() ([]models.Article, error) {
	var list []models.Article
	err := r.db.Where("published = ?", true).Order("id desc").Find(&list).Error
	return list, err
}

func (r *ArticlesRepo) Create(article *models.Article) error {
	return r.db.Create(article).Error
}

func (r *ArticlesRepo) Update(article *models.Article) error {
	return r.db.Save(article).Error
}

func (r *ArticlesRepo) Delete(id uint) error {
	article := &models.Article{}
	err := r.db.First(article, id).Error
	if err != nil {
		return err
	}
	return r.db.Delete(article).Error
}
