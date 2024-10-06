package view

import (
	"database/sql"
	"framer/internal/auth"
	"framer/internal/database"
	"framer/internal/view/layout"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

var email = "email"
var password = "password"

func HandleGetRegister(w http.ResponseWriter, r *http.Request) {
	user := auth.GetUser(r)

	if user != nil {
		http.Redirect(w, r, IndexPath, http.StatusSeeOther)
	}

	layout.Unauthenticated(RegisterPath, LoginPath, registerPage()).Render(r.Context(), w)
}

func HandleRegister(w http.ResponseWriter, r *http.Request) {
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

	http.Redirect(w, r, IndexPath, http.StatusSeeOther)
}
