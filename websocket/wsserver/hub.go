package wsserver

import (
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

type Hub struct {
	connections map[*WsConnection]bool
	register    chan *WsConnection
	unregister  chan *WsConnection
	broadcast   chan []byte
}

type WsConnection struct {
	conn     *websocket.Conn
	sendChan chan []byte // 防止存在conn的写阻塞，导致整个hub广播阻塞
}

func NewHub() *Hub {
	return &Hub{
		connections: make(map[*WsConnection]bool),
		register:    make(chan *WsConnection),
		unregister:  make(chan *WsConnection),
		broadcast:   make(chan []byte),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case conn := <-h.register:
			h.connections[conn] = true
			println("new connection registered")
		case message := <-h.broadcast:
			for conn := range h.connections {
				select {
				case conn.sendChan <- message:
				default:
					// 如果发送通道已满，则关闭连接
					// 并踢掉
					close(conn.sendChan)
					delete(h.connections, conn)
				}
			}
		case conn := <-h.unregister:
			println("connection unregistered")
			if _, ok := h.connections[conn]; ok {
				close(conn.sendChan)
				delete(h.connections, conn)
				conn.conn.Close()
			}
		}
	}
}

const (
	maxMessageSize = 512                 // 最大消息大小
	pongWait       = 60 * time.Second    // 心跳超时时间
	pingPeriod     = (pongWait * 9) / 10 // 心跳间隔时间，略小于超时时间
)

func (c *WsConnection) readPump(hub *Hub) {
	defer func() {
		hub.unregister <- c
	}()

	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPingHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			println("read error:", err)
			break
		}
		hub.broadcast <- message
	}
}

// unimplemented
func (c *WsConnection) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.sendChan:
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			if err := c.conn.WriteMessage(websocket.TextMessage, message); err != nil {
				println("write error:", err)
				return
			}
		case <-ticker.C:
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				println("ping error:", err)
				return
			}
		}
	}
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  maxMessageSize,
	WriteBufferSize: maxMessageSize,
	CheckOrigin: func(r *http.Request) bool {
		return true // 允许跨域连接
	},
}

func ServeWs(hub *Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		println("upgrade error: ", err.Error())
		return
	}
	wsConn := &WsConnection{
		conn:     conn,
		sendChan: make(chan []byte, maxMessageSize), // with buffer, avoid blocking
	}

	hub.register <- wsConn
	go wsConn.readPump(hub)
	go wsConn.writePump()
}
