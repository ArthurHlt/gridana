package gridana

import (
	"github.com/ArthurHlt/gridana/emitter"
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
	"time"
)

const (
	// Time allowed to write the file to the client.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the client.
	pongWait = 60 * time.Second

	// Send pings to client with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10
)

func notifyReader(c *websocket.Conn) {
	entry := log.WithField("user_addr", c.RemoteAddr())
	c.SetReadLimit(512)
	c.SetReadDeadline(time.Now().Add(pongWait))
	c.SetPongHandler(
		func(string) error {
			c.SetReadDeadline(time.Now().Add(pongWait))
			entry.Debug("Pong.")
			return nil
		})
	for {
		_, _, err := c.ReadMessage()
		if err != nil {
			break
		}
	}
}
func notifyWriter(c *websocket.Conn, close <-chan bool) {
	entry := log.WithField("user_addr", c.RemoteAddr())
	pingTicker := time.NewTicker(pingPeriod)
	listener := emitter.On()
	for {
		select {
		case event := <-listener:
			alert := emitter.ToAlert(event)
			c.SetWriteDeadline(time.Now().Add(writeWait))
			err := c.WriteJSON(alert)
			if err != nil {
				emitter.Off(listener)
				entry.Debug(err)
				break
			}
			entry.Debug("Event sent to client")
		case <-pingTicker.C:
			c.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				return
			}
			entry.Debug("Ping.")
		case <-close:
			emitter.Off(listener)
			return
		}
	}
}
