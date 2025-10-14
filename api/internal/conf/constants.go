package conf

import (
	"log"
	"os"
	"strconv"
)

const (
	AGE_AC_TOKEN = 30 * 60          // 30 minutes
	AGE_RF_TOKEN = 7 * 24 * 60 * 60 // 7 days
	AGE_COOKIE   = 7 * 24 * 60 * 60 // 7 days

	USER_CTX_KEY string = "user_session"
)

var (
	HOST     string
	DB_URI   string
	RD_URI   string
	RBMQ_URI string
	JWT_SEC  string

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
