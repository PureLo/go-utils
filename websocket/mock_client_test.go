package websocket

import (
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// MockClient represents a test WebSocket client
type MockClient struct {
	Conn        *websocket.Conn
	ReceivedMsg []byte
	mu          sync.Mutex
	done        chan struct{}
}

// NewMockClient creates a new mock client connected to the given URL
func NewMockClient(url string) (*MockClient, error) {
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		return nil, err
	}

	client := &MockClient{
		Conn: conn,
		done: make(chan struct{}),
	}

	// Start listening for messages
	go client.readPump()

	return client, nil
}

// readPump reads messages from the WebSocket connection
func (c *MockClient) readPump() {
	defer close(c.done)
	defer c.Conn.Close()

	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			return
		}

		c.mu.Lock()
		c.ReceivedMsg = message
		c.mu.Unlock()
	}
}

// GetLastMessage returns the last message received by the client
func (c *MockClient) GetLastMessage() []byte {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.ReceivedMsg
}

// Close closes the WebSocket connection
func (c *MockClient) Close() error {
	err := c.Conn.WriteMessage(websocket.CloseMessage,
		websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	if err != nil {
		return err
	}

	select {
	case <-c.done:
	case <-time.After(time.Second):
	}

	return c.Conn.Close()
}
