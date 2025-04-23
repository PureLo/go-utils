package websocket

import "sync"

var (
	once sync.Once
	ins  *RoomManager
)

type RoomManager struct {
	rooms map[string]*Room
	sync.RWMutex
}

func GetRoomManagerInstance() *RoomManager {
	once.Do(func() {
		ins = &RoomManager{
			rooms: make(map[string]*Room),
		}
	})
	return ins
}

func (rm *RoomManager) GetRoom(id string) *Room {
	rm.RLock()
	defer rm.RUnlock()

	return rm.rooms[id]
}

func (rm *RoomManager) AddRoom(room *Room) {
	rm.Lock()
	defer rm.Unlock()

	rm.rooms[room.ID] = room
}
