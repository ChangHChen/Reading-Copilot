package main

import (
	"path/filepath"
	"text/template"
)

func newHtmlTemplateCache(staticDir string) (map[string]*template.Template, error) {
	htmlCache := map[string]*template.Template{}
	pages, err := filepath.Glob(staticDir + "/html/pages")
	if err != nil {
		return nil, err
	}
	for _, page := range pages {
		name := filepath.Base(page)

		files := []string{
			staticDir + "/html/html/base.tmpl",
			staticDir + "/html/html/partials/nav.tmpl",
			page,
		}

		ts, err := template.ParseFiles(files...)
		if err != nil {
			return nil, err
		}
		htmlCache[name] = ts
	}
	return htmlCache, nil
}
