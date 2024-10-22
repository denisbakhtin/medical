package controllers

import (
	"net/http"
	"path"

	"github.com/denisbakhtin/medical/system"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

// SetupRouter establishes web routes
func SetupRouter() (router *gin.Engine) {
	router = gin.Default()
	store := cookie.NewStore([]byte(system.GetConfig().SessionSecret))
	router.Use(sessions.Sessions("gin-session", store))
	router.SetHTMLTemplate(system.GetTemplates())
	router.NoRoute(Error404)
	router.StaticFS("/public", http.Dir("public"))
	router.StaticFile("/robots.txt", path.Join(system.GetConfig().Public, "robots.txt"))
	router.GET("/", Home)
	router.GET("/signin", SignInGet)
	router.POST("/signin", SignInPost)
	router.GET("/logout", LogOut)
	if system.GetConfig().SignupEnabled {
		router.GET("/signup", SignUpGet)
		router.POST("/signup", SignUpPost)
	}

	router.GET("/pages/:idslug", PageShow)
	router.GET("/exercises", ExercisesIndex)
	router.GET("/exercises/:idslug", ExerciseShow)
	router.GET("/articles", ArticlesIndex)
	router.GET("/articles/:idslug", ArticleShow)
	router.GET("/comments/:id", CommentsIndex)
	router.GET("/info/:idslug", InfoShow)
	router.GET("/reviews", ReviewsIndex)
	router.GET("/reviews/:idslug", ReviewShow)
	router.POST("/new_request", RequestCreatePost)
	router.POST("/new_comment", CommentCreatePost)
	// http.Handle("/edit_comment", Default(CommentPublicUpdate))
	router.GET("/new_review", ReviewCreateGet)
	router.POST("/new_review", ReviewCreatePost)
	router.GET("/edit_review", ReviewUpdateGet)
	router.POST("/edit_review", ReviewUpdatePost)

	authorized := router.Group("/admin", Authenticated())
	{
		authorized.GET("/", Dashboard)
		authorized.GET("/users", UsersAdminIndex)
		authorized.GET("/new_user", UserAdminCreateGet)
		authorized.POST("/new_user", UserAdminCreatePost)
		authorized.GET("/edit_user/:id", UserAdminUpdateGet)
		authorized.POST("/edit_user/:id", UserAdminUpdatePost)
		authorized.POST("/delete_user", UserAdminDelete)

		authorized.GET("/pages", PagesAdminIndex)
		authorized.GET("/new_page", PageAdminCreateGet)
		authorized.POST("/new_page", PageAdminCreatePost)
		authorized.GET("/edit_page/:id", PageAdminUpdateGet)
		authorized.POST("/edit_page/:id", PageAdminUpdatePost)
		authorized.POST("/delete_page", PageAdminDelete)

		authorized.GET("/exercises", ExercisesAdminIndex)
		authorized.GET("/new_exercise", ExerciseAdminCreateGet)
		authorized.POST("/new_exercise", ExerciseAdminCreatePost)
		authorized.GET("/edit_exercise/:id", ExerciseAdminUpdateGet)
		authorized.POST("/edit_exercise/:id", ExerciseAdminUpdatePost)
		authorized.POST("/delete_exercise", ExerciseAdminDelete)

		authorized.GET("/articles", ArticlesAdminIndex)
		authorized.GET("/new_article", ArticleAdminCreateGet)
		authorized.POST("/new_article", ArticleAdminCreatePost)
		authorized.GET("/edit_article/:id", ArticleAdminUpdateGet)
		authorized.POST("/edit_article/:id", ArticleAdminUpdatePost)
		authorized.POST("/delete_article", ArticleAdminDelete)

		authorized.GET("/info", InfoAdminIndex)
		authorized.GET("/new_info", InfoAdminCreateGet)
		authorized.POST("/new_info", InfoAdminCreatePost)
		authorized.GET("/edit_info/:id", InfoAdminUpdateGet)
		authorized.POST("/edit_info/:id", InfoAdminUpdatePost)
		authorized.POST("/delete_info", InfoAdminDelete)

		authorized.GET("/comments", CommentsAdminIndex)
		authorized.GET("/edit_comment/:id", CommentAdminUpdateGet)
		authorized.POST("/edit_comment/:id", CommentAdminUpdatePost)
		authorized.POST("/delete_comment", CommentAdminDelete)

		authorized.GET("/reviews", ReviewsAdminIndex)
		authorized.GET("/new_review", ReviewAdminCreateGet)
		authorized.POST("/new_review", ReviewAdminCreatePost)
		authorized.GET("/edit_review/:id", ReviewAdminUpdateGet)
		authorized.POST("/edit_review/:id", ReviewAdminUpdatePost)
		authorized.POST("/delete_review", ReviewAdminDelete)

		authorized.POST("/ckupload", CkUpload)
	}
	return router
}
