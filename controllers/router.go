package controllers

import (
	"net/http"
	"path"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

// SetupRouter establishes web routes
func (app *Application) SetupRouter() *gin.Engine {
	router := gin.Default()

	store := cookie.NewStore([]byte(app.Config.SessionSecret))
	router.Use(sessions.Sessions("gin-session", store))

	router.SetHTMLTemplate(app.Template)

	// Static files
	router.StaticFS("/public", http.Dir(app.Config.Public))
	router.StaticFile("/robots.txt", path.Join(app.Config.Public, "robots.txt"))

	// Routes
	router.NoRoute(app.Error404)
	router.GET("/", app.Home)

	router.GET("/signin", app.SignInGet)
	router.POST("/signin", app.SignInPost)
	router.GET("/logout", app.LogOut)
	if app.Config.SignupEnabled {
		router.GET("/signup", app.SignUpGet)
		router.POST("/signup", app.SignUpPost)
	}

	router.GET("/pages/:idslug", app.PageShow)
	router.GET("/exercises", app.ExercisesIndex)
	router.GET("/exercises/:idslug", app.ExerciseShow)
	router.GET("/articles", app.ArticlesIndex)
	router.GET("/articles/:idslug", app.ArticleShow)
	router.GET("/comments/:id", app.CommentsIndex)
	router.GET("/info/:idslug", app.InfoShow)
	router.GET("/reviews", app.ReviewsIndex)
	router.GET("/reviews/:idslug", app.ReviewShow)
	router.POST("/new_request", app.RequestCreatePost)
	router.POST("/new_comment", app.CommentCreatePost)
	// http.Handle("/edit_comment", Default(CommentPublicUpdate))
	router.GET("/new_review", app.ReviewCreateGet)
	router.POST("/new_review", app.ReviewCreatePost)
	router.GET("/edit_review", app.ReviewUpdateGet)
	router.POST("/edit_review", app.ReviewUpdatePost)

	authorized := router.Group("/admin", app.FilterAuthenticated())
	{
		authorized.GET("/", app.Dashboard)
		authorized.GET("/users", app.UsersAdminIndex)
		authorized.GET("/new_user", app.UserAdminCreateGet)
		authorized.POST("/new_user", app.UserAdminCreatePost)
		authorized.GET("/edit_user/:id", app.UserAdminUpdateGet)
		authorized.POST("/edit_user/:id", app.UserAdminUpdatePost)
		authorized.POST("/delete_user", app.UserAdminDelete)

		authorized.GET("/pages", app.PagesAdminIndex)
		authorized.GET("/new_page", app.PageAdminCreateGet)
		authorized.POST("/new_page", app.PageAdminCreatePost)
		authorized.GET("/edit_page/:id", app.PageAdminUpdateGet)
		authorized.POST("/edit_page/:id", app.PageAdminUpdatePost)
		authorized.POST("/delete_page", app.PageAdminDelete)

		authorized.GET("/exercises", app.ExercisesAdminIndex)
		authorized.GET("/new_exercise", app.ExerciseAdminCreateGet)
		authorized.POST("/new_exercise", app.ExerciseAdminCreatePost)
		authorized.GET("/edit_exercise/:id", app.ExerciseAdminUpdateGet)
		authorized.POST("/edit_exercise/:id", app.ExerciseAdminUpdatePost)
		authorized.POST("/delete_exercise", app.ExerciseAdminDelete)

		authorized.GET("/articles", app.ArticlesAdminIndex)
		authorized.GET("/new_article", app.ArticleAdminCreateGet)
		authorized.POST("/new_article", app.ArticleAdminCreatePost)
		authorized.GET("/edit_article/:id", app.ArticleAdminUpdateGet)
		authorized.POST("/edit_article/:id", app.ArticleAdminUpdatePost)
		authorized.POST("/delete_article", app.ArticleAdminDelete)

		authorized.GET("/info", app.InfoAdminIndex)
		authorized.GET("/new_info", app.InfoAdminCreateGet)
		authorized.POST("/new_info", app.InfoAdminCreatePost)
		authorized.GET("/edit_info/:id", app.InfoAdminUpdateGet)
		authorized.POST("/edit_info/:id", app.InfoAdminUpdatePost)
		authorized.POST("/delete_info", app.InfoAdminDelete)

		authorized.GET("/comments", app.CommentsAdminIndex)
		authorized.GET("/edit_comment/:id", app.CommentAdminUpdateGet)
		authorized.POST("/edit_comment/:id", app.CommentAdminUpdatePost)
		authorized.POST("/delete_comment", app.CommentAdminDelete)

		authorized.GET("/reviews", app.ReviewsAdminIndex)
		authorized.GET("/new_review", app.ReviewAdminCreateGet)
		authorized.POST("/new_review", app.ReviewAdminCreatePost)
		authorized.GET("/edit_review/:id", app.ReviewAdminUpdateGet)
		authorized.POST("/edit_review/:id", app.ReviewAdminUpdatePost)
		authorized.POST("/delete_review", app.ReviewAdminDelete)

		authorized.POST("/ckupload", app.CkUpload)
	}
	return router
}
