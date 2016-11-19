package models

//User type contains user info
type Login struct {
	Email    string `form:"email" binding:"required"`
	Password string `form:"password" binding:"required"`
}

type Register struct {
	Email    string `form:"email" binding:"required"`
	Password string `form:"password" binding:"required"`
}
