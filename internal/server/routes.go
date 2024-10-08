package server

import (
	"encoding/json"
	"framer/internal/database"
	"log"
	"log/slog"
	"net/http"
	"os"

	"fmt"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/httplog/v2"

	_ "framer/docs"

	"github.com/coder/websocket"
	httpSwagger "github.com/swaggo/http-swagger"
)

func RegisterRoutes() http.Handler {
	r := chi.NewRouter()
	logger := NewLogger()
	r.Use(httplog.RequestLogger(logger))

	r.Group(viewRouteHandler)
	r.Group(apiRouteHandler)

	r.Get("/health", healthHandler)

	r.Get("/websocket", websocketHandler)

	r.Get("/swagger/*", httpSwagger.WrapHandler)

	return r
}

func NewLogger() *httplog.Logger {
	return httplog.NewLogger("framer-server-logger", httplog.Options{
		// JSON:             true,
		LogLevel:         slog.LevelDebug,
		Concise:          true,
		RequestHeaders:   true,
		MessageFieldName: "message",
		// TimeFieldFormat: time.RFC850,
		Tags: map[string]string{
			"version": "1.0.0",
			"env":     os.Getenv("APP_ENV"),
		},
		QuietDownRoutes: []string{
			"/",
		},
		QuietDownPeriod: 10 * time.Second,
		// SourceFieldName: "source",
	})
}

func HelloWorldHandler(w http.ResponseWriter, r *http.Request) {
	resp := make(map[string]string)
	resp["message"] = "Hello World"

	jsonResp, err := json.Marshal(resp)
	if err != nil {
		log.Fatalf("error handling JSON marshal. Err: %v", err)
	}

	_, _ = w.Write(jsonResp)
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	jsonResp, _ := json.Marshal(database.Health())
	_, _ = w.Write(jsonResp)
}

func websocketHandler(w http.ResponseWriter, r *http.Request) {
	socket, err := websocket.Accept(w, r, nil)

	if err != nil {
		log.Printf("could not open websocket: %v", err)
		_, _ = w.Write([]byte("could not open websocket"))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	defer socket.Close(websocket.StatusGoingAway, "server closing websocket")

	ctx := r.Context()
	socketCtx := socket.CloseRead(ctx)

	for {
		payload := fmt.Sprintf("server timestamp: %d", time.Now().UnixNano())
		err := socket.Write(socketCtx, websocket.MessageText, []byte(payload))
		if err != nil {
			break
		}
		time.Sleep(time.Second * 2)
	}
}
