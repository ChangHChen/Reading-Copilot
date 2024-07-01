package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/ChangHChen/Reading-Copilot/webGateway/internal/models"
	"github.com/ChangHChen/Reading-Copilot/webGateway/internal/validator"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r, nil)
	app.render(w, r, http.StatusOK, "home", data)

}

func (app *application) about(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r, nil)
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
	data := app.newTemplateData(r, nil)
	data.Form = userLoginForm{}
	app.render(w, r, http.StatusOK, "login", data)
}

func (app *application) loginPost(w http.ResponseWriter, r *http.Request) {
	var form userLoginForm
	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form.CheckField(validator.NotBlank(form.Email), "email", "Email cannot be blank")
	form.CheckField(validator.Matches(form.Email, validator.EmailRX), "email", "This field must be a valid email address")
	form.CheckField(validator.NotBlank(form.PWD), "pwd", "Password cannot be blank")

	if !form.Valid() {
		data := app.newTemplateData(r, form)
		app.render(w, r, http.StatusUnprocessableEntity, "login", data)
		return
	}

	id, username, err := app.users.Authenticate(form.Email, form.PWD)
	if err != nil {
		if errors.Is(err, models.ErrInvalidCredentials) {
			form.AddNonFieldError("Email or password is incorrect")
			data := app.newTemplateData(r, form)
			app.render(w, r, http.StatusUnprocessableEntity, "login", data)
		} else {
			app.serverError(w, r, err)
		}
		return
	}
	err = app.sessionManager.RenewToken(r.Context())
	if err != nil {
		app.serverError(w, r, err)
		return
	}
	app.sessionManager.Put(r.Context(), "flash", "Successfully logged in!")
	app.sessionManager.Put(r.Context(), "authenticatedUserID", id)
	app.sessionManager.Put(r.Context(), "authenticatedUserName", username)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (app *application) signUp(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r, nil)
	data.Form = userSignupForm{}
	app.render(w, r, http.StatusOK, "signup", data)
}

func (app *application) signUpPost(w http.ResponseWriter, r *http.Request) {
	var form userSignupForm
	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	form.CheckField(validator.NotBlank(form.UserName), "username", "User Name cannot be blank")
	form.CheckField(validator.NotBlank(form.Email), "email", "Email cannot be blank")
	form.CheckField(validator.Matches(form.Email, validator.EmailRX), "email", "This field must be a valid email address")
	form.CheckField(validator.NotBlank(form.PWD), "pwd", "Password cannot be blank")
	form.CheckField(validator.MinChars(form.PWD, 8), "pwd", "Password must be at least 8 characters long")
	form.CheckField(validator.Repeat(form.PWD, form.PWDConfirm), "pwdconfirm", "Passwords must match")
	if !form.Valid() {
		data := app.newTemplateData(r, form)
		app.render(w, r, http.StatusUnprocessableEntity, "signup", data)
		return
	}
	err = app.users.Insert(form.UserName, form.Email, form.PWD)
	if err != nil {
		if errors.Is(err, models.ErrDuplicateEmail) {
			form.AddFieldError("email", "Email address is already in use")
			data := app.newTemplateData(r, form)
			app.render(w, r, http.StatusUnprocessableEntity, "signup", data)
		} else if errors.Is(err, models.ErrDuplicateUserName) {
			form.AddFieldError("username", "Username is already in use")
			data := app.newTemplateData(r, form)
			app.render(w, r, http.StatusUnprocessableEntity, "signup", data)
		} else {
			app.serverError(w, r, err)
		}
		return
	}

	app.sessionManager.Put(r.Context(), "flash", "User sucessfully signed up, please log in.")
	http.Redirect(w, r, "/user/login", http.StatusSeeOther)
}

func (app *application) logoutPost(w http.ResponseWriter, r *http.Request) {
	err := app.sessionManager.RenewToken(r.Context())
	if err != nil {
		app.serverError(w, r, err)
		return
	}
	app.sessionManager.Remove(r.Context(), "authenticatedUserID")
	app.sessionManager.Remove(r.Context(), "authenticatedUserName")
	app.sessionManager.Put(r.Context(), "flash", "You've been logged out successfully!")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
