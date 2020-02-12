package controllers

import (
	"fmt"
	"strings"

	"github.com/denisbakhtin/medical/helpers"
	"github.com/denisbakhtin/medical/models"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

//ExerciseShow handles GET /exercises/:id-slug route
func ExerciseShow(c *gin.Context) {
	db := models.GetDB()
	session := sessions.Default(c)

	idslug := c.Param("idslug")
	id := helpers.Atouint(strings.Split(idslug, "-")[0])
	exercise := &models.Exercise{}
	db.First(exercise, id)
	if exercise.ID == 0 || !exercise.Published {
		c.HTML(404, "errors/404", nil)
		return
	}
	//redirect to canonical url
	if c.Request.URL.Path != exercise.URL() {
		c.Redirect(301, exercise.URL())
		return
	}

	flashes := session.Flashes()
	session.Save()
	c.HTML(200, "exercises/show", gin.H{
		"Exercise":        exercise,
		"Title":           exercise.Name,
		"Active":          "/exercises",
		"MetaDescription": exercise.MetaDescription,
		"MetaKeywords":    exercise.MetaKeywords,
		"Ogheadprefix":    "og: http://ogp.me/ns# fb: http://ogp.me/ns/fb# article: http://ogp.me/ns/article#",
		"Ogtitle":         exercise.Name,
		"Ogurl":           fmt.Sprintf("http://%s%s", c.Request.Host, exercise.URL()),
		"Ogtype":          "video.movie",
		"Ogdescription":   exercise.MetaDescription,
		"Ogimage":         exercise.Image,
		"Flash":           flashes,
		"Authenticated":   (session.Get("user_id") != nil),
	})
}

//ExercisesIndex handles GET /exercises route
func ExercisesIndex(c *gin.Context) {
	db := models.GetDB()
	session := sessions.Default(c)

	var list []models.Exercise
	if err := db.Where("published = ?", true).Order("id asc").Find(&list).Error; err != nil {
		c.HTML(500, "errors/500", helpers.ErrorData(err))
		return
	}

	c.HTML(200, "exercises/index", gin.H{
		"Title":           "Упражнения для пациентов",
		"Active":          c.Request.RequestURI,
		"List":            list,
		"MetaDescription": "Упражнения пациентам для самостоятельного домашнего выполнения...",
		"MetaKeywords":    "кинезиология, упражнения, реабилитация, прикладная",
		"Authenticated":   (session.Get("user_id") != nil),
	})
}

//ExercisesAdminIndex handles GET /admin/exercises route
func ExercisesAdminIndex(c *gin.Context) {
	db := models.GetDB()

	var list []models.Exercise
	if err := db.Order("published desc, id desc").Find(&list).Error; err != nil {
		c.HTML(500, "errors/500", helpers.ErrorData(err))
		return
	}
	c.HTML(200, "exercises/admin/index", gin.H{
		"Title":  "Упражнения",
		"Active": "exercises",
		"List":   list,
	})
}

//ExerciseAdminCreateGet handles /admin/new_exercise route
func ExerciseAdminCreateGet(c *gin.Context) {
	session := sessions.Default(c)
	flashes := session.Flashes()
	session.Save()

	c.HTML(200, "exercises/admin/form", gin.H{
		"Title":    "Новое упражнение",
		"Active":   "exercises",
		"Flash":    flashes,
		"Exercise": models.Exercise{Published: true},
	})
}

//ExerciseAdminCreatePost handles /admin/new_exercise post request
func ExerciseAdminCreatePost(c *gin.Context) {
	session := sessions.Default(c)
	db := models.GetDB()
	vuri := ""
	iuri := ""

	err := c.Request.ParseMultipartForm(32 << 24) // ~500MB
	if err != nil {
		c.HTML(500, "errors/500", helpers.ErrorData(err))
		return
	}

	vmpartFile, vmpartHeader, err := c.Request.FormFile("videofile")
	if err == nil {
		defer vmpartFile.Close()

		vuri, err = saveFile(vmpartHeader, vmpartFile)
		if err != nil {
			c.HTML(500, "errors/500", helpers.ErrorData(err))
			return
		}
	}

	impartFile, impartHeader, err := c.Request.FormFile("imagefile")
	if err == nil {
		defer impartFile.Close()

		iuri, err = saveFile(impartHeader, impartFile)
		if err != nil {
			c.HTML(500, "errors/500", helpers.ErrorData(err))
			return
		}
	}

	exercise := &models.Exercise{}
	if c.Bind(exercise) == nil {
		exercise.Video = vuri
		exercise.Image = iuri
		if err := db.Create(exercise).Error; err != nil {
			session.AddFlash(err.Error())
			session.Save()
			c.Redirect(303, "/admin/new_exercise")
			return
		}
		c.Redirect(303, "/admin/exercises")
	} else {
		session.AddFlash("Ошибка! Проверьте внимательно заполнение всех полей!")
		session.Save()
		c.Redirect(303, "/admin/new_exercise")
	}
}

//ExerciseAdminUpdateGet handles /admin/edit_exercise/:id get request
func ExerciseAdminUpdateGet(c *gin.Context) {
	session := sessions.Default(c)
	flashes := session.Flashes()
	session.Save()
	db := models.GetDB()

	id := c.Param("id")
	exercise := &models.Exercise{}
	db.First(exercise, id)
	if exercise.ID == 0 {
		c.HTML(404, "errors/404", nil)
		return
	}

	c.HTML(200, "exercises/admin/form", gin.H{
		"Title":    "Редактировать упражнение",
		"Active":   "exercises",
		"Exercise": exercise,
		"Flash":    flashes,
	})
}

//ExerciseAdminUpdatePost handles /admin/edit_exercise/:id post request
func ExerciseAdminUpdatePost(c *gin.Context) {
	session := sessions.Default(c)
	db := models.GetDB()
	vuri := ""
	iuri := ""

	err := c.Request.ParseMultipartForm(32 << 24) // ~500MB
	if err != nil {
		c.HTML(500, "errors/500", helpers.ErrorData(err))
		return
	}

	vmpartFile, vmpartHeader, err := c.Request.FormFile("videofile")
	if err == nil {
		defer vmpartFile.Close()

		vuri, err = saveFile(vmpartHeader, vmpartFile)
		if err != nil {
			c.HTML(500, "errors/500", helpers.ErrorData(err))
			return
		}
	}

	impartFile, impartHeader, err := c.Request.FormFile("imagefile")
	if err == nil {
		defer impartFile.Close()

		iuri, err = saveFile(impartHeader, impartFile)
		if err != nil {
			c.HTML(500, "errors/500", helpers.ErrorData(err))
			return
		}
	}

	exercise := &models.Exercise{}
	if c.Bind(exercise) == nil {
		if len(vuri) > 0 {
			exercise.Video = vuri
		}
		if len(iuri) > 0 {
			exercise.Image = iuri
		}
		if err := db.Save(exercise).Error; err != nil {
			session.AddFlash(err.Error())
			session.Save()
			c.Redirect(303, c.Request.RequestURI)
			return
		}
		c.Redirect(303, "/admin/exercises")
	} else {
		session.AddFlash("Ошибка! Проверьте внимательно заполнение всех полей!")
		session.Save()
		c.Redirect(303, c.Request.RequestURI)
	}
}

//ExerciseAdminDelete handles /admin/delete_exercise route
func ExerciseAdminDelete(c *gin.Context) {
	db := models.GetDB()

	exercise := &models.Exercise{}
	db.First(exercise, c.Request.PostFormValue("id"))
	if exercise.ID == 0 {
		c.HTML(404, "errors/404", nil)
		return
	}

	if err := db.Delete(exercise).Error; err != nil {
		c.HTML(500, "errors/500", helpers.ErrorData(err))
		return
	}
	c.Redirect(303, "/admin/exercises")
}
