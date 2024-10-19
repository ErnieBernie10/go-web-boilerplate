package server

import (
	"context"
	"framer/internal/core"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

func OptionalUserMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Step 1: Get the JWT token from the "token" tokenCookie.
		refreshStr := getCookieValue(r, string(core.RefreshContextKey))
		tokenStr := getTokenString(r)

		if tokenStr == "" {
			ctx := context.WithValue(r.Context(), core.UserContextKey, nil)
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}

		// Step 2: Parse and validate the JWT token.
		claims := &core.Claims{}
		token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
			return core.JwtSecret, nil
		})
		if err != nil || !token.Valid {
			ctx := context.WithValue(r.Context(), core.UserContextKey, nil)
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}

		// Step 3: Store the user information in the request context.
		ctx := context.WithValue(r.Context(), core.UserContextKey, claims)
		ctx = context.WithValue(ctx, core.TokenContextKey, tokenStr)
		ctx = context.WithValue(ctx, core.RefreshContextKey, refreshStr)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func AuthGuardMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Step 1: Get the JWT token from the "token" tokenCookie.
		refreshStr := getCookieValue(r, string(core.RefreshContextKey))
		tokenStr := getTokenString(r)

		if tokenStr == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Step 2: Parse and validate the JWT token.
		claims := &core.Claims{}
		token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
			return core.JwtSecret, nil
		})
		if err != nil || !token.Valid {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Step 3: Store the user information in the request context.
		ctx := context.WithValue(r.Context(), core.UserContextKey, claims)
		ctx = context.WithValue(ctx, core.TokenContextKey, tokenStr)
		ctx = context.WithValue(ctx, core.RefreshContextKey, refreshStr)
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
	tokenCookie, err := r.Cookie(string(core.TokenContextKey))
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
		accessToken, _ := r.Context().Value(core.TokenContextKey).(string)
		refreshToken, _ := r.Context().Value(core.RefreshContextKey).(string)

		return accessToken, refreshToken
	}
}
