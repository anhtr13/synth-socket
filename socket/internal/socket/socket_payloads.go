package socket

type SEvent string

const (
	EVENT_ERROR        SEvent = "error"
	EVENT_MESSAGE      SEvent = "message"
	EVENT_NOTIFICATION SEvent = "notification"
	EVENT_ROOM_IO      SEvent = "room_io"
)

type RoomIoType string

const (
	ROOM_IN  RoomIoType = "room_in"
	ROOM_OUT RoomIoType = "room_out"
)

type NotificationType string

const (
	REQ_FRIEND NotificationType = "friend_request"
	REQ_ROOM   NotificationType = "room_request"
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

type SNotification struct {
	NotificationId string           `json:"notification_id"`
	UserId         string           `json:"user_id"`
	Message        string           `json:"message"`
	Type           NotificationType `json:"type"`
	IdRef          string           `json:"id_ref"`
	Seen           bool             `json:"seen"`
	CreatedAt      string           `json:"created_at"`
}

type SRoomIo struct {
	UserId string     `json:"user_id"`
	RoomId string     `json:"room_id"`
	Type   RoomIoType `json:"type"`
}

type SPayload struct {
	Event SEvent `json:"event"`
	Data  any    `json:"data"` // SMessage | SError | SRoomIo | SNotification
}
