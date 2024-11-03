package auth

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"framer/internal/pkg"
)

var (
	ErrInvalidLogin = errors.New("invalid username or password")
)

func refreshTokens(accessToken, refreshToken string) (string, string, error) {
	expirationTime := time.Now().Add(5 * time.Minute) // Set token expiration time.
	if os.Getenv("APP_ENV") == string(pkg.Development) {
		expirationTime = time.Now().Add(1 * time.Hour)
	}
	claims := &pkg.Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	refreshClaim := &pkg.Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	// Create the JWT using the claims and the secret key.
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	refresh := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaim)

	tokenString, err := token.SignedString(pkg.JwtSecret)
	if err != nil {
		return "", "", ErrInvalidLogin
	}

	refreshTokenString, err := refresh.SignedString(pkg.JwtSecret)
	if err != nil {
		return "", "", ErrInvalidLogin
	}

	return tokenString, refreshTokenString, nil
}

func login(id uuid.UUID, email, password, hash string) (string, string, error) {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		return "", "", ErrInvalidLogin
	}

	expirationTime := time.Now().Add(5 * time.Minute) // Set token expiration time.
	if os.Getenv("APP_ENV") == string(pkg.Development) {
		expirationTime = time.Now().Add(1 * time.Hour)
	}
	claims := &pkg.Claims{
		Email: email,
		ID:    id,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	refreshClaim := &pkg.Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	// Create the JWT using the claims and the secret key.
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	refresh := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaim)

	tokenString, err := token.SignedString(pkg.JwtSecret)
	if err != nil {
		return "", "", ErrInvalidLogin
	}

	refreshTokenString, err := refresh.SignedString(pkg.JwtSecret)
	if err != nil {
		return "", "", ErrInvalidLogin
	}

	return tokenString, refreshTokenString, nil
}

func hashString(password string) ([]byte, error) {

	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}
