package socket

type Event string

const (
	EVENT_MESSAGE      Event = "message"
	EVENT_NOTIFICATION Event = "notification"
)

type NotificationType string

const (
	REQ_FRIEND  NotificationType = "friend_request"
	ROOM_INVITE NotificationType = "room_invite"
)

type Message struct {
	SenderId   string `json:"sender_id"`
	ReceiverId string `json:"receiver_id"`
	Text       string `json:"text"`
	MediaUrl   string `json:"media_url"`
}

type Notification struct {
	NotificationId string           `json:"notification_id"`
	UserId         string           `json:"user_id"`
	Message        string           `json:"message"`
	Type           NotificationType `json:"type"`
	IdRef          string           `json:"id_ref"`
	Seen           bool             `json:"seen"`
	CreatedAt      string           `json:"created_at"`
}

type BroadcastPayload struct {
	Event Event `json:"event"`
	Data  any   `json:"data"` // Message | Notification
}
