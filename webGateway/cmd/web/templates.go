package main

import (
	"path/filepath"
	"strings"
	"text/template"
	"time"
)

type templateData struct {
	CurYear int
	Flash   string
}

func humanTime(t time.Time) string {
	return t.Format("02 Jan 2006 at 15:04")
}

var templateFunctions = template.FuncMap{
	"humanTime": humanTime,
}

func newHtmlTemplateCache(staticDir string) (map[string]*template.Template, error) {
	htmlCache := map[string]*template.Template{}
	pages, err := filepath.Glob(staticDir + "/html/pages/*.tmpl")
	if err != nil {
		return nil, err
	}
	for _, page := range pages {
		pageName := strings.TrimSuffix(filepath.Base(page), ".tmpl")

		ts, err := template.New(pageName).Funcs(templateFunctions).ParseFiles("./ui/html/base.tmpl")
		if err != nil {
			return nil, err
		}
		ts, err = ts.ParseGlob("./ui/html/partials/*.tmpl")
		if err != nil {
			return nil, err
		}

		ts, err = ts.ParseFiles(page)
		if err != nil {
			return nil, err
		}

		htmlCache[pageName] = ts
	}
	return htmlCache, nil
}
