package queue

type (
	RoomIoType       string
	FriendIoType     string
	NotificationType string
)

const (
	ROOM_IN  RoomIoType = "room_in"
	ROOM_OUT RoomIoType = "room_out"

	FRIEND_IN  FriendIoType = "friend_in"
	FRIEND_OUT FriendIoType = "friend_out"

	REQ_FRIEND  NotificationType = "friend_request"
	ROOM_INVITE NotificationType = "room_invite"
)

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

type Notification struct {
	NotificationId string           `json:"notification_id"`
	UserId         string           `json:"user_id"`
	Message        string           `json:"message"`
	Type           NotificationType `json:"type"`
	IdRef          string           `json:"id_ref"`
	Seen           bool             `json:"seen"`
	CreatedAt      string           `json:"created_at"`
}

type SaveMessage struct {
	SenderId   string `json:"sender_id"`
	ReceiverId string `json:"receiver_id"`
	Text       string `json:"text"`
	MediaUrl   string `json:"media_url"`
}
