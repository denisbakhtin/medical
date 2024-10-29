package controllers

import (
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/denisbakhtin/medical/config"
)

// createTokenFromId creates secure token for id
func createTokenFromID(ID uint) string {
	digest := sha1.New().Sum([]byte(fmt.Sprintf("%d-%s", ID, config.GetConfig().Salt)))
	return base64.URLEncoding.EncodeToString(
		[]byte(fmt.Sprintf("%d-%x", ID, digest)),
	)
}

// getIDFromToken deciphers token and returns review ID. Returns empty string if error
func getIDFromToken(token string) string {
	idDigest, err := base64.URLEncoding.DecodeString(token)
	if err != nil {
		return ""
	}
	if sl := strings.Split(string(idDigest), "-"); len(sl) == 2 {
		digest := sha1.New().Sum([]byte(fmt.Sprintf("%s-%s", sl[0], config.GetConfig().Salt)))
		if fmt.Sprintf("%x", digest) == sl[1] {
			return sl[0]
		}
	}
	return ""
}
