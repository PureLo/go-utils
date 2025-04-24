package websocket

import (
	"context"
	"sync"

	"github.com/panjf2000/ants"
)

type Room struct {
	sync.RWMutex

	ID   string
	crew map[string]*Member
	ctx  context.Context
	pool *ants.Pool
}

func NewRoom(ctx context.Context, id string, psize int) *Room {
	p, err := ants.NewPool(psize)
	if err != nil {
		return nil
	}
	room := &Room{
		ID:   id,
		crew: make(map[string]*Member, 0),
		ctx:  ctx,
		pool: p,
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
		r.pool.Submit(func() {
			if err := m.Publish(message); err != nil {
				r.LeaveRoomSafe(m)
			}
		})
	}
	return nil
}

func (r *Room) LeaveRoomSafe(Member *Member) error {
	if Member == nil {
		return ErrorInvalidMember
	}
	r.Lock()
	defer r.Unlock()
	if _, ok := r.crew[Member.uuid]; ok {
		delete(r.crew, Member.uuid)
		Member.wsCoon.Close()
	}
	return nil
}
