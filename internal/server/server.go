package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	_ "github.com/joho/godotenv/autoload"

	"framer/internal/database"
)

var port int

func NewServer(ctx context.Context) *http.Server {
	port, _ := strconv.Atoi(os.Getenv("PORT"))

	databaseUrl := os.Getenv("DATABASE_URL")
	database.NewDb(databaseUrl)
	// Declare Server config
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", port),
		Handler:      RegisterRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return server
}

func Shutdown() {
	database.Shutdown()
}
