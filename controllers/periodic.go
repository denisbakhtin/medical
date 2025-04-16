package controllers

import (
	"path"
	"time"

	"github.com/claudiu/gocron"
	"github.com/denisbakhtin/sitemap"
)

// StartPeriodicTasks initializes background periodic tasks
func (app *Application) StartPeriodicTasks() {
	gocron.Every(1).Day().Do(app.createXMLSitemap) // refresh daily
	gocron.Start()
}

// createXMLSitemap creates xml sitemap for search engines, and saves in public/sitemap folder
func (app *Application) createXMLSitemap() {
	app.Logger.Info("Starting XML sitemap generation")
	folder := path.Join(app.Config.Public, "sitemap")
	domain := app.Config.FullDomain
	now := time.Now()
	items := make([]sitemap.Page, 0, 500)

	// Home page
	items = append(items, sitemap.Page{
		Loc:        domain,
		LastMod:    now,
		Changefreq: sitemap.Daily,
		Priority:   1,
	})

	// Articles
	items = append(items, sitemap.Page{
		Loc:        app.fullURL("/articles"),
		LastMod:    now,
		Changefreq: sitemap.Monthly,
		Priority:   0.9,
	})

	articles, err := app.ArticlesRepo.GetAllPublished()
	if err != nil {
		app.Logger.Error(err)
	} else {
		for i := range articles {
			items = append(items, sitemap.Page{
				Loc:        app.fullURL(articles[i].URL()),
				LastMod:    articles[i].UpdatedAt,
				Changefreq: sitemap.Weekly,
				Priority:   0.9,
			})
		}
	}

	infos, err := app.InfosRepo.GetAllPublished()
	if err != nil {
		app.Logger.Error(err)
	} else {
		for i := range infos {
			items = append(items, sitemap.Page{
				Loc:        app.fullURL(infos[i].URL()),
				LastMod:    infos[i].UpdatedAt,
				Changefreq: sitemap.Weekly,
				Priority:   0.9,
			})
		}
	}

	// Static pages
	pages, err := app.PagesRepo.GetAllPublished()
	if err != nil {
		app.Logger.Error(err)
	} else {
		for i := range pages {
			items = append(items, sitemap.Page{
				Loc:        app.fullURL(pages[i].URL()),
				LastMod:    pages[i].UpdatedAt,
				Changefreq: sitemap.Monthly,
				Priority:   0.8,
			})
		}
	}

	// Reviews
	items = append(items, sitemap.Page{
		Loc:        app.fullURL("/reviews"),
		LastMod:    now,
		Changefreq: sitemap.Monthly,
		Priority:   0.7,
	})

	reviews, err := app.ReviewsRepo.GetAllPublished()
	if err != nil {
		app.Logger.Error(err)
	} else {
		for i := range reviews {
			items = append(items, sitemap.Page{
				Loc:        app.fullURL(reviews[i].URL()),
				LastMod:    reviews[i].UpdatedAt,
				Changefreq: sitemap.Monthly,
				Priority:   0.7,
			})
		}
	}

	if err := sitemap.SiteMap(path.Join(folder, "sitemap1.xml.gz"), items); err != nil {
		app.Logger.Error(err)
		return
	}
	if err := sitemap.SiteMapIndex(folder, "sitemap_index.xml", domain+"/public/sitemap/"); err != nil {
		app.Logger.Error(err)
		return
	}
	app.Logger.Infof("XML sitemap has been generated in %s\n", folder)
}
