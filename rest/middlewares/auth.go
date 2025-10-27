package middlewares

import (
	"context"
	"log" 
	"net/http"
	"strconv"
	"strings"

	"ecoscan.com/utils"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "application/json") // Ensure JSON errors

		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			log.Println("AuthMiddleware: Missing Authorization header")
			http.Error(w, `{"message": "Authorization header required"}`, http.StatusUnauthorized)
			return
		}

		// Ensure the header starts with "Bearer " (case-insensitive)
		parts := strings.Fields(authHeader) // Splits by whitespace
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
             log.Printf("AuthMiddleware: Invalid token format: %s", authHeader)
			 http.Error(w, `{"message": "Invalid token format"}`, http.StatusUnauthorized)
			return
		}
		tokenString := parts[1]


		claims, err := utils.ValidateAccessToken(tokenString)
		if err != nil {
            // Log the specific validation error
			log.Printf("AuthMiddleware: Invalid token: %v", err)
			// Send specific error message if token expired
			if strings.Contains(err.Error(), "token is expired") {
                 http.Error(w, `{"message": "Token has expired"}`, http.StatusUnauthorized)
            } else {
                 http.Error(w, `{"message": "Invalid token"}`, http.StatusUnauthorized)
            }
			return
		}

		// Robust user ID extraction (handle various numeric types)
        var userID int64
        userIDClaim, ok := claims["user_id"]
        if !ok {
            log.Println("AuthMiddleware: user_id claim missing")
            http.Error(w, `{"message": "Invalid token claims (user_id missing)"}`, http.StatusUnauthorized)
            return
        }

        switch v := userIDClaim.(type) {
        case float64:
            userID = int64(v)
        case int64:
            userID = v
        case int:
             userID = int64(v)
        case string: // Handle if user_id was stored as a string
             parsedID, parseErr := strconv.ParseInt(v, 10, 64)
             if parseErr != nil {
                  log.Printf("AuthMiddleware: Could not parse user_id claim string: %v", parseErr)
                  http.Error(w, `{"message": "Invalid token claims (user_id format)"}`, http.StatusUnauthorized)
                  return
             }
             userID = parsedID
        default:
             log.Printf("AuthMiddleware: Unexpected type for user_id claim: %T", v)
             http.Error(w, `{"message": "Invalid token claims (user_id type)"}`, http.StatusUnauthorized)
             return
        }
        
        // Ensure userID is positive
        if userID <= 0 {
             log.Printf("AuthMiddleware: Invalid user_id value: %d", userID)
             http.Error(w, `{"message": "Invalid token claims (user_id value)"}`, http.StatusUnauthorized)
             return
        }


		ctx := context.WithValue(r.Context(), "userID", userID)
		log.Printf("AuthMiddleware: User %d authorized for %s %s", userID, r.Method, r.URL.Path)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}