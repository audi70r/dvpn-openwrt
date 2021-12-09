package node

import (
	"github.com/gorilla/websocket"
	"net/http"
	"sync"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type Connection struct {
	Socket *websocket.Conn
	mu     sync.Mutex
}

func (c *Connection) Send(msg []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.Socket.WriteMessage(1, msg)
}

func (n *Node) Handle(w http.ResponseWriter, r *http.Request) {
	n.SocketConn.Socket, _ = upgrader.Upgrade(w, r, nil) // error ignored for sake of simplicity

	for {
		// Read message from browser
		msgType, msg, err := n.SocketConn.Socket.ReadMessage()
		if err != nil {
			return
		}

		// Write message back to browser
		if err = n.SocketConn.Socket.WriteMessage(msgType, msg); err != nil {
			return
		}
	}
}
