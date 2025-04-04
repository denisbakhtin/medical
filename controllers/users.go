package controllers

import (
	"fmt"

	"github.com/denisbakhtin/medical/helpers"
	"github.com/denisbakhtin/medical/models"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// UsersAdminIndex handles GET /admin/users route
func (app *Application) UsersAdminIndex(c *gin.Context) {
	list, err := app.UsersRepo.GetAll()
	if err != nil {
		app.Error(c, err)
		return
	}
	c.HTML(200, "users/admin/index", gin.H{
		"Title":  "Пользователи",
		"Active": "users",
		"List":   list,
	})
}

// UserAdminCreateGet handles /admin/new_user get request
func (app *Application) UserAdminCreateGet(c *gin.Context) {
	session := sessions.Default(c)
	flashes := session.Flashes()
	_ = session.Save()

	c.HTML(200, "users/admin/form", gin.H{
		"Title":  "Новый пользователь",
		"Active": "users",
		"Flash":  flashes,
	})
}

// UserAdminCreatePost handles /admin/new_user post request
func (app *Application) UserAdminCreatePost(c *gin.Context) {
	session := sessions.Default(c)

	user := &models.User{}
	if c.Bind(user) == nil {
		if err := app.UsersRepo.Create(user); err != nil {
			session.AddFlash(err.Error())
			_ = session.Save()
			c.Redirect(303, "/admin/new_user")
			return
		}
	} else {
		session.AddFlash("Ошибка! Проверьте заполнение всех полей!")
		_ = session.Save()
		c.Redirect(303, "/admin/new_user")
		return
	}
	c.Redirect(303, "/admin/users")
}

// UserAdminUpdateGet handles /admin/edit_user/:id get request
func (app *Application) UserAdminUpdateGet(c *gin.Context) {
	session := sessions.Default(c)
	flashes := session.Flashes()
	_ = session.Save()

	id := helpers.Atouint(c.Param("id"))
	user, err := app.UsersRepo.Get(id)
	if err != nil {
		app.Error(c, err)
		return
	}

	c.HTML(200, "users/admin/form", gin.H{
		"Title":  "Редактировать пользователя",
		"Active": "users",
		"User":   user,
		"Flash":  flashes,
	})
}

// UserAdminUpdatePost handles /admin/edit_user/:id post request
func (app *Application) UserAdminUpdatePost(c *gin.Context) {
	session := sessions.Default(c)

	user := &models.User{}
	id := helpers.Atouint(c.Param("id"))
	if c.Bind(user) == nil {
		if err := app.UsersRepo.Update(user); err != nil {
			session.AddFlash(err.Error())
			_ = session.Save()
			c.Redirect(303, fmt.Sprintf("/admin/edit_user/%v", id))
			return
		}
	} else {
		session.AddFlash("Ошибка! Проверьте заполнение всех полей!")
		_ = session.Save()
		c.Redirect(303, fmt.Sprintf("/admin/edit_user/%v", id))
		return
	}
	c.Redirect(303, "/admin/users")
}

// UserAdminDelete handles /admin/delete_user route
func (app *Application) UserAdminDelete(c *gin.Context) {
	id := helpers.Atouint(c.Request.PostFormValue("id"))

	if err := app.UsersRepo.Delete(id); err != nil {
		app.Error(c, err)
		return
	}
	c.Redirect(303, "/admin/users")
}
