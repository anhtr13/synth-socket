package cache

const (
	KEY_USER           string = "user"
	USER_FRIENDS       string = "friends" // user:<user_id>:friends = set[friend_ids]
	USER_ROOMS         string = "rooms"   // user:<user_id>:rooms 	= set[room_ids]
	USER_STATUS        string = "status"  // user:<user_id>:status  = "online" | <last_active>
	USER_STATUS_ONLINE string = "online"
)
