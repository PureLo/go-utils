package echo

import (
	"net/url"
	"testing"
	"time"

	"github.com/gorilla/websocket"
)

func TestWsClient(t *testing.T) {
	uri := "ws://localhost:9001/echo"
	u, err := url.Parse(uri)
	if err != nil {
		t.Fatal("parse:", err)
	}

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		t.Fatal("dial:", err)
	}
	defer c.Close()

	done := make(chan struct{})

	// read message
	go func() {
		defer close(done)
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				println("read err:", err)
				return
			}
			println("recv: %s", message)
		}
	}()

	// ping handler
	c.SetPingHandler(func(appData string) error {
		println("client recv ping:", appData)
		return c.WriteControl(websocket.PongMessage, []byte(appData), time.Now().Add(time.Second))
	})

	// pong handler
	c.SetPongHandler(func(appData string) error {
		println("client recv pong:", appData)
		return nil
	})

	// close handler
	c.SetCloseHandler(func(code int, text string) error {
		println("client recv close:", code, text)
		return nil
	})

	// heartbeat
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-done:
			return
		case t := <-ticker.C:
			c.WriteMessage(websocket.PingMessage, []byte(t.String()))
			println("---------------------------------\n client send ping:", t.String())
		}
	}
}
