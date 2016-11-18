package controllers

import (
	"bytes"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/denisbakhtin/medical/helpers"
	"github.com/denisbakhtin/medical/models"
	"github.com/denisbakhtin/medical/system"
	"gopkg.in/gomail.v2"
)

//ReviewShow handles /reviews/:id route
func ReviewShow(w http.ResponseWriter, r *http.Request) {
	tmpl := helpers.Template(r)
	data := helpers.DefaultData(r)
	db := models.GetDB()
	if r.Method == "GET" {

		id := r.URL.Path[len("/reviews/"):]
		review := &models.Review{}
		db.First(review, id)
		if review.ID == 0 || !review.Published {
			w.WriteHeader(404)
			tmpl.Lookup("errors/404").Execute(w, nil)
			return
		}
		data["Review"] = review
		data["Title"] = "Отзыв о работе кинезиолога" + ". " + review.AuthorName
		data["Active"] = "/reviews"
		data["MetaDescription"] = review.MetaDescription
		data["MetaKeywords"] = review.MetaKeywords
		tmpl.Lookup("reviews/show").Execute(w, data)

	} else {
		err := fmt.Errorf("Method %q not allowed", r.Method)
		log.Printf("ERROR: %s\n", err)
		w.WriteHeader(405)
		tmpl.Lookup("errors/405").Execute(w, helpers.ErrorData(err))
	}
}

//ReviewPublicIndex handles GET /reviews route
func ReviewPublicIndex(w http.ResponseWriter, r *http.Request) {
	tmpl := helpers.Template(r)
	data := helpers.DefaultData(r)
	session := helpers.Session(r)
	db := models.GetDB()
	if r.Method == "GET" {

		var list []models.Review
		db.Where("published = ?", true).Order("id desc").Find(&list)
		data["Title"] = "Кинезиология - отзывы пациентов"
		data["Active"] = r.RequestURI
		data["List"] = list
		data["Flash"] = session.Flashes("reviews")
		session.Save(r, w)
		tmpl.Lookup("reviews/public-index").Execute(w, data)

	} else {
		err := fmt.Errorf("Method %q not allowed", r.Method)
		log.Printf("ERROR: %s\n", err)
		w.WriteHeader(405)
		tmpl.Lookup("errors/405").Execute(w, helpers.ErrorData(err))
	}
}

//ReviewIndex handles GET /admin/reviews route
func ReviewIndex(w http.ResponseWriter, r *http.Request) {
	tmpl := helpers.Template(r)
	data := helpers.DefaultData(r)
	db := models.GetDB()
	if r.Method == "GET" {

		var list []models.Review
		db.Order("id desc").Find(&list)
		data["Title"] = "Отзывы"
		data["Active"] = "reviews"
		data["List"] = list
		tmpl.Lookup("reviews/index").Execute(w, data)

	} else {
		err := fmt.Errorf("Method %q not allowed", r.Method)
		log.Printf("ERROR: %s\n", err)
		w.WriteHeader(405)
		tmpl.Lookup("errors/405").Execute(w, helpers.ErrorData(err))
	}
}

//ReviewPublicCreate handles /new_review route
func ReviewPublicCreate(w http.ResponseWriter, r *http.Request) {
	session := helpers.Session(r)
	tmpl := helpers.Template(r)
	data := helpers.DefaultData(r)
	db := models.GetDB()
	if r.Method == "GET" {

		data["Title"] = "Новый отзыв"
		data["Active"] = "reviews"
		data["Flash"] = session.Flashes("reviews")
		session.Save(r, w)
		tmpl.Lookup("reviews/public-form").Execute(w, data)

	} else if r.Method == "POST" {

		r.ParseMultipartForm(32 << 20)
		//simple captcha check
		captcha, err := base64.StdEncoding.DecodeString(r.FormValue("captcha"))
		if err != nil {
			w.WriteHeader(500)
			tmpl.Lookup("errors/500").Execute(w, helpers.ErrorData(err))
			return
		}
		if string(captcha) != "100.00" {
			w.WriteHeader(400)
			tmpl.Lookup("errors/400").Execute(w, nil)
			return
		}

		review := &models.Review{
			AuthorName:  r.FormValue("author_name"),
			AuthorEmail: r.FormValue("author_email"),
			Content:     r.FormValue("content"),
			Published:   false, //reviews are published by admin via dashboard
		}

		if mpartFile, mpartHeader, err := r.FormFile("image"); err == nil {
			defer mpartFile.Close()
			review.Image, err = saveFile(mpartHeader, mpartFile)
			if err != nil {
				log.Println(err.Error())
				w.WriteHeader(500)
				tmpl.Lookup("errors/500").Execute(w, helpers.ErrorData(err))
				return
			}
		}

		if err := db.Create(review).Error; err != nil {
			log.Printf("ERROR: %s\n", err)
			w.WriteHeader(400)
			tmpl.Lookup("errors/400").Execute(w, helpers.ErrorData(err))
			return
		}
		notifyAdminOfReview(r, review)
		session.AddFlash("Спасибо! Ваш отзыв будет опубликован после проверки.", "reviews")
		session.Save(r, w)
		http.Redirect(w, r, "/reviews", 303)

	} else {
		err := fmt.Errorf("Method %q not allowed", r.Method)
		log.Printf("ERROR: %s\n", err)
		w.WriteHeader(405)
		tmpl.Lookup("errors/405").Execute(w, helpers.ErrorData(err))
	}
}

//ReviewCreate handles /admin/new_review route
func ReviewCreate(w http.ResponseWriter, r *http.Request) {
	session := helpers.Session(r)
	tmpl := helpers.Template(r)
	data := helpers.DefaultData(r)
	db := models.GetDB()
	if r.Method == "GET" {

		var articles []models.Article
		db.Where("published = ?", true).Find(&articles)
		data["Title"] = "Новый отзыв"
		data["Active"] = "reviews"
		data["Articles"] = articles
		data["Flash"] = session.Flashes("reviews")
		session.Save(r, w)
		tmpl.Lookup("reviews/form").Execute(w, data)

	} else if r.Method == "POST" {

		r.ParseMultipartForm(32 << 20)
		review := &models.Review{
			ArticleID:       helpers.Atouintr(r.FormValue("article_id")),
			AuthorName:      r.FormValue("author_name"),
			AuthorEmail:     r.FormValue("author_email"),
			Content:         r.FormValue("content"),
			Published:       helpers.Atob(r.FormValue("published")),
			Video:           r.FormValue("video"),
			MetaDescription: r.FormValue("meta_description"),
			MetaKeywords:    r.FormValue("meta_keywords"),
		}

		if mpartFile, mpartHeader, err := r.FormFile("image"); err == nil {
			defer mpartFile.Close()
			review.Image, err = saveFile(mpartHeader, mpartFile)
			if err != nil {
				log.Println(err.Error())
				w.WriteHeader(500)
				tmpl.Lookup("errors/500").Execute(w, helpers.ErrorData(err))
				return
			}
		}

		if err := db.Create(review).Error; err != nil {
			log.Printf("ERROR: %s\n", err)
			w.WriteHeader(400)
			tmpl.Lookup("errors/400").Execute(w, helpers.ErrorData(err))
			return
		}
		http.Redirect(w, r, "/admin/reviews", 303)

	} else {
		err := fmt.Errorf("Method %q not allowed", r.Method)
		log.Printf("ERROR: %s\n", err)
		w.WriteHeader(405)
		tmpl.Lookup("errors/405").Execute(w, helpers.ErrorData(err))
	}
}

//ReviewUpdate handles /admin/edit_review/:id route
func ReviewUpdate(w http.ResponseWriter, r *http.Request) {
	tmpl := helpers.Template(r)
	session := helpers.Session(r)
	data := helpers.DefaultData(r)
	db := models.GetDB()
	if r.Method == "GET" {

		id := r.URL.Path[len("/admin/edit_review/"):]
		review := &models.Review{}
		db.First(review, id)
		if review.ID == 0 {
			w.WriteHeader(404)
			tmpl.Lookup("errors/404").Execute(w, nil)
			return
		}

		var articles []models.Article
		db.Where("published = ?", true).Find(&articles)
		data["Title"] = "Редактировать отзыв"
		data["Active"] = "reviews"
		data["Review"] = review
		data["Articles"] = articles
		data["Flash"] = session.Flashes("reviews")
		session.Save(r, w)
		tmpl.Lookup("reviews/form").Execute(w, data)

	} else if r.Method == "POST" {

		r.ParseMultipartForm(32 << 20)
		rev := &models.Review{}
		db.First(rev, r.FormValue("id"))
		if rev.ID == 0 {
			w.WriteHeader(400)
			tmpl.Lookup("errors/400").Execute(w, nil)
			return
		}

		review := &models.Review{
			ID:              rev.ID,
			ArticleID:       helpers.Atouintr(r.FormValue("article_id")),
			AuthorName:      r.FormValue("author_name"),
			AuthorEmail:     r.FormValue("author_email"),
			Content:         r.FormValue("content"),
			Image:           rev.Image,
			Published:       helpers.Atob(r.FormValue("published")),
			Video:           r.FormValue("video"),
			MetaDescription: r.FormValue("meta_description"),
			MetaKeywords:    r.FormValue("meta_keywords"),
		}
		if mpartFile, mpartHeader, err := r.FormFile("image"); err == nil {
			defer mpartFile.Close()
			review.Image, err = saveFile(mpartHeader, mpartFile)
			if err != nil {
				log.Println(err.Error())
				w.WriteHeader(500)
				tmpl.Lookup("errors/500").Execute(w, helpers.ErrorData(err))
				return
			}
		}

		if err := db.Save(review).Error; err != nil {
			session.AddFlash(err.Error(), "reviews")
			session.Save(r, w)
			http.Redirect(w, r, r.RequestURI, 303)
			return
		}
		http.Redirect(w, r, "/admin/reviews", 303)

	} else {
		err := fmt.Errorf("Method %q not allowed", r.Method)
		log.Printf("ERROR: %s\n", err)
		w.WriteHeader(405)
		tmpl.Lookup("errors/405").Execute(w, helpers.ErrorData(err))
	}
}

//ReviewPublicUpdate handles /edit_review?token=:secure_token route
func ReviewPublicUpdate(w http.ResponseWriter, r *http.Request) {
	tmpl := helpers.Template(r)
	session := helpers.Session(r)
	data := helpers.DefaultData(r)
	db := models.GetDB()
	if r.Method == "GET" {

		id := getIDFromToken(r.FormValue("token"))
		review := &models.Review{}
		db.First(review, id)
		if review.ID == 0 || review.Published {
			err := fmt.Errorf("Отзыв не найден или уже был опубликован и не подлежит редактированию")
			w.WriteHeader(404)
			tmpl.Lookup("errors/404").Execute(w, helpers.ErrorData(err))
			return
		}

		var articles []models.Article
		db.Where("published = ?", true).Find(&articles)
		review.Published = true //set default to true
		data["Title"] = "Редактировать отзыв"
		data["Articles"] = articles
		data["Active"] = "reviews"
		data["Review"] = review
		data["SecureEdit"] = true
		data["Flash"] = session.Flashes("reviews")
		session.Save(r, w)
		tmpl.Lookup("reviews/public-form").Execute(w, data)

	} else if r.Method == "POST" {

		r.ParseMultipartForm(32 << 20)
		rev := &models.Review{}
		db.First(rev, r.FormValue("id"))
		if rev.ID == 0 {
			w.WriteHeader(400)
			tmpl.Lookup("errors/400").Execute(w, nil)
			return
		}

		review := &models.Review{
			ID:          rev.ID,
			ArticleID:   helpers.Atouintr(r.FormValue("article_id")),
			AuthorName:  r.FormValue("author_name"),
			AuthorEmail: r.FormValue("author_email"),
			Content:     r.FormValue("content"),
			Image:       rev.Image,
			Published:   helpers.Atob(r.FormValue("published")),
			Video:       rev.Video,
		}
		if mpartFile, mpartHeader, err := r.FormFile("image"); err == nil {
			defer mpartFile.Close()
			review.Image, err = saveFile(mpartHeader, mpartFile)
			if err != nil {
				log.Println(err.Error())
				w.WriteHeader(500)
				tmpl.Lookup("errors/500").Execute(w, helpers.ErrorData(err))
				return
			}
		}

		if err := db.Save(review).Error; err != nil {
			session.AddFlash(err.Error(), "reviews")
			session.Save(r, w)
			http.Redirect(w, r, r.RequestURI, 303)
			return
		}
		session.AddFlash("Отзыв был успешно сохранен", "reviews")
		session.Save(r, w)
		http.Redirect(w, r, "/reviews", 303)

	} else {
		err := fmt.Errorf("Method %q not allowed", r.Method)
		log.Printf("ERROR: %s\n", err)
		w.WriteHeader(405)
		tmpl.Lookup("errors/405").Execute(w, helpers.ErrorData(err))
	}
}

//ReviewDelete handles /admin/delete_review route
func ReviewDelete(w http.ResponseWriter, r *http.Request) {
	tmpl := helpers.Template(r)
	db := models.GetDB()

	if r.Method == "POST" {

		review := &models.Review{}
		db.First(review, r.PostFormValue("id"))
		if review.ID == 0 {
			w.WriteHeader(404)
			tmpl.Lookup("errors/404").Execute(w, nil)
		}

		if err := db.Delete(review).Error; err != nil {
			log.Printf("ERROR: %s\n", err)
			w.WriteHeader(500)
			tmpl.Lookup("errors/500").Execute(w, helpers.ErrorData(err))
			return
		}
		http.Redirect(w, r, "/admin/reviews", 303)

	} else {
		err := fmt.Errorf("Method %q not allowed", r.Method)
		log.Printf("ERROR: %s\n", err)
		w.WriteHeader(405)
		tmpl.Lookup("errors/405").Execute(w, helpers.ErrorData(err))
	}
}

func notifyAdminOfReview(r *http.Request, review *models.Review) {
	//closure is needed here, as r may be released by the time func finishes
	tmpl := helpers.Template(r)
	go func() {
		data := map[string]interface{}{
			"Review": review,
			"Token":  createTokenFromID(review.ID),
		}
		var b bytes.Buffer
		if err := tmpl.Lookup("emails/review").Execute(&b, data); err != nil {
			log.Printf("ERROR: %s\n", err)
			return
		}

		smtp := system.GetConfig().SMTP
		msg := gomail.NewMessage()
		msg.SetHeader("From", smtp.From)
		msg.SetHeader("To", smtp.To)
		if len(smtp.Cc) > 0 {
			msg.SetHeader("Cc", smtp.Cc)
		}
		msg.SetHeader("Subject", fmt.Sprintf("Новый отзыв на сайте www.miobalans.ru: %s", review.AuthorName))
		msg.SetBody(
			"text/html",
			b.String(),
		)

		port, _ := strconv.Atoi(smtp.Port)
		dialer := gomail.NewPlainDialer(smtp.SMTP, port, smtp.User, smtp.Password)
		sender, err := dialer.Dial()
		if err != nil {
			log.Printf("ERROR: %s\n", err)
			return
		}
		if err := gomail.Send(sender, msg); err != nil {
			log.Printf("ERROR: %s\n", err)
			return
		}
	}()
}

//createTokenFromId creates secure token for id
func createTokenFromID(ID uint) string {
	digest := sha1.New().Sum([]byte(fmt.Sprintf("%d-%s", ID, system.GetConfig().Salt)))
	return base64.URLEncoding.EncodeToString(
		[]byte(fmt.Sprintf("%d-%x", ID, digest)),
	)
}

//getIDFromToken deciphers token and returns review ID. Returns empty string if error
func getIDFromToken(token string) string {
	idDigest, err := base64.URLEncoding.DecodeString(token)
	if err != nil {
		return ""
	}
	if sl := strings.Split(string(idDigest), "-"); len(sl) == 2 {
		digest := sha1.New().Sum([]byte(fmt.Sprintf("%s-%s", sl[0], system.GetConfig().Salt)))
		if fmt.Sprintf("%x", digest) == sl[1] {
			return sl[0]
		}
	}
	return ""
}
