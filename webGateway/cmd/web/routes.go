package main

import (
	"net/http"

	"github.com/justinas/alice"
)

func (app *application) routes(staticDir string) http.Handler {
	router := http.NewServeMux()
	fileServer := http.FileServer(http.Dir(staticDir))
	router.Handle("GET /static/", fileServer)

	dynamic := alice.New(app.sessionManager.LoadAndSave, noSurf)
	router.Handle("GET /{$}", dynamic.ThenFunc(app.home))
	router.Handle("GET /about", dynamic.ThenFunc(app.about))
	router.Handle("GET /book/view/{id}", dynamic.ThenFunc(app.bookView))

	guestChain := dynamic.Append(app.requireGuest)
	router.Handle("GET /user/login", guestChain.ThenFunc(app.login))
	router.Handle("POST /user/login", guestChain.ThenFunc(app.loginPost))
	router.Handle("GET /user/signup", guestChain.ThenFunc(app.signUp))
	router.Handle("POST /user/signup", guestChain.ThenFunc(app.signUpPost))

	authenticatedChain := dynamic.Append(app.requireAuthentication)
	router.Handle("POST /user/logout", authenticatedChain.ThenFunc(app.logoutPost))

	commonMiddleware := alice.New(app.recoverPanic, app.logRequest, commonHeaders)
	return commonMiddleware.Then(router)
}
