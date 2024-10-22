package system

import (
	"fmt"
	"io"
	"os"

	"github.com/denisbakhtin/medical/models"
	"github.com/gin-gonic/gin"
)

// Application mode
const (
	DebugMode   = "debug"
	ReleaseMode = "release"
	TestMode    = "test"
)

// Init initializes core system elements (DB, sessions, templates, et al)
func Init(mode string) {
	loadConfig(mode)
	loadTemplates()
	connection := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable", config.Database.Host, config.Database.User, config.Database.Password, config.Database.Name)
	models.InitDB(connection)
	setupGin(mode)
}

func setupGin(mode string) {
	gin.SetMode(mode)
	gin.DisableConsoleColor()
	f, _ := os.Create("logs/gin.txt")
	gin.DefaultWriter = io.MultiWriter(f, os.Stdout)
}
