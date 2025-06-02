package controllers

import (
	"strings"

	"github.com/denisbakhtin/medical/helpers"
	"github.com/denisbakhtin/medical/models"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// ExerciseShow handles GET /exercises/:id-slug route
func (app *Application) ExerciseShow(c *gin.Context) {
	session := sessions.Default(c)

	idslug := c.Param("idslug")
	id := helpers.Atouint(strings.Split(idslug, "-")[0])
	exercise, err := app.ExercisesRepo.GetPublished(id)
	if err != nil {
		app.Error(c, err)
		return
	}
	// redirect to canonical url
	if c.Request.URL.Path != exercise.URL() {
		c.Redirect(301, exercise.URL())
		return
	}

	flashes := session.Flashes()
	_ = session.Save()
	c.HTML(200, "exercises/show", gin.H{
		"Exercise":        exercise,
		"Title":           exercise.Name,
		"Active":          "/exercises",
		"MetaDescription": exercise.MetaDescription,
		"MetaKeywords":    exercise.MetaKeywords,
		"Ogheadprefix":    "og: http://ogp.me/ns# fb: http://ogp.me/ns/fb# article: http://ogp.me/ns/article#",
		"Ogtitle":         exercise.Name,
		"Ogurl":           app.fullURL(exercise.URL()),
		"Ogtype":          "video.movie",
		"Ogdescription":   exercise.MetaDescription,
		"Ogimage":         exercise.Image,
		"Flash":           flashes,
		"Authenticated":   app.authenticated(session),
	})
}

// ExercisesIndex handles GET /exercises route
func (app *Application) ExercisesIndex(c *gin.Context) {
	session := sessions.Default(c)

	list, err := app.ExercisesRepo.GetAllPublished()
	if err != nil {
		app.Error(c, err)
		return
	}

	c.HTML(200, "exercises/index", gin.H{
		"Title":           "Упражнения для пациентов",
		"Active":          c.Request.RequestURI,
		"List":            list,
		"MetaDescription": "Упражнения пациентам для самостоятельного домашнего выполнения...",
		"MetaKeywords":    "кинезиология, упражнения, реабилитация, прикладная",
		"Authenticated":   app.authenticated(session),
	})
}

// ExercisesAdminIndex handles GET /admin/exercises route
func (app *Application) ExercisesAdminIndex(c *gin.Context) {
	list, err := app.ExercisesRepo.GetAll()
	if err != nil {
		app.Error(c, err)
		return
	}
	c.HTML(200, "exercises/admin/index", gin.H{
		"Title":  "Упражнения",
		"Active": "exercises",
		"List":   list,
	})
}

// ExerciseAdminCreateGet handles /admin/new_exercise route
func (app *Application) ExerciseAdminCreateGet(c *gin.Context) {
	session := sessions.Default(c)
	flashes := session.Flashes()
	_ = session.Save()

	c.HTML(200, "exercises/admin/form", gin.H{
		"Title":    "Новое упражнение",
		"Active":   "exercises",
		"Flash":    flashes,
		"Exercise": models.Exercise{Published: true},
	})
}

// ExerciseAdminCreatePost handles /admin/new_exercise post request
func (app *Application) ExerciseAdminCreatePost(c *gin.Context) {
	session := sessions.Default(c)
	//video uri
	vuri := ""
	//image uri
	iuri := ""

	err := c.Request.ParseMultipartForm(32 << 24) // ~500MB
	if err != nil {
		app.Error(c, err)
		return
	}

	vmpartFile, vmpartHeader, err := c.Request.FormFile("videofile")
	if err == nil {
		defer func() {
			err := vmpartFile.Close()
			if err != nil {
				app.Logger.Errorf("Error closing videofile multi-part %v", err)
			}
		}()

		vuri, err = app.saveFile(vmpartHeader, vmpartFile)
		if err != nil {
			app.Error(c, err)
			return
		}
	}

	impartFile, impartHeader, err := c.Request.FormFile("imagefile")
	if err == nil {
		defer func() {
			err := impartFile.Close()
			if err != nil {
				app.Logger.Errorf("Error closing imagefile multi-part %v", err)
			}
		}()

		iuri, err = app.saveFile(impartHeader, impartFile)
		if err != nil {
			app.Error(c, err)
			return
		}
	}

	exercise := &models.Exercise{}
	if c.Bind(exercise) == nil {
		exercise.Video = vuri
		exercise.Image = iuri
		if err := app.ExercisesRepo.Create(exercise); err != nil {
			session.AddFlash(err.Error())
			_ = session.Save()
			c.Redirect(303, "/admin/new_exercise")
			return
		}
		c.Redirect(303, "/admin/exercises")
	} else {
		session.AddFlash("Ошибка! Проверьте внимательно заполнение всех полей!")
		_ = session.Save()
		c.Redirect(303, "/admin/new_exercise")
	}
}

// ExerciseAdminUpdateGet handles /admin/edit_exercise/:id get request
func (app *Application) ExerciseAdminUpdateGet(c *gin.Context) {
	session := sessions.Default(c)
	flashes := session.Flashes()
	_ = session.Save()

	id := helpers.Atouint(c.Param("id"))
	exercise, err := app.ExercisesRepo.Get(id)
	if err != nil {
		app.Error(c, err)
		return
	}

	c.HTML(200, "exercises/admin/form", gin.H{
		"Title":    "Редактировать упражнение",
		"Active":   "exercises",
		"Exercise": exercise,
		"Flash":    flashes,
	})
}

// ExerciseAdminUpdatePost handles /admin/edit_exercise/:id post request
func (app *Application) ExerciseAdminUpdatePost(c *gin.Context) {
	session := sessions.Default(c)
	//video uri
	vuri := ""
	//image uri
	iuri := ""

	err := c.Request.ParseMultipartForm(32 << 24) // ~500MB
	if err != nil {
		app.Error(c, err)
		return
	}

	vmpartFile, vmpartHeader, err := c.Request.FormFile("videofile")
	if err == nil {
		defer func() {
			err := vmpartFile.Close()
			if err != nil {
				app.Logger.Errorf("Error closing videofile multi-part %v", err)
			}
		}()

		vuri, err = app.saveFile(vmpartHeader, vmpartFile)
		if err != nil {
			app.Error(c, err)
			return
		}
	}

	impartFile, impartHeader, err := c.Request.FormFile("imagefile")
	if err == nil {
		defer func() {
			err := impartFile.Close()
			if err != nil {
				app.Logger.Errorf("Error closing imagefile multi-part %v", err)
			}
		}()

		iuri, err = app.saveFile(impartHeader, impartFile)
		if err != nil {
			app.Error(c, err)
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
		if err := app.ExercisesRepo.Update(exercise); err != nil {
			session.AddFlash(err.Error())
			_ = session.Save()
			c.Redirect(303, c.Request.RequestURI)
			return
		}
		c.Redirect(303, "/admin/exercises")
	} else {
		session.AddFlash("Ошибка! Проверьте внимательно заполнение всех полей!")
		_ = session.Save()
		c.Redirect(303, c.Request.RequestURI)
	}
}

// ExerciseAdminDelete handles /admin/delete_exercise route
func (app *Application) ExerciseAdminDelete(c *gin.Context) {
	id := helpers.Atouint(c.Request.PostFormValue("id"))

	if err := app.ExercisesRepo.Delete(id); err != nil {
		app.Error(c, err)
		return
	}
	c.Redirect(303, "/admin/exercises")
}
