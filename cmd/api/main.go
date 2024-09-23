package main

import (
	"context"
	"fmt"
	"framer/internal/server"
)

func main() {
	ctx := context.Background()

	server := server.NewServer(ctx)

	err := server.ListenAndServe()
	if err != nil {
		panic(fmt.Sprintf("cannot start server: %s", err))
	}

	server.Shutdown(ctx)
}
