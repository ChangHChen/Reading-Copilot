package main

import (
	"html/template"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

func chatHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, "Could not open websocket connection", http.StatusBadRequest)
		return
	}
	defer conn.Close()

	for {
		msgType, msg, err := conn.ReadMessage()
		if err != nil {
			return
		}
		msg = append(msg, []byte("Hello")...)
		if err = conn.WriteMessage(msgType, msg); err != nil {
			return
		}
	}
}

func serveTemplate(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("./chatTest.tmpl")
	if err != nil {
		http.Error(w, "Error loading template", http.StatusInternalServerError)
		return
	}
	tmpl.ExecuteTemplate(w, "chat", nil)
}

func main() {
	http.HandleFunc("/chat", serveTemplate)
	http.HandleFunc("/echo", chatHandler)
	http.ListenAndServe(":8080", nil)
}
