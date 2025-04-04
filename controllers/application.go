package controllers

import (
	"fmt"
	"html/template"

	"github.com/denisbakhtin/medical/config"
	"github.com/denisbakhtin/medical/models"
	"github.com/denisbakhtin/medical/repos"
	"github.com/denisbakhtin/medical/services"
	"github.com/gin-gonic/gin"
)

type Application struct {
	Config        *config.Config
	Template      *template.Template
	Emailer       services.Emailer
	Logger        services.Logger
	ArticlesRepo  repos.Articles
	ReviewsRepo   repos.Reviews
	CommentsRepo  repos.Comments
	InfosRepo     repos.Infos
	UsersRepo     repos.Users
	ExercisesRepo repos.Exercises
	PagesRepo     repos.Pages
	MenusRepo     repos.Menus
}

func NewApplication(mode string) *Application {
	conf := config.LoadConfig(mode)
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable",
		conf.Database.Host,
		conf.Database.User,
		conf.Database.Password,
		conf.Database.Name)

	//gin settings
	gin.SetMode(mode)
	gin.DisableConsoleColor()

	db := models.InitDB(dsn)
	logger := services.NewStdLogger()

	app := &Application{
		Config:        conf,
		Emailer:       services.NewGmailer(conf, logger),
		Logger:        logger,
		ArticlesRepo:  repos.NewArticlesRepo(db),
		CommentsRepo:  repos.NewCommentsRepo(db),
		ExercisesRepo: repos.NewExercisesRepo(db),
		InfosRepo:     repos.NewInfosRepo(db),
		PagesRepo:     repos.NewPagesRepo(db),
		ReviewsRepo:   repos.NewReviewsRepo(db),
		UsersRepo:     repos.NewUsersRepo(db),
		MenusRepo:     repos.NewMenusRepo(db),
	}
	//load html templates
	app.Template = app.loadTemplate()

	return app
}
