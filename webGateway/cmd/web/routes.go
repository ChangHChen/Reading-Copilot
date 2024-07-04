package main

import (
	"net/http"

	"github.com/ChangHChen/Reading-Copilot/webGateway/ui"
	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler {
	router := http.NewServeMux()

	cacheDir := http.Dir("./cache")
	cacheFileServer := http.FileServer(cacheDir)
	router.Handle("/cache/", app.noDirListing(http.StripPrefix("/cache/", app.noDirListing(cacheFileServer))))
	router.Handle("GET /static/", app.noDirListing(http.FileServerFS(ui.Files)))

	dynamic := alice.New(app.sessionManager.LoadAndSave, noSurf, app.validAuthentication)
	router.Handle("GET /{$}", dynamic.ThenFunc(app.home))
	router.Handle("GET /about", dynamic.ThenFunc(app.about))
	router.Handle("GET /search", dynamic.ThenFunc(app.search))
	router.Handle("GET /book/view/{id}", dynamic.ThenFunc(app.bookView))

	guestChain := dynamic.Append(app.requireGuest)
	router.Handle("GET /user/login", guestChain.ThenFunc(app.login))
	router.Handle("POST /user/login", guestChain.ThenFunc(app.loginPost))
	router.Handle("GET /user/signup", guestChain.ThenFunc(app.signUp))
	router.Handle("POST /user/signup", guestChain.ThenFunc(app.signUpPost))

	authenticatedChain := dynamic.Append(app.requireAuthentication)
	router.Handle("POST /user/logout", authenticatedChain.ThenFunc(app.logoutPost))
	router.Handle("GET /user/profile", authenticatedChain.ThenFunc(app.profile))
	router.Handle("GET /user/password", authenticatedChain.ThenFunc(app.updatePWD))
	router.Handle("POST /user/password", authenticatedChain.ThenFunc(app.updatePWDPost))

	commonMiddleware := alice.New(app.recoverPanic, app.logRequest, commonHeaders)
	return commonMiddleware.Then(router)
}
