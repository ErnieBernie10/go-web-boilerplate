package auth

import (
	"context"
	"database/sql"
	"errors"
	"framer/internal/core"
	"framer/internal/database"
	pb "framer/internal/proto"
)

type AuthController struct {
	pb.AppUserServiceServer
	Db database.DbService
}

func (c *AuthController) Login(ctx context.Context, req *pb.LoginRequest) (*pb.AppUser, error) {
	user, err := c.Db.Queries.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return nil, errors.Join(core.ErrUnauthorized, err)
	}

	if !user.PasswordHash.Valid {
		return nil, errors.Join(core.ErrUnauthorized, errors.New("invalid username or password"))
	}

	tokenString, refreshTokenString, err := login(
		user.ID,
		req.Email,
		req.Password,
		user.PasswordHash.String,
	)

	if err != nil {
		return nil, errors.Join(core.ErrUnauthorized, err)
	}

	return &pb.AppUser{
		AccessToken:  tokenString,
		RefreshToken: refreshTokenString,
	}, nil
}

func (c *AuthController) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.EmptyResponse, error) {
	_, err := c.Db.Queries.GetUserByEmail(ctx, req.Email)
	if err == nil {
		return nil, errors.Join(core.ErrValidation, errors.New("user with e-mail already exists"))
	}

	hashedPassword, err := hashString(req.Password)
	if err != nil {
		return nil, err
	}

	_, err = c.Db.Queries.Register(ctx, database.RegisterParams{
		Email:        req.Email,
		PasswordHash: sql.NullString{Valid: true, String: string(hashedPassword)},
	})

	return &pb.EmptyResponse{}, err
}
