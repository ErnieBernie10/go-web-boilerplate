package auth

import (
	"database/sql"
	"encoding/json"
	"framer/internal/api"
	"framer/internal/database"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func AuthApiResourceHandler(r chi.Router) {
	r.Post(api.LoginApiPath, handleApiPostLogin)
	r.Post(api.RegisterApiPath, handleApiPostRegister)
}

func handleApiPostLogin(w http.ResponseWriter, r *http.Request) {
	body := &loginCommandDto{}
	if err := json.NewDecoder(r.Body).Decode(body); err != nil {
		api.HandleError(r, w, err)
		return
	}

	user, err := database.Service.GetUserByEmail(r.Context(), body.Email)
	if err != nil {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	if !user.PasswordHash.Valid {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	tokenString, refreshTokenString, err := login(
		body.Email,
		body.Password,
		user.PasswordHash.String,
	)

	if err != nil {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	encoded, err := json.Marshal(&loginResponseDto{
		AccessToken:  tokenString,
		RefreshToken: refreshTokenString,
	})

	if err != nil {
		api.HandleError(r, w, err)
		return
	}

	w.Write(encoded)
	w.WriteHeader(http.StatusCreated)
}

func handleApiPostRegister(w http.ResponseWriter, r *http.Request) {
	body := &registerCommandDto{}
	if err := json.NewDecoder(r.Body).Decode(body); err != nil {
		api.HandleError(r, w, err)
	}

	_, err := database.Service.GetUserByEmail(r.Context(), body.Email)
	if err == nil {
		http.Error(w, "User with e-mail already exists", http.StatusBadRequest)
		return
	}

	hashedPassword, err := hashString(body.Password)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	if err := database.Service.Register(r.Context(), database.RegisterParams{
		Email:        body.Email,
		PasswordHash: sql.NullString{Valid: true, String: string(hashedPassword)},
	}); err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

type loginCommandDto struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type loginResponseDto struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

type registerCommandDto struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
