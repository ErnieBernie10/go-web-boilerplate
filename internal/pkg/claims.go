package pkg

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type Claims struct {
	Email string    `json:"username"`
	ID    uuid.UUID `json:"id"`
	jwt.RegisteredClaims
}

func (c *Claims) Valid() error {
	if !time.Now().Before(c.ExpiresAt.Time) {
		return fmt.Errorf("the token has expired")
	}
	return nil
}

var JwtSecret = []byte(os.Getenv("JWT_SECRET"))
