package websocket

import "github.com/gorilla/websocket"

type member struct {
	uuid   string
	wsCoon *websocket.Conn
	room   *Room
}

func NewMember(uuid string, wsCoon *websocket.Conn, room *Room) *member {
	if uuid == "" || wsCoon == nil || room == nil {
		return nil
	}

	return &member{
		uuid:   uuid,
		wsCoon: wsCoon,
		room:   room,
	}
}

func (m *member) UUID() string {
	return m.uuid
}

func (m *member) Publish(msg []byte) error {
	return m.wsCoon.WriteMessage(websocket.TextMessage, msg)
}
