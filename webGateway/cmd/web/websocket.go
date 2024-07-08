package main

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gorilla/websocket"
)

type GeneralMessage struct {
	Type string `json:"type"`
}
type ChatMessage struct {
	Message string `json:"message"`
	Page    int    `json:"page"`
	Model   string `json:"model"`
}
type ProgressMessage struct {
	Page int `json:"page"`
}

func (app *application) bookWebSocketHandler(w http.ResponseWriter, r *http.Request) {
	app.logger.Debug("Starting websocket")
	bookID, err := strconv.Atoi(r.PathValue("id"))

	if err != nil || bookID < 1 {
		app.clientError(w, http.StatusNotFound)
		return
	}

	var upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	userID := app.sessionManager.GetInt(r.Context(), "authenticatedUserID")
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		app.serverError(w, r, err)
		return
	}
	defer conn.Close()

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				app.serverError(w, r, err)
			}
			break
		}

		var generalMsg GeneralMessage
		if err := json.Unmarshal(message, &generalMsg); err != nil {
			app.serverError(w, r, err)
			continue
		}
		switch generalMsg.Type {
		case "chat":
			var msg ChatMessage
			if err := json.Unmarshal(message, &msg); err != nil {
				app.serverError(w, r, err)
				return
			}

			app.logger.Debug("Received message:", slog.String("message", msg.Message), slog.Int("page", msg.Page))

			response, err := processWithLLM(msg, bookID)
			if err != nil {
				app.serverError(w, r, err)
				return
			}
			app.logger.Debug("Sending response:", slog.String("response", response))
			if err = conn.WriteMessage(websocket.TextMessage, []byte(response)); err != nil {
				app.serverError(w, r, err)
				return
			}
		case "progress":
			var curPageNum ProgressMessage
			if err := json.Unmarshal(message, &curPageNum); err != nil {
				app.serverError(w, r, err)
				return
			}
			app.logger.Debug("Reading progress update", slog.Int("user", userID), slog.Int("book", bookID), slog.Int("page", curPageNum.Page))

			if err = app.users.UpdateReadingProgress(userID, bookID, curPageNum.Page); err != nil {
				app.serverError(w, r, err)
				return
			}
		}
	}
}
