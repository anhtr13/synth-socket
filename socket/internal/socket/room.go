package socket

import (
	"sync"

	"github.com/google/uuid"
)

type Room struct {
	RoomId  uuid.UUID
	Members map[uuid.UUID]*User
	Mtx     sync.RWMutex
}

func NewRoom(room_id uuid.UUID) *Room {
	return &Room{
		RoomId:  room_id,
		Members: map[uuid.UUID]*User{},
		Mtx:     sync.RWMutex{},
	}
}

func (r *Room) AddMember(members ...*User) {
	r.Mtx.Lock()
	defer r.Mtx.Unlock()
	for _, u := range members {
		r.Members[u.UserId] = u
	}
}

func (r *Room) RemoveMember(members ...*User) {
	r.Mtx.Lock()
	defer r.Mtx.Unlock()
	for _, c := range members {
		delete(r.Members, c.UserId)
	}
}

func (r *Room) GetMember(user_id uuid.UUID) *User {
	r.Mtx.RLock()
	defer r.Mtx.RUnlock()
	client := r.Members[user_id]
	return client
}

func (r *Room) GetAllMembers() []*User {
	res := []*User{}
	r.Mtx.RLock()
	defer r.Mtx.RUnlock()
	for _, u := range r.Members {
		res = append(res, u)
	}
	return res
}

func (r *Room) CountMembers() int {
	r.Mtx.RLock()
	defer r.Mtx.RUnlock()
	res := len(r.Members)
	return res
}

func (r *Room) Broadcast(msg BroadcastPayload) {
	r.Mtx.RLock()
	defer r.Mtx.RUnlock()
	for _, u := range r.Members {
		u.Broadcast(msg)
	}
}

// A collection of active Rooms
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
