package auth

import (
	"framer/internal/api"
	"framer/internal/core"
	"framer/internal/view"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
)

var email = "email"
var password = "password"
var confirm = "confirm"

func AuthResourceHandler(r chi.Router) {
	r.Get(view.LoginPath, handleGetLogin)
	r.Post(view.LoginPath, handlePostLogin)
	r.Get(view.RegisterPath, handleGetRegister)
	r.Post(view.RegisterPath, handlePostRegister)
}

func handleGetLogin(w http.ResponseWriter, r *http.Request) {
	user := api.GetUser(r)

	if user != nil {
		http.Redirect(w, r, view.IndexPath, http.StatusSeeOther)
		return
	}

	view.Render(w, r, loginPage(), nil)
}

func handleGetRegister(w http.ResponseWriter, r *http.Request) {
	user := api.GetUser(r)

	if user != nil {
		http.Redirect(w, r, view.IndexPath, http.StatusSeeOther)
	}

	view.Render(w, r, registerPage(), nil)
}

func handlePostRegister(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
		return
	}

	email := r.FormValue(email)
	pw := r.FormValue(password)
	confirm := r.FormValue(confirm)

	if pw != confirm {
		view.Message(view.Error, "Password and confirmation do not match").Render(r.Context(), w)
		return
	}

	if status, err := api.ApiClient.Request("POST",
		api.RegisterApiPath,
		registerCommandDto{
			Email:    email,
			Password: pw,
		}, nil); err != nil {
		switch status {
		case http.StatusBadRequest:
			view.Message(view.Error, "User with E-mail already exists").Render(r.Context(), w)
			return
		default:
			break
		}
	}

	w.Header().Set("HX-Redirect", view.IndexPath)
	w.WriteHeader(http.StatusSeeOther)
}

func handlePostLogin(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		view.Message(view.Error, "Unable to parse form").Render(r.Context(), w)
		view.GetLogger(r).Error(err.Error())
		return
	}

	email := r.FormValue(email)
	pw := r.FormValue(password)

	response := api.LoginResponseDto{}

	status, err := api.ApiClient.Request("POST", api.LoginApiPath, loginCommandDto{
		Email:    email,
		Password: pw,
	}, &response)

	if err != nil {
		switch status {
		case http.StatusUnauthorized:
			view.Message(view.Error, "E-mail or password do not match").Render(r.Context(), w)
			return
		}
		view.Message(view.Error, "Something went wrong").Render(r.Context(), w)
		view.GetLogger(r).Error(err.Error())
	}

	// Step 3: Set the JWT token in an HTTP-only cookie.
	http.SetCookie(w, &http.Cookie{
		Name:     string(core.TokenContextKey),  // Cookie name
		Value:    response.AccessToken,          // JWT token value
		Expires:  time.Now().Add(time.Hour * 1), // Cookie expiration time (same as JWT)
		HttpOnly: true,                          // Make the cookie HTTP-only
		Secure:   false,                         // Set to true if using HTTPS
		Path:     "/",
	})

	http.SetCookie(w, &http.Cookie{
		Name:     string(core.RefreshContextKey), // Cookie name
		Value:    response.RefreshToken,          // JWT token value
		Expires:  time.Now().Add(time.Hour * 72), // Cookie expiration time (same as JWT)
		HttpOnly: true,                           // Make the cookie HTTP-only
		Secure:   false,                          // Set to true if using HTTPS
		Path:     "/",
	})

	w.Header().Set("HX-Redirect", view.IndexPath)
	w.WriteHeader(http.StatusSeeOther)
}
