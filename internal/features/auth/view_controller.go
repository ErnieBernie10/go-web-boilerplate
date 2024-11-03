package auth

import (
	"framer/internal/pkg"
	"framer/internal/pkg/api"
	"framer/internal/pkg/rpc"
	"framer/internal/pkg/view"
	"net/http"
	"time"

	pb "framer/internal/proto"

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

	_, err := rpc.Client.Auth.Register(r.Context(), &pb.RegisterRequest{
		Email:    email,
		Password: pw,
		Name:     "Unimplemented",
	})

	if err != nil {
		view.Message(view.Error, "Something went wrong").Render(r.Context(), w)
		view.GetLogger(r).Error(err.Error())
		return
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

	response, err := rpc.Client.Auth.Login(r.Context(), &pb.LoginRequest{
		Email:    email,
		Password: pw,
	})

	if err != nil {
		view.Message(view.Error, "Invalid credentials").Render(r.Context(), w)
		view.GetLogger(r).Error(err.Error())
		return
	}

	// Step 3: Set the JWT token in an HTTP-only cookie.
	http.SetCookie(w, &http.Cookie{
		Name:     string(pkg.TokenContextKey),   // Cookie name
		Value:    response.AccessToken,          // JWT token value
		Expires:  time.Now().Add(time.Hour * 1), // Cookie expiration time (same as JWT)
		HttpOnly: true,                          // Make the cookie HTTP-only
		Secure:   false,                         // Set to true if using HTTPS
		Path:     "/",
	})

	http.SetCookie(w, &http.Cookie{
		Name:     string(pkg.RefreshContextKey),  // Cookie name
		Value:    response.RefreshToken,          // JWT token value
		Expires:  time.Now().Add(time.Hour * 72), // Cookie expiration time (same as JWT)
		HttpOnly: true,                           // Make the cookie HTTP-only
		Secure:   false,                          // Set to true if using HTTPS
		Path:     "/",
	})

	w.Header().Set("HX-Redirect", view.IndexPath)
	w.WriteHeader(http.StatusSeeOther)
}
