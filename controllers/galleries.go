package controllers

import (
	"net/http"

	"../views"
)

type Galleries struct {
	Gallery *views.View
}

func NewGalleries() *Galleries {
	return &Galleries{
		Gallery: views.NewView("bootstrap", "galleries/gallery"),
	}
}

// New is used to render the form where a user can
// create a new gallery
//
// POST /signup
func (g *Galleries) Create(w http.ResponseWriter, r *http.Request) {
	if err := g.Gallery.Render(w, nil); err != nil {
		panic(err)
	}
}
