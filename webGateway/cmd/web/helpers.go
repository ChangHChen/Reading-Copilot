package main

import (
	"database/sql"
	"log/slog"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

func (app *application) setup(cfg config) {
	var err error
	app.logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))
	app.db, err = openDB(cfg.dsn)
	if err != nil {
		app.fatalError("Errors occured when connecting to the DB", err)
	}
	app.router = app.routes(cfg.staticDir)
	app.htmlTemplateCache, err = newHtmlTemplateCache(cfg.staticDir)
	if err != nil {
		app.fatalError("Errors occured when preparing html pages", err)
	}
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
