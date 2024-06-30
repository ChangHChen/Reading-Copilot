package main

import (
	"net/http"

	"github.com/justinas/alice"
)

func (app *application) routes(staticDir string) http.Handler {
	router := http.NewServeMux()
	fileServer := http.FileServer(http.Dir(staticDir))
	router.Handle("GET /static/", fileServer)

	dynamic := alice.New(app.sessionManager.LoadAndSave)
	router.Handle("GET /{$}", dynamic.ThenFunc(app.home))
	router.Handle("GET /about", dynamic.ThenFunc(app.about))
	router.Handle("GET /user/login", dynamic.ThenFunc(app.login))
	router.Handle("POST /user/login", dynamic.ThenFunc(app.loginPost))
	router.Handle("GET /book/view/{id}", dynamic.ThenFunc(app.bookView))
	router.Handle("GET /user/signup", dynamic.ThenFunc(app.signUp))
	router.Handle("POST /user/signup", dynamic.ThenFunc(app.signUpPost))
	router.Handle("POST /user/logout", dynamic.ThenFunc(app.logoutPost))

	commonMiddleware := alice.New(app.recoverPanic, app.logRequest, commonHeaders)
	return commonMiddleware.Then(router)
}
