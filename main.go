package main

import (
	"flag"
	"log"
	"net/http"
	"path"

	"github.com/claudiu/gocron"
	"github.com/denisbakhtin/medical/controllers"
	"github.com/denisbakhtin/medical/system"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

func init() {
	log.SetFlags(log.Lshortfile)
}

func main() {
	mode := flag.String("mode", "debug", "Application mode: debug, release, test")
	flag.Parse()

	system.SetMode(mode)
	system.Init()

	//Periodic tasks
	if system.GetMode() == system.ReleaseMode {
		system.CreateXMLSitemap() //refresh sitemap now
	}
	gocron.Every(1).Day().Do(system.CreateXMLSitemap) //refresh daily
	gocron.Start()

	gin.SetMode(system.GetMode())
	router := gin.Default()
	store := cookie.NewStore([]byte(system.GetConfig().SessionSecret))
	router.Use(sessions.Sessions("gin-session", store))
	router.SetHTMLTemplate(system.GetTemplates())
	router.NoRoute(controllers.Error404)
	router.StaticFS("/public", http.Dir("public"))
	router.StaticFile("/robots.txt", path.Join(system.GetConfig().Public, "robots.txt"))
	router.GET("/", controllers.Home)
	router.GET("/signin", controllers.SignInGet)
	router.POST("/signin", controllers.SignInPost)
	router.GET("/logout", controllers.LogOut)
	if system.GetConfig().SignupEnabled {
		router.GET("/signup", controllers.SignUpGet)
		router.POST("/signup", controllers.SignUpPost)
	}

	router.GET("/pages/:idslug", controllers.PageShow)
	router.GET("/exercises", controllers.ExercisesIndex)
	router.GET("/exercises/:idslug", controllers.ExerciseShow)
	router.GET("/articles", controllers.ArticlesIndex)
	router.GET("/articles/:idslug", controllers.ArticleShow)
	router.GET("/comments/:id", controllers.CommentsIndex)
	router.GET("/info/:idslug", controllers.InfoShow)
	router.GET("/reviews", controllers.ReviewsIndex)
	router.GET("/reviews/:idslug", controllers.ReviewShow)
	router.POST("/new_request", controllers.RequestCreatePost)
	router.POST("/new_comment", controllers.CommentCreatePost)
	//http.Handle("/edit_comment", Default(controllers.CommentPublicUpdate))
	router.GET("/new_review", controllers.ReviewCreateGet)
	router.POST("/new_review", controllers.ReviewCreatePost)
	router.GET("/edit_review", controllers.ReviewUpdateGet)
	router.POST("/edit_review", controllers.ReviewUpdatePost)

	authorized := router.Group("/admin", controllers.Authenticated())
	{
		authorized.GET("/", controllers.Dashboard)
		authorized.GET("/users", controllers.UsersAdminIndex)
		authorized.GET("/new_user", controllers.UserAdminCreateGet)
		authorized.POST("/new_user", controllers.UserAdminCreatePost)
		authorized.GET("/edit_user/:id", controllers.UserAdminUpdateGet)
		authorized.POST("/edit_user/:id", controllers.UserAdminUpdatePost)
		authorized.POST("/delete_user", controllers.UserAdminDelete)

		authorized.GET("/pages", controllers.PagesAdminIndex)
		authorized.GET("/new_page", controllers.PageAdminCreateGet)
		authorized.POST("/new_page", controllers.PageAdminCreatePost)
		authorized.GET("/edit_page/:id", controllers.PageAdminUpdateGet)
		authorized.POST("/edit_page/:id", controllers.PageAdminUpdatePost)
		authorized.POST("/delete_page", controllers.PageAdminDelete)

		authorized.GET("/exercises", controllers.ExercisesAdminIndex)
		authorized.GET("/new_exercise", controllers.ExerciseAdminCreateGet)
		authorized.POST("/new_exercise", controllers.ExerciseAdminCreatePost)
		authorized.GET("/edit_exercise/:id", controllers.ExerciseAdminUpdateGet)
		authorized.POST("/edit_exercise/:id", controllers.ExerciseAdminUpdatePost)
		authorized.POST("/delete_exercise", controllers.ExerciseAdminDelete)

		authorized.GET("/articles", controllers.ArticlesAdminIndex)
		authorized.GET("/new_article", controllers.ArticleAdminCreateGet)
		authorized.POST("/new_article", controllers.ArticleAdminCreatePost)
		authorized.GET("/edit_article/:id", controllers.ArticleAdminUpdateGet)
		authorized.POST("/edit_article/:id", controllers.ArticleAdminUpdatePost)
		authorized.POST("/delete_article", controllers.ArticleAdminDelete)

		authorized.GET("/info", controllers.InfoAdminIndex)
		authorized.GET("/new_info", controllers.InfoAdminCreateGet)
		authorized.POST("/new_info", controllers.InfoAdminCreatePost)
		authorized.GET("/edit_info/:id", controllers.InfoAdminUpdateGet)
		authorized.POST("/edit_info/:id", controllers.InfoAdminUpdatePost)
		authorized.POST("/delete_info", controllers.InfoAdminDelete)

		authorized.GET("/comments", controllers.CommentsAdminIndex)
		authorized.GET("/edit_comment/:id", controllers.CommentAdminUpdateGet)
		authorized.POST("/edit_comment/:id", controllers.CommentAdminUpdatePost)
		authorized.POST("/delete_comment", controllers.CommentAdminDelete)

		authorized.GET("/reviews", controllers.ReviewsAdminIndex)
		authorized.GET("/new_review", controllers.ReviewAdminCreateGet)
		authorized.POST("/new_review", controllers.ReviewAdminCreatePost)
		authorized.GET("/edit_review/:id", controllers.ReviewAdminUpdateGet)
		authorized.POST("/edit_review/:id", controllers.ReviewAdminUpdatePost)
		authorized.POST("/delete_review", controllers.ReviewAdminDelete)

		authorized.POST("/ckupload", controllers.CkUpload)
	}

	log.Fatal(router.Run(":8010"))
}
