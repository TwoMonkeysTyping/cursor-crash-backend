package auth

import (
	"cursor-crash-backend/database"
	"cursor-crash-backend/models"
	"encoding/json"
	"log"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

type LoginRequest struct {
	Email		string `json:"email"`
	Password	string `json:"password"`
}

type LoginResponse struct {
	Message 	string `json:"message"`
}

type RegisterRequest struct {
	Email 			string `json:"email"`
	Password		string `json:"password"`
}


type RegisterResponse struct {
	Message 		string `json:"message"`
}

// LoginHandler godoc
// @Summary      User login
// @Description  Authenticates a user with email and password
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request body LoginRequest true "Login credentials"
// @Success      200  {object}  LoginResponse
// @Failure      400  {string}  string "Bad Request"
// @Failure      401  {object}  LoginResponse
// @Failure      405  {string}  string "Method Not Allowed"
// @Router       /login [post]
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	var loginReq LoginRequest
	err := json.NewDecoder(r.Body).Decode(&loginReq)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	// Find user by email
	var user models.User
	result := database.DB.Where("email = ?", loginReq.Email).First(&user)
	if result.Error != nil {
		// Don't reveal whether a user exists or not for security reasons
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(LoginResponse{Message: "Invalid credentials"})
		return
	}

	// Compare the stored hashed password with the provided password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginReq.Password))
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(LoginResponse{Message: "Invalid credentials"})
		return
	}

	// JWT is yet to be implemented
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(LoginResponse{
		Message: "Login successful",
	})
}



// RegisterHandler godoc
// @Summary      User registration
// @Description  Registers a new user with email and password
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request body RegisterRequest true "Registration information"
// @Success      201  {object}  RegisterResponse
// @Failure      400  {string}  string "Bad Request"
// @Failure      405  {string}  string "Method Not Allowed"
// @Failure      409  {string}  string "Email already registered"
// @Failure      500  {string}  string "Internal Server Error"
// @Router       /register [post]
func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	var registerReq RegisterRequest
	err := json.NewDecoder(r.Body).Decode(&registerReq)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	// Check if email already exists
	var count int64
	result := database.DB.Model(&models.User{}).Where("email = ?", registerReq.Email).Count(&count)
	if result.Error != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Println("Database error:", result.Error)
		return
	}
	exists := count > 0
	if exists {
		http.Error(w, "Email already registered", http.StatusConflict)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(registerReq.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Println("Password hashing error:", err)
		return
	}

	// Create new user with GORM
	newUser := models.User{
		Email:    registerReq.Email,
		Password: string(hashedPassword),
	}
	result = database.DB.Create(&newUser)
	if result.Error != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Println("Database insert error:", result.Error)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(RegisterResponse{Message: "User registered successfully"})
}