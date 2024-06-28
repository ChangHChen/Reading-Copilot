package main

import "net/http"

func (app *application) routes() *http.ServeMux {
	router := http.NewServeMux()
	fileServer := http.FileServer(http.Dir(app.cfg.staticDir))
	router.Handle("GET /static/", fileServer)

	router.HandleFunc("GET /{$}", app.home)
	router.HandleFunc("GET /user/login", app.login)
	router.HandleFunc("POST /user/login", app.loginPost)
	router.HandleFunc("GET /book/view/{id}", app.bookView)
	router.HandleFunc("GET /user/signup", app.signUp)
	router.HandleFunc("POST /user/signup", app.signUpPost)
	router.HandleFunc("POST /user/logout", app.logoutPost)
	return router
}
