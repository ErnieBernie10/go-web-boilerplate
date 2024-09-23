package server

import (
	"encoding/json"
	"framer/internal/database"
	"log"
	"net/http"

	"fmt"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/coder/websocket"
)

func RegisterRoutes() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Get("/", HelloWorldHandler)

	r.Get("/health", healthHandler)

	r.Get("/websocket", websocketHandler)

	r.Route("/frame", frameResourceHandler)

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
