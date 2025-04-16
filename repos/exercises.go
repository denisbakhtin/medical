package repos

import (
	"github.com/denisbakhtin/medical/models"
	"github.com/jinzhu/gorm"
)

// Exercises is an interface for exercises repository
type Exercises interface {
	GetPublished(id uint) (*models.Exercise, error)
	GetAllPublished() ([]models.Exercise, error)
	Get(id uint) (*models.Exercise, error)
	GetAll() ([]models.Exercise, error)
	Create(ex *models.Exercise) error
	Update(ex *models.Exercise) error
	Delete(id uint) error
}

// ExercisesRepo implements Exercises repository interface
type ExercisesRepo struct {
	db *gorm.DB
}

// NewExercisesRepo creates an instance of ExercisesRepo
func NewExercisesRepo(db *gorm.DB) *ExercisesRepo {
	return &ExercisesRepo{db: db}
}

// GetPublished returns a published exercise by its ID
func (r *ExercisesRepo) GetPublished(id uint) (*models.Exercise, error) {
	ex := &models.Exercise{}
	err := r.db.Where("published = ?", true).First(ex, id).Error
	return ex, err
}

// GetAllPublished returns all published exercises
func (r *ExercisesRepo) GetAllPublished() ([]models.Exercise, error) {
	var list []models.Exercise
	err := r.db.Where("published = ?", true).Order("sort_ord asc").Find(&list).Error
	return list, err
}

// Get returns an exercise by its ID
func (r *ExercisesRepo) Get(id uint) (*models.Exercise, error) {
	ex := &models.Exercise{}
	err := r.db.First(ex, id).Error
	return ex, err
}

// GetAll returns all exercises
func (r *ExercisesRepo) GetAll() ([]models.Exercise, error) {
	var list []models.Exercise
	err := r.db.Order("sort_ord asc").Find(&list).Error
	return list, err
}

// Create creates a new exercise in db
func (r *ExercisesRepo) Create(ex *models.Exercise) error {
	return r.db.Create(ex).Error
}

// Update updates an exercise in db, ID must be non-zero
func (r *ExercisesRepo) Update(ex *models.Exercise) error {
	return r.db.Save(ex).Error
}

// Delete removes an exercise by ID from db
func (r *ExercisesRepo) Delete(id uint) error {
	ex := &models.Exercise{}
	err := r.db.First(ex, id).Error
	if err != nil {
		return err
	}
	return r.db.Delete(ex).Error
}
