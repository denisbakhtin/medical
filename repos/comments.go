package repos

import (
	"github.com/denisbakhtin/medical/models"
	"github.com/jinzhu/gorm"
)

type Comments interface {
	Get(id uint) (*models.Comment, error)
	GetTopByArticle(articleID uint) ([]models.Comment, error)
	GetUnpublishedByArticle(articleID uint) ([]models.Comment, error)
	GetCountByArticle(articleID uint) (int, error)
	GetByArticlePage(articleID uint, page, perPage int) ([]models.Comment, error)
	GetAll() ([]models.Comment, error)
	Create(*models.Comment) error
	Update(*models.Comment) error
	Delete(id uint) error
}

type CommentsRepo struct {
	db *gorm.DB
}

func NewCommentsRepo(db *gorm.DB) *CommentsRepo {
	return &CommentsRepo{db: db}
}

func (r *CommentsRepo) Get(id uint) (*models.Comment, error) {
	comment := &models.Comment{}
	err := r.db.First(comment, id).Error
	return comment, err
}

func (r *CommentsRepo) GetTopByArticle(articleID uint) ([]models.Comment, error) {
	var comments []models.Comment
	err := r.db.
		Where("published = ? AND article_id = ? AND author_city = ?", true, articleID, "Москва").
		Order("id desc").
		Limit(10).
		Find(&comments).
		Error
	return comments, err
}

func (r *CommentsRepo) GetUnpublishedByArticle(articleID uint) ([]models.Comment, error) {
	var moscow []models.Comment
	var nonmoscow []models.Comment
	err := r.db.
		Where("published = ? AND answer <> ? AND article_id = ? AND author_city = ?", false, "", articleID, "Москва").
		Order("id desc").
		Find(&moscow).
		Error
	err2 := r.db.
		Where("published = ? AND answer <> ? AND article_id = ? AND author_city <> ?", false, "", articleID, "Москва").
		Order("id desc").
		Find(&nonmoscow).
		Error
	//still ok, can use multierror package though
	if err2 != nil {
		err = err2
	}
	return append(moscow, nonmoscow...), err
}

func (r *CommentsRepo) GetCountByArticle(articleID uint) (int, error) {
	var totalCount int
	err := r.db.Model(models.Comment{}).Where("article_id = ?", articleID).Count(&totalCount).Error
	return totalCount, err
}

func (r *CommentsRepo) GetByArticlePage(articleID uint, page, perPage int) ([]models.Comment, error) {
	var list []models.Comment
	err := r.db.Model(models.Comment{}).
		Where("article_id = ?", articleID).
		Limit(perPage).Offset((page - 1) * perPage).
		Order("answer desc, id desc").Find(&list).
		Error
	return list, err
}

func (r *CommentsRepo) GetAll() ([]models.Comment, error) {
	var list []models.Comment
	err := r.db.Order("id desc").Find(&list).Error
	return list, err
}

func (r *CommentsRepo) Create(comment *models.Comment) error {
	return r.db.Create(comment).Error
}

func (r *CommentsRepo) Update(comment *models.Comment) error {
	return r.db.Save(comment).Error
}

func (r *CommentsRepo) Delete(id uint) error {
	comment := &models.Comment{}
	err := r.db.First(comment, id).Error
	if err != nil {
		return err
	}
	return r.db.Delete(comment).Error
}
