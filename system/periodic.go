package system

import (
	"fmt"
	"log"
	"path"
	"time"

	"github.com/denisbakhtin/medical/models"
	"github.com/denisbakhtin/sitemap"
)

//CreateXMLSitemap creates xml sitemap for search engines, and saves in public/sitemap folder
func CreateXMLSitemap() {
	log.Printf("INFO: Starting XML sitemap generation\n")
	db := models.GetDB()
	folder := path.Join(GetConfig().Public, "sitemap")
	domain := "http://www." + GetConfig().Domain
	now := time.Now()
	items := make([]sitemap.Item, 0, 500)

	//Home page
	items = append(items, sitemap.Item{
		Loc:        fmt.Sprintf("%s", domain),
		LastMod:    now,
		Changefreq: "daily",
		Priority:   1,
	})

	//Articles
	items = append(items, sitemap.Item{
		Loc:        fmt.Sprintf("%s%s", domain, "/articles"),
		LastMod:    now,
		Changefreq: "monthly",
		Priority:   0.9,
	})

	var articles []models.Article
	db.Where("published = ?", true).Order("id desc").Find(&articles)
	for i := range articles {
		items = append(items, sitemap.Item{
			Loc:        fmt.Sprintf("%s%s", domain, articles[i].URL()),
			LastMod:    articles[i].UpdatedAt,
			Changefreq: "weekly",
			Priority:   0.9,
		})
	}

	var infos []models.Info
	db.Where("published = ?", true).Order("id desc").Find(&infos)
	for i := range infos {
		items = append(items, sitemap.Item{
			Loc:        fmt.Sprintf("%s%s", domain, infos[i].URL()),
			LastMod:    infos[i].UpdatedAt,
			Changefreq: "weekly",
			Priority:   0.9,
		})
	}

	var exercises []models.Exercise
	db.Where("published = ?", true).Order("id desc").Find(&exercises)
	for i := range exercises {
		items = append(items, sitemap.Item{
			Loc:        fmt.Sprintf("%s%s", domain, exercises[i].URL()),
			LastMod:    infos[i].UpdatedAt,
			Changefreq: "weekly",
			Priority:   0.9,
		})
	}

	//Static pages
	var pages []models.Page
	db.Where("published = ?", true).Order("id desc").Find(&pages)
	for i := range pages {
		items = append(items, sitemap.Item{
			Loc:        fmt.Sprintf("%s%s", domain, pages[i].URL()),
			LastMod:    pages[i].UpdatedAt,
			Changefreq: "monthly",
			Priority:   0.8,
		})
	}

	//Reviews
	items = append(items, sitemap.Item{
		Loc:        fmt.Sprintf("%s%s", domain, "/reviews"),
		LastMod:    now,
		Changefreq: "monthly",
		Priority:   0.7,
	})

	var reviews []models.Review
	db.Where("published = ?", true).Order("id desc").Find(&reviews)
	for i := range reviews {
		items = append(items, sitemap.Item{
			Loc:        fmt.Sprintf("%s%s", domain, reviews[i].URL()),
			LastMod:    reviews[i].UpdatedAt,
			Changefreq: "monthly",
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
