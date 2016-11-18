package system

import (
	"fmt"
	"log"

	"github.com/denisbakhtin/medical/models"
)

//Application mode
const (
	DebugMode   = "debug"
	ReleaseMode = "release"
	TestMode    = "test"
)

var mode string //application mode: debug, release, test

//Init initializes core system elements (DB, sessions, templates, et al)
func Init() {
	loadConfig()
	loadI18n()
	createSession()
	loadTemplates()
	connection := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable", config.Database.Host, config.Database.User, config.Database.Password, config.Database.Name)
	models.InitDB(connection)
}

//SetMode sets application mode
func SetMode(flagmode *string) {
	switch *flagmode {
	case DebugMode:
		mode = DebugMode
	case ReleaseMode:
		mode = ReleaseMode
	case TestMode:
		mode = TestMode
	default:
		log.Fatalf("Unknown application mode: %q\n", *flagmode)
	}
}

//GetMode returns application mode
func GetMode() string {
	return mode
}
