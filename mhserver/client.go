package main

import (
	"encoding/json"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 2 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 10000
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  4096,
	WriteBufferSize: 4096,
}

// readPump pumps messages from the websocket connection to the hub.
//
// The application runs readPump in a per-connection goroutine. The application
// ensures that there is at most one reader on a connection by executing all
// reads from this goroutine.
func (c *Client) readPump() {
	defer func() {
		log.Println("closing websocket read")
		ti := time.Now()
		note := rawmsg{Type: "notification", Value: ti.Format(time.Stamp) + ": " + c.device.Desc + " has disconnected"}

		jnote, err := json.Marshal(note)
		if err != nil {
			log.Printf("Marshal error: %s", err.Error())
		}
		mesg := message{Type: "control", Targets: []string{"text"}, Msg: jnote}
		c.hub.broadcast <- mesg

		c.hub.unregister <- c
		c.conn.Close()
	}()
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })

	for {
		mtype, msg, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}

		if mtype == websocket.TextMessage {
			var messg message
			err = json.Unmarshal(msg, &messg)
			if err != nil {
				log.Printf("unmarshal error: %s", err.Error())
			}
			c.msg = &messg
			processchan <- c

		}
	}
}

// writePump pumps messages from the hub to the websocket connection.
//
// A goroutine running writePump is started for each connection. The
// application ensures that there is at most one writer to a connection by
// executing all writes from this goroutine.
func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		log.Println("closing websocket write")
		ticker.Stop()
		c.conn.Close()
	}()
	for {
		select {
		case msg, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				log.Println("hub closed channel") // The hub closed the channel.
				ti := time.Now()
				note := rawmsg{Type: "notification", Value: ti.Format(time.Stamp) + ": " + c.device.Desc + " was disconnected by the Hub "}

				jnote, err := json.Marshal(note)
				if err != nil {
					log.Printf("Marshal error: %s", err.Error())
				}
				mesg := message{Type: "control", Targets: []string{"text"}, Msg: jnote}
				c.hub.broadcast <- mesg
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				log.Printf("couldn't get next writer: %s \n", err.Error())
				return
			}

			msgstr, _ := json.Marshal(msg)

			//	websocket.WriteJSON(c.conn, msg)
			w.Write(msgstr)

			// Add queued chat messages to the current websocket message.
			n := len(c.send)
			for i := 0; i < n; i++ {

				websocket.WriteJSON(c.conn, <-c.send)
			}

			if err := w.Close(); err != nil {
				return
			}

		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))

			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				log.Println("couldn't send ping")
				return
			}
			//	default:
			//		fmt.Println("dropped message")
		}
	}
}

// serveWs handles websocket requests from the peer.
func serveWs(hub *Hub, w http.ResponseWriter, r *http.Request) {
	log.Println(r)
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {

		log.Println(err)
		return
	}
	host, _, _ := net.SplitHostPort(r.RemoteAddr)

	device := &device{Type: "", Name: "", Desc: "", State: "", Address: host}
	client := &Client{hub: hub, conn: conn, send: make(chan message), device: device}
	client.hub.register <- client

	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go client.writePump()
	go client.readPump()

}
