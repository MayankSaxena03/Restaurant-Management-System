package middleware

import (
	"encoding/json"
	"net/http"

	"github.com/MayankSaxena03/Restaurant-Management-System/helpers"
)

// Middleware for mux router
func Authentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get the bearer token from the header
		token := r.Header.Get("token")
		if token == "" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode("Unauthorized")
			return
		}
		// Validate the token
		claims, err := helpers.ValidateToken(token)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode("Unauthorized")
			return
		}

		r.Header.Set("userId", claims.UserID.Hex())
		r.Header.Set("email", claims.Email)
		r.Header.Set("username", claims.Username)
		next.ServeHTTP(w, r)
	})
}
