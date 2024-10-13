package database

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func Seed() (uuid.UUID, error) {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte("test"), bcrypt.DefaultCost)
	if err != nil {
		return uuid.UUID{}, err
	}

	Service.Register(context.Background(), RegisterParams{
		Email:        "test@test.com",
		PasswordHash: sql.NullString{String: string(passwordHash), Valid: true},
	})

	user, err := Service.GetUserByEmail(context.Background(), "test@test.com")
	if err != nil {
		return uuid.UUID{}, err
	}

	return user.ID, nil
}
