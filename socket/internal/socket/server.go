package socket

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/coder/websocket"
	"github.com/google/uuid"

	"github.com/anhtr13/synth-socket/socket-service/internal/cache"
	"github.com/anhtr13/synth-socket/socket-service/internal/conf"
	"github.com/anhtr13/synth-socket/socket-service/internal/util"
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

		user_rooms, err := conf.RD_Client.SMembers(
			ctx,
			fmt.Sprintf("%s:%s:%s", cache.KEY_USER, user.UserId.String(), cache.USER_ROOMS),
		).Result()
		if err != nil {
			return
		}

		for _, id := range user_rooms {
			room_id, err := uuid.Parse(id)
			if err != nil {
				continue
			}
			room := s.RoomPool.GetRoom(room_id)
			if room == nil {
				continue
			}
			room.RemoveMember(user)
			if room.CountMembers() == 0 {
				s.RoomPool.RemoveRoom(room)
			}
		}
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
