package main

import (
	"context"
	"fmt"
	srv "framer/internal/server"
)

func main() {
	ctx := context.Background()

	server := srv.NewServer(ctx)

	err := server.ListenAndServe()
	if err != nil {
		panic(fmt.Sprintf("cannot start server: %s", err))
	}
	srv.Shutdown()
}
