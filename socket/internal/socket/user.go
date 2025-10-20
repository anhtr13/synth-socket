package socket

import (
	"sync"

	"github.com/google/uuid"
)

type User struct {
	UserId  uuid.UUID
	Clients map[uuid.UUID]*Client
	Mtx     sync.RWMutex
}

func NewUser(user_id uuid.UUID) *User {
	return &User{
		UserId:  user_id,
		Clients: map[uuid.UUID]*Client{},
		Mtx:     sync.RWMutex{},
	}
}

func (u *User) AddClient(clients ...*Client) {
	u.Mtx.Lock()
	defer u.Mtx.Unlock()
	for _, c := range clients {
		u.Clients[c.ClientId] = c
	}
}

func (u *User) RemoveClient(clients ...*Client) {
	u.Mtx.Lock()
	defer u.Mtx.Unlock()
	for _, c := range clients {
		delete(u.Clients, c.ClientId)
	}
}

func (u *User) GetClient(client_id uuid.UUID) *Client {
	u.Mtx.RLock()
	defer u.Mtx.RUnlock()
	client := u.Clients[client_id]
	return client
}

func (u *User) GetAllClient() []*Client {
	res := []*Client{}
	u.Mtx.RLock()
	defer u.Mtx.RUnlock()
	for _, c := range u.Clients {
		res = append(res, c)
	}
	return res
}

func (u *User) CountClients() int {
	u.Mtx.RLock()
	defer u.Mtx.RUnlock()
	res := len(u.Clients)
	return res
}

func (u *User) Broadcast(msg BroadcastPayload) {
	clients := u.GetAllClient()
	for _, c := range clients {
		c.WriteMsg(msg)
	}
}
