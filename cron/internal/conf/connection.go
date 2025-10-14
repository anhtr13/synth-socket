package conf

import (
	"context"
	"log"
	"time"

	"github.com/anhtr13/synth-socket/api/pkgs/database"
	"github.com/anhtr13/synth-socket/api/pkgs/queue"
	"github.com/jackc/pgx/v5/pgxpool"
	amqp "github.com/rabbitmq/amqp091-go"
)

var (
	DB_Connection *pgxpool.Pool
	DB_Queries    *database.Queries

	RBMQ_Connection *amqp.Connection
	RBMQ_Channel    *amqp.Channel

	Queue_NewMessages amqp.Queue
)

func InitConnection() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

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
	q_new_messages, err := ch.QueueDeclare(
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
	Queue_NewMessages = q_new_messages
	err = ch.QueueBind(
		Queue_NewMessages.Name,
		queue.ROUTE_ROOM_IO,
		queue.EXCHANGE_API_SOCKET,
		false,
		nil,
	)
	if err != nil {
		log.Fatal("Error when bind queue: ", err)
	}
}

func CloseAllConnections() {
	RBMQ_Channel.Close()
	RBMQ_Connection.Close()
	DB_Connection.Close()
}
