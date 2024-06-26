package main

import (
	"net/http"
	"webGateway/ui"

	"github.com/justinas/alice"
)

func (app *application) router() http.Handler {
	router := http.NewServeMux()
	router.Handle("GET /static/", http.FileServerFS(ui.Files))

	dynamicChain := alice.New(app.sessionManager.LoadAndSave, noSurf)

	router.Handle("GET /{$}", dynamicChain.ThenFunc(app.home))
	router.Handle("GET /about", dynamicChain.ThenFunc(app.about))

	standardChain := alice.New(app.recoverPanic, app.logRequest, commonHeaders)
	return standardChain.Then(router)
}
