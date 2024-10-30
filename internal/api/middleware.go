package api

import (
	"context"
	"framer/internal/pkg"
	"framer/internal/pkg/view"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

func OptionalUserMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Step 1: Get the JWT token from the "token" tokenCookie.
		refreshStr := getCookieValue(r, string(pkg.RefreshContextKey))
		tokenStr := getTokenString(r)

		if tokenStr == "" {
			ctx := context.WithValue(r.Context(), pkg.UserContextKey, nil)
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}

		// Step 2: Parse and validate the JWT token.
		claims := &pkg.Claims{}
		token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
			return pkg.JwtSecret, nil
		})
		if err != nil || !token.Valid {
			ctx := context.WithValue(r.Context(), pkg.UserContextKey, nil)
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}

		// Step 3: Store the user information in the request context.
		ctx := context.WithValue(r.Context(), pkg.UserContextKey, claims)
		ctx = context.WithValue(ctx, pkg.TokenContextKey, tokenStr)
		ctx = context.WithValue(ctx, pkg.RefreshContextKey, refreshStr)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func AuthGuardMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Step 1: Get the JWT token from the "token" tokenCookie.
		refreshStr := getCookieValue(r, string(pkg.RefreshContextKey))
		tokenStr := getTokenString(r)

		if tokenStr == "" {
			if r.Header.Get("Content-Type") != "application/json" {
				http.Redirect(w, r, view.LoginPath, http.StatusSeeOther)
				return
			}
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Step 2: Parse and validate the JWT token.
		claims := &pkg.Claims{}
		token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
			return pkg.JwtSecret, nil
		})
		if err != nil || !token.Valid {
			if r.Header.Get("Content-Type") != "application/json" {
				http.Redirect(w, r, view.LoginPath, http.StatusSeeOther)
				return
			}
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Step 3: Store the user information in the request context.
		ctx := context.WithValue(r.Context(), pkg.UserContextKey, claims)
		ctx = context.WithValue(ctx, pkg.TokenContextKey, tokenStr)
		ctx = context.WithValue(ctx, pkg.RefreshContextKey, refreshStr)
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
	tokenCookie, err := r.Cookie(string(pkg.TokenContextKey))
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
