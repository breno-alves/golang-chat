package router

type ActionType string

const (
	CreateRoom  ActionType = "CREATE_ROOM"
	JoinRoom               = "JOIN_MESSAGE"
	LeaveRoom              = "LEAVE_MESSAGE"
	SendMessage            = "SEND_MESSAGE"
)

type WebsocketMessage struct {
	Action ActionType `json:"action"`
	Token  string     `json:"token"`
	//Payload struct{}   `json:"payload"`
}
