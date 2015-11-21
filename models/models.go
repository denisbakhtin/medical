package models

import (
	"github.com/fiam/gounidecode/unidecode"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" //postgresql driver, don't remove
	"regexp"
	"strings"
)

var db *sqlx.DB

//InitDB establishes connection to database and saves its handler into db *sqlx.DB
func InitDB(connection string) {
	db = sqlx.MustConnect("postgres", connection)
}

//GetDB returns database handler
func GetDB() *sqlx.DB {
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
