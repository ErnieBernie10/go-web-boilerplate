package pkg

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var JwtSecret = []byte(os.Getenv("JWT_SECRET"))

type ContextKey string

const UserContextKey ContextKey = "User"
const TokenContextKey ContextKey = "Token"
const RefreshContextKey ContextKey = "Refresh"

func OptionalUserMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Step 1: Get the JWT token from the "token" tokenCookie.
		refreshStr := getCookieValue(r, string(RefreshContextKey))
		tokenStr := getTokenString(r)

		if tokenStr == "" {
			ctx := context.WithValue(r.Context(), UserContextKey, nil)
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}

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
		ctx = context.WithValue(ctx, TokenContextKey, tokenStr)
		ctx = context.WithValue(ctx, RefreshContextKey, refreshStr)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func AuthGuardMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Step 1: Get the JWT token from the "token" tokenCookie.
		refreshStr := getCookieValue(r, string(RefreshContextKey))
		tokenStr := getTokenString(r)

		if tokenStr == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

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
		ctx = context.WithValue(ctx, TokenContextKey, tokenStr)
		ctx = context.WithValue(ctx, RefreshContextKey, refreshStr)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func getCookieValue(r *http.Request, key string) string {
	cookie, err := r.Cookie(key)
	if err != nil {
		return ""
	}
	return cookie.Value
}

func getTokenString(r *http.Request) string {
	tokenCookie, err := r.Cookie(string(TokenContextKey))
	if err != nil {
		if err == http.ErrNoCookie {
			h := r.Header.Get("Authorization")
			if h != "" {
				return strings.Split(h, " ")[1]
			}
		}
		return ""
	}
	return tokenCookie.Value
}

func GetTokens(r *http.Request) func() (string, string) {
	return func() (string, string) {
		accessToken, _ := r.Context().Value(TokenContextKey).(string)
		refreshToken, _ := r.Context().Value(RefreshContextKey).(string)

		return accessToken, refreshToken
	}
}

func GetUser(r *http.Request) *Claims {
	user, ok := r.Context().Value(UserContextKey).(*Claims)

	if !ok {
		return nil
	}

	return user
}

type Claims struct {
	Email string    `json:"username"`
	ID    uuid.UUID `json:"id"`
	jwt.RegisteredClaims
}

func (c *Claims) Valid() error {
	if !time.Now().Before(c.ExpiresAt.Time) {
		return fmt.Errorf("the token has expired")
	}
	return nil
}
