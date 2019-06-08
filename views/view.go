package views

import (
	"bytes"
	"html/template"
	"io"
	"net/http"
	"path/filepath"
)

// LayoutDir sets the path to layout files
var LayoutDir = "views/layouts/"

// TemplateDir sets the path to template files
var TemplateDir = "views/"

// TemplateExt sets the file extension for template files
var TemplateExt = ".gohtml"

// View struct used by most view methods
type View struct {
	Template *template.Template
	Layout   string
}

// Render method used to render templates into web pages
func (v *View) Render(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "text/html")

	switch data.(type) {
	case Data:
		// Data struct - do nothing, View processing this
	default:
		// not a Data struct - pass the data argument in a Data struct
		data = Data{
			Yield: data,
		}
	}

	var buf bytes.Buffer
	err := v.Template.ExecuteTemplate(&buf, v.Layout, data)
	if err != nil {
		http.Error(w, "Something went wrong. If the problem persists, please email support@exercise.com",
			http.StatusInternalServerError)
		return
	}

	// template executed without error, copy the buffer to w, and done
	io.Copy(w, &buf)
}

func (v *View) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	v.Render(w, nil)
}

func layoutFiles() []string {
	files, err := filepath.Glob(LayoutDir + "*" + TemplateExt)
	if err != nil {
		panic(err)
	}
	return files
}

// NewView creates a new View based on the layout and template files arguments
func NewView(layout string, files ...string) *View {
	addTemplatePath(files)
	addTemplateExt(files)
	files = append(files, layoutFiles()...)

	t, err := template.ParseFiles(files...)
	if err != nil {
		panic(err)
	}

	return &View{
		Template: t,
		Layout:   layout,
	}
}

// addTemplatePath takes in a slice of strings
// representing file paths for templates, and it prepends
// the TemplateDir directory to each string in the slice
//
// E.g., the input {"home"} would result in the output
// {"views/home"} if TemplateDir == "views/"
func addTemplatePath(files []string) {
	for i, f := range files {
		files[i] = TemplateDir + f
	}
}

// addTemplateExt takes in a slice of strings
// representing file paths for templates, and it appends
// the TemplateExt extension to each string in the slice
//
// E.g., the input {"home"} would result in the output
// {"home.gohtml"} if TemplateExt == ".gohtml"
func addTemplateExt(files []string) {
	for i, f := range files {
		files[i] = f + TemplateExt
	}
}
