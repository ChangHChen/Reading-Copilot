package main

import (
	"log/slog"
	"os"
	"time"

	"github.com/alexedwards/scs/mysqlstore"
	"github.com/alexedwards/scs/v2"
)

func setupApplication(dsn string) *application {
	logger := slog.New(slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelDebug, AddSource: true}))
	tc, err := newTemplateCache()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
	app := &application{
		templateCache: tc,
		logger:        logger,
	}
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
	sessionManager := scs.New()
	sessionManager.Store = mysqlstore.New(app.snippets.DB)
	sessionManager.Lifetime = 12 * time.Hour
	sessionManager.Cookie.Secure = true
	app.sessionManager = sessionManager
	return app
}
