package models

import (
	"log"
	"regexp"
	"strings"

	"github.com/fiam/gounidecode/unidecode"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

var db *gorm.DB

//InitDB establishes connection to database and saves its handler into db *sqlx.DB
func InitDB(connection string) {
	var err error
	db, err = gorm.Open("postgres", connection)
	if err != nil {
		log.Fatal(err)
	}
	//automigrate
	db.AutoMigrate(&Article{}, &Comment{}, &Page{}, &Review{}, &User{})
}

//GetDB returns database handler
func GetDB() *gorm.DB {
	return db
}

//utility functions

//truncate truncates string to n runes
func truncate(s string, n int) string {
	runes := []rune(s)
	if len(runes) > n {
		return string(runes[:n])
	}
	return s
}

//createSlug makes url slug out of string
func createSlug(s string) string {
	s = strings.ToLower(unidecode.Unidecode(s))                     //transliterate if it is not in english
	s = regexp.MustCompile("[^a-z0-9\\s]+").ReplaceAllString(s, "") //spaces
	s = regexp.MustCompile("\\s+").ReplaceAllString(s, "-")         //spaces
	return s
}
