package internal

import (
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

const (
	redisTTL = 60 * time.Second
	readWait = 60 * time.Second

	pongWait = 40 * time.Second

	pingPeriod = (pongWait * 7) / 10

	maxMessageSize = 512
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type Client struct {
	room *Room

	conn *websocket.Conn

	send chan []byte
}

func (c *Client) readPump() {

	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	var msg []byte
	defer func() {
		c.room.unregister <- &ClientMessage{client: c, msg: msg}
		c.conn.Close()
	}()

	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {

			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				return
			}
			msg = []byte("opponent left the game")
			c.room.broadcast <- &ClientMessage{client: c, msg: msg}
			return

		}
		// fmt.Println("readpump default:", string(message))
		c.room.broadcast <- &ClientMessage{client: c, msg: message}
	}
}

func (c *Client) writePump() {
	pingTicker := time.NewTicker(pingPeriod)
	defer func() {
		pingTicker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			if strings.Compare(string(message), "ROOM_UNAVAILABLE") == 0 {
				c.writeMsg(message)
				return
			} else if strings.HasPrefix(string(message), "Game Over") {
				c.writeMsg(message)
				return
			} else if strings.HasPrefix(string(message), "GAME_EVENT\nFORFEITED") {
				c.writeMsg(message)
				return
			} else if strings.HasSuffix(string(message), "took too long to make a move") {
				c.writeMsg(message)
				return
			} else if strings.HasPrefix(string(message), "TIME_OUT") {
				c.writeMsg(message)
				return
			} else {
				if err := c.writeMsg(message); err != nil {
					return
				}

			}
		case <-pingTicker.C:
			c.conn.SetWriteDeadline(time.Now().Add(pingPeriod))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}

	}

}

func (c *Client) writeMsg(msg []byte) error {
	w, err := c.conn.NextWriter(websocket.TextMessage)
	if err != nil {
		return err
	}
	w.Write(msg)
	if err := w.Close(); err != nil {
		return err
	}

	return nil
}
