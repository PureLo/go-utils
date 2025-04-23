How to use 
```go
package main

func main() {
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
		client1, _ := NewMockClient(wsURL)
		defer client1.Close()

		// Add members to the room using the mock clients
		member1 := NewMember("user1", client1.Conn, room)
		_ = room.JoinRoomSafe(member1)

		// Broadcast a test message
		testMessage := []byte("test broadcast message")
		_ = room.Broadcast(testMessage)
}
```