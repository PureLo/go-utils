package websocket

import "github.com/gorilla/websocket"

type Member struct {
	uuid   string
	wsCoon *websocket.Conn
	room   *Room
}

func NewMember(uuid string, wsCoon *websocket.Conn, room *Room) *Member {
	if uuid == "" || wsCoon == nil || room == nil {
		return nil
	}

	return &Member{
		uuid:   uuid,
		wsCoon: wsCoon,
		room:   room,
	}
}

func (m *Member) UUID() string {
	return m.uuid
}

func (m *Member) Publish(msg []byte) error {
	if m.wsCoon == nil {
		return nil
	}

	return m.wsCoon.WriteMessage(websocket.TextMessage, msg)
}

func (m *Member) Close() error {
	if m.wsCoon != nil {
		return m.wsCoon.Close()
	}

	return nil
}

func (m *Member) SetCloseHandler(h func(code int, text string) error) {
	m.wsCoon.SetCloseHandler(h)
}
