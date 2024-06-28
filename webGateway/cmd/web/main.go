package main

import (
	"flag"
	"log/slog"
	"net/http"
	"os"
)

type config struct {
	addr      string
	staticDir string
}

type application struct {
	cfg    config
	logger *slog.Logger
}

func main() {
	app := application{}
	app.logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))
	flag.StringVar(&app.cfg.addr, "addr", "4000", "HTTP network address")
	flag.StringVar(&app.cfg.staticDir, "static-dir", "./ui", "Path to static files")
	flag.Parse()

	app.logger.Info("starting server", slog.String("port", app.cfg.addr))

	err := http.ListenAndServe(":"+app.cfg.addr, app.routes())
	app.logger.Error(err.Error())
	os.Exit(1)
}
