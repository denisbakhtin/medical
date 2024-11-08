package system

import (
	"fmt"
	"io"
	"os"

	"github.com/denisbakhtin/medical/config"
	"github.com/denisbakhtin/medical/models"
	"github.com/denisbakhtin/medical/views"
	"github.com/gin-gonic/gin"
)

// Current app mode for reference
var appmode string

func GetMode() string {
	return appmode
}

// Init initializes core system elements (DB, sessions, templates, et al)
func Init(mode string) {
	//store mode
	appmode = mode

	conf := config.LoadConfig(mode)
	views.Load()
	connection := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable",
		conf.Database.Host,
		conf.Database.User,
		conf.Database.Password,
		conf.Database.Name)
	models.InitDB(connection)
	setupGin(mode)
}

func setupGin(mode string) {
	gin.SetMode(mode)
	gin.DisableConsoleColor()
	f, _ := os.Create("logs/gin.txt")
	gin.DefaultWriter = io.MultiWriter(f, os.Stdout)
}
