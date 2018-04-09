package models

//Login view model
type Login struct {
	Email    string `form:"email" binding:"required"`
	Password string `form:"password" binding:"required"`
}

//Register view model
type Register struct {
	Email    string `form:"email" binding:"required"`
	Password string `form:"password" binding:"required"`
}
