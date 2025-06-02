// Package views contains html templates
package views

import (
	"embed"
)

// TemplateFiles contains embedded html template files from /views/html directory
//
//go:embed html/*
var TemplateFiles embed.FS
