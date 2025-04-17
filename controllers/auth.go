package controllers

import (
	"errors"

	"github.com/denisbakhtin/medical/models"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

// SignInGet handles /signin get request
func (app *Application) SignInGet(c *gin.Context) {
	session := sessions.Default(c)
	flashes := session.Flashes()
	_ = session.Save()
	c.HTML(200, "auth/signin", gin.H{
		"Title":  "Вход в систему",
		"Active": "signin",
		"Flash":  flashes,
	})
}

// SignInPost handles /signin post request
func (app *Application) SignInPost(c *gin.Context) {
	session := sessions.Default(c)

	login := &models.Login{}
	if c.Bind(login) == nil {
		user, err := app.UsersRepo.GetByEmail(login.Email)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			app.Logger.Errorf("Login failed, IP: %s, Email: %s\n", c.ClientIP(), login.Email)
			session.AddFlash("Эл. адрес или пароль указаны неверно")
			_ = session.Save()
			c.Redirect(303, "/signin")
			return
		}
		if err != nil {
			app.Error(c, err)
			return
		}

		if err := user.ComparePassword(login.Password); err != nil {
			app.Logger.Errorf("Login failed, IP: %s, Email: %s\n", c.ClientIP(), login.Email)
			session.AddFlash("Эл. адрес или пароль указаны неверно")
			_ = session.Save()
			c.Redirect(303, "/signin")
			return
		}

		session.Set(userIDkey, user.ID)
		_ = session.Save()
		c.Redirect(303, "/")
	}
}

// LogOut handles logout request
func (app *Application) LogOut(c *gin.Context) {
	session := sessions.Default(c)
	session.Delete(userIDkey)
	_ = session.Save()
	c.Redirect(303, "/")
}

// SignUpGet handles /signup get request
func (app *Application) SignUpGet(c *gin.Context) {
	session := sessions.Default(c)
	flashes := session.Flashes()
	_ = session.Save()
	c.HTML(200, "auth/signup", gin.H{
		"Title":  "Регистрация в системе",
		"Active": "signup",
		"Flash":  flashes,
	})
}

// SignUpPost handles /signup post request
func (app *Application) SignUpPost(c *gin.Context) {
	session := sessions.Default(c)

	register := &models.Register{}
	if c.Bind(register) == nil {
		existing, err := app.UsersRepo.GetByEmail(register.Email)
		//already exists
		if existing.ID != 0 {
			session.AddFlash("Пользователь с таким эл. адресом уже существует")
			_ = session.Save()
			c.Redirect(303, "/signup")
			return
		}
		if err != nil && !gorm.IsRecordNotFoundError(err) {
			app.Error(c, err)
			return
		}

		// create user
		user := &models.User{
			Email:    register.Email,
			Password: register.Password,
		}

		if err := app.UsersRepo.Create(user); err != nil {
			session.AddFlash("Ошибка регистрации пользователя")
			_ = session.Save()
			app.Logger.Errorf("Ошибка регистрации пользователя: %v", err)
			c.Redirect(303, "/signup")
			return
		}
		session.Set(userIDkey, user.ID)
		_ = session.Save()
		c.Redirect(303, "/")
	}
}
