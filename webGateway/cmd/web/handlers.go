package main

import (
	"fmt"
	"net/http"
	"strconv"
	"text/template"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	files := []string{
		"./ui/html/base.tmpl",
		"./ui/html/pages/home.tmpl",
		"./ui/html/partials/nav.tmpl",
	}

	ts, err := template.ParseFiles(files...)
	if err != nil {
		app.serverError(w, r, err)
		return

	}
	err = ts.ExecuteTemplate(w, "base", nil)
	if err != nil {
		app.serverError(w, r, err)
	}

}

func (app *application) about(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Something about this applicatioin")
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
