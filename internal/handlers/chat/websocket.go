package chat

import (
	"chatroom/internal/models"
	"context"
	"encoding/json"
	"github.com/gorilla/websocket"
	"log/slog"
	"net/http"
	"strconv"
	"sync"
)

const (
	LeaveRoom   = "LEAVE_MESSAGE"
	SendMessage = "SEND_MESSAGE"
)

type Client struct {
	conn     *websocket.Conn
	username string
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
type JoinRoomPayload struct {
	RoomId uint `json:"room_id"`
}

type WebsocketActionJoinRoom struct {
	WebsocketMessage
	Payload JoinRoomPayload `json:"payload"`
}

var connections = make(map[string]*websocket.Conn)
var mutex = &sync.Mutex{}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (h *Handler) WebsocketConnect(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	token := r.URL.Query().Get("token")
	if token == "" {
		slog.Error("token required")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	ctx = context.WithValue(ctx, "token", token)

	user, err := h.userService.ValidateUserToken(ctx, token)
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
	room, err := h.roomService.FindByID(ctx, uint(roomId))
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
	go h.handleConnection(ctx, conn)
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

func (h *Handler) broadcastMessage(ctx context.Context, message *models.Message) error {
	slog.Debug("broadcasting message")

	tokens, err := h.roomService.GetCurrentUserTokensInRoom(ctx)
	if err != nil {
		return err
	}
	for _, token := range tokens {
		conn, ok := connections[token]
		if !ok {
			slog.Error("could not find connection for token", token)
			continue
		}
		err := conn.WriteJSON(message)
		if err != nil {
			slog.Error("error sending message", err.Error())
		}
	}
	return nil
}

func (h *Handler) handleConnection(ctx context.Context, conn *websocket.Conn) {
	defer closeClientConnection(ctx, conn)
	token := ctx.Value("token").(string)

	mutex.Lock()
	connections[token] = conn
	mutex.Unlock()

	err := h.roomService.UserJoinRoom(ctx)
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

		user, err := h.userService.ValidateUserToken(ctx, msg.Token)
		if err != nil {
			slog.Error("error checking token:", err.Error())
			break
		}

		switch msg.Action {
		case SendMessage:
			sendMessagePayload := &WebsocketMessageSendMessage{}
			err := json.Unmarshal(message, sendMessagePayload)
			if err != nil {
				slog.Error("error unmarshalling message:", err.Error())
				break
			}
			clientMessage, err := h.messageService.CreateMessage(ctx, sendMessagePayload.Payload.Room, user.Username, sendMessagePayload.Payload.Content)
			if err != nil {
				slog.Error("error creating message:", err.Error())
				break
			}
			_ = h.broadcastMessage(ctx, clientMessage)
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
