package auth

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"framer/internal/api"
	"framer/internal/database"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func AuthApiResourceHandler(r chi.Router) {
	r.Post(api.LoginApiPath, handleApiPostLogin)
	r.Post(api.RefreshApiPath, handleApiPostRefresh)
	r.Post(api.RegisterApiPath, handleApiPostRegister)
}

func handleApiPostRefresh(w http.ResponseWriter, r *http.Request) {
	body := &api.LoginResponseDto{}
	if err := json.NewDecoder(r.Body).Decode(body); err != nil {
		api.HandleError(r, w, err, http.StatusBadRequest)
		return
	}

}

func handleApiPostLogin(w http.ResponseWriter, r *http.Request) {
	body := &loginCommandDto{}
	if err := json.NewDecoder(r.Body).Decode(body); err != nil {
		api.HandleError(r, w, err, http.StatusBadRequest)
		return
	}

	user, err := database.Service.GetUserByEmail(r.Context(), body.Email)
	if err != nil {
		api.HandleError(r, w, fmt.Errorf("invalid username or password"), http.StatusUnauthorized)
		return
	}

	if !user.PasswordHash.Valid {
		api.HandleError(r, w, fmt.Errorf("invalid username or password"), http.StatusUnauthorized)
		return
	}

	tokenString, refreshTokenString, err := login(
		body.Email,
		body.Password,
		user.PasswordHash.String,
	)

	if err != nil {
		api.HandleError(r, w, fmt.Errorf("invalid username or password"), http.StatusUnauthorized)
		return
	}

	encoded, err := json.Marshal(&api.LoginResponseDto{
		AccessToken:  tokenString,
		RefreshToken: refreshTokenString,
	})

	if err != nil {
		api.HandleError(r, w, err, http.StatusInternalServerError)
		return
	}

	w.Write(encoded)
	w.WriteHeader(http.StatusCreated)
}

func handleApiPostRegister(w http.ResponseWriter, r *http.Request) {
	body := &registerCommandDto{}
	if err := json.NewDecoder(r.Body).Decode(body); err != nil {
		api.HandleError(r, w, err, http.StatusInternalServerError)
		return
	}

	_, err := database.Service.GetUserByEmail(r.Context(), body.Email)
	if err == nil {
		api.HandleError(r, w, fmt.Errorf("user with e-mail already exists"), http.StatusBadRequest)
		return
	}

	hashedPassword, err := hashString(body.Password)
	if err != nil {
		api.HandleError(r, w, err, http.StatusInternalServerError)
		return
	}

	if err := database.Service.Register(r.Context(), database.RegisterParams{
		Email:        body.Email,
		PasswordHash: sql.NullString{Valid: true, String: string(hashedPassword)},
	}); err != nil {
		api.HandleError(r, w, err, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

type loginCommandDto struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type registerCommandDto struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
