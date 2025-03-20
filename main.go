package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Error upgrading to websocket:", err)
		return
	}
	defer conn.Close()

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println("Websocket read error:", err)
			break
		}
		log.Println("Received:", string(msg))
	}
}

func saveDocument(w http.ResponseWriter, r *http.Request) {
	var doc Document
	err := json.NewDecoder(r.Body).Decode(&doc)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	DB.Create(&doc)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(doc)
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "CursorCrash backend is running!")
}

func main() {
	ConnectDatabase()
	
	http.HandleFunc("/", handler)
	http.HandleFunc("/ws", handleWebSocket)
	http.HandleFunc("/save", saveDocument)

	fmt.Println("CursorCrash backend is running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
