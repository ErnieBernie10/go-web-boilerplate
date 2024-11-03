package rpc

import (
	"context"
	"framer/internal/pkg"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func IsAuthenticated(ctx context.Context) (*pkg.Claims, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "no metadata found")
	}

	// Extract the "authorization" token
	token := md["authorization"]
	if len(token) == 0 {
		return nil, status.Errorf(codes.Unauthenticated, "no auth token provided")
	}

	// Validate the token (this is just an example, you would call a real JWT validator here)
	claims := &pkg.Claims{}
	_, err := jwt.ParseWithClaims(strings.Split(token[0], " ")[1], claims, func(token *jwt.Token) (interface{}, error) {
		return pkg.JwtSecret, nil
	})
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "invalid token")
	}

	// Continue with the actual request handling
	return claims, nil
}
