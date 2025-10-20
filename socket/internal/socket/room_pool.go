package socket

import (
	"sync"

	"github.com/google/uuid"
)

type RoomPool struct {
	Pool map[uuid.UUID]*Room
	Mtx  sync.RWMutex
}

func NewRoomPool() *RoomPool {
	return &RoomPool{
		Pool: map[uuid.UUID]*Room{},
		Mtx:  sync.RWMutex{},
	}
}

func (pool *RoomPool) AddRoom(rooms ...*Room) {
	pool.Mtx.Lock()
	defer pool.Mtx.Unlock()
	for _, r := range rooms {
		pool.Pool[r.RoomId] = r
	}
}

func (pool *RoomPool) RemoveRoom(rooms ...*Room) {
	pool.Mtx.Lock()
	defer pool.Mtx.Unlock()
	for _, r := range rooms {
		delete(pool.Pool, r.RoomId)
	}
}

func (pool *RoomPool) GetRoom(room_id uuid.UUID) *Room {
	pool.Mtx.RLock()
	defer pool.Mtx.RUnlock()
	room := pool.Pool[room_id]
	return room
}

func (pool *RoomPool) GetAllRooms() []*Room {
	res := []*Room{}
	pool.Mtx.RLock()
	defer pool.Mtx.RUnlock()
	for _, r := range pool.Pool {
		res = append(res, r)
	}
	return res
}
