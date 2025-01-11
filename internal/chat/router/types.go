package router

import "github.com/gorilla/websocket"

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

type WebsocketActionJoinRoom struct {
	WebsocketMessage
	Payload JoinRoomPayload `json:"payload"`
}

type JoinRoomPayload struct {
	RoomId uint `json:"room_id"`
}
