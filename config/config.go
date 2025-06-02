// Package config contains data structures for reading web service config.json file
package config

import (
	_ "embed"
	"encoding/json"
	"os"
	"path"
)

// Application mode
const (
	DebugMode   = "debug"
	ReleaseMode = "release"
	TestMode    = "test"
)

// Configs contains application configurations for all application modes
type Configs struct {
	Debug   Config
	Release Config
	Test    Config
}

// Config contains application configuration for active application mode
type Config struct {
	Public        string `json:"public"`
	Uploads       string `json:"-"`
	Domain        string `json:"domain"`
	FullDomain    string `json:"full_domain"`
	SessionSecret string `json:"session_secret"`
	CsrfSecret    string `json:"csrf_secret"`
	Ssl           bool   `json:"ssl"`
	SignupEnabled bool   `json:"signup_enabled"` // always set to false in release mode (config.json)
	Salt          string `json:"salt"`           // sha salt for generation of review & comment tokens
	Database      DatabaseConfig
	SMTP          SMTPConfig
}

// DatabaseConfig contains database connection info
type DatabaseConfig struct {
	Host     string
	Name     string // database name
	User     string
	Password string
}

// SMTPConfig contains smtp mailer info
type SMTPConfig struct {
	From     string // from email
	To       string // to email
	Cc       string // cc email
	SMTP     string // smtp server address
	Port     string // smtp port
	User     string // smtp user login
	Password string // smtp user password
}

//go:embed config.json
var cfg []byte

// LoadConfig unmarshals config for current application mode
func LoadConfig(mode string) *Config {
	configs := &Configs{}
	config := &Config{}
	if err := json.Unmarshal(cfg, configs); err != nil {
		panic(err)
	}
	switch mode {
	case ReleaseMode:
		config = &configs.Release
	case TestMode:
		config = &configs.Test
	default:
		config = &configs.Debug
	}
	config.Public = getPublicDir(config)
	config.Uploads = path.Join(config.Public, "uploads")
	return config
}

// getPublicDir returns absulute public dir path.
// the trick is in debug mode, I can't reliably use os.Getwd() there, because
// different IDEs use different wd, so is air watcher. Some use projectdir,
// some projectdir/cmd, so I have to figure it out below
// panics are ok here becaus this all is run at startup
func getPublicDir(config *Config) string {
	if path.IsAbs(config.Public) {
		//already ok
		return config.Public
	}
	workingDir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	testPath := path.Join(workingDir, config.Public)
	if dirExists(testPath) {
		return testPath
	}
	testPath = path.Join(workingDir, "..", config.Public)
	if dirExists(testPath) {
		return testPath
	}
	panic("Public dir not found")
}

func dirExists(dir string) bool {
	_, err := os.Stat(dir)
	return err == nil
}
