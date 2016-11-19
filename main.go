package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/claudiu/gocron"
	"github.com/denisbakhtin/medical/controllers"
	"github.com/denisbakhtin/medical/system"
	"github.com/gin-gonic/contrib/sessions"
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
	//system.CreateXMLSitemap()                         //refresh sitemap now
	gocron.Every(1).Day().Do(system.CreateXMLSitemap) //refresh daily
	gocron.Start()

	router := gin.Default()
	store := sessions.NewCookieStore([]byte(system.GetConfig().SessionSecret))
	router.Use(sessions.Sessions("gin-session", store))
	router.SetHTMLTemplate(system.GetTemplates())
	router.GET("/", controllers.Home)
	router.StaticFS("/public", http.Dir("public"))
	router.GET("/signin", controllers.SignInGet)
	router.POST("/signin", controllers.SignInPost)
	router.GET("/logout", controllers.LogOut)
	if system.GetConfig().SignupEnabled {
		router.GET("/signup", controllers.SignUpGet)
		router.POST("/signup", controllers.SignUpPost)
	}

	router.GET("/pages/:idslug", controllers.PageShow)
	router.GET("/articles", controllers.ArticlesIndex)
	router.GET("/articles/:idslug", controllers.ArticleShow)
	router.GET("/reviews", controllers.ReviewsIndex)
	router.GET("/reviews/:id", controllers.ReviewShow)
	router.POST("/new_request", controllers.RequestCreatePost)
	router.POST("/new_comment", controllers.CommentCreatePost)
	//http.Handle("/edit_comment", Default(controllers.CommentPublicUpdate))
	router.GET("/new_review", controllers.ReviewCreateGet)
	router.POST("/new_review", controllers.ReviewCreatePost)
	router.GET("/edit_review", controllers.ReviewUpdateGet)
	router.POST("/edit_review", controllers.ReviewUpdatePost)

	router.Use(system.Authenticated())
	{
		router.GET("/admin", controllers.Dashboard)
		router.GET("/admin/users", controllers.UsersAdminIndex)
		router.GET("/admin/new_user", controllers.UserAdminCreateGet)
		router.POST("/admin/new_user", controllers.UserAdminCreatePost)
		router.GET("/admin/edit_user/:id", controllers.UserAdminUpdateGet)
		router.POST("/admin/edit_user/:id", controllers.UserAdminUpdatePost)
		router.POST("/admin/delete_user", controllers.UserAdminDelete)

		router.GET("/admin/pages", controllers.PagesAdminIndex)
		router.GET("/admin/new_page", controllers.PageAdminCreateGet)
		router.POST("/admin/new_page", controllers.PageAdminCreatePost)
		router.GET("/admin/edit_page/:id", controllers.PageAdminUpdateGet)
		router.POST("/admin/edit_page/:id", controllers.PageAdminUpdatePost)
		router.POST("/admin/delete_page", controllers.PageAdminDelete)

		router.GET("/admin/articles", controllers.ArticlesAdminIndex)
		router.GET("/admin/new_article", controllers.ArticleAdminCreateGet)
		router.POST("/admin/new_article", controllers.ArticleAdminCreatePost)
		router.GET("/admin/edit_article/:id", controllers.ArticleAdminUpdateGet)
		router.POST("/admin/edit_article/:id", controllers.ArticleAdminUpdatePost)
		router.POST("/admin/delete_article", controllers.ArticleAdminDelete)

		router.GET("/admin/comments", controllers.CommentsAdminIndex)
		router.GET("/admin/edit_comment/:id", controllers.CommentAdminUpdateGet)
		router.POST("/admin/edit_comment/:id", controllers.CommentAdminUpdatePost)
		router.POST("/admin/delete_comment", controllers.CommentAdminDelete)

		router.GET("/admin/reviews", controllers.ReviewsAdminIndex)
		router.GET("/admin/new_review", controllers.ReviewAdminCreateGet)
		router.POST("/admin/new_review", controllers.ReviewAdminCreatePost)
		router.GET("/admin/edit_review/:id", controllers.ReviewAdminUpdateGet)
		router.POST("/admin/edit_review/:id", controllers.ReviewAdminUpdatePost)
		router.POST("/admin/delete_review", controllers.ReviewAdminDelete)

		router.POST("/admin/ckupload", controllers.CkUpload)
	}

	log.Fatal(router.Run(":8010"))
}
