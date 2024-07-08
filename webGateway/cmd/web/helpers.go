package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"time"

	apis "github.com/ChangHChen/Reading-Copilot/webGateway/internal/APIs"
	"github.com/ChangHChen/Reading-Copilot/webGateway/internal/models"
	"github.com/alexedwards/scs/mysqlstore"
	"github.com/alexedwards/scs/v2"
	"github.com/go-playground/form/v4"
	_ "github.com/go-sql-driver/mysql"
	"github.com/justinas/nosurf"
)

func setup(cfg config) *application {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))
	db, err := openDB(cfg.dsn)
	if err != nil {
		fatalError(logger, "Errors occured when connecting to the DB", err)
	}
	htmlTemplateCache, err := newHtmlTemplateCache()
	if err != nil {
		fatalError(logger, "Errors occured when preparing html pages", err)
	}
	sessionManager := scs.New()
	sessionManager.Store = mysqlstore.New(db)
	sessionManager.Lifetime = 12 * time.Hour
	sessionManager.Cookie.Secure = true

	app := &application{
		logger:            logger,
		db:                db,
		users:             &models.UserModel{DB: db},
		books:             &models.BookModel{DB: db},
		htmlTemplateCache: htmlTemplateCache,
		formDecoder:       form.NewDecoder(),
		sessionManager:    sessionManager,
	}
	app.router = app.routes()
	return app
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		db.Close()
		return nil, err
	}
	return db, nil
}

func (app *application) render(w http.ResponseWriter, r *http.Request, status int, pageName string, data templateData) {
	ts, ok := app.htmlTemplateCache[pageName]
	if !ok {
		err := fmt.Errorf("page template %s does not exist", pageName)
		app.serverError(w, r, err)
		return
	}

	buf := new(bytes.Buffer)
	err := ts.ExecuteTemplate(buf, "base", data)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	w.WriteHeader(status)
	buf.WriteTo(w)
}

func (app *application) newTemplateData(r *http.Request, form any) templateData {
	newData := templateData{
		CurYear:         time.Now().Year(),
		Flash:           app.sessionManager.PopString(r.Context(), "flash"),
		Form:            form,
		IsAuthenticated: app.isAuthenticated(r),
		CSRFToken:       nosurf.Token(r),
		BookList:        bookList{},
		Book:            models.BookMeta{},
		CurPage:         1,
		APIKeys:         apiKeys{},
	}
	if newData.IsAuthenticated {
		newData.UserName = app.sessionManager.GetString(r.Context(), "authenticatedUserName")
	}
	return newData
}

func (app *application) decodePostForm(r *http.Request, dst any) error {
	err := r.ParseForm()
	if err != nil {
		return err
	}
	err = app.formDecoder.Decode(dst, r.PostForm)
	if err != nil {
		var invalidDecoderError *form.InvalidDecoderError
		if errors.As(err, &invalidDecoderError) {
			panic(err)
		}
		return err
	}
	return nil
}

func (app *application) isAuthenticated(r *http.Request) bool {
	isAuthenticated, ok := r.Context().Value(isAuthenticatedContextKey).(bool)
	if !ok {
		return false
	}
	return isAuthenticated

}

func (app *application) redirectToLastURL(w http.ResponseWriter, r *http.Request) {
	lastURL := app.sessionManager.GetString(r.Context(), "lastURL")
	if lastURL == "" {
		lastURL = "/"
	}
	http.Redirect(w, r, lastURL, http.StatusSeeOther)
}

func processWithLLM(msg ChatMessage, bookID int) (string, error) {
	requestData := apis.ChatRequest{
		Model:       msg.Model,
		BookID:      bookID,
		Progress:    msg.Page,
		UserMessage: msg.Message,
	}

	requestJson, err := json.Marshal(requestData)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", "https://localhost:4010/chat", bytes.NewBuffer(requestJson))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	responseJson, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	var responseData apis.ChatResponse
	if err := json.Unmarshal([]byte(responseJson), &responseData); err != nil {
		return "", err
	}

	if responseData.Error != "" {
		return "", errors.New(responseData.Error)
	}
	return responseData.ResponseMessage, nil
}

func buildUpBook(bookID int) (string, error) {
	requestData := apis.BuildBookRequest{
		BookID: bookID,
	}

	requestJson, err := json.Marshal(requestData)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", "https://localhost:4010/build", bytes.NewBuffer(requestJson))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	responseJson, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	var responseData apis.BuildBookResponse
	if err := json.Unmarshal([]byte(responseJson), &responseData); err != nil {
		return "", err
	}
	return responseData.Msg, nil
}
