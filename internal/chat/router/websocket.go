package router

import (
	"chatroom/api/chat/model"
	"chatroom/internal/chat/service"
	"context"
	"encoding/json"
	"github.com/gorilla/websocket"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"log/slog"
	"net/http"
	"sync"
)

// key => room:uuid value => { ...user tokens }
// key => user:token => { username, activeRoom, ... }

const (
	JoinRoom    = "JOIN_ROOM"
	LeaveRoom   = "LEAVE_MESSAGE"
	SendMessage = "SEND_MESSAGE"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var connections = make(map[*websocket.Conn]bool)
var clients = make(map[string]*Client)
var mutex = &sync.Mutex{}

var rooms = make(map[uint][]*Client)
var roomBroadcast = make(map[uint]chan []byte)

// Start websocket connection steps
// - Send token in query parameters to authorize
// - if token is valid we upgrade connection, if not refuse
// - Change token property in redis to connected but no room so far

// Connect to room
// At this point we could re-validate token, but I don't think it's necessary
// Update room list (watch out for race conditions in this resource) <-- this could be a Set to prevent duplication
// Update user active room property

// Leave room
// Update room list
// Update user active room property

// Send message
// - Validate token
// - Retrieve all connected users in room and broadcast message
// - update ui

func wsHandler(ctx context.Context, db *gorm.DB, cache *redis.Client, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		slog.Error("error upgrading connection", err.Error())
		return
	}
	go handleConnection(ctx, db, cache, conn)
}

func closeClientConnection(conn *websocket.Conn) {
	mutex.Lock()
	connections[conn] = false
	mutex.Unlock()
	err := conn.Close()
	if err != nil {
		slog.Error("error closing connection", err.Error())
	}
}

func broadcastMessage(message *model.Message) {
	for conn, active := range connections {
		if active {
			err := conn.WriteJSON(message)
			if err != nil {
				slog.Error("error sending message to client", err.Error())
			}
		}
	}
}

func handleConnection(ctx context.Context, db *gorm.DB, cache *redis.Client, conn *websocket.Conn) {
	defer closeClientConnection(conn)

	mutex.Lock()
	connections[conn] = true
	mutex.Unlock()

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			slog.Error("error reading message", err.Error())
			break
		}

		slog.Debug("received message", string(message))

		msg, err := parseAction(message)
		if err != nil {
			slog.Error("error parsing action:", err.Error())
			break
		}

		// if token is invalid/expired we should disconnect user
		username, err := service.CheckToken(ctx, cache, msg.Token)
		if err != nil {
			slog.Error("error checking token:", err.Error())
			break
		}

		switch msg.Action {
		case JoinRoom:
		case LeaveRoom:
		case SendMessage:
			sendMessagePayload := &WebsocketMessageSendMessage{}
			err := json.Unmarshal(message, sendMessagePayload)
			if err != nil {
				slog.Error("error unmarshalling message:", err.Error())
				break
			}
			clientMessage, err := service.CreateMessage(db, sendMessagePayload.Payload.Room, username, sendMessagePayload.Payload.Content)
			if err != nil {
				slog.Error("error creating message:", err.Error())
				break
			}
			broadcastMessage(clientMessage)
			break
		}

	}
}

func parseAction(message []byte) (*WebsocketMessage, error) {
	data := &WebsocketMessage{}
	err := json.Unmarshal(message, data)
	if err != nil {
		return nil, err
	}
	return data, nil
}
