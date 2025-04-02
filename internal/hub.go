package internal

import (
	"context"
	"hunaidsav/gamesite/games"
	"sync"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/redis/go-redis/v9"
)

const (
	createRoomMsg        = "CREATE_ROOM"
	roomExistsMsg        = "ROOM_EXISTS"
	createClientMsg      = "CREATE_CLIENT"
	deleteRoomMsg        = "DELETE_ROOM"
	terminateAllRoomsMsg = "TERMINATE_ALL_ROOMS"
)

type HubMessage interface {
	getTitle() string
}

type CreateRoomMsg struct {
	GameConstructor func() games.Game
	GameType        string
	ReceiveChan     chan string
}

func (m *CreateRoomMsg) getTitle() string {
	return createRoomMsg
}

type RoomExistsMsg struct {
	RoomId      string
	ReceiveChan chan string
}

func (m *RoomExistsMsg) getTitle() string {
	return roomExistsMsg
}

type CreateClientMsg struct {
	RoomId      string
	Conn        *websocket.Conn
	ReceiveChan chan bool
}

func (m *CreateClientMsg) getTitle() string {
	return createClientMsg
}

type DeleteRoomMsg struct {
	RoomId      string
	ReceiveChan chan bool
}

func (m *DeleteRoomMsg) getTitle() string {
	return deleteRoomMsg
}

type TerminateAllRoomsMsg struct {
	ReceiveChan chan bool
}

func (m *TerminateAllRoomsMsg) getTitle() string {
	return terminateAllRoomsMsg
}

var (
	HubChan  = make(chan HubMessage)
	doerChan = make(chan HubMessage, 1024)
	rooms    = make(map[string]struct {
		room          *Room
		terminateChan chan bool
	})
	hubWg  = &sync.WaitGroup{}
	mainWg *sync.WaitGroup
	rdb    *redis.Client
	sID    string
)

func runListener(terminateChan chan struct{}) {
	var queue []HubMessage

	for {
		select {
		case <-terminateChan:
			close(doerChan)
			return
		case msg := <-HubChan:
			queue = append(queue, msg)
			for len(queue) > 0 {
				select {
				case doerChan <- queue[0]:
					queue = queue[1:]
				case x, ok := <-HubChan:
					if !ok {
						for _, x := range queue {
							doerChan <- x
						}
						close(doerChan)
						return
					}
					queue = append(queue, x)
				}
			}
		}
	}
}

func runDoer() {
	for msg := range doerChan {
		switch msg.getTitle() {
		case createRoomMsg:
			m := msg.(*CreateRoomMsg)
			m.ReceiveChan <- createRoom(m.GameConstructor, m.GameType)
		case roomExistsMsg:
			m := msg.(*RoomExistsMsg)
			m.ReceiveChan <- roomExists(m.RoomId)
		case createClientMsg:
			m := msg.(*CreateClientMsg)
			m.ReceiveChan <- createClient(m.RoomId, m.Conn)
		case deleteRoomMsg:
			m := msg.(*DeleteRoomMsg)
			m.ReceiveChan <- deleteRoom(m.RoomId)
		}
	}
	shutdownAllRooms()
}

func StartHub(redisClient *redis.Client, serverID string, wg *sync.WaitGroup, terminateChan chan struct{}) {
	rdb = redisClient
	sID = serverID
	mainWg = wg
	mainWg.Add(2)

	go func() {
		defer mainWg.Done()
		runListener(terminateChan)
	}()

	go func() {
		defer mainWg.Done()
		runDoer()
	}()
}

func createRoom(gameConstructor func() games.Game, gameType string) string {

	id := uuid.New().String()
	terminateChan := make(chan bool)
	// fmt.Println("room created with id", id)

	hubWg.Add(1)
	room := newRoom(gameConstructor, gameType, id, terminateChan)
	// log.Println("room created")
	rooms[id] = struct {
		room          *Room
		terminateChan chan bool
	}{room, terminateChan}
	// log.Println("size of rooms: ", len(rooms))

	// insert the roomid to the redis
	rdb.Set(context.Background(), id, sID, redisTTL).Err()

	go room.run()

	return id
}

func roomExists(roomId string) string {
	if roomstruct, found := rooms[roomId]; found {
		return roomstruct.room.gameType
	} else {
		return ""
	}
}

func createClient(roomId string, conn *websocket.Conn) bool {
	if roomStruct, found := rooms[roomId]; found {
		client := &Client{room: roomStruct.room, conn: conn, send: make(chan []byte)}
		// log.Println("in the hub.createClient before regestering client")
		client.room.register <- client
		go client.readPump()
		go client.writePump()
		// log.Println("in the hub.createClient after running off readpump and writepump")
		return found
	} else {
		return found
	}
}

func deleteRoom(roomId string) bool {
	roomStruct, ok := rooms[roomId]
	if ok {
		close(roomStruct.terminateChan)
		delete(rooms, roomId)
		return true
	}
	return false

}

func shutdownAllRooms() bool {
	for _, val := range rooms {
		val.terminateChan <- true
	}
	hubWg.Wait()
	rooms = nil
	return true
}
