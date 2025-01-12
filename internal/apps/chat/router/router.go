package router

import (
	"chatroom/internal/handlers"
	"chatroom/internal/pkg/http/middlewares"
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"log/slog"
	"net/http"
)

type RequestHandlerFunction func(ctx context.Context, db *gorm.DB, cache *redis.Client, w http.ResponseWriter, r *http.Request)

type Router struct {
	Router  *mux.Router
	handler *handlers.Handler
}

func NewRouter(db *gorm.DB, cache *redis.Client) *Router {
	r := &Router{
		Router:  mux.NewRouter(),
		handler: handlers.NewHandler(db, cache),
	}
	r.initialize()
	return r
}

func (r *Router) initialize() {
	r.Router.Use(middlewares.SetCORS)

	// HEALTH ROUTES
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// AUTH ROUTES
	r.Post("/auth/login", r.handler.Login)

	// USER ROUTES
	r.Post("/user", r.handler.SignUp)

	// ROOM ROUTES
	r.Get("/rooms", r.handler.ListRooms)
	r.Post("/rooms", r.handler.CreateRoom)

	// WS ROUTES
	r.Ws("/ws", r.handler.WebsocketConnect)
}

// Get wraps the router for GET method
func (r *Router) Get(path string, f func(w http.ResponseWriter, r *http.Request)) {
	slog.Info(fmt.Sprintf("GET %s", path))
	r.Router.HandleFunc(path, f).Methods("GET")
}

// Post wraps the router for POST method
func (r *Router) Post(path string, f func(w http.ResponseWriter, r *http.Request)) {
	slog.Info(fmt.Sprintf("POST %s", path))
	r.Router.HandleFunc(path, f).Methods("POST")
}

// Put wraps the router for PUT method
func (r *Router) Put(path string, f func(w http.ResponseWriter, r *http.Request)) {
	slog.Info(fmt.Sprintf("PUT %s", path))
	r.Router.HandleFunc(path, f).Methods("PUT")
}

// Delete wraps the router for DELETE method
func (r *Router) Delete(path string, f func(w http.ResponseWriter, r *http.Request)) {
	slog.Info(fmt.Sprintf("DELETE %s", path))
	r.Router.HandleFunc(path, f).Methods("DELETE")
}

func (r *Router) Ws(path string, f func(w http.ResponseWriter, r *http.Request)) {
	r.Router.HandleFunc(path, f)
}
