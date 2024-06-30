package main

import (
	"fmt"
	"net/http"
	"strconv"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	app.render(w, r, http.StatusOK, "home", data)

}

func (app *application) about(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	app.render(w, r, http.StatusOK, "about", data)
}

func (app *application) bookView(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id < 1 {
		app.clientError(w, http.StatusNotFound)
		return
	}
	fmt.Fprintf(w, "Reading book %d\n", id)

}

func (app *application) login(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Page for users to log in")
}

func (app *application) loginPost(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Actual log in postinng")
}

func (app *application) signUp(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Page for users to sign up")
}

func (app *application) signUpPost(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Actual sign up posting")
}

func (app *application) logoutPost(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Log out posting")
}
