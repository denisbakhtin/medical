package repos

import (
	"errors"

	"github.com/denisbakhtin/medical/models"
	"github.com/jinzhu/gorm"
)

// Users is an interface for users repository
type Users interface {
	Get(id uint) (*models.User, error)
	GetByEmail(email string) (*models.User, error)
	GetAll() ([]models.User, error)
	Create(*models.User) error
	Update(*models.User) error
	Delete(id uint) error
}

// UsersRepo implements Users repository interface
type UsersRepo struct {
	db *gorm.DB
}

// NewUsersRepo creates an instance of UsersRepo
func NewUsersRepo(db *gorm.DB) *UsersRepo {
	return &UsersRepo{db: db}
}

// Get returns a user by ID
func (r *UsersRepo) Get(id uint) (*models.User, error) {
	user := &models.User{}
	err := r.db.First(user, id).Error
	return user, err
}

// GetByEmail returns a user by email
func (r *UsersRepo) GetByEmail(email string) (*models.User, error) {
	user := &models.User{}
	err := r.db.Where("lower(email) = lower(?)", email).First(user).Error
	return user, err
}

// GetAll returns a list of all users
func (r *UsersRepo) GetAll() ([]models.User, error) {
	var users []models.User
	err := r.db.Find(&users).Error
	return users, err
}

// Create creates a new user in db
func (r *UsersRepo) Create(user *models.User) error {
	return r.db.Create(user).Error
}

// Update updates a user in db, ID must be non-zero
func (r *UsersRepo) Update(user *models.User) error {
	return r.db.Save(user).Error
}

// Delete removes a user from db by ID
func (r *UsersRepo) Delete(id uint) error {
	count := 0
	r.db.Model(&models.User{}).Count(&count)
	if count == 1 {
		return errors.New("невозможно удалить последнего пользователя")
	}

	user := &models.User{}
	err := r.db.First(user, id).Error
	if err != nil {
		return err
	}

	return r.db.Delete(user).Error
}
