package main

import (
	"log/slog"
	"net/http"
	"os"
)

func (app *application) serverError(w http.ResponseWriter, r *http.Request, err error) {
	method, uri := r.Method, r.URL.RequestURI()
	app.logger.Error(err.Error(), slog.String("method", method), slog.String("uri", uri))
	http.Error(w, err.Error(), http.StatusInternalServerError)
}

func (app *application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

func fatalError(logger *slog.Logger, msg string, err error) {
	logger.Error(msg, slog.String("error msg", err.Error()))
	os.Exit(1)
}
