package controllers

import (
	"net/http"

	"github.com/peterpla/webdevgo/views"
)

// Galleries ... [add documentation]
type Galleries struct {
	Gallery *views.View
}

// NewGalleries ... [add documentation]
func NewGalleries() *Galleries {
	return &Galleries{
		Gallery: views.NewView("bootstrap", "galleries/gallery"),
	}
}

// Create is used to render the form where a user can
// create a new gallery
//
// POST /signup
func (g *Galleries) Create(w http.ResponseWriter, r *http.Request) {
	if err := g.Gallery.Render(w, nil); err != nil {
		panic(err)
	}
}
