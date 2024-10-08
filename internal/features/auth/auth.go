package auth

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type Token string

const TokenContextKey Token = "Token"
const RefreshContextKey Token = "Refresh"

var JwtSecret = []byte(os.Getenv("JWT_SECRET"))

var (
	ErrInvalidLogin = errors.New("invalid username or password")
)

type Claims struct {
	Email string `json:"username"`
	jwt.RegisteredClaims
}

func (c *Claims) Valid() error {
	return nil
}

func login(email, password, hash string) (string, string, error) {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		return "", "", ErrInvalidLogin
	}

	expirationTime := time.Now().Add(5 * time.Minute) // Set token expiration time.
	claims := &Claims{
		Email: email,
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
		return "", "", ErrInvalidLogin
	}

	refreshTokenString, err := refresh.SignedString(JwtSecret)
	if err != nil {
		return "", "", ErrInvalidLogin
	}

	return tokenString, refreshTokenString, nil
}

func hashString(password string) ([]byte, error) {

	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}
