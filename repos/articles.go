package repos

import (
	"github.com/denisbakhtin/medical/models"
	"github.com/jinzhu/gorm"
)

// Articles is an interface for articles repository
type Articles interface {
	Get(id uint) (*models.Article, error)
	GetPublished(id uint) (*models.Article, error)
	GetAll() ([]models.Article, error)
	GetAllPublished() ([]models.Article, error)
	Create(*models.Article) error
	Update(*models.Article) error
	Delete(id uint) error
}

// ArticlesRepo implements Articles repository interface
type ArticlesRepo struct {
	db *gorm.DB
}

// NewArticlesRepo creates an instance of ArticlesRepo
func NewArticlesRepo(db *gorm.DB) *ArticlesRepo {
	return &ArticlesRepo{db: db}
}

// Get returns an article by its ID
func (r *ArticlesRepo) Get(id uint) (*models.Article, error) {
	article := &models.Article{}
	err := r.db.First(article, id).Error
	return article, err
}

// GetPublished returns a published article by its ID
func (r *ArticlesRepo) GetPublished(id uint) (*models.Article, error) {
	article := &models.Article{}
	err := r.db.Where("published = ?", true).First(article, id).Error
	return article, err
}

// GetAll returns all articles
func (r *ArticlesRepo) GetAll() ([]models.Article, error) {
	var list []models.Article
	err := r.db.Order("published desc, id desc").Find(&list).Error
	return list, err
}

// GetAllPublished returns all published articles
func (r *ArticlesRepo) GetAllPublished() ([]models.Article, error) {
	var list []models.Article
	err := r.db.Where("published = ?", true).Order("id desc").Find(&list).Error
	return list, err
}

// Create creates a new article in db
func (r *ArticlesRepo) Create(article *models.Article) error {
	return r.db.Create(article).Error
}

// Update updates an article. Article ID should be set.
func (r *ArticlesRepo) Update(article *models.Article) error {
	return r.db.Save(article).Error
}

// Delete deletes an article by ID
func (r *ArticlesRepo) Delete(id uint) error {
	article := &models.Article{}
	err := r.db.First(article, id).Error
	if err != nil {
		return err
	}
	return r.db.Delete(article).Error
}
