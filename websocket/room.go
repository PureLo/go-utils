package websocket

import (
	"context"
	"sync"
)

type Room struct {
	sync.RWMutex

	ID   string
	crew map[string]*member
	ctx  context.Context
}

func NewRoom(id string, ctx context.Context) *Room {
	room := &Room{
		ID:   id,
		crew: make(map[string]*member, 0),
		ctx:  ctx,
	}

	GetRoomManagerInstance().AddRoom(room)
	return room
}

func (r *Room) JoinRoomSafe(member *member) error {
	if member == nil {
		return ErrorInvalidMember
	}
	r.Lock()
	defer r.Unlock()

	if prev, ok := r.crew[member.uuid]; ok {
		prev.wsCoon.Close()
	}

	r.crew[member.uuid] = member
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
