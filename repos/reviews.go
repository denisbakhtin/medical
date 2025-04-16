package repos

import (
	"github.com/denisbakhtin/medical/models"
	"github.com/jinzhu/gorm"
)

// Reviews is an interface for clients reviews repository
type Reviews interface {
	Get(id uint) (*models.Review, error)
	GetPublished(id uint) (*models.Review, error)
	GetPublishedByArticle(articleID uint) ([]models.Review, error)
	GetAll() ([]models.Review, error)
	GetAllPublished() ([]models.Review, error)
	GetLastPublished() ([]models.Review, error)
	Create(review *models.Review) error
	Update(review *models.Review) error
	Delete(id uint) error
}

// ReviewsRepo implements Reviews repository interface
type ReviewsRepo struct {
	db *gorm.DB
}

// NewReviewsRepo creates an instance of ReviewsRepo
func NewReviewsRepo(db *gorm.DB) *ReviewsRepo {
	return &ReviewsRepo{db: db}
}

// Get returns a review by ID
func (r *ReviewsRepo) Get(id uint) (*models.Review, error) {
	review := &models.Review{}
	err := r.db.First(review, id).Error
	return review, err
}

// GetPublished returns a published review by ID
func (r *ReviewsRepo) GetPublished(id uint) (*models.Review, error) {
	review := &models.Review{}
	err := r.db.Where("published=?", true).First(review, id).Error
	return review, err
}

// GetPublishedByArticle returns a list of published reviews by article ID
func (r *ReviewsRepo) GetPublishedByArticle(articleID uint) ([]models.Review, error) {
	var reviews []models.Review
	err := r.db.Where("published = ? and article_id = ?", true, articleID).Order("created_at desc").Find(&reviews).Error
	return reviews, err
}

// GetAll returns a list of all reviews
func (r *ReviewsRepo) GetAll() ([]models.Review, error) {
	var reviews []models.Review
	err := r.db.Order("id desc").Find(&reviews).Error
	return reviews, err
}

// GetAllPublished returns a list of all published reviews
func (r *ReviewsRepo) GetAllPublished() ([]models.Review, error) {
	var reviews []models.Review
	err := r.db.Where("published = ?", true).Order("id desc").Find(&reviews).Error
	return reviews, err
}

// GetLastPublished returns a list of last (5) published reviews
func (r *ReviewsRepo) GetLastPublished() ([]models.Review, error) {
	const limit = 5
	var reviews []models.Review
	err := r.db.Where("published = ?", true).Order("id desc").Limit(limit).Find(&reviews).Error
	return reviews, err
}

// Create creates a new review in db
func (r *ReviewsRepo) Create(review *models.Review) error {
	return r.db.Create(review).Error
}

// Update updates a review in db, ID must be non-zero
func (r *ReviewsRepo) Update(review *models.Review) error {
	//update only non empty fields
	return r.db.Model(&models.Review{}).Updates(review).Error
}

// Delete removes a review from db by ID
func (r *ReviewsRepo) Delete(id uint) error {
	review := &models.Review{}
	err := r.db.First(review, id).Error
	if err != nil {
		return err
	}
	return r.db.Delete(review).Error
}
