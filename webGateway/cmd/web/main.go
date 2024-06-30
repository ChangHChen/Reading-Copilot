package main

import (
	"database/sql"
	"flag"
	"log/slog"
	"net/http"
	"os"
	"text/template"
)

type config struct {
	addr      string
	staticDir string
	dsn       string
}

type application struct {
	logger            *slog.Logger
	db                *sql.DB
	router            *http.ServeMux
	htmlTemplateCache map[string]*template.Template
}

func main() {
	var cfg config
	flag.StringVar(&cfg.addr, "addr", "4000", "HTTP network address")
	flag.StringVar(&cfg.staticDir, "static-dir", "./ui", "Path to static files")
	flag.StringVar(&cfg.dsn, "dsn", "web:pass@/readingcopilot?parseTime=true", "MySQL data source name")

	flag.Parse()
	app := &application{}
	app.setup(cfg)

	defer app.db.Close()
	app.logger.Info("starting server", slog.String("port", cfg.addr))
	err := http.ListenAndServe(":"+cfg.addr, app.router)
	app.logger.Error(err.Error())
	os.Exit(1)
}
