package websocket

import (
	"context"
	"sync"
)

type Room struct {
	sync.RWMutex

	ID   string
	crew map[string]*Member
	ctx  context.Context
}

func NewRoom(id string, ctx context.Context) *Room {
	room := &Room{
		ID:   id,
		crew: make(map[string]*Member, 0),
		ctx:  ctx,
	}

	GetRoomManagerInstance().AddRoom(room)
	return room
}

func (r *Room) JoinRoomSafe(Member *Member) error {
	if Member == nil {
		return ErrorInvalidMember
	}
	r.Lock()
	defer r.Unlock()

	if prev, ok := r.crew[Member.uuid]; ok {
		prev.wsCoon.Close()
	}

	r.crew[Member.uuid] = Member
	return nil
}

func (r *Room) Broadcast(message []byte) error {
	r.RWMutex.RLock()
	defer r.RWMutex.RUnlock()
	for _, m := range r.crew {
		if err := m.Publish(message); err != nil {
			delete(r.crew, m.uuid)
		}
	}
	return nil
}
