package conf

import (
	"log"
	"os"
	"strconv"
)

const (
	RBMQ_EXCHANGE_API_SOCKET string = "exchange_api_socket"
	RBMQ_KEY_ROOM_IO         string = "room_io_route"
	RBMQ_KEY_NOTIFICATION    string = "notification_route"

	// Redis key constraint: PREFIX:[user_id | room_id]:SUFFIX
	// ex: user:6ae32cad-b464-4298-b8da-1f08b0e0d6c9:status
	RDKEY_PREFIX_USER        string = "user"
	RDKEY_SUFFIX_STATUS      string = "status"
	RDKEY_SUFFIX_FRIENDSHIPS string = "friends"
	RDKEY_SUFFIX_GROUPS      string = "groups"

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
