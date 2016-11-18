package controllers

import (
	"fmt"
	"log"
	"net/http"
	"regexp"

	"github.com/denisbakhtin/medical/helpers"
	"github.com/denisbakhtin/medical/models"
)

//ArticleShow handles GET /articles/:id-slug route
func ArticleShow(w http.ResponseWriter, r *http.Request) {
	tmpl := helpers.Template(r)
	session := helpers.Session(r)
	data := helpers.DefaultData(r)
	db := models.GetDB()
	if r.Method == "GET" {

		re := regexp.MustCompile("^[0-9]+")
		id := re.FindString(r.URL.Path[len("/articles/"):])
		article := &models.Article{}
		db.First(article, id)
		if article.ID == 0 || !article.Published {
			w.WriteHeader(404)
			tmpl.Lookup("errors/404").Execute(w, nil)
			return
		}
		//redirect to canonical url
		if r.URL.Path != article.URL() {
			http.Redirect(w, r, article.URL(), http.StatusSeeOther)
			return
		}
		var testimonials []models.Review
		db.Where("published = ? and article_id = ?", true, article.ID).Order("created_at desc").Find(&testimonials)
		topComments := models.GetTopComments(article.ID)
		comments := models.GetComments(article.ID)
		article.Comments = append(topComments, comments...)
		data["Article"] = article
		data["Testimonials"] = testimonials
		data["Title"] = article.Name
		data["Active"] = "/articles"
		data["MetaDescription"] = article.MetaDescription
		data["MetaKeywords"] = article.MetaKeywords
		//Facebook open graph meta tags
		data["Ogheadprefix"] = "og: http://ogp.me/ns# fb: http://ogp.me/ns/fb# article: http://ogp.me/ns/article#"
		data["Ogtitle"] = article.Name
		data["Ogurl"] = fmt.Sprintf("http://%s/articles/%d", r.Host, article.ID)
		data["Ogtype"] = "article"
		data["Ogdescription"] = article.Excerpt
		if img := article.GetImage(); len(img) > 0 {
			data["Ogimage"] = fmt.Sprintf("http://%s%s", r.Host, img)
		}
		//flashes
		data["Flash"] = session.Flashes("comments")
		session.Save(r, w)
		tmpl.Lookup("articles/show").Execute(w, data)

	} else {
		err := fmt.Errorf("Method %q not allowed", r.Method)
		log.Printf("ERROR: %s\n", err)
		w.WriteHeader(405)
		tmpl.Lookup("errors/405").Execute(w, helpers.ErrorData(err))
	}
}

//ArticlePublicIndex handles GET /articles route
func ArticlePublicIndex(w http.ResponseWriter, r *http.Request) {
	tmpl := helpers.Template(r)
	data := helpers.DefaultData(r)
	T := helpers.T(r)
	db := models.GetDB()
	if r.Method == "GET" {

		var list []models.Article
		if err := db.Where("published = ?", true).Order("id desc").Find(&list).Error; err != nil {
			log.Printf("ERROR: %s\n", err)
			w.WriteHeader(500)
			tmpl.Lookup("errors/500").Execute(w, helpers.ErrorData(err))
			return
		}
		data["Title"] = T("kinesiology_in_practice")
		data["Active"] = r.RequestURI
		data["List"] = list
		tmpl.Lookup("articles/public-index").Execute(w, data)

	} else {
		err := fmt.Errorf("Method %q not allowed", r.Method)
		log.Printf("ERROR: %s\n", err)
		w.WriteHeader(405)
		tmpl.Lookup("errors/405").Execute(w, helpers.ErrorData(err))
	}
}

//ArticleIndex handles GET /admin/articles route
func ArticleIndex(w http.ResponseWriter, r *http.Request) {
	tmpl := helpers.Template(r)
	data := helpers.DefaultData(r)
	T := helpers.T(r)
	db := models.GetDB()
	if r.Method == "GET" {

		var list []models.Article
		if err := db.Order("published desc, id desc").Find(&list).Error; err != nil {
			log.Printf("ERROR: %s\n", err)
			w.WriteHeader(500)
			tmpl.Lookup("errors/500").Execute(w, helpers.ErrorData(err))
			return
		}
		data["Title"] = T("articles")
		data["Active"] = "articles"
		data["List"] = list
		tmpl.Lookup("articles/index").Execute(w, data)

	} else {
		err := fmt.Errorf("Method %q not allowed", r.Method)
		log.Printf("ERROR: %s\n", err)
		w.WriteHeader(405)
		tmpl.Lookup("errors/405").Execute(w, helpers.ErrorData(err))
	}
}

//ArticleCreate handles /admin/new_article route
func ArticleCreate(w http.ResponseWriter, r *http.Request) {
	tmpl := helpers.Template(r)
	session := helpers.Session(r)
	data := helpers.DefaultData(r)
	T := helpers.T(r)
	db := models.GetDB()
	if r.Method == "GET" {

		data["Title"] = T("new_article")
		data["Active"] = "articles"
		data["Flash"] = session.Flashes()
		session.Save(r, w)
		tmpl.Lookup("articles/form").Execute(w, data)

	} else if r.Method == "POST" {

		r.ParseForm()
		article := &models.Article{
			Name:            r.PostFormValue("name"),
			Slug:            r.PostFormValue("slug"),
			Content:         r.PostFormValue("content"),
			Excerpt:         r.PostFormValue("excerpt"),
			SellingPreface:  r.PostFormValue("selling_preface"),
			MetaDescription: r.PostFormValue("meta_description"),
			MetaKeywords:    r.PostFormValue("meta_keywords"),
			Published:       helpers.Atob(r.PostFormValue("published")),
		}

		if err := db.Create(article).Error; err != nil {
			session.AddFlash(err.Error())
			session.Save(r, w)
			http.Redirect(w, r, "/admin/new_article", 303)
			return
		}
		http.Redirect(w, r, "/admin/articles", 303)

	} else {
		err := fmt.Errorf("Method %q not allowed", r.Method)
		log.Printf("ERROR: %s\n", err)
		w.WriteHeader(405)
		tmpl.Lookup("errors/405").Execute(w, helpers.ErrorData(err))
	}
}

//ArticleUpdate handles /admin/edit_article/:id route
func ArticleUpdate(w http.ResponseWriter, r *http.Request) {
	tmpl := helpers.Template(r)
	session := helpers.Session(r)
	data := helpers.DefaultData(r)
	T := helpers.T(r)
	db := models.GetDB()
	if r.Method == "GET" {

		id := r.URL.Path[len("/admin/edit_article/"):]
		article := &models.Article{}
		db.First(article, id)
		if article.ID == 0 {
			w.WriteHeader(404)
			tmpl.Lookup("errors/404").Execute(w, nil)
			return
		}

		data["Title"] = T("edit_article")
		data["Active"] = "articles"
		data["Article"] = article
		data["Flash"] = session.Flashes()
		session.Save(r, w)
		tmpl.Lookup("articles/form").Execute(w, data)

	} else if r.Method == "POST" {

		r.ParseForm()
		article := &models.Article{
			ID:              helpers.Atouint(r.PostFormValue("id")),
			Name:            r.PostFormValue("name"),
			Slug:            r.PostFormValue("slug"),
			Content:         r.PostFormValue("content"),
			Excerpt:         r.PostFormValue("excerpt"),
			SellingPreface:  r.PostFormValue("selling_preface"),
			MetaDescription: r.PostFormValue("meta_description"),
			MetaKeywords:    r.PostFormValue("meta_keywords"),
			Published:       helpers.Atob(r.PostFormValue("published")),
		}

		if err := db.Save(article).Error; err != nil {
			session.AddFlash(err.Error())
			session.Save(r, w)
			http.Redirect(w, r, r.RequestURI, 303)
			return
		}
		http.Redirect(w, r, "/admin/articles", 303)

	} else {
		err := fmt.Errorf("Method %q not allowed", r.Method)
		log.Printf("ERROR: %s\n", err)
		w.WriteHeader(405)
		tmpl.Lookup("errors/405").Execute(w, helpers.ErrorData(err))
	}
}

//ArticleDelete handles /admin/delete_article route
func ArticleDelete(w http.ResponseWriter, r *http.Request) {
	tmpl := helpers.Template(r)
	db := models.GetDB()

	if r.Method == "POST" {

		article := &models.Article{}
		db.First(article, r.PostFormValue("id"))
		if article.ID == 0 {
			w.WriteHeader(404)
			tmpl.Lookup("errors/404").Execute(w, nil)
			return
		}

		if err := db.Delete(article).Error; err != nil {
			log.Printf("ERROR: %s\n", err)
			w.WriteHeader(500)
			tmpl.Lookup("errors/500").Execute(w, helpers.ErrorData(err))
			return
		}
		http.Redirect(w, r, "/admin/articles", 303)

	} else {
		err := fmt.Errorf("Method %q not allowed", r.Method)
		log.Printf("ERROR: %s\n", err)
		w.WriteHeader(405)
		tmpl.Lookup("errors/405").Execute(w, helpers.ErrorData(err))
	}
}
