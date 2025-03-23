package main

import (
	"cursor-crash-backend/auth"
	"cursor-crash-backend/database"
	"cursor-crash-backend/docs"
	"cursor-crash-backend/models"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	httpSwagger "github.com/swaggo/http-swagger"
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
	var doc models.Document
	err := json.NewDecoder(r.Body).Decode(&doc)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	database.DB.Create(&doc)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(doc)
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "CursorCrash backend is running!")
}


func CORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		
		// Handle preflight OPTIONS requests
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		
		next.ServeHTTP(w, r)
	})
}



// @title           My API
// @version         1.0
// @description     This is a sample API using Swagger in Go.
// @host           localhost:8080
// @BasePath       /
func main() {
	database.ConnectDatabase()
	

	docs.SwaggerInfo.Title = "Swagger Example API"
	docs.SwaggerInfo.Description = "This is a sample server using Go standard library."
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = "localhost:8080"
	docs.SwaggerInfo.BasePath = "/api"
	docs.SwaggerInfo.Schemes = []string{"http", "https"}

	http.HandleFunc("/", handler)
	http.HandleFunc("/ws", handleWebSocket)
	http.HandleFunc("/save", saveDocument)
	http.Handle("/api/register", CORSMiddleware(http.HandlerFunc(auth.RegisterHandler)))
	http.Handle("/api/login", CORSMiddleware(http.HandlerFunc(auth.LoginHandler)))

	// Serve Swagger documentation
	http.HandleFunc("/swagger/", httpSwagger.WrapHandler)

	fmt.Println("CursorCrash backend is running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
