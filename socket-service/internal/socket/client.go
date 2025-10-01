package socket

import (
	"context"
	"encoding/json"

	"github.com/coder/websocket"
	"github.com/google/uuid"
)

// A user can connect on multiple devices, each connection represented by a "Client"
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
	data, _ := json.Marshal(msg)
	err := c.Conn.Write(context.Background(), websocket.MessageText, []byte(data))
	return err
}
