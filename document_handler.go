package main


import (

	"cursor-crash-backend/database"
	"cursor-crash-backend/models"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
)

// DocumentInitRequest only includes the language
type DocumentInitRequest struct {
	Language string `json:"language"`
}

// DocumentInitResponse returns the document details
type DocumentInitResponse struct {
	Message  string          `json:"message"`
	Document models.Document `json:"document"`
}

// CreateDocumentHandler handles new doc creation with language only

// CreateDocumentHandler godoc
// @Summary      Create a new document
// @Description  Initializes a new document with the selected language. A temporary name is assigned ("Untitled").
// @Tags         documents
// @Accept       json
// @Produce      json
// @Param        request body DocumentInitRequest true "Document language"
// @Success      201 {object} DocumentInitResponse
// @Failure      400 {string} string "Bad Request"
// @Failure      405 {string} string "Method Not Allowed"
// @Failure      500 {string} string "Internal Server Error"
// @Router       /documents [post]
func CreateDocumentHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	var req DocumentInitRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Bad Request: "+err.Error(), http.StatusBadRequest)
		return
	}

	if req.Language == "" {
		http.Error(w, "Missing language", http.StatusBadRequest)
		return
	}

	now := time.Now().Format(time.RFC3339)
	newDoc := models.Document{
		ID:        uuid.New().String(),
		Name:      "Untitled", // Optionally add numbering on frontend
		Content:   "",
		Language:  req.Language,
		IsSaved:   false,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if result := database.DB.Create(&newDoc); result.Error != nil {
		log.Println("Database error:", result.Error)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(DocumentInitResponse{
		Message:  "Document created successfully",
		Document: newDoc,
	})
}






