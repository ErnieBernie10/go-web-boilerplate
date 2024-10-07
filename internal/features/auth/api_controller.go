package auth

import (
	"database/sql"
	"encoding/json"
	"framer/internal/api"
	"framer/internal/core"
	"framer/internal/database"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func AuthApiResourceHandler(r chi.Router) {
	r.Post(api.LoginApiPath, handleApiPostLogin)
	r.Post(api.RegisterApiPath, handleApiPostRegister)
}

func handleApiPostLogin(w http.ResponseWriter, r *http.Request) {
	body := &loginCommandDto{}
	if err := json.NewDecoder(r.Body).Decode(body); err != nil {
		core.HandleError(w, err)
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

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash.String), []byte(body.Password))
	if err != nil {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	expirationTime := time.Now().Add(5 * time.Minute) // Set token expiration time.
	claims := &Claims{
		Email: body.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	refreshClaim := &Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	// Create the JWT using the claims and the secret key.
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	refresh := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaim)

	tokenString, err := token.SignedString(JwtSecret)
	if err != nil {
		core.HandleError(w, err)
		return
	}

	refreshTokenString, err := refresh.SignedString(JwtSecret)
	if err != nil {
		core.HandleError(w, err)
		return
	}

	encoded, err := json.Marshal(&loginResponseDto{AccessToken: tokenString, RefreshToken: refreshTokenString})
	if err != nil {
		core.HandleError(w, err)
		return
	}

	w.Write(encoded)
	w.WriteHeader(http.StatusCreated)
}

func handleApiPostRegister(w http.ResponseWriter, r *http.Request) {
	body := &registerCommandDto{}
	if err := json.NewDecoder(r.Body).Decode(body); err != nil {
		core.HandleError(w, err)
	}

	_, err := database.Service.GetUserByEmail(r.Context(), body.Email)
	if err == nil {
		http.Error(w, "User with e-mail already exists", http.StatusBadRequest)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(body.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
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
