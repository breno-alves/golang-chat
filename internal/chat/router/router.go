package router

import (
	"chatroom/api/chat/handler"
	"context"
	"github.com/gorilla/mux"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"net/http"
)

type RequestHandlerFunction func(ctx context.Context, db *gorm.DB, cache *redis.Client, w http.ResponseWriter, r *http.Request)

type Router struct {
	Router *mux.Router
	Db     *gorm.DB
	Cache  *redis.Client
}

func NewRouter(db *gorm.DB, cache *redis.Client) *Router {
	r := &Router{
		Router: mux.NewRouter(),
		Db:     db,
		Cache:  cache,
	}
	r.initialize()
	return r
}

func (router *Router) initialize() {

	// HEALTH ROUTES
	router.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// AUTH ROUTES
	router.Post("/auth/login", router.handleRestRequest(handler.Login))

	// USER ROUTES
	router.Post("/user", router.handleRestRequest(handler.SignUp))

	// ROOM ROUTES
	router.Get("/rooms", router.handleRestRequest(handler.ListRooms))

	// WS ROUTES
	router.Ws("/ws", router.handleWsRequest(wsHandler))
}
