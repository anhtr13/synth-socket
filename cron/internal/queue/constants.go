package queue

const (
	EXCHANGE_API_TO_SOCKET string = "exchange_api_to_socket" // exchange to push notifications from api service to socket service
	ROUTE_ROOM_IO          string = "route_room_io"
	ROUTE_FRIEND_IO        string = "route_friend_io"
	ROUTE_NOTIFICATION     string = "route_notification"

	EXCHANGE_SOCKET_TO_CRON string = "exchange_socket_to_cron" // exchange to push message from socket service to cron
	ROUTE_SAVE_MESSAGE      string = "route_save_message"
)
