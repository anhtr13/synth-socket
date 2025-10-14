package conf

import (
	"context"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	redis "github.com/redis/go-redis/v9"

	"github.com/anhtr13/synth-socket/api-service/pkgs/queue"
)

var (
	RBMQ_Connection *amqp.Connection
	RBMQ_Channel    *amqp.Channel

	Queue_RoomIO       amqp.Queue
	Queue_Notification amqp.Queue

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

	q_conn, err := amqp.Dial(RBMQ_URI)
	if err != nil {
		log.Fatal("Cannot connect to queue:", err)
	}
	RBMQ_Connection = q_conn
	ch, err := q_conn.Channel()
	if err != nil {
		log.Fatal("Error when open amqp channel:", err)
	}
	RBMQ_Channel = ch
	log.Println("Connected to queue")

	// ========== ExchangeDeclare ==========
	err = ch.ExchangeDeclare(
		queue.EXCHANGE_API_SOCKET,
		"direct",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatal("Error when declare exchange:", err)
	}

	// ========== QueueDeclare ==========
	q_room_io, err := ch.QueueDeclare(
		"",
		false,
		false,
		true,
		false,
		nil,
	)
	if err != nil {
		log.Fatal("Error when declare queue: ", err)
	}
	Queue_RoomIO = q_room_io

	q_notification, err := ch.QueueDeclare(
		"",
		false,
		false,
		true,
		false,
		nil,
	)
	if err != nil {
		log.Fatal("Error when declare queue: ", err)
	}
	Queue_Notification = q_notification

	// ========== QueueBind ==========
	err = ch.QueueBind(
		Queue_RoomIO.Name,
		queue.ROUTE_ROOM_IO,
		queue.EXCHANGE_API_SOCKET,
		false,
		nil,
	)
	if err != nil {
		log.Fatal("Error when bind queue: ", err)
	}
	err = ch.QueueBind(
		Queue_Notification.Name,
		queue.ROUTE_NOTIFICATION,
		queue.EXCHANGE_API_SOCKET,
		false,
		nil,
	)
	if err != nil {
		log.Fatal("Error when bind queue: ", err)
	}
}

func CloseAllConnections() {
	RD_Client.Close()
	RBMQ_Channel.Close()
	RBMQ_Connection.Close()
}
