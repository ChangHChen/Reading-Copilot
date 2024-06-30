package main

import (
	"database/sql"
	"flag"
	"log/slog"
	"net/http"
	"os"
	"text/template"

	"github.com/alexedwards/scs/v2"
)

type config struct {
	addr      string
	staticDir string
	dsn       string
}

type application struct {
	logger            *slog.Logger
	db                *sql.DB
	router            http.Handler
	htmlTemplateCache map[string]*template.Template
	sessionManager    *scs.SessionManager
}

func main() {
	var cfg config
	flag.StringVar(&cfg.addr, "addr", ":4000", "HTTP network address")
	flag.StringVar(&cfg.staticDir, "static-dir", "./ui", "Path to static files")
	flag.StringVar(&cfg.dsn, "dsn", "web:pass@/readingcopilot?parseTime=true", "MySQL data source name")

	flag.Parse()
	app := setup(cfg)

	defer app.db.Close()

	srv := &http.Server{
		Addr:     cfg.addr,
		Handler:  app.router,
		ErrorLog: slog.NewLogLogger(app.logger.Handler(), slog.LevelError),
	}
	app.logger.Info("starting server", slog.String("port", cfg.addr))

	err := srv.ListenAndServeTLS("./tls/cert.pem", "./tls/key.pem")

	app.logger.Error(err.Error())
	os.Exit(1)
}
