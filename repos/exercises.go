package repos

import (
	"github.com/denisbakhtin/medical/models"
	"github.com/jinzhu/gorm"
)

type Exercises interface {
	GetPublished(id uint) (*models.Exercise, error)
	GetAllPublished() ([]models.Exercise, error)
	Get(id uint) (*models.Exercise, error)
	GetAll() ([]models.Exercise, error)
	Create(ex *models.Exercise) error
	Update(ex *models.Exercise) error
	Delete(id uint) error
}

type ExercisesRepo struct {
	db *gorm.DB
}

func NewExercisesRepo(db *gorm.DB) *ExercisesRepo {
	return &ExercisesRepo{db: db}
}

func (r *ExercisesRepo) GetPublished(id uint) (*models.Exercise, error) {
	ex := &models.Exercise{}
	err := r.db.Where("published = ?", true).First(ex, id).Error
	return ex, err
}

func (r *ExercisesRepo) GetAllPublished() ([]models.Exercise, error) {
	var list []models.Exercise
	err := r.db.Where("published = ?", true).Order("sort_ord asc").Find(&list).Error
	return list, err
}

func (r *ExercisesRepo) Get(id uint) (*models.Exercise, error) {
	ex := &models.Exercise{}
	err := r.db.First(ex, id).Error
	return ex, err
}

func (r *ExercisesRepo) GetAll() ([]models.Exercise, error) {
	var list []models.Exercise
	err := r.db.Order("sort_ord asc").Find(&list).Error
	return list, err
}

func (r *ExercisesRepo) Create(ex *models.Exercise) error {
	return r.db.Create(ex).Error
}

func (r *ExercisesRepo) Update(ex *models.Exercise) error {
	return r.db.Save(ex).Error
}

func (r *ExercisesRepo) Delete(id uint) error {
	ex := &models.Exercise{}
	err := r.db.First(ex, id).Error
	if err != nil {
		return err
	}
	return r.db.Delete(ex).Error
}
