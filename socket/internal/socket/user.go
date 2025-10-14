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

func (u *User) Broadcast(msg SPayload) {
	u.Mtx.RLock()
	defer u.Mtx.RUnlock()
	for _, c := range u.Clients {
		c.WriteMsg(msg)
	}
}

// A collection of active users
type UserPool struct {
	Pool map[uuid.UUID]*User
	Mtx  sync.RWMutex
}

func NewUserPool() *UserPool {
	return &UserPool{
		Pool: map[uuid.UUID]*User{},
		Mtx:  sync.RWMutex{},
	}
}

func (pool *UserPool) AddUser(users ...*User) {
	pool.Mtx.Lock()
	defer pool.Mtx.Unlock()
	for _, u := range users {
		pool.Pool[u.UserId] = u
	}
}

func (pool *UserPool) RemoveUser(users ...*User) {
	pool.Mtx.Lock()
	defer pool.Mtx.Unlock()
	for _, u := range users {
		delete(pool.Pool, u.UserId)
	}
}

func (pool *UserPool) GetUser(user_id uuid.UUID) *User {
	pool.Mtx.RLock()
	defer pool.Mtx.RUnlock()
	user := pool.Pool[user_id]
	return user
}

func (pool *UserPool) GetAllUsers() []*User {
	res := []*User{}
	pool.Mtx.RLock()
	defer pool.Mtx.RUnlock()
	for _, u := range pool.Pool {
		res = append(res, u)
	}
	return res
}
