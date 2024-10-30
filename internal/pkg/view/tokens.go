package view

import (
	"context"
	"framer/internal/pkg"

	"google.golang.org/grpc/metadata"
)

func GetTokens(c context.Context) func() (string, string) {
	return func() (string, string) {
		accessToken, _ := c.Value(pkg.TokenContextKey).(string)
		refreshToken, _ := c.Value(pkg.RefreshContextKey).(string)

		return accessToken, refreshToken
	}
}

func ContextWithToken(ctx context.Context, getTokens func() (string, string)) context.Context {
	accessToken, _ := getTokens()
	ctx = metadata.AppendToOutgoingContext(ctx, "authorization", "Bearer "+accessToken)
	return ctx
}
