package middleware

import (
	"car-mobile-project/internal/models"
	"encoding/json"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"strings"
)

func AuthJwtMiddleware(secretKey string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var responseObject *models.ResponseObject
			tokenString := r.Header.Get("Authorization")

			if tokenString == "" {
				w.WriteHeader(http.StatusUnauthorized)
				responseObject = models.NewResponseObject(false, "invalid token", nil)
				jsonResponse, _ := json.Marshal(responseObject)
				w.Write(jsonResponse)
				return
			}

			splitToken := strings.Split(tokenString, "Bearer ")
			reqToken := splitToken[1]

			token, err := jwt.Parse(reqToken, func(token *jwt.Token) (interface{}, error) {
				return []byte(secretKey), nil
			})

			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				responseObject = models.NewResponseObject(false, "invalid token", nil)
				jsonResponse, _ := json.Marshal(responseObject)
				w.Write(jsonResponse)
				return
			}

			if !token.Valid {
				w.WriteHeader(http.StatusUnauthorized)
				responseObject = models.NewResponseObject(false, "invalid token", nil)
				jsonResponse, _ := json.Marshal(responseObject)
				w.Write(jsonResponse)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
