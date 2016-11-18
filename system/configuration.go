package system

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path"
)

//Configs contains application configurations for all application modes
type Configs struct {
	Debug   Config
	Release Config
	Test    Config
}

//Config contains application configuration for active application mode
type Config struct {
	Public        string `json:"public"`
	Uploads       string `json:"-"`
	Domain        string `json:"domain"`
	SessionSecret string `json:"session_secret"`
	CsrfSecret    string `json:"csrf_secret"`
	Ssl           bool   `json:"ssl"`
	SignupEnabled bool   `json:"signup_enabled"` //always set to false in release mode (config.json)
	Salt          string `json:"salt"`           //sha salt for generation of review & comment tokens
	Database      DatabaseConfig
	SMTP          SMTPConfig
}

//DatabaseConfig contains database connection info
type DatabaseConfig struct {
	Host     string
	Name     string //database name
	User     string
	Password string
}

//SMTPConfig contains smtp mailer info
type SMTPConfig struct {
	From     string //from email
	To       string //to email
	Cc       string //cc email
	SMTP     string //smtp server address
	Port     string //smtp port
	User     string //smtp user login
	Password string //smtp user password
}

var (
	config *Config
)

//loadConfig unmarshals config for current application mode
func loadConfig() {
	data, err := ioutil.ReadFile("./config/config.json")
	if err != nil {
		panic(err)
	}
	configs := &Configs{}
	if err := json.Unmarshal(data, configs); err != nil {
		panic(err)
	}
	switch GetMode() {
	case DebugMode:
		config = &configs.Debug
	case ReleaseMode:
		config = &configs.Release
	case TestMode:
		config = &configs.Test
	}
	if !path.IsAbs(config.Public) {
		workingDir, err := os.Getwd()
		if err != nil {
			panic(err)
		}
		config.Public = path.Join(workingDir, config.Public)
	}
	config.Uploads = path.Join(config.Public, "uploads")
}

//GetConfig returns actual config
func GetConfig() *Config {
	return config
}
