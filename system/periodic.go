package system

import (
	"fmt"
	"log"
	"path"
	"time"

	"github.com/claudiu/gocron"
	"github.com/denisbakhtin/medical/config"
	"github.com/denisbakhtin/medical/models"
	"github.com/denisbakhtin/sitemap"
)

// SetupPeriodicTasks initializes periodic tasks
func SetupPeriodicTasks(mode string) {
	// Periodic tasks
	if mode == config.ReleaseMode {
		createXMLSitemap() // refresh sitemap now
	}
	gocron.Every(1).Day().Do(createXMLSitemap) // refresh daily
	gocron.Start()
}

// createXMLSitemap creates xml sitemap for search engines, and saves in public/sitemap folder
func createXMLSitemap() {
	log.Printf("INFO: Starting XML sitemap generation\n")
	db := models.GetDB()
	folder := path.Join(config.GetConfig().Public, "sitemap")
	domain := "http://www." + config.GetConfig().Domain
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
		Loc:        fmt.Sprintf("%s%s", domain, "/articles"),
		LastMod:    now,
		Changefreq: sitemap.Monthly,
		Priority:   0.9,
	})

	var articles []models.Article
	db.Where("published = ?", true).Order("id desc").Find(&articles)
	for i := range articles {
		items = append(items, sitemap.Page{
			Loc:        fmt.Sprintf("%s%s", domain, articles[i].URL()),
			LastMod:    articles[i].UpdatedAt,
			Changefreq: sitemap.Weekly,
			Priority:   0.9,
		})
	}

	var infos []models.Info
	db.Where("published = ?", true).Order("id desc").Find(&infos)
	for i := range infos {
		items = append(items, sitemap.Page{
			Loc:        fmt.Sprintf("%s%s", domain, infos[i].URL()),
			LastMod:    infos[i].UpdatedAt,
			Changefreq: sitemap.Weekly,
			Priority:   0.9,
		})
	}

	// Static pages
	var pages []models.Page
	db.Where("published = ?", true).Order("id desc").Find(&pages)
	for i := range pages {
		items = append(items, sitemap.Page{
			Loc:        fmt.Sprintf("%s%s", domain, pages[i].URL()),
			LastMod:    pages[i].UpdatedAt,
			Changefreq: sitemap.Monthly,
			Priority:   0.8,
		})
	}

	// Reviews
	items = append(items, sitemap.Page{
		Loc:        fmt.Sprintf("%s%s", domain, "/reviews"),
		LastMod:    now,
		Changefreq: sitemap.Monthly,
		Priority:   0.7,
	})

	var reviews []models.Review
	db.Where("published = ?", true).Order("id desc").Find(&reviews)
	for i := range reviews {
		items = append(items, sitemap.Page{
			Loc:        fmt.Sprintf("%s%s", domain, reviews[i].URL()),
			LastMod:    reviews[i].UpdatedAt,
			Changefreq: sitemap.Monthly,
			Priority:   0.7,
		})
	}

	if err := sitemap.SiteMap(path.Join(folder, "sitemap1.xml.gz"), items); err != nil {
		log.Printf("ERROR: %s\n", err)
		return
	}
	if err := sitemap.SiteMapIndex(folder, "sitemap_index.xml", domain+"/public/sitemap/"); err != nil {
		log.Printf("ERROR: %s\n", err)
		return
	}
	log.Printf("INFO: XML sitemap has been generated in %s\n", folder)
}
