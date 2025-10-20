package middlewares

import (
	"context"
	"net/http"
	"strings"

	"ecoscan.com/utils"
)

func AuthMiddleware(next http.Handler) http.Handler{

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
		// reading authorization header : bearer <token>
		authHeader := r.Header.Get("Authorization")

		if authHeader == ""{
			http.Error(w, "Authorization header required", http.StatusUnauthorized)
			return 
		}

		// separate the bearer and token
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader{ // if there is no 'Bearer '
			http.Error(w, "Invalid token format", http.StatusUnauthorized)
			return 
		}

		// validation of token and returns the payload
		claims, err := utils.ValidateAccessToken(tokenString)
		if err != nil{
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		// extract user id

		userIDFloat, ok := claims["user_id"].(float64) //jwt pass numbers as float64

		if !ok{
			http.Error(w, "invalid token claims (user_id not found)", http.StatusUnauthorized)
			return 
		}
		// context to the new handler as userID 
		ctx := context.WithValue(r.Context(), "userID", int64(userIDFloat))

		next.ServeHTTP(w, r.WithContext(ctx))






	})
}