package socket

type SEvent string

const (
	EVENT_ERROR        SEvent = "error"
	EVENT_MESSAGE      SEvent = "message"
	EVENT_NOTIFICATION SEvent = "notification"
	EVENT_ROOM_IO      SEvent = "room_io"
)

type SError struct {
	Error string `json:"error"`
}

type SMessage struct {
	SenderId   string `json:"sender_id"`
	ReceiverId string `json:"receiver_id"`
	Text       string `json:"text"`
	MediaUrl   string `json:"media_url"`
}

type SPayload struct {
	Event SEvent `json:"event"`
	Data  any    `json:"data"` // SMessage | SError | conf.QueueMsg_RoomIo | conf.QueueMsg_Notification
}
