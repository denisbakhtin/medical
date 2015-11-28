package controllers

import (
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"github.com/denisbakhtin/medical/helpers"
	"github.com/denisbakhtin/medical/models"
	"github.com/denisbakhtin/medical/system"
	"gopkg.in/gomail.v2"
	"log"
	"net/http"
	"strconv"
	"strings"
)

//ReviewShow handles /reviews/:id route
func ReviewShow(w http.ResponseWriter, r *http.Request) {
	tmpl := helpers.Template(r)
	data := helpers.DefaultData(r)
	T := helpers.T(r)
	if r.Method == "GET" {

		id := r.URL.Path[len("/reviews/"):]
		review, err := models.GetReview(id)
		if err != nil || !review.Published {
			w.WriteHeader(404)
			tmpl.Lookup("errors/404").Execute(w, nil)
			return
		}
		data["Review"] = review
		data["Title"] = T("testimonial_title") + ". " + review.AuthorName
		data["Active"] = "/reviews"
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
	T := helpers.T(r)
	if r.Method == "GET" {

		list, err := models.GetPublishedReviews()
		if err != nil {
			log.Printf("ERROR: %s\n", err)
			w.WriteHeader(500)
			tmpl.Lookup("errors/500").Execute(w, helpers.ErrorData(err))
			return
		}
		data["Title"] = T("testimonials_title")
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
	T := helpers.T(r)
	if r.Method == "GET" {

		list, err := models.GetReviews()
		if err != nil {
			w.WriteHeader(500)
			tmpl.Lookup("errors/500").Execute(w, helpers.ErrorData(err))
			return
		}
		data["Title"] = T("reviews")
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
	T := helpers.T(r)
	if r.Method == "GET" {

		data["Title"] = T("new_review")
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

		if err := review.Insert(); err != nil {
			log.Printf("ERROR: %s\n", err)
			w.WriteHeader(400)
			tmpl.Lookup("errors/400").Execute(w, helpers.ErrorData(err))
			return
		}
		notifyAdminOfReview(r, review)
		session.AddFlash(T("thank_you_for_posting_review"), "reviews")
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
	T := helpers.T(r)
	if r.Method == "GET" {

		data["Title"] = T("new_review")
		data["Active"] = "reviews"
		data["Flash"] = session.Flashes("reviews")
		session.Save(r, w)
		tmpl.Lookup("reviews/form").Execute(w, data)

	} else if r.Method == "POST" {

		r.ParseMultipartForm(32 << 20)
		review := &models.Review{
			AuthorName:  r.FormValue("author_name"),
			AuthorEmail: r.FormValue("author_email"),
			Content:     r.FormValue("content"),
			Published:   helpers.Atob(r.FormValue("published")),
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

		if err := review.Insert(); err != nil {
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
	T := helpers.T(r)
	if r.Method == "GET" {

		id := r.URL.Path[len("/admin/edit_review/"):]
		review, err := models.GetReview(id)
		if err != nil {
			w.WriteHeader(404)
			tmpl.Lookup("errors/404").Execute(w, helpers.ErrorData(err))
			return
		}

		data["Title"] = T("edit_review")
		data["Active"] = "reviews"
		data["Review"] = review
		data["Flash"] = session.Flashes("reviews")
		session.Save(r, w)
		tmpl.Lookup("reviews/form").Execute(w, data)

	} else if r.Method == "POST" {

		r.ParseMultipartForm(32 << 20)
		rev, _ := models.GetReview(r.FormValue("id"))
		if rev.ID == 0 {
			w.WriteHeader(400)
			tmpl.Lookup("errors/400").Execute(w, nil)
			return
		}

		review := &models.Review{
			ID:          rev.ID,
			AuthorName:  r.FormValue("author_name"),
			AuthorEmail: r.FormValue("author_email"),
			Content:     r.FormValue("content"),
			Image:       rev.Image,
			Published:   helpers.Atob(r.FormValue("published")),
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

		if err := review.Update(); err != nil {
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
	T := helpers.T(r)
	if r.Method == "GET" {

		id := getIDFromToken(r.FormValue("token"))
		review, err := models.GetReview(id)
		if err != nil || review.Published {
			err := fmt.Errorf(T("review_not_found_or_already_published"))
			w.WriteHeader(404)
			tmpl.Lookup("errors/404").Execute(w, helpers.ErrorData(err))
			return
		}

		review.Published = true //set default to true
		data["Title"] = T("edit_review")
		data["Active"] = "reviews"
		data["Review"] = review
		data["SecureEdit"] = true
		data["Flash"] = session.Flashes("reviews")
		session.Save(r, w)
		tmpl.Lookup("reviews/public-form").Execute(w, data)

	} else if r.Method == "POST" {

		r.ParseMultipartForm(32 << 20)
		rev, _ := models.GetReview(r.FormValue("id"))
		if rev.ID == 0 {
			w.WriteHeader(400)
			tmpl.Lookup("errors/400").Execute(w, nil)
			return
		}

		review := &models.Review{
			ID:          rev.ID,
			AuthorName:  r.FormValue("author_name"),
			AuthorEmail: r.FormValue("author_email"),
			Content:     r.FormValue("content"),
			Image:       rev.Image,
			Published:   helpers.Atob(r.FormValue("published")),
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

		if err := review.Update(); err != nil {
			session.AddFlash(err.Error(), "reviews")
			session.Save(r, w)
			http.Redirect(w, r, r.RequestURI, 303)
			return
		}
		session.AddFlash(T("review_has_been_successfully_updated"), "reviews")
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

	if r.Method == "POST" {

		review, err := models.GetReview(r.PostFormValue("id"))
		if err != nil {
			log.Printf("ERROR: %s\n", err)
			w.WriteHeader(404)
			tmpl.Lookup("errors/404").Execute(w, helpers.ErrorData(err))
		}

		if err := review.Delete(); err != nil {
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
	T := helpers.T(r)
	host := r.Host
	go func() {

		link := fmt.Sprintf("http://%s/edit_review?token=%s", host, createTokenFromID(review.ID))

		smtp := system.GetConfig().SMTP
		msg := gomail.NewMessage()
		msg.SetHeader("From", smtp.From)
		msg.SetHeader("To", smtp.To)
		if len(smtp.Cc) > 0 {
			msg.SetHeader("Cc", smtp.Cc)
		}
		msg.SetHeader("Subject", T("new_review_has_been_created", map[string]interface{}{"Name": review.AuthorName}))
		msg.SetBody(
			"text/html",
			T("new_review_email_text", map[string]interface{}{"Link": link}),
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
func createTokenFromID(ID int64) string {
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
