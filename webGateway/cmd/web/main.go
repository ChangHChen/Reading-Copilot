package main

import (
	"flag"
	"log/slog"
	"net/http"
	"os"
	"text/template"

	"github.com/alexedwards/scs/v2"
)

type config struct {
	port string
}

type application struct {
	logger         *slog.Logger
	templateCache  map[string]*template.Template
	sessionManager *scs.SessionManager
}

func main() {
	var cfg config
	flag.StringVar(&cfg.port, "port", "8000", "HTTP network address")
	flag.Parse()

	app := &application{}


	srv := &http.Server{
		Addr:     ":" + cfg.port,
		Handler:  app.router(),
		ErrorLog: slog.NewLogLogger(app.logger.Handler(), slog.LevelError),
	}
	app.logger.Info("Starting Server...", slog.Any("port", cfg.port))
	err := srv.ListenAndServeTLS("./tls/cert.pem", "./tls/key.pem")
	app.logger.Error(err.Error())
	os.Exit(1)
}
