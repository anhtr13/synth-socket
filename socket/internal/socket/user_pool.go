package socket

import (
	"sync"

	"github.com/google/uuid"
)

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
