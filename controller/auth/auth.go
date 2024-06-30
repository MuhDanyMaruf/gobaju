package auth

import (
	"belajar/database"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"belajar/model/users"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

var jwtKey = []byte("your_secret_key")

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

func Registration(w http.ResponseWriter, r *http.Request) {
	var creds users.Users
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		log.Printf("Error decoding request body: %v", err)
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	var existingUser users.Users
	err = database.DB.QueryRow("SELECT username FROM users WHERE username = ?", creds.Username).Scan(&existingUser.Username)

	if err != nil && err != sql.ErrNoRows {
		log.Printf("Error querying database: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	if existingUser.Username != "" {
		http.Error(w, "Username already exists", http.StatusConflict)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(creds.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Error hashing password: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	_, err = database.DB.Exec("INSERT INTO users (username, email, password) VALUES (?, ?, ?)", creds.Username, creds.Email, hashedPassword)
	if err != nil {
		log.Printf("Error inserting into database: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	response := map[string]interface{}{
		"message": "User registered successfully",
	}
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Printf("Error encoding response: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

func Login(w http.ResponseWriter, r *http.Request) {
	var creds users.Users
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		log.Printf("Error decoding request body: %v", err)
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	var user users.Users
	err = database.DB.QueryRow("SELECT users_id, username, password, email FROM users WHERE username = ?", creds.Username).Scan(&user.UsersId, &user.Username, &user.Password, &user.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "User not found", http.StatusUnauthorized)
			return
		}
		log.Printf("Error querying database: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(creds.Password))
	if err != nil {
		http.Error(w, "Invalid password", http.StatusUnauthorized)
		return
	}

	expirationTime := time.Now().Add(120 * time.Minute)
	claims := &Claims{
		Username: creds.Username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		log.Printf("Error signing token: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	response := map[string]interface{}{
		"message": "Login successful",
		"token":   tokenString,
	}
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Printf("Error encoding response: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

func ValidateToken(tokenString string) (bool, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil {
		return false, err
	}

	return token.Valid, nil
}

func JWTAuth(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		bearerToken := r.Header.Get("Authorization")
		sttArr := strings.Split(bearerToken, " ")
		if len(sttArr) == 2 {
			isValid, err := ValidateToken(sttArr[1])
			if err != nil {
				log.Printf("Token validation error: %v", err)
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
			if isValid {
				next.ServeHTTP(w, r)
			} else {
				log.Printf("Invalid token")
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
			}
		} else {
			log.Printf("Invalid Authorization header format")
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
		}
	})
}
