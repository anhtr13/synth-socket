package socket

type (
	Event            string
	NotificationType string
	RoomIoType       string
	FriendIoType     string
)

const (
	EVENT_MESSAGE      Event = "message"
	EVENT_NOTIFICATION Event = "notification"
	EVENT_ROOM_IO      Event = "room_io"
	EVENT_FIREND_IO    Event = "friend_io"

	REQ_FRIEND  NotificationType = "friend_request"
	ROOM_INVITE NotificationType = "room_invite"

	ROOM_IN  RoomIoType = "room_in"
	ROOM_OUT RoomIoType = "room_out"

	FRIEND_IN  FriendIoType = "friend_in"
	FRIEND_OUT FriendIoType = "friend_out"
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

type RoomIO struct {
	UserId string     `json:"user_id"`
	RoomId string     `json:"room_id"`
	Type   RoomIoType `json:"type"`
}

type FriendIO struct {
	User1Id string       `json:"user1_id"`
	User2Id string       `json:"user2_id"`
	Type    FriendIoType `json:"type"`
}

type BroadcastPayload struct {
	Event Event `json:"event"`
	Data  any   `json:"data"` // Message | Notification | RoomIO
}
