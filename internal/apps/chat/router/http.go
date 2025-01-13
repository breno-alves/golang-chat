package router

import (
	"chatroom/internal/pkg/http/middlewares"
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"log/slog"
	"net/http"
)

type RequestHandlerFunction func(ctx context.Context, db *gorm.DB, cache *redis.Client, w http.ResponseWriter, r *http.Request)

func (r *Router) InitHttpRouter() {
	r.HttpRouter.Use(middlewares.SetCORS)

	// HEALTH ROUTES
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// AUTH ROUTES
	r.Post("/auth/login", r.handler.Login)

	// USER ROUTES
	r.Post("/user", r.handler.SignUp)

	// ROOM ROUTES
	r.Post("/rooms/leave", r.handler.LeaveRoom)
	r.Get("/rooms", r.handler.ListRooms)
	r.Post("/rooms", r.handler.CreateRoom)

	// MESSAGES ROUTES
	r.Get("/messages", r.handler.ListMessages)

	// WS ROUTES
	r.Ws("/ws", r.handler.WebsocketConnect)
}

// Get wraps the router for GET method
func (r *Router) Get(path string, f func(w http.ResponseWriter, r *http.Request)) {
	slog.Info(fmt.Sprintf("GET %s", path))
	r.HttpRouter.HandleFunc(path, f).Methods("GET")
}

// Post wraps the router for POST method
func (r *Router) Post(path string, f func(w http.ResponseWriter, r *http.Request)) {
	slog.Info(fmt.Sprintf("POST %s", path))
	r.HttpRouter.HandleFunc(path, f).Methods("POST")
}

// Put wraps the router for PUT method
func (r *Router) Put(path string, f func(w http.ResponseWriter, r *http.Request)) {
	slog.Info(fmt.Sprintf("PUT %s", path))
	r.HttpRouter.HandleFunc(path, f).Methods("PUT")
}

// Delete wraps the router for DELETE method
func (r *Router) Delete(path string, f func(w http.ResponseWriter, r *http.Request)) {
	slog.Info(fmt.Sprintf("DELETE %s", path))
	r.HttpRouter.HandleFunc(path, f).Methods("DELETE")
}

func (r *Router) Ws(path string, f func(w http.ResponseWriter, r *http.Request)) {
	r.HttpRouter.HandleFunc(path, f)
}
