package socket

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/coder/websocket"
	"github.com/google/uuid"

	"github.com/anhtr13/synth-socket/socket-service/internal/conf"
	"github.com/anhtr13/synth-socket/socket-service/internal/util"
)

func (s *SocketServer) HandleSocketConnection(w http.ResponseWriter, r *http.Request) {
	user_id, err := s.getUserId(r)
	if err != nil {
		util.WriteError(w, 400, err.Error())
		return
	}
	client_conn, err := websocket.Accept(w, r, &websocket.AcceptOptions{
		InsecureSkipVerify: true,
	})
	if err != nil {
		util.WriteError(w, 400, fmt.Sprintf("Failed to accept connection: %s", err.Error()))
		return
	}
	defer client_conn.Close(websocket.StatusNormalClosure, "")

	// Create client
	client, user, err := s.initClient(client_conn, user_id)
	if err != nil {
		util.WriteError(w, 500, fmt.Sprintf("Failed to accept connection: %s", err.Error()))
		return
	}
	if user.CountClients() == 1 {
		s.updateUserActiveStatus(user, conf.USER_STATUS_ONLINE)
	}

	ctx := context.Background()

	// Handle websocket messages
	for {
		msgType, data, err := client_conn.Read(ctx)

		if err != nil {
			// Read failed => client disconnected
			s.cleanUpClient(client, user)
			return
		}
		if msgType != websocket.MessageText {
			continue
		}

		msg := Message{}
		err = json.Unmarshal(data, &msg)
		if err != nil {
			user.Broadcast(BroadcastPayload{
				Event: EVENT_ERROR,
				Data: SocketError{
					Message: err.Error(),
				},
			})
			continue
		}

		room_uuid, err := uuid.Parse(msg.ReceiverId)
		if err != nil {
			user.Broadcast(BroadcastPayload{
				Event: EVENT_ERROR,
				Data: SocketError{
					Message: err.Error(),
				},
			})
			continue
		}

		room := s.RoomPool.GetRoom(room_uuid)
		if room == nil || room.GetMember(user_id) == nil {
			user.Broadcast(BroadcastPayload{
				Event: EVENT_ERROR,
				Data: SocketError{
					Message: "Room not found",
				},
			})
			continue
		}

		msg.SenderId = user.UserId.String()
		room.Broadcast(BroadcastPayload{
			Event: EVENT_MESSAGE,
			Data:  msg,
		})
	}
}

func (s *SocketServer) HandleQueue_RoomIo(done chan bool) {
	delivery, err := conf.RBMQ_Channel.Consume(
		conf.Queue_RoomIO.Name, // queue
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
	for {
		select {
		case <-done:
			return
		case d := <-delivery:
			room_io_msg := conf.QueueMsg_RoomIo{}
			err = json.Unmarshal(d.Body, &room_io_msg)
			if err != nil {
				continue
			}

			user_uuid, err := uuid.Parse(room_io_msg.UserId)
			if err != nil {
				continue
			}
			user := s.UserPool.GetUser(user_uuid)
			if user == nil {
				continue
			}

			room_uuid, err := uuid.Parse(room_io_msg.RoomId)
			if err != nil {
				continue
			}
			room := s.RoomPool.GetRoom(room_uuid)
			if room == nil {
				switch room_io_msg.Type {
				case conf.ROOM_IN:
					room = NewRoom(room_uuid)
					s.RoomPool.AddRoom(room)
				case conf.ROOM_OUT:
					continue
				}
			}

			switch room_io_msg.Type {
			case conf.ROOM_IN:
				room.AddMember(user)
				room_in_msg := BroadcastPayload{
					Event: EVENT_CHAT_NOTIFICATION,
					Data: ChatNotification{
						ChatId:  room.RoomId.String(),
						Message: fmt.Sprintf("%s has joined the room.", user.UserId.String()),
					},
				}
				room.Broadcast(room_in_msg)
			case conf.ROOM_OUT:
				room.RemoveMember(user)
				room_out_msg := BroadcastPayload{
					Event: EVENT_CHAT_NOTIFICATION,
					Data: ChatNotification{
						ChatId:  room.RoomId.String(),
						Message: fmt.Sprintf("%s has left the room.", user.UserId.String()),
					},
				}
				room.Broadcast(room_out_msg)
			}
		}
	}
}

func (s *SocketServer) HandleQueue_Notification(done chan bool) {
	delivery, err := conf.RBMQ_Channel.Consume(
		conf.Queue_Notification.Name, // queue
		"",                           // consumer
		true,                         // auto-ack
		false,                        // exclusive
		false,                        // no-local
		false,                        // no-wait
		nil,                          // args
	)
	if err != nil {
		panic(err)
	}
	for {
		select {
		case <-done:
			return
		case d := <-delivery:
			notification_msg := conf.QueueMsg_Notification{}
			err = json.Unmarshal(d.Body, &notification_msg)
			if err != nil {
				continue
			}
			user_uuid, err := uuid.Parse(notification_msg.UserId)
			if err != nil {
				continue
			}

			user := s.UserPool.GetUser(user_uuid)
			if user == nil {
				continue
			}

			noti_msg := BroadcastPayload{
				Event: EVENT_NOTIFICATION,
				Data:  notification_msg,
			}

			user.Broadcast(noti_msg)
		}
	}
}
