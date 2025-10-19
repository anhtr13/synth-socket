package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os/signal"
	"syscall"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"

	"github.com/anhtr13/synth-socket/cron/internal/database"
	"github.com/anhtr13/synth-socket/cron/internal/queue"
)

func init() {
	LoadEnvs()
	InitConnection()
}

func main() {
	delivery, err := RBMQ_Channel.Consume(
		Queue_NewMessages.Name, // queue
		"",                     // consumer
		true,                   // auto-ack
		false,                  // exclusive
		false,                  // no-local
		false,                  // no-wait
		nil,                    // args
	)
	if err != nil {
		panic(err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	log.Println("Cron job is running...")

	for {
		select {
		case <-ctx.Done():
			log.Println("Cron job is shutting down...")
			CloseAllConnections()
			log.Println("Cron job shutdown complete.")
			return
		case d := <-delivery:
			qmsg := queue.SaveMessage{}
			err := json.Unmarshal(d.Body, &qmsg)
			if err != nil {
				log.Println("Cannot unmarshal queue message:", err.Error())
				continue
			}
			go insertDb(qmsg)
		}
	}
}

func insertDb(qmsg queue.SaveMessage) error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	sender_uuid, err := uuid.Parse(qmsg.SenderId)
	room_uuid, err := uuid.Parse(qmsg.ReceiverId)
	if err != nil {
		return err
	}
	msg, err := DB_Queries.CreateMessage(
		ctx,
		database.CreateMessageParams{
			RoomID:   pgtype.UUID{Bytes: room_uuid, Valid: true},
			SenderID: pgtype.UUID{Bytes: sender_uuid, Valid: true},
			Text:     qmsg.Text,
			MediaUrl: pgtype.Text{String: qmsg.MediaUrl},
		},
	)
	if err != nil {
		log.Println("Cannot save message:", err.Error())
	}
	err = DB_Queries.UpdateRoomLastMessage(
		ctx,
		database.UpdateRoomLastMessageParams{
			RoomID:      pgtype.UUID{Bytes: room_uuid, Valid: true},
			LastMessage: msg.MessageID,
			UpdatedAt:   pgtype.Timestamp{Time: time.Now(), InfinityModifier: pgtype.Finite, Valid: true},
		},
	)
	if err != nil {
		return err
	}
	return nil
}
