package echo

import (
	"flag"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

func WsServerRun() {
	go func() {
		http.HandleFunc("/echo", echo)
		log.Fatal(http.ListenAndServe(*addr, nil))
	}()
}

var addr = flag.String("addr", "localhost:9001", "http service address")

var upgrader = websocket.Upgrader{} // use default options

func echo(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()

	c.SetPingHandler(func(appData string) error {
		println("server recv ping:", appData)
		return c.WriteControl(websocket.PongMessage, []byte(appData), time.Now().Add(time.Second))
	})

	// pong handler
	c.SetPongHandler(func(appData string) error {
		println("server recv pong:", appData)
		return nil
	})

	// close handler
	c.SetCloseHandler(func(code int, text string) error {
		println("server recv close:", code, text)
		return c.WriteControl(websocket.CloseMessage, []byte(text), time.Now().Add(time.Second))
	})

	for {
		mt, message, err := c.ReadMessage()
		if err != nil {
			println("read:", err)
			break
		}

		// echo
		println("server read:", string(message))
		err = c.WriteMessage(mt, message)
		if err != nil {
			println("write:", err)
			break
		}
	}
}
