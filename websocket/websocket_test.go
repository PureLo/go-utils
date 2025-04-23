package websocket

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
)

func TestRoom(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		upgrader := websocket.Upgrader{}
		_, _ = upgrader.Upgrade(w, r, nil)
	}))
	defer server.Close()
	wsURL := "ws" + server.URL[4:]

	t.Run("NewRoom", func(t *testing.T) {
		room := NewRoom("test-room", t.Context())
		assert.NotNil(t, room)
		assert.Equal(t, "test-room", room.ID)
		assert.NotNil(t, room.crew)
	})

	t.Run("JoinRoomSafe", func(t *testing.T) {
		room := NewRoom("test-room", t.Context())
		ws, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
		if err != nil {
			t.Fatalf("无法创建WebSocket连接: %v", err)
		}
		defer ws.Close()

		member := NewMember("user1", ws, room)
		err = room.JoinRoomSafe(member)
		assert.Nil(t, err)
	})

	t.Run("Broadcast", func(t *testing.T) {
		// Create a test WebSocket server that will echo messages back
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			upgrader := websocket.Upgrader{}
			conn, err := upgrader.Upgrade(w, r, nil)
			if err != nil {
				t.Fatalf("Failed to upgrade connection: %v", err)
				return
			}

			// Echo any received messages back to the client
			go func() {
				for {
					mt, message, err := conn.ReadMessage()
					if err != nil {
						break
					}
					err = conn.WriteMessage(mt, message)
					if err != nil {
						break
					}
				}
			}()
		}))
		defer server.Close()

		wsURL := "ws" + server.URL[4:]
		room := NewRoom("test-room", t.Context())

		// Create two mock clients
		client1, err := NewMockClient(wsURL)
		if err != nil {
			t.Fatalf("Failed to create mock client 1: %v", err)
		}
		defer client1.Close()

		client2, err := NewMockClient(wsURL)
		if err != nil {
			t.Fatalf("Failed to create mock client 2: %v", err)
		}
		defer client2.Close()

		// Add members to the room using the mock clients
		member1 := NewMember("user1", client1.Conn, room)
		member2 := NewMember("user2", client2.Conn, room)

		_ = room.JoinRoomSafe(member1)
		_ = room.JoinRoomSafe(member2)

		// Broadcast a test message
		testMessage := []byte("test broadcast message")
		err = room.Broadcast(testMessage)
		assert.Nil(t, err)

		// Give some time for the message to be delivered
		time.Sleep(100 * time.Millisecond)

		// Check if both clients received the message
		assert.Equal(t, testMessage, client1.GetLastMessage())
		assert.Equal(t, testMessage, client2.GetLastMessage())
	})
}

func TestMember(t *testing.T) {
	room := NewRoom("test-room", t.Context())
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		upgrader := websocket.Upgrader{}
		_, _ = upgrader.Upgrade(w, r, nil)
	}))
	defer server.Close()

	wsURL := "ws" + server.URL[4:]
	ws, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatalf("无法创建WebSocket连接: %v", err)
	}
	defer ws.Close()

	t.Run("NewMember", func(t *testing.T) {
		member := NewMember("user1", ws, room)
		assert.NotNil(t, member)
		assert.Equal(t, "user1", member.UUID())
	})

	t.Run("Publish", func(t *testing.T) {
		member := NewMember("user1", ws, room)
		err := member.Publish([]byte("test message"))
		assert.Nil(t, err)
	})
}

func TestRoomManager(t *testing.T) {
	t.Run("GetRoomManagerInstance", func(t *testing.T) {
		manager1 := GetRoomManagerInstance()
		manager2 := GetRoomManagerInstance()
		assert.Equal(t, manager1, manager2)
	})

	t.Run("RoomManagement", func(t *testing.T) {
		manager := GetRoomManagerInstance()
		room := NewRoom("test-room", t.Context())
		gotRoom := manager.GetRoom("test-room")
		assert.Equal(t, room, gotRoom)

		newRoom := &Room{ID: "new-room"}
		manager.AddRoom(newRoom)
		gotNewRoom := manager.GetRoom("new-room")
		assert.Equal(t, newRoom, gotNewRoom)
	})
}
