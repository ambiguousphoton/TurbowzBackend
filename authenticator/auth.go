package authenticator

import (
	"net/http"
	"GoServer/models"
	"context"
	"log"
	"github.com/golang-jwt/jwt/v5"
)


var jwtKey = []byte("om namo bhagwate vaudevay") // TODO: move to environment variable


func RequireAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

			tokenString := r.Header.Get("Authorization")

			if tokenString == "" {
				http.Error(w, "missing Authentication Token", http.StatusUnauthorized)
				return
			}

			claims := &models.Claims{}
			token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
				return jwtKey, nil
			})

			if err != nil || !token.Valid {
				log.Printf("RequireAuth: invalid or expired token - %v", err)
				http.Error(w, "invalid token", http.StatusUnauthorized)
				return
			}


			allowed := context.WithValue(r.Context(), "userID", claims.UserID)
			// Continue if valid
			log.Println("Authenticating Request successful for userID: ", claims.UserID)
			next.ServeHTTP(w, r.WithContext(allowed))
		}
}