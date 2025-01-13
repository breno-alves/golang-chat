package router

import (
	"chatroom/internal/handlers/chat"
	"chatroom/internal/pkg/broker"
	"github.com/gorilla/mux"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type Router struct {
	HttpRouter *mux.Router
	handler    *chat.Handler
	broker     *broker.Broker
}

func NewRouter(db *gorm.DB, cache *redis.Client, broker *broker.Broker) *Router {
	r := &Router{
		HttpRouter: mux.NewRouter(),
		handler:    chat.NewHandler(db, cache),
	}
	r.InitHttpRouter()
	r.InitRmqRouter(broker)
	return r
}
