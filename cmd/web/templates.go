package main

import "github.com/rozbeikonazar/gosnippetbox/internal/models"

type templateData struct {
	Snippet  *models.Snippet
	Snippets []*models.Snippet
}
