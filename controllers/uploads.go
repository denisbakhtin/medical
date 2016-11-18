package controllers

import (
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/denisbakhtin/medical/system"
)

//CkUpload handles POST /admin/ckupload route
func CkUpload(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {

		err := r.ParseMultipartForm(32 << 20)
		if err != nil {
			log.Printf("ERROR: %s\n", err)
			http.Error(w, err.Error(), 500)
			return
		}
		mpartFile, mpartHeader, err := r.FormFile("upload")
		if err != nil {
			log.Println(err.Error())
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		defer mpartFile.Close()
		uri, err := saveFile(mpartHeader, mpartFile)
		if err != nil {
			log.Println(err.Error())
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		CKEdFunc := r.FormValue("CKEditorFuncNum")
		fmt.Fprintln(w, "<script>window.parent.CKEDITOR.tools.callFunction("+CKEdFunc+", \""+uri+"\");</script>")

	} else {
		err := fmt.Errorf("Method %q not allowed", r.Method)
		log.Printf("ERROR: %s\n", err)
		http.Error(w, err.Error(), 405)
	}
}

//saveFile saves file to disc and returns its relative uri
func saveFile(fh *multipart.FileHeader, f multipart.File) (string, error) {
	fileExt := strings.ToLower(filepath.Ext(fh.Filename))
	if !regexp.MustCompile("^\\.(jpe?g|bmp|gif|png)$").MatchString(fileExt) {
		return "", fmt.Errorf("File is not an image")
	}
	newName := fmt.Sprint(time.Now().Unix()) + fileExt //unique file name ;D
	uri := "/public/uploads/" + newName
	fullName := filepath.Join(system.GetConfig().Uploads, newName)

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
