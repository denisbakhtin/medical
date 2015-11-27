package controllers

import (
	"fmt"
	"log"
	"net/http"
	"regexp"

	"github.com/denisbakhtin/medical/controllers/oauth"
	"github.com/denisbakhtin/medical/helpers"
	"github.com/denisbakhtin/medical/models"
)

//ArticleShow handles GET /articles/:id-slug route
func ArticleShow(w http.ResponseWriter, r *http.Request) {
	tmpl := helpers.Template(r)
	session := helpers.Session(r)
	data := helpers.DefaultData(r)
	if r.Method == "GET" {

		re := regexp.MustCompile("^[0-9]+")
		id := re.FindString(r.URL.Path[len("/articles/"):])
		article, err := models.GetArticle(id)
		if err != nil || !article.Published {
			w.WriteHeader(404)
			tmpl.Lookup("errors/404").Execute(w, nil)
			return
		}
		//redirect to canonical url
		if r.URL.Path != article.URL() {
			http.Redirect(w, r, article.URL(), http.StatusSeeOther)
			return
		}
		data["Article"] = article
		data["Title"] = article.Name
		data["Active"] = "/articles"
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
	if r.Method == "GET" {

		list, err := models.GetPublishedArticles()
		if err != nil {
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
	if r.Method == "GET" {

		list, err := models.GetArticles()
		if err != nil {
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
	if r.Method == "GET" {

		data["Title"] = T("new_article")
		data["Active"] = "articles"
		data["Flash"] = session.Flashes()
		session.Save(r, w)
		tmpl.Lookup("articles/form").Execute(w, data)

	} else if r.Method == "POST" {

		r.ParseForm()
		article := &models.Article{
			Name:      r.PostFormValue("name"),
			Slug:      r.PostFormValue("slug"),
			Content:   r.PostFormValue("content"),
			Excerpt:   r.PostFormValue("excerpt"),
			Published: helpers.Atob(r.PostFormValue("published")),
		}

		if err := article.Insert(); err != nil {
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
	if r.Method == "GET" {

		id := r.URL.Path[len("/admin/edit_article/"):]
		article, err := models.GetArticle(id)
		if err != nil {
			w.WriteHeader(404)
			tmpl.Lookup("errors/404").Execute(w, helpers.ErrorData(err))
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
			ID:        helpers.Atoi64(r.PostFormValue("id")),
			Name:      r.PostFormValue("name"),
			Slug:      r.PostFormValue("slug"),
			Content:   r.PostFormValue("content"),
			Excerpt:   r.PostFormValue("excerpt"),
			Published: helpers.Atob(r.PostFormValue("published")),
		}

		if err := article.Update(); err != nil {
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

	if r.Method == "POST" {

		article, err := models.GetArticle(r.PostFormValue("id"))
		if err != nil {
			log.Printf("ERROR: %s\n", err)
			w.WriteHeader(404)
			tmpl.Lookup("errors/404").Execute(w, helpers.ErrorData(err))
			return
		}

		if err := article.Delete(); err != nil {
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

//PostOnFacebook publishes article preview on facebook page wall
func PostOnFacebook(w http.ResponseWriter, r *http.Request) {

	if r.Method == "POST" {

		article, err := models.GetArticle(r.FormValue("id"))
		if err != nil {
			log.Printf("ERROR: %s\n", err)
			w.WriteHeader(404)
			return
		}
		err = oauth.PostOnFacebook(
			fmt.Sprintf("http://%s/articles/%d", r.Host, article.ID),
			article.Name,
		)
		if err != nil {
			log.Printf("ERROR: %s\n", err)
			w.WriteHeader(500)
			return
		}
		w.WriteHeader(200)

	} else {
		err := fmt.Errorf("Method %q not allowed", r.Method)
		log.Printf("ERROR: %s\n", err)
		w.WriteHeader(405)
	}
}
