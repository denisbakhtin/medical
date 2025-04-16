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

// Application represents the core of medical web-site. All methods on this struct are used as or by http handlers.
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

// NewApplication creates a pointer to new Application struct, respecting current app mode
// It setups database handler, loads config, html templates and initializes db repositories
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
	menus := repos.NewMenusRepo(db)
	reviews := repos.NewReviewsRepo(db)

	app := &Application{
		Config:        conf,
		Emailer:       services.NewGmailer(conf, logger),
		Logger:        logger,
		ArticlesRepo:  repos.NewArticlesRepo(db),
		CommentsRepo:  repos.NewCommentsRepo(db),
		ExercisesRepo: repos.NewExercisesRepo(db),
		InfosRepo:     repos.NewInfosRepo(db),
		PagesRepo:     repos.NewPagesRepo(db),
		ReviewsRepo:   reviews,
		UsersRepo:     repos.NewUsersRepo(db),
		MenusRepo:     menus,
		Template:      loadTemplate(menus, reviews, conf),
	}

	return app
}
