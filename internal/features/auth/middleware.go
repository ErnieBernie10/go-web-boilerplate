package auth

import (
	"context"
	"net/http"

	"github.com/golang-jwt/jwt"
)

func AuthGuardMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Step 1: Get the JWT token from the "token" cookie.
		cookie, err := r.Cookie("token")
		if err != nil {
			if err == http.ErrNoCookie {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}
		tokenStr := cookie.Value

		// Step 2: Parse and validate the JWT token.
		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
			return JwtSecret, nil
		})
		if err != nil || !token.Valid {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Step 3: Store the user information in the request context.
		ctx := context.WithValue(r.Context(), TokenContextKey, claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func OptionalUserMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Step 1: Get the JWT token from the "token" cookie.
		cookie, err := r.Cookie(string(TokenContextKey))
		if err != nil {
			ctx := context.WithValue(r.Context(), TokenContextKey, nil)
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}
		tokenStr := cookie.Value

		// Step 2: Parse and validate the JWT token.
		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
			return JwtSecret, nil
		})
		if err != nil || !token.Valid {
			ctx := context.WithValue(r.Context(), TokenContextKey, nil)
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}

		// Step 3: Store the user information in the request context.
		ctx := context.WithValue(r.Context(), TokenContextKey, claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetUser(r *http.Request) *Claims {
	user, ok := r.Context().Value(TokenContextKey).(*Claims)

	if !ok {
		return nil
	}

	return user
}
