package controllers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/denisbakhtin/medical/controllers/oauth"
	"github.com/denisbakhtin/medical/helpers"
	"github.com/denisbakhtin/medical/models"
	"github.com/gorilla/context"
	"gopkg.in/guregu/null.v3"
)

//PostShow handles GET /posts/:id route
func PostShow(w http.ResponseWriter, r *http.Request) {
	tmpl := helpers.Template(r)
	session := helpers.Session(r)
	data := helpers.DefaultData(r)
	if r.Method == "GET" {

		id := r.URL.Path[len("/posts/"):]
		post, err := models.GetPost(id)
		if err != nil || !post.Published {
			tmpl.Lookup("errors/404").Execute(w, nil)
			return
		}
		data["Post"] = post
		data["Title"] = post.Name
		data["Active"] = fmt.Sprintf("posts/%s", id)
		data["OauthName"] = session.Values["oauth_name"]
		//Facebook open graph meta tags
		data["Ogheadprefix"] = "og: http://ogp.me/ns# fb: http://ogp.me/ns/fb# article: http://ogp.me/ns/article#"
		data["Ogtitle"] = post.Name
		data["Ogurl"] = fmt.Sprintf("http://%s/posts/%d", r.Host, post.ID)
		data["Ogtype"] = "article"
		data["Ogdescription"] = post.Excerpt()
		if img := post.GetImage(); len(img) > 0 {
			data["Ogimage"] = fmt.Sprintf("http://%s%s", r.Host, img)
		}
		//flashes
		data["Flash"] = session.Flashes("comments")
		session.Save(r, w)
		tmpl.Lookup("posts/show").Execute(w, data)

	} else {
		err := fmt.Errorf("Method %q not allowed", r.Method)
		log.Printf("ERROR: %s\n", err)
		tmpl.Lookup("errors/405").Execute(w, helpers.ErrorData(err))
	}
}

//PostIndex handles GET /admin/posts route
func PostIndex(w http.ResponseWriter, r *http.Request) {
	tmpl := helpers.Template(r)
	data := helpers.DefaultData(r)
	T := helpers.T(r)
	if r.Method == "GET" {

		list, err := models.GetPosts()
		if err != nil {
			w.WriteHeader(500)
			tmpl.Lookup("errors/500").Execute(w, helpers.ErrorData(err))
			return
		}
		data["Title"] = T("posts")
		data["Active"] = "posts"
		data["List"] = list
		tmpl.Lookup("posts/index").Execute(w, data)

	} else {
		err := fmt.Errorf("Method %q not allowed", r.Method)
		log.Printf("ERROR: %s\n", err)
		w.WriteHeader(405)
		tmpl.Lookup("errors/405").Execute(w, helpers.ErrorData(err))
	}
}

//PostCreate handles /admin/new_post route
func PostCreate(w http.ResponseWriter, r *http.Request) {
	tmpl := helpers.Template(r)
	session := helpers.Session(r)
	data := helpers.DefaultData(r)
	T := helpers.T(r)
	if r.Method == "GET" {

		tags, err := models.GetTags()
		if err != nil {
			w.WriteHeader(404)
			tmpl.Lookup("errors/404").Execute(w, nil)
			return
		}
		data["Title"] = T("new_post")
		data["Active"] = "posts"
		data["Tags"] = tags
		data["Flash"] = session.Flashes()
		session.Save(r, w)
		tmpl.Lookup("posts/form").Execute(w, data)

	} else if r.Method == "POST" {

		r.ParseForm()
		post := &models.Post{
			Name:      r.PostFormValue("name"),
			Content:   r.PostFormValue("content"),
			Published: helpers.Atob(r.PostFormValue("published")),
			Tags:      r.Form["tags"], //PostFormValue returns only first value
		}

		if user := context.Get(r, "user"); user != nil {
			post.UserID = null.NewInt(user.(*models.User).ID, user.(*models.User).ID > 0)
		}
		if err := post.Insert(); err != nil {
			session.AddFlash(err.Error())
			session.Save(r, w)
			http.Redirect(w, r, "/admin/new_post", 303)
			return
		}
		http.Redirect(w, r, "/admin/posts", 303)

	} else {
		err := fmt.Errorf("Method %q not allowed", r.Method)
		log.Printf("ERROR: %s\n", err)
		w.WriteHeader(405)
		tmpl.Lookup("errors/405").Execute(w, helpers.ErrorData(err))
	}
}

//PostUpdate handles /admin/edit_post/:id route
func PostUpdate(w http.ResponseWriter, r *http.Request) {
	tmpl := helpers.Template(r)
	session := helpers.Session(r)
	data := helpers.DefaultData(r)
	T := helpers.T(r)
	if r.Method == "GET" {

		id := r.URL.Path[len("/admin/edit_post/"):]
		post, err := models.GetPost(id)
		if err != nil {
			w.WriteHeader(404)
			tmpl.Lookup("errors/404").Execute(w, helpers.ErrorData(err))
			return
		}
		tags, err := models.GetTags()
		if err != nil {
			w.WriteHeader(404)
			tmpl.Lookup("errors/404").Execute(w, nil)
			return
		}

		data["Title"] = T("edit_post")
		data["Active"] = "posts"
		data["Post"] = post
		data["Tags"] = tags
		data["Flash"] = session.Flashes()
		session.Save(r, w)
		tmpl.Lookup("posts/form").Execute(w, data)

	} else if r.Method == "POST" {

		r.ParseForm()
		post := &models.Post{
			ID:        helpers.Atoi64(r.PostFormValue("id")),
			Name:      r.PostFormValue("name"),
			Content:   r.PostFormValue("content"),
			Published: helpers.Atob(r.PostFormValue("published")),
			Tags:      r.Form["tags"], //PostFormValue returns only first value
		}

		if err := post.Update(); err != nil {
			session.AddFlash(err.Error())
			session.Save(r, w)
			http.Redirect(w, r, r.RequestURI, 303)
			return
		}
		http.Redirect(w, r, "/admin/posts", 303)

	} else {
		err := fmt.Errorf("Method %q not allowed", r.Method)
		log.Printf("ERROR: %s\n", err)
		w.WriteHeader(405)
		tmpl.Lookup("errors/405").Execute(w, helpers.ErrorData(err))
	}
}

//PostDelete handles /admin/delete_post route
func PostDelete(w http.ResponseWriter, r *http.Request) {
	tmpl := helpers.Template(r)

	if r.Method == "POST" {

		post, err := models.GetPost(r.PostFormValue("id"))
		if err != nil {
			log.Printf("ERROR: %s\n", err)
			w.WriteHeader(404)
			tmpl.Lookup("errors/404").Execute(w, helpers.ErrorData(err))
			return
		}

		if err := post.Delete(); err != nil {
			log.Printf("ERROR: %s\n", err)
			w.WriteHeader(500)
			tmpl.Lookup("errors/500").Execute(w, helpers.ErrorData(err))
			return
		}
		http.Redirect(w, r, "/admin/posts", 303)

	} else {
		err := fmt.Errorf("Method %q not allowed", r.Method)
		log.Printf("ERROR: %s\n", err)
		w.WriteHeader(405)
		tmpl.Lookup("errors/405").Execute(w, helpers.ErrorData(err))
	}
}

//PostOnFacebook publishes post preview on facebook page wall
func PostOnFacebook(w http.ResponseWriter, r *http.Request) {

	if r.Method == "POST" {

		post, err := models.GetPost(r.FormValue("id"))
		if err != nil {
			log.Printf("ERROR: %s\n", err)
			w.WriteHeader(404)
			return
		}
		err = oauth.PostOnFacebook(
			fmt.Sprintf("http://%s/posts/%d", r.Host, post.ID),
			post.Name,
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
