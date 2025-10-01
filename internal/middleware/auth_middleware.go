package middleware

import (
	"context"
	"net/http"
	"os"
	"strings"
	"strconv"
	"log"

	"github.com/dgrijalva/jwt-go"
)

type Claims struct {
	UserID int    `json:"user_id"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	jwt.StandardClaims
}
// getEnvAsInt obtiene una variable de entorno como entero
func getEnvAsInt(name string, defaultVal int) int {
	valStr := os.Getenv(name)
	if val, err := strconv.Atoi(valStr); err == nil {
		return val
	}
	return defaultVal
}

type contextKey string

const ContextUserKey contextKey = "user"

// AuthMiddleware verifica el header Authorization: Bearer <token>
// y añade las claims al contexto de la request.
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header missing", http.StatusUnauthorized)
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			http.Error(w, "Authorization header format must be: Bearer {token}", http.StatusUnauthorized)
			return
		}
		secret := os.Getenv("JWT_SECRET")
		if secret == "" {
			log.Println("⚠️ JWT_SECRET no está seteado en el middleware")
		} else {
			log.Println("🔑 JWT_SECRET cargado en middleware:", secret)
		}


		tokenStr := parts[1]
		claims := &Claims{}
		var token *jwt.Token
		var err error
		// Parse token using the JWT secret available at request time.
		// Handle error and invalid token separately to avoid calling err.Error() when err is nil.
		// 'secret' was already read above and is reused here.
		token, err = jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(secret), nil
		})
		if err != nil {
			http.Error(w, "Token inválido: "+err.Error(), http.StatusUnauthorized)
			return
		}
		if !token.Valid {
			http.Error(w, "Token inválido", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), ContextUserKey, claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// GetClaims extrae las claims del contexto (útil desde handlers)
func GetClaims(r *http.Request) (*Claims, bool) {
	c, ok := r.Context().Value(ContextUserKey).(*Claims)
	return c, ok
}
func authorizeRole(next http.HandlerFunc, roles ...string) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        claims, ok := r.Context().Value("claims").(*Claims)
        if !ok {
            http.Error(w, "Unauthorized", http.StatusUnauthorized)
            return
        }

        for _, role := range roles {
            if claims.Role == role {
                next.ServeHTTP(w, r)
                return
            }
        }

        http.Error(w, "Forbidden", http.StatusForbidden)
    }
}
