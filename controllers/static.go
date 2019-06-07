package controllers

import "github.com/peterpla/webdevgo/views"

// Static ... [add documentation]
type Static struct {
	Home    *views.View
	Contact *views.View
	Faq     *views.View
}

// NewStatic ... [add documentation]
func NewStatic() *Static {
	return &Static{
		Home:    views.NewView("bootstrap", "static/home"),
		Contact: views.NewView("bootstrap", "static/contact"),
		Faq:     views.NewView("bootstrap", "static/faq"),
	}
}
