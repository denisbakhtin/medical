package models

import (
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"
)

//User type contains user info
type User struct {
	ID        uint      `json:"id" form:"id"`
	Email     string    `json:"email" form:"email" binding:"required"`
	Name      string    `json:"name" form:"name"`
	Password  string    `json:"password" form:"password" binding:"required"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

//HashPassword substitutes User.Password with its bcrypt hash
func (user *User) HashPassword() error {
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hash)
	return nil
}

//ComparePassword compares User.Password hash with raw password
func (user *User) ComparePassword(password string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return err
	}
	return nil
}

func (user *User) BeforeDelete() (err error) {
	count := 0
	db.Model(&User{}).Count(&count)
	if count == 1 {
		return fmt.Errorf("Невозможно удалить последнего пользователя")
	}
	return
}

func (user *User) BeforeCreate() (err error) {
	return user.HashPassword()
}

func (user *User) BeforeSave() (err error) {
	return user.HashPassword()
}
