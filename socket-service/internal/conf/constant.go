package conf

import (
	"log"
	"os"
	"strconv"
)

const (
	RBMQ_EXCHANGE_API_SOCKET string = "exchange_api_socket"
	RBMQ_KEY_ROOM_IO         string = "route_room_io"
	RBMQ_KEY_NOTIFICATION    string = "route_notification"

	// Redis keys
	KEY_USER           string = "user"
	USER_FRIENDS       string = "friends" // user:<user_id>:friends = set[friend_ids]
	USER_ROOMS         string = "rooms"   // user:<user_id>:rooms 	= set[room_ids]
	USER_STATUS        string = "status"  // user:<user_id>:status  = "online" | <last_active>
	USER_STATUS_ONLINE string = "online"
)

var (
	DB_URI   string
	RD_URI   string
	RBMQ_URI string
	JWT_SEC  string
	HOST     string

	PORT int = 8080
)

func LoadEnvs() {
	dbUri := os.Getenv("DB_URI")
	if dbUri == "" {
		log.Fatal("Cannot find env: DB_URI")
	}
	DB_URI = dbUri

	qUri := os.Getenv("RBMQ_URI")
	if qUri == "" {
		log.Fatal("Cannot find env: RBMQ_URI")
	}
	RBMQ_URI = qUri

	rdUri := os.Getenv("RD_URI")
	if rdUri == "" {
		log.Fatal("Cannot find env: RD_URI")
	}
	RD_URI = rdUri

	port := os.Getenv("PORT")
	p, err := strconv.Atoi(port)
	if err == nil {
		PORT = p
	}

	jWTSec := os.Getenv("JWT_SECRET")
	if jWTSec == "" {
		log.Fatal("Cannot find env: JWT_SECRET")
	}
	JWT_SEC = jWTSec

	host := os.Getenv("HOST")
	if host == "" {
		log.Fatal("Cannot find env: HOST")
	}
	HOST = host
}

func init() {
	LoadEnvs()
}
