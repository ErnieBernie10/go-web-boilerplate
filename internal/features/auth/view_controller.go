package auth

import (
	"database/sql"
	"framer/internal/database"
	"framer/internal/view"
	"framer/internal/view/layout"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var email = "email"
var password = "password"

func FrameResourceHandler(r chi.Router) {
	r.Get(view.LoginPath, handleGetLogin)
	r.Post(view.LoginPath, handlePostLogin)
	r.Get(view.RegisterPath, handleGetRegister)
	r.Post(view.RegisterPath, handlePostRegister)
}

func handleGetLogin(w http.ResponseWriter, r *http.Request) {
	user := GetUser(r)

	if user != nil {
		http.Redirect(w, r, view.IndexPath, http.StatusSeeOther)
	}

	layout.Unauthenticated(view.RegisterPath, view.LoginPath, loginPage()).Render(r.Context(), w)
}

func handleGetRegister(w http.ResponseWriter, r *http.Request) {
	user := GetUser(r)

	if user != nil {
		http.Redirect(w, r, view.IndexPath, http.StatusSeeOther)
	}

	layout.Unauthenticated(view.RegisterPath, view.LoginPath, registerPage()).Render(r.Context(), w)
}

func handlePostRegister(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
		return
	}

	email := r.FormValue(email)
	pw := r.FormValue(password)

	_, err := database.Service.GetUserByEmail(r.Context(), email)
	if err == nil {
		http.Error(w, "User with e-mail already exists", http.StatusBadRequest)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(pw), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	if err := database.Service.Register(r.Context(), database.RegisterParams{
		Email:        email,
		PasswordHash: sql.NullString{Valid: true, String: string(hashedPassword)},
	}); err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, view.IndexPath, http.StatusSeeOther)
}

func handlePostLogin(w http.ResponseWriter, r *http.Request) {
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
	claims := &Claims{
		Email: email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	// Create the JWT using the claims and the secret key.
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(JwtSecret)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Error creating token", http.StatusInternalServerError)
		return
	}

	// Step 3: Set the JWT token in an HTTP-only cookie.
	http.SetCookie(w, &http.Cookie{
		Name:     UserContextKey, // Cookie name
		Value:    tokenString,    // JWT token value
		Expires:  expirationTime, // Cookie expiration time (same as JWT)
		HttpOnly: true,           // Make the cookie HTTP-only
		Secure:   false,          // Set to true if using HTTPS
		Path:     "/",
	})

	w.Header().Set("HX-Redirect", view.IndexPath)
	w.WriteHeader(http.StatusSeeOther)
}
