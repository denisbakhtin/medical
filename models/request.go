package models

type Request struct {
	Name    string `form:"name" binding:"required"`
	Phone   string `form:"phone" binding:"required"`
	Comment string `form:"comment" binding:"required"`
	Captcha string `json:"captcha" form:"captcha" gorm:"-" db:"-"`
}
