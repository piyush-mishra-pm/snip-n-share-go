package main

import "github.com/piyush-mishra-pm/snip-n-share-go/internal/models"

type templateData struct {
	Snip  *models.Snip
	Snips []*models.Snip
}
