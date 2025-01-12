package router

import (
	"chatroom/internal/models"
	"chatroom/internal/services"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"log/slog"
	"net/http"
	"strconv"
	"sync"
)

// key => room:uuid value => { ...user tokens }
// key => user:token => { username, activeRoom, ... }

const (
	LeaveRoom   = "LEAVE_MESSAGE"
	SendMessage = "SEND_MESSAGE"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var connections = make(map[string]*websocket.Conn)
var mutex = &sync.Mutex{}

// Start websocket connection steps
// - Send token in query parameters to authorize
// - if token is valid we upgrade connection, if not refuse
// - Change token property in redis to connected but no room so far

// Connect to room
// At this point we could re-validate token, but I don't think it's necessary
// room:id:user:username
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
	token := r.URL.Query().Get("token")
	if token == "" {
		slog.Error("token required")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	ctx = context.WithValue(ctx, "token", token)

	user, err := services.ValidateUserToken(ctx, cache, token)
	if err != nil {
		slog.Error("could not validate token")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	ctx = context.WithValue(ctx, "user", user)

	a := r.URL.Query().Get("room_id")
	roomId, err := strconv.Atoi(a)
	if err != nil {
		slog.Error("could not parse room_id")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	room, err := services.FindRoomById(db, uint(roomId))
	if err != nil {
		slog.Error("could not find room")
		w.WriteHeader(http.StatusNotFound)
		return
	}
	ctx = context.WithValue(ctx, "room", room)

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		slog.Error("error upgrading connection", err.Error())
		return
	}
	go handleConnection(ctx, db, cache, conn)
}

func closeClientConnection(ctx context.Context, conn *websocket.Conn) {
	token := ctx.Value("token").(string)
	mutex.Lock()
	connections[token] = conn
	mutex.Unlock()
	err := conn.Close()
	if err != nil {
		slog.Error("error closing connection", err.Error())
	}
}

func broadcastMessage(ctx context.Context, cache *redis.Client, message *models.Message) error {
	slog.Debug("broadcasting message")
	iter := cache.Scan(ctx, 0, fmt.Sprintf("room:%d:user:*", message.RoomId), 0).Iterator()
	for iter.Next(ctx) {
		k := iter.Val()
		token, _ := cache.Get(ctx, k).Result()
		conn := connections[token]
		err := conn.WriteJSON(message)
		if err != nil {
			slog.Error("error sending message", err.Error())
		}
	}
	if err := iter.Err(); err != nil {
		return err
	}
	return nil
}

func handleConnection(ctx context.Context, db *gorm.DB, cache *redis.Client, conn *websocket.Conn) {
	defer closeClientConnection(ctx, conn)
	token := ctx.Value("token").(string)

	mutex.Lock()
	connections[token] = conn
	mutex.Unlock()

	err := services.JoinRoom(ctx, db, cache)
	if err != nil {
		slog.Error("error joining room", err.Error())
		return
	}

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

		user, err := services.ValidateUserToken(ctx, cache, msg.Token)
		if err != nil {
			slog.Error("error checking token:", err.Error())
			break
		}

		switch msg.Action {
		case LeaveRoom:
		case SendMessage:
			sendMessagePayload := &WebsocketMessageSendMessage{}
			err := json.Unmarshal(message, sendMessagePayload)
			if err != nil {
				slog.Error("error unmarshalling message:", err.Error())
				break
			}
			clientMessage, err := services.CreateMessage(db, sendMessagePayload.Payload.Room, user.Username, sendMessagePayload.Payload.Content)
			if err != nil {
				slog.Error("error creating message:", err.Error())
				break
			}
			_ = broadcastMessage(ctx, cache, clientMessage)
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
