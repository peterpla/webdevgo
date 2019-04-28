package main

import (
	"html/template"
	"os"
)

func main() {
	t, err := template.ParseFiles("hello.gohtml")
	if err != nil {
		panic(err)
	}

	type Salutation struct {
		Name  string
		Title string
		Born  int
	}

	var data = []Salutation{
		{"John Smith", "Explorer", 1580},
		{"Batman", "Superhero", 1939},
	}

	for _, r := range data {
		err = t.Execute(os.Stdout, r)
		if err != nil {
			panic(err)
		}
	}
}
