package socket

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/coder/websocket"
	"github.com/google/uuid"
	amqp "github.com/rabbitmq/amqp091-go"

	"github.com/anhtr13/synth-socket/socket/internal/cache"
	"github.com/anhtr13/synth-socket/socket/internal/conf"
	"github.com/anhtr13/synth-socket/socket/internal/queue"
	"github.com/anhtr13/synth-socket/socket/internal/util"
)

type SocketServer struct {
	UserPool *UserPool
	RoomPool *RoomPool
}

func NewSocketServer() *SocketServer {
	return &SocketServer{
		UserPool: NewUserPool(),
		RoomPool: NewRoomPool(),
	}
}

func (s *SocketServer) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
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
		s.updateUserActiveStatus(user, cache.USER_STATUS_ONLINE)
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
			continue
		}

		room_uuid, err := uuid.Parse(msg.ReceiverId)
		if err != nil {
			continue
		}

		room := s.RoomPool.GetRoom(room_uuid)
		if room == nil {
			continue
		}
		if room.GetMember(user_id) == nil {
			continue
		}

		msg.SenderId = user.UserId.String()
		room.Broadcast(BroadcastPayload{
			Event: EVENT_MESSAGE,
			Data:  msg,
		})

		err = conf.RBMQ_Channel.Publish(
			queue.EXCHANGE_SOCKET_TO_CRON, // exchange
			queue.ROUTE_SAVE_MESSAGE,      // routing key
			false,                         // mandatory
			false,                         // immediate
			amqp.Publishing{
				ContentType: "text/plain",
				Body:        data,
			},
		)
	}
}

func (s *SocketServer) HandleQueue_RoomIO(done chan bool) {
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
			q_msg := queue.RoomIO{}
			err = json.Unmarshal(d.Body, &q_msg)
			if err != nil {
				continue
			}

			user_uuid, err := uuid.Parse(q_msg.UserId)
			if err != nil {
				continue
			}
			user := s.UserPool.GetUser(user_uuid)
			if user == nil {
				continue
			}

			room_uuid, err := uuid.Parse(q_msg.RoomId)
			if err != nil {
				continue
			}
			room := s.RoomPool.GetRoom(room_uuid)
			if room == nil {
				switch q_msg.Type {
				case queue.ROOM_IN:
					room = NewRoom(room_uuid)
					s.RoomPool.AddRoom(room)
				case queue.ROOM_OUT:
					continue
				}
			}

			switch q_msg.Type {
			case queue.ROOM_IN:
				room.AddMember(user)
				room_in := BroadcastPayload{
					Event: EVENT_ROOM_IO,
					Data: RoomIO{
						UserId: q_msg.UserId,
						RoomId: q_msg.RoomId,
						Type:   ROOM_IN,
					},
				}
				room.Broadcast(room_in)
			case queue.ROOM_OUT:
				room.RemoveMember(user)
				room_out := BroadcastPayload{
					Event: EVENT_ROOM_IO,
					Data: RoomIO{
						UserId: q_msg.UserId,
						RoomId: q_msg.RoomId,
						Type:   ROOM_OUT,
					},
				}
				room.Broadcast(room_out)
			}
		}
	}
}

func (s *SocketServer) HandleQueue_FriendIO(done chan bool) {
	delivery, err := conf.RBMQ_Channel.Consume(
		conf.Queue_FriendIO.Name, // queue
		"",                       // consumer
		true,                     // auto-ack
		false,                    // exclusive
		false,                    // no-local
		false,                    // no-wait
		nil,                      // args
	)
	if err != nil {
		panic(err)
	}
	for {
		select {
		case <-done:
			return
		case d := <-delivery:
			q_msg := queue.FriendIO{}
			err = json.Unmarshal(d.Body, &q_msg)
			if err != nil {
				continue
			}
			user1_uuid, err := uuid.Parse(q_msg.User1Id)
			user2_uuid, err := uuid.Parse(q_msg.User2Id)
			if err != nil {
				continue
			}
			msg := BroadcastPayload{
				Event: EVENT_FIREND_IO,
				Data: FriendIO{
					User1Id: q_msg.User1Id,
					User2Id: q_msg.User2Id,
					Type:    FriendIoType(q_msg.Type),
				},
			}
			user1 := s.UserPool.GetUser(user1_uuid)
			user2 := s.UserPool.GetUser(user2_uuid)
			if user1 != nil {
				user1.Broadcast(msg)
			}
			if user2 != nil {
				user2.Broadcast(msg)
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
			q_msg := queue.Notification{}
			err = json.Unmarshal(d.Body, &q_msg)
			if err != nil {
				continue
			}
			user_uuid, err := uuid.Parse(q_msg.UserId)
			if err != nil {
				continue
			}

			user := s.UserPool.GetUser(user_uuid)
			if user == nil {
				continue
			}

			noti_msg := BroadcastPayload{
				Event: EVENT_NOTIFICATION,
				Data:  q_msg,
			}

			user.Broadcast(noti_msg)
		}
	}
}

func (s *SocketServer) getUserId(r *http.Request) (uuid.UUID, error) {
	access_token := r.URL.Query().Get("access_token")
	if access_token == "" {
		cookie, err := r.Cookie("access_token")
		if err == nil {
			access_token = cookie.Value
		}
	}
	if access_token == "" {
		return [16]byte{}, fmt.Errorf("access_token not found")
	}
	user_claim, err := util.VerifyJWT(access_token)
	if err != nil {
		return [16]byte{}, err
	}
	user_id, ok := user_claim["id"].(string)
	if !ok {
		return [16]byte{}, fmt.Errorf("client_id not found")
	}
	user_uuid, err := uuid.Parse(user_id)
	if err != nil {
		return [16]byte{}, err
	}
	return user_uuid, nil
}

func (s *SocketServer) initClient(client_conn *websocket.Conn, user_id uuid.UUID) (*Client, *User, error) {
	client := NewClient(user_id, client_conn)
	user := s.UserPool.GetUser(user_id)
	if user == nil {
		user = NewUser(user_id)
		s.UserPool.AddUser(user)

		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		user_rooms, err := conf.RD_Client.SMembers(
			ctx,
			fmt.Sprintf("%s:%s:%s", cache.KEY_USER, user_id.String(), cache.USER_ROOMS),
		).Result()
		if err != nil {
			return nil, nil, err
		}

		for _, id := range user_rooms {
			room_id, err := uuid.Parse(id)
			if err != nil {
				continue
			}
			room := s.RoomPool.GetRoom(room_id)
			if room == nil {
				room = NewRoom(room_id)
				s.RoomPool.AddRoom(room)
			}
			room.AddMember(user)
		}
	}

	user.AddClient(client)

	return client, user, nil
}

func (s *SocketServer) cleanUpClient(client *Client, user *User) {
	user.RemoveClient(client)

	if user.CountClients() == 0 {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		room_ids, err := conf.RD_Client.SMembers(
			ctx,
			fmt.Sprintf("%s:%s:%s", cache.KEY_USER, user.UserId.String(), cache.USER_ROOMS),
		).Result()
		if err != nil {
			return
		}

		for _, r_id := range room_ids {
			r_uuid, err := uuid.Parse(r_id)
			if err != nil {
				continue
			}
			room := s.RoomPool.GetRoom(r_uuid)
			if room == nil {
				continue
			}
			room.RemoveMember(user)
			if room.CountMembers() == 0 {
				s.RoomPool.RemoveRoom(room)
			}
		}

		s.UserPool.RemoveUser(user)

		conf.RD_Client.Set(
			ctx,
			fmt.Sprintf("%s:%s:%s", cache.KEY_USER, user.UserId.String(), cache.USER_STATUS),
			time.Now().String(),
			0,
		).Err()
	}
}

// Update user active status
func (s *SocketServer) updateUserActiveStatus(user *User, active_status string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	err := conf.RD_Client.Set(
		ctx,
		fmt.Sprintf("%s:%s:%s", cache.KEY_USER, user.UserId.String(), cache.USER_STATUS),
		active_status,
		0,
	).Err()
	if err != nil {
		return err
	}
	return nil
}
