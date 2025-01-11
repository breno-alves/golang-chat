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

type Client struct {
	conn *websocket.Conn
	user *model.User
}

type ActionType string

type WebsocketMessage struct {
	Action ActionType `json:"action"`
	Token  string     `json:"token"`
}

type WebsocketMessageSendMessage struct {
	WebsocketMessage
	Payload SendMessagePayload `json:"payload"`
}

type SendMessagePayload struct {
	Content string `json:"content"`
	Room    uint   `json:"room"`
}

const (
	CreateRoom  ActionType = "CREATE_ROOM"
	JoinRoom               = "JOIN_ROOM"
	LeaveRoom              = "LEAVE_MESSAGE"
	SendMessage            = "SEND_MESSAGE"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var clients = make(map[*websocket.Conn]bool)
var mutex = &sync.Mutex{}

var rooms = make(map[uint][]*Client)
var roomBroadcast = make(map[uint]chan []byte)

func (router *Router) Ws(path string, f func(w http.ResponseWriter, r *http.Request)) {
	router.Router.HandleFunc(path, f)
}

func wsHandler(ctx context.Context, db *gorm.DB, cache *redis.Client, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		slog.Error("error upgrading connection", err.Error())
		return
	}

	go handleConnection(ctx, db, cache, conn)
}

func handleConnection(ctx context.Context, db *gorm.DB, cache *redis.Client, conn *websocket.Conn) {
	defer func() {
		err := conn.Close()
		if err != nil {
			slog.Error("error closing connection", err.Error())
		}
		delete(clients, conn)
	}()

	mutex.Lock()
	clients[conn] = true
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
		case CreateRoom:
			_, err = service.CreateRoom(db)
			if err != nil {
				slog.Error("error creating room:", err.Error())
				break
			}
			break
		case JoinRoom:
		case LeaveRoom:
		case SendMessage:
			sendMessagePayload := &WebsocketMessageSendMessage{}
			err := json.Unmarshal(message, sendMessagePayload)
			if err != nil {
				slog.Error("error unmarshalling message:", err.Error())
				break
			}
			_, err = service.CreateMessage(db, sendMessagePayload.Payload.Room, username, sendMessagePayload.Payload.Content)
			if err != nil {
				slog.Error("error creating message:", err.Error())
				break
			}
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
