package conf

import (
	"context"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	amqp "github.com/rabbitmq/amqp091-go"
	redis "github.com/redis/go-redis/v9"

	"github.com/anhtr13/synth-socket/api-service/pkgs/database"
	"github.com/anhtr13/synth-socket/api-service/pkgs/queue"
)

var (
	DB_Connection *pgxpool.Pool
	DB_Queries    *database.Queries

	RBMQ_Connection *amqp.Connection
	RBMQ_Channel    *amqp.Channel

	RD_Client *redis.Client
)

func InitConnection() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	rd_opt, err := redis.ParseURL(RD_URI)
	if err != nil {
		log.Fatal("Redis failed:", err)
	}
	RD_Client = redis.NewClient(rd_opt)
	_, err = RD_Client.Ping(ctx).Result()
	if err != nil {
		log.Fatal("Could not connect to Redis:", err)
	}
	log.Println("Connected to redis")

	db_conn, err := pgxpool.New(ctx, DB_URI)
	if err != nil {
		log.Fatal("Cannot connect to database", err)
	}
	DB_Connection = db_conn
	DB_Queries = database.New(db_conn)
	log.Println("Connected to database")

	q_conn, err := amqp.Dial(RBMQ_URI)
	if err != nil {
		log.Fatal("Cannot connect to queue:", err.Error())
	}
	RBMQ_Connection = q_conn
	ch, err := q_conn.Channel()
	if err != nil {
		log.Fatal("Error when open amqp channel:", err.Error())
	}
	RBMQ_Channel = ch
	log.Println("Connected to queue")

	err = ch.ExchangeDeclare(
		queue.EXCHANGE_API_SOCKET, // name
		"direct",                  // type
		true,                      // durable
		false,                     // auto-deleted
		false,                     // internal
		false,                     // no-wait
		nil,                       // arguments
	)
	if err != nil {
		log.Fatal("Error when declare exchange:", err.Error())
	}
}

func CloseAllConnections() {
	RD_Client.Close()
	RBMQ_Channel.Close()
	RBMQ_Connection.Close()
	DB_Connection.Close()
}
