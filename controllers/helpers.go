package controllers

import (
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/gin-contrib/sessions"
)

// createTokenFromId creates secure token for id
func (app *Application) createTokenFromID(ID uint) string {
	digest := sha1.New().Sum([]byte(fmt.Sprintf("%d-%s", ID, app.Config.Salt)))
	return base64.URLEncoding.EncodeToString(
		[]byte(fmt.Sprintf("%d-%x", ID, digest)),
	)
}

// getIDFromToken deciphers token and returns review ID. Returns empty string if error
func (app *Application) getIDFromToken(token string) string {
	idDigest, err := base64.URLEncoding.DecodeString(token)
	if err != nil {
		return ""
	}
	if sl := strings.Split(string(idDigest), "-"); len(sl) == 2 {
		digest := sha1.New().Sum([]byte(fmt.Sprintf("%s-%s", sl[0], app.Config.Salt)))
		if fmt.Sprintf("%x", digest) == sl[1] {
			return sl[0]
		}
	}
	return ""
}

// authenticated checks if user is authenticated
func (app *Application) authenticated(session sessions.Session) bool {
	return session.Get(userIDkey) != nil
}

// fullURL transforms relative url into full with schema and domain name
func (app *Application) fullURL(relativeURL string) string {
	return fmt.Sprintf("%s%s", app.Config.FullDomain, relativeURL)
}
