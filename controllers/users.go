package controllers

import (
	"fmt"

	"github.com/denisbakhtin/medical/helpers"
	"github.com/denisbakhtin/medical/models"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// UsersAdminIndex handles GET /admin/users route
func UsersAdminIndex(c *gin.Context) {
	db := models.GetDB()

	var list []models.User
	db.Find(&list)
	c.HTML(200, "users/admin/index", gin.H{
		"Title":  "Пользователи",
		"Active": "users",
		"List":   list,
	})
}

// UserAdminCreateGet handles /admin/new_user get request
func UserAdminCreateGet(c *gin.Context) {
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
func UserAdminCreatePost(c *gin.Context) {
	db := models.GetDB()
	session := sessions.Default(c)

	user := &models.User{}
	if c.Bind(user) == nil {
		if err := db.Create(user).Error; err != nil {
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
func UserAdminUpdateGet(c *gin.Context) {
	db := models.GetDB()
	session := sessions.Default(c)
	flashes := session.Flashes()
	_ = session.Save()

	id := helpers.Atouint(c.Param("id"))
	user := &models.User{}
	db.First(user, id)
	if user.ID == 0 {
		c.HTML(404, "errors/404", nil)
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
func UserAdminUpdatePost(c *gin.Context) {
	db := models.GetDB()
	session := sessions.Default(c)

	user := &models.User{}
	id := helpers.Atouint(c.Param("id"))
	if c.Bind(user) == nil {
		if err := db.Save(user).Error; err != nil {
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
func UserAdminDelete(c *gin.Context) {
	db := models.GetDB()

	id := helpers.Atouint(c.Request.PostFormValue("id"))
	user := &models.User{}
	db.First(user, id)
	if user.ID == 0 {
		c.HTML(404, "errors/404", nil)
	}

	if err := db.Delete(user).Error; err != nil {
		c.HTML(500, "errors/500", helpers.ErrorData(err))
		return
	}
	c.Redirect(303, "/admin/users")
}
