package controllers

import (
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// CkUpload handles POST /admin/ckupload route
func (app *Application) CkUpload(c *gin.Context) {
	err := c.Request.ParseMultipartForm(32 << 20)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	mpartFile, mpartHeader, err := c.Request.FormFile("upload")
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	defer mpartFile.Close()
	uri, err := app.saveFile(mpartHeader, mpartFile)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"uploaded": true, "url": uri})
}

// saveFile saves file to disc and returns its relative uri
func (app *Application) saveFile(fh *multipart.FileHeader, f multipart.File) (string, error) {
	fileExt := strings.ToLower(filepath.Ext(fh.Filename))
	if !regexp.MustCompile(`^\.(jpe?g|bmp|gif|png|mp4)$`).MatchString(fileExt) {
		return "", errors.New("file is not an image or .mp4 video")
	}
	newName := fmt.Sprint(time.Now().Unix()) + fileExt // unique file name ;D
	uri := "/public/uploads/" + newName
	fullName := filepath.Join(app.Config.Uploads, newName)

	file, err := os.OpenFile(fullName, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return "", err
	}
	defer file.Close()
	_, err = io.Copy(file, f)
	if err != nil {
		return "", err
	}
	return uri, nil
}
