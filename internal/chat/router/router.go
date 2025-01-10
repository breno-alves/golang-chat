package router

import (
	"chatroom/api/chat/handler"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"log/slog"
	"net/http"
)

type RequestHandlerFunction func(db *gorm.DB, w http.ResponseWriter, r *http.Request)

type Router struct {
	Router *mux.Router
	Db     *gorm.DB
}

func NewRouter(db *gorm.DB, cache *redis.Client) *Router {
	r := &Router{
		Router: mux.NewRouter(),
		Db:     db,
	}
	r.initialize()
	return r
}

func (router *Router) initialize() {
	router.Ws("/ws", wsHandler)

	// AUTH ROUTES
	router.Post("/auth/login", router.handleRequest(handler.Login))

	// USER ROUTES
	router.Post("/user", router.handleRequest(handler.SignUp))

	// ROOM ROUTES
	router.Get("/room", router.handleRequest(handler.ListRooms))
	router.Post("/room", router.handleRequest(handler.CreateRoom))
	router.Post("/room/join", router.handleRequest(handler.JoinRoom))

	// MESSAGES ROUTES
	router.Get("/room/{room_id}", router.handleRequest(handler.ListMessage))
	router.Post("/room/message", router.handleRequest(handler.CreateMessage))
}

// Get wraps the router for GET method
func (router *Router) Get(path string, f func(w http.ResponseWriter, r *http.Request)) {
	slog.Info(fmt.Sprintf("GET %s", path))
	router.Router.HandleFunc(path, f).Methods("GET")
}

// Post wraps the router for POST method
func (router *Router) Post(path string, f func(w http.ResponseWriter, r *http.Request)) {
	slog.Info(fmt.Sprintf("POST %s", path))
	router.Router.HandleFunc(path, f).Methods("POST")
}

// Put wraps the router for PUT method
func (router *Router) Put(path string, f func(w http.ResponseWriter, r *http.Request)) {
	slog.Info(fmt.Sprintf("PUT %s", path))
	router.Router.HandleFunc(path, f).Methods("PUT")
}

// Delete wraps the router for DELETE method
func (router *Router) Delete(path string, f func(w http.ResponseWriter, r *http.Request)) {
	slog.Info(fmt.Sprintf("DELETE %s", path))
	router.Router.HandleFunc(path, f).Methods("DELETE")
}

func (router *Router) Ws(path string, f func(w http.ResponseWriter, r *http.Request)) {
	router.Router.HandleFunc(path, f)
}

func (router *Router) handleRequest(handler RequestHandlerFunction) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		handler(router.Db, w, r)
	}
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		slog.Error(err.Error())
		return
	}
	defer func() {
		err := conn.Close()
		if err != nil {
			fmt.Println("Error closing connection:", err)
		}
	}()

	// Listen for incoming messages
	for {
		// Read message from the client
		_, message, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("Error reading message:", err)
			break
		}
		fmt.Printf("Received: %s\\n", message)
		// Echo the message back to the client
		if err := conn.WriteMessage(websocket.TextMessage, message); err != nil {
			fmt.Println("Error writing message:", err)
			break
		}
	}
}
