package main

import (
	"bytes"
	"database/sql"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/mysqlstore"
	"github.com/alexedwards/scs/v2"
	_ "github.com/go-sql-driver/mysql"
)

func setup(cfg config) *application {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))
	db, err := openDB(cfg.dsn)
	if err != nil {
		fatalError(logger, "Errors occured when connecting to the DB", err)
	}
	htmlTemplateCache, err := newHtmlTemplateCache(cfg.staticDir)
	if err != nil {
		fatalError(logger, "Errors occured when preparing html pages", err)
	}
	sessionManager := scs.New()
	sessionManager.Store = mysqlstore.New(db)
	sessionManager.Lifetime = 12 * time.Hour
	sessionManager.Cookie.Secure = true

	app := &application{
		logger:            logger,
		db:                db,
		htmlTemplateCache: htmlTemplateCache,
		sessionManager:    sessionManager,
	}
	app.router = app.routes(cfg.staticDir)
	return app
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		db.Close()
		return nil, err
	}
	return db, nil
}

func (app *application) render(w http.ResponseWriter, r *http.Request, status int, pageName string, data templateData) {
	ts, ok := app.htmlTemplateCache[pageName]
	if !ok {
		err := fmt.Errorf("page template %s does not exist", pageName)
		app.serverError(w, r, err)
		return
	}

	buf := new(bytes.Buffer)
	err := ts.ExecuteTemplate(buf, "base", data)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	w.WriteHeader(status)
	buf.WriteTo(w)
}

func (app *application) newTemplateData(r *http.Request) templateData {
	newData := templateData{
		CurYear: time.Now().Year(),
		Flash:   app.sessionManager.PopString(r.Context(), "flash"),
	}
	return newData
}
