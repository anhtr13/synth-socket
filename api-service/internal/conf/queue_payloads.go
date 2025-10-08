package conf

type RoomIoType string

const (
	ROOM_IN  RoomIoType = "room_in"
	ROOM_OUT RoomIoType = "room_out"
)

type QueueMsg_RoomIo struct {
	UserId string     `json:"user_id"`
	RoomId string     `json:"room_id"`
	Type   RoomIoType `json:"type"`
}

type NotificationType string

const (
	REQ_FRIEND NotificationType = "friend_request"
	REQ_ROOM   NotificationType = "room_request"
)

type QueueMsg_Notification struct {
	NotificationId string           `json:"notification_id"`
	UserId         string           `json:"user_id"`
	Message        string           `json:"message"`
	Type           NotificationType `json:"type"`
	IdRef          string           `json:"id_ref"`
	Seen           bool             `json:"seen"`
	CreatedAt      string           `json:"created_at"`
}
