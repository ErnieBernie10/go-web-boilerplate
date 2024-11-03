package main

import (
	"context"
	"fmt"
	"framer/internal/api"
	"framer/internal/pkg/database"
	"framer/internal/pkg/rpc"
	"net/http"
	"os"
	"strconv"
	"time"
)

func main() {
	ctx := context.Background()

	server := NewServer(ctx)
	defer Shutdown()

	err := server.ListenAndServe()
	if err != nil {
		panic(fmt.Sprintf("cannot start server: %s", err))
	}
}

func NewServer(ctx context.Context) *http.Server {
	port, _ := strconv.Atoi(os.Getenv("API_PORT"))
	rpcPort := os.Getenv("GRPC_PORT")

	databaseUrl := os.Getenv("DATABASE_URL")
	database.NewDb(databaseUrl)

	err := rpc.NewClient("localhost:" + rpcPort)
	if err != nil {
		panic(fmt.Sprintf("cannot connect to rpc server: %s", err))
	}

	// Declare Server config
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", port),
		Handler:      api.RegisterRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return server
}

func Shutdown() {
	database.Shutdown()
}
