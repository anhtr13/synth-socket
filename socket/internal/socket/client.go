package socket

import (
	"context"
	"encoding/json"
	"time"

	"github.com/coder/websocket"
	"github.com/google/uuid"
)

type Client struct {
	ClientId uuid.UUID
	UserId   uuid.UUID
	Conn     *websocket.Conn
}

func NewClient(user_id uuid.UUID, client_conn *websocket.Conn) *Client {
	return &Client{
		ClientId: uuid.New(),
		UserId:   user_id,
		Conn:     client_conn,
	}
}

func (c *Client) WriteMsg(msg BroadcastPayload) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()
	data, _ := json.Marshal(msg)
	err := c.Conn.Write(ctx, websocket.MessageText, []byte(data))
	return err
}
