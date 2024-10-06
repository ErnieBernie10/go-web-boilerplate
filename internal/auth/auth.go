package auth

import (
	"context"
	"net/http"
	"os"

	"github.com/golang-jwt/jwt/v5"
)

const UserContextKey = "User"

var JwtSecret = []byte(os.Getenv("JWT_SECRET"))

type Claims struct {
	Email string `json:"username"`
	jwt.RegisteredClaims
}

func (c *Claims) Valid() error {
	return nil
}
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
		ctx := context.WithValue(r.Context(), UserContextKey, claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func OptionalUserMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Step 1: Get the JWT token from the "token" cookie.
		cookie, err := r.Cookie(UserContextKey)
		if err != nil {
			ctx := context.WithValue(r.Context(), UserContextKey, nil)
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
			ctx := context.WithValue(r.Context(), UserContextKey, nil)
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}

		// Step 3: Store the user information in the request context.
		ctx := context.WithValue(r.Context(), UserContextKey, claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetUser(r *http.Request) *Claims {
	user, ok := r.Context().Value(UserContextKey).(*Claims)

	if !ok {
		return nil
	}

	return user
}
