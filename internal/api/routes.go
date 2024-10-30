package api

import (
	"encoding/json"
	"log"
	"net/http"

	"fmt"
	"time"

	"github.com/go-chi/chi/v5"

	_ "framer/docs"
	"framer/internal/pkg/database"
	"framer/internal/pkg/logger"

	"github.com/coder/websocket"
	httpSwagger "github.com/swaggo/http-swagger"
)

func RegisterRoutes() http.Handler {
	r := chi.NewRouter()
	r.Use(logger.HttpLoggingMiddleware(logger.NewLogger()))

	r.Group(viewRouteHandler)
	r.Group(apiRouteHandler)

	r.Get("/health", healthHandler)

	r.Get("/websocket", websocketHandler)

	r.Get("/swagger/*", httpSwagger.WrapHandler)

	return r
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
