package socket

type SocketEvent string

const (
	EVENT_ERROR             SocketEvent = "error"
	EVENT_MESSAGE           SocketEvent = "message"
	EVENT_NOTIFICATION      SocketEvent = "user_notification"
	EVENT_CHAT_NOTIFICATION SocketEvent = "chat_notification"
)

type BroadcastPayload struct {
	Event SocketEvent `json:"event"`
	Data  any         `json:"data"` // Message | ChatNotification | Error | conf.Notification
}

type SocketError struct {
	Message string `json:"message"`
}

type ChatNotification struct {
	ChatId  string `json:"chat_id"`
	Message string `json:"message"`
}

type Message struct {
	SenderId   string `json:"sender_id"`
	ReceiverId string `json:"receiver_id"`
	IsRoomMsg  bool   `json:"is_room_msg"`
	Text       string `json:"text"`
	MediaUrl   string `json:"media_url"`
}
