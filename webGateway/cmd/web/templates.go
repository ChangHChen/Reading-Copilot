package main

import (
	"io/fs"
	"path/filepath"
	"strings"
	"text/template"
	"time"

	"github.com/ChangHChen/Reading-Copilot/webGateway/internal/models"
	"github.com/ChangHChen/Reading-Copilot/webGateway/ui"
)

type templateData struct {
	CurYear         int
	Flash           string
	Form            any
	IsAuthenticated bool
	UserName        string
	CSRFToken       string
	User            models.User
	BookList        bookList
	Book            models.BookMeta
	CurPage         int
	APIKeys         apiKeys
}
type apiKeys struct {
	OpenAIKeyReady    bool
	AnthropicKeyReady bool
	GoogleKeyReady    bool
}
type bookList struct {
	Error          string
	SearchKeyWords string
	Books          []models.BookMeta
}

func humanTime(t time.Time) string {
	return t.Format("02 Jan 2006 at 15:04")
}

var templateFunctions = template.FuncMap{
	"humanTime": humanTime,
}

func newHtmlTemplateCache() (map[string]*template.Template, error) {
	htmlCache := map[string]*template.Template{}
	pages, err := fs.Glob(ui.Files, "html/pages/*.tmpl")
	if err != nil {
		return nil, err
	}
	for _, page := range pages {
		pageName := strings.TrimSuffix(filepath.Base(page), ".tmpl")
		patterns := []string{
			"html/base.tmpl",
			"html/partials/*.tmpl",
			page,
		}

		ts, err := template.New(pageName).Funcs(templateFunctions).ParseFS(ui.Files, patterns...)
		if err != nil {
			return nil, err
		}
		htmlCache[pageName] = ts
	}
	return htmlCache, nil
}
