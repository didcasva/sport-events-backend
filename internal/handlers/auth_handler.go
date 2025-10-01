package handlers

import (
	"encoding/json"
	"net/http"
	"os"
	"time"
	"strconv"
	"log"

	"github.com/dgrijalva/jwt-go"
	"sport-events-backend/internal/services"
)
// getEnvAsInt obtiene una variable de entorno como entero
func getEnvAsInt(name string, defaultVal int) int {
	valStr := os.Getenv(name)
	if val, err := strconv.Atoi(valStr); err == nil {
		return val
	}
	return defaultVal
}
	

// Avoid reading JWT_SECRET at import time; read at runtime when creating tokens.
// jwtKey is intentionally not a package-level variable to ensure changes to the
// environment (or loading order) are respected.

type Claims struct {
	UserID int    `json:"user_id"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	jwt.StandardClaims
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
		Role     string `json:"role"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	err := services.RegisterUser(input.Name, input.Email, input.Password, input.Role)
	if err != nil {
		http.Error(w, "Error registrando usuario: "+err.Error(), 500)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Usuario creado"})
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	user, err := services.AuthenticateUser(input.Email, input.Password)
	if err != nil {
		http.Error(w, "Credenciales inválidas", 401)
		return
	}

	// crear token JWT
	expiration := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		UserID: user.ID,
		Email:  user.Email,
		Role:   user.Role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expiration.Unix(),
		},
	}
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		log.Println("⚠️ JWT_SECRET authn no está seteado en auth handler")
	} else {
		log.Println("🔑 JWT_SECRET cargado en auth handler: <redacted>")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// Use the runtime secret to sign the token so it matches validation in middleware.
	tokenStr, err := token.SignedString([]byte(secret))
	if err != nil {
		http.Error(w, "Error creando token", 500)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
	"token": tokenStr,
	"user": map[string]interface{}{
		"id":    user.ID,
		"email": user.Email,
		"role":  user.Role,
	},
})

}
