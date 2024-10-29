package auth

import (
	"database/sql"
	"encoding/json"
	"errors"
	"framer/internal/api"
	"framer/internal/core"
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
		api.HandleError(r, w, errors.Join(core.ErrMalformedRequest, err))
		return
	}

}

func handleApiPostLogin(w http.ResponseWriter, r *http.Request) {
	body := &loginCommandDto{}
	if err := json.NewDecoder(r.Body).Decode(body); err != nil {
		api.HandleError(r, w, errors.Join(core.ErrMalformedRequest, err))
		return
	}

	user, err := database.Service.Queries.GetUserByEmail(r.Context(), body.Email)
	if err != nil {
		api.HandleError(r, w, errors.Join(core.ErrUnauthorized, err))
		return
	}

	if !user.PasswordHash.Valid {
		api.HandleError(r, w, errors.Join(core.ErrUnauthorized, errors.New("invalid username or password")))
		return
	}

	tokenString, refreshTokenString, err := login(
		user.ID,
		body.Email,
		body.Password,
		user.PasswordHash.String,
	)

	if err != nil {
		api.HandleError(r, w, errors.Join(core.ErrUnauthorized, errors.New("invalid username or password")))
		return
	}

	encoded, err := json.Marshal(&api.LoginResponseDto{
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
		api.HandleError(r, w, errors.Join(core.ErrMalformedRequest, err))
		return
	}

	_, err := database.Service.Queries.GetUserByEmail(r.Context(), body.Email)
	if err == nil {
		api.HandleError(r, w, errors.Join(core.ErrValidation, errors.New("user with e-mail already exists")))
		return
	}

	hashedPassword, err := hashString(body.Password)
	if err != nil {
		api.HandleError(r, w, err)
		return
	}

	id, err := database.Service.Queries.Register(r.Context(), database.RegisterParams{
		Email:        body.Email,
		PasswordHash: sql.NullString{Valid: true, String: string(hashedPassword)},
	})
	if err != nil {
		api.HandleError(r, w, err)
		return
	}

	api.WriteCreatedResponse(w, api.LoginApiPath, api.CreatedResponseDto{Id: id.String()})
}

type loginCommandDto struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type registerCommandDto struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
