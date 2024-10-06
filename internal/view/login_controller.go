package view

import (
	"framer/internal/auth"
	"framer/internal/database"
	"framer/internal/view/layout"
	"log"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func HandleGetLogin(w http.ResponseWriter, r *http.Request) {
	user := auth.GetUser(r)

	if user != nil {
		http.Redirect(w, r, IndexPath, http.StatusSeeOther)
	}

	layout.Unauthenticated(RegisterPath, LoginPath, loginPage()).Render(r.Context(), w)
}

func HandleLogin(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
		return
	}

	email := r.FormValue(email)
	pw := r.FormValue(password)

	user, err := database.Service.GetUserByEmail(r.Context(), email)
	if err != nil {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	if !user.PasswordHash.Valid {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash.String), []byte(pw))
	if err != nil {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	expirationTime := time.Now().Add(5 * time.Minute) // Set token expiration time.
	claims := &auth.Claims{
		Email: email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	// Create the JWT using the claims and the secret key.
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(auth.JwtSecret)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Error creating token", http.StatusInternalServerError)
		return
	}

	// Step 3: Set the JWT token in an HTTP-only cookie.
	http.SetCookie(w, &http.Cookie{
		Name:     auth.UserContextKey, // Cookie name
		Value:    tokenString,         // JWT token value
		Expires:  expirationTime,      // Cookie expiration time (same as JWT)
		HttpOnly: true,                // Make the cookie HTTP-only
		Secure:   false,               // Set to true if using HTTPS
		Path:     "/",
	})

	w.Header().Set("HX-Redirect", IndexPath)
	w.WriteHeader(http.StatusSeeOther)
}
