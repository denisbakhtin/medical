package controllers

import (
	"log"

	"github.com/denisbakhtin/medical/models"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

//SignInGet handles /signin get request
func SignInGet(c *gin.Context) {
	session := sessions.Default(c)
	flashes := session.Flashes()
	session.Save()
	c.HTML(200, "auth/signin", gin.H{
		"Title":  "Вход в систему",
		"Active": "signin",
		"Flash":  flashes,
	})
}

//SignInPost handles /signin post request
func SignInPost(c *gin.Context) {
	db := models.GetDB()
	session := sessions.Default(c)

	login := &models.Login{}
	if c.Bind(login) == nil {
		user := &models.User{}
		db.Where("lower(email) = lower(?)", login.Email).First(user)
		if user.ID == 0 {
			log.Printf("ERROR: Login failed, IP: %s, Email: %s\n", c.ClientIP(), login.Email)
			session.AddFlash("Эл. адрес или пароль указаны неверно")
			session.Save()
			c.Redirect(303, "/signin")
			return
		}
		//create user
		if err := user.ComparePassword(login.Password); err != nil {
			log.Printf("ERROR: Login failed, IP: %s, Email: %s\n", c.ClientIP(), login.Email)
			session.AddFlash("Эл. адрес или пароль указаны неверно")
			session.Save()
			c.Redirect(303, "/signin")
			return
		}

		session.Set("user_id", user.ID)
		session.Save()
		c.Redirect(303, "/")
	}
}

//LogOut handles logout request
func LogOut(c *gin.Context) {
	session := sessions.Default(c)
	session.Delete("user_id")
	session.Save()
	c.Redirect(303, "/")
}

//SignUpGet handles /signup get request
func SignUpGet(c *gin.Context) {
	session := sessions.Default(c)
	flashes := session.Flashes()
	session.Save()
	c.HTML(200, "auth/signup", gin.H{
		"Title":  "Регистрация в системе",
		"Active": "signup",
		"Flash":  flashes,
	})
}

//SignUpPost handles /signup post request
func SignUpPost(c *gin.Context) {
	session := sessions.Default(c)
	db := models.GetDB()

	register := &models.Register{}
	if c.Bind(register) == nil {
		user := &models.User{}
		db.Where("lower(email) = lower(?)", register.Email).First(user)
		if user.ID != 0 {
			session.AddFlash("Пользователь с таким эл. адресом уже существует")
			session.Save()
			c.Redirect(303, "/signup")
			return
		}
		//create user
		user.Email = register.Email
		user.Password = register.Password
		if err := db.Create(user).Error; err != nil {
			session.AddFlash("Ошибка регистрации пользователя")
			session.Save()
			log.Printf("ERROR: ошибка регистрации пользователя: %v", err)
			c.Redirect(303, "/signup")
			return
		}
		session.Set("user_id", user.ID)
		session.Save()
		c.Redirect(303, "/")
	}
}
