package internal

import (
	"fmt"
	"hunaidsav/gamesite/games"
	"strconv"
	"strings"
	"time"
)

type ClientMessage struct {
	client *Client
	msg    []byte
}

type Room struct {
	clients       []*Client
	broadcast     chan *ClientMessage
	register      chan *Client
	unregister    chan *ClientMessage
	players       map[*Client]int
	gameType      string
	game          games.Game
	roomId        string
	terminateChan chan bool
}

func newRoom(gameConstructor func() games.Game, gametype string, roomId string, terminateChan chan bool) *Room {

	return &Room{
		broadcast:  make(chan *ClientMessage),
		register:   make(chan *Client),
		unregister: make(chan *ClientMessage),
		//clients read and write messages to the users(browsers)
		clients: make([]*Client, 0),
		//players is used to make sure that only the game's current player is making the move
		//and also to know who the winner or loser is
		players:  make(map[*Client]int),
		gameType: gametype,
		game:     gameConstructor(),
		roomId:   roomId,
		//send signal to close the room
		terminateChan: terminateChan,
	}
}

func (r *Room) run() {

	terminateSignal := false

	ticker := time.NewTicker(readWait)
	unregisterCount := 0
	currentPlayer := "X"

	defer func() {
		ticker.Stop()
		if !terminateSignal {
			ch := make(chan bool)
			HubChan <- &DeleteRoomMsg{RoomId: r.roomId, ReceiveChan: ch}
			<-ch
		}
		hubWg.Done()
	}()

	for {
		select {
		case <-r.terminateChan:
			terminateSignal = true
			return
		case <-ticker.C:
			if len(r.players) == 1 {
				r.clients[0].send <- []byte("TIME_OUT\nNo one joined the game")
			} else if len(r.players) >= 2 {
				i := 0
				for i < len(r.clients) {
					select {
					case r.clients[i].send <- []byte(fmt.Sprintf("TIME_OUT\n%s took too long to make a move", currentPlayer)):
						i++
					default:
						i++
					}

				}
			} else {
				ticker.Reset(readWait)
			}

		case client := <-r.register:
			if len(r.players) >= 2 {
				client.send <- []byte("ROOM_UNAVAILABLE")
			} else {
				ticker.Reset(readWait)

				if len(r.players) == 1 {
					r.players[client] = games.PLAYER2
					r.clients = append(r.clients, client)
					initData := r.game.GetInitData()
					if initData != "" {
						r.clients[0].send <- []byte(fmt.Sprintf("GAME_START\nYou are player %s\n%s", r.game.GetCurrentPlayerSymbol(), initData))
						r.clients[1].send <- []byte(fmt.Sprintf("GAME_START\nYou are player %s\n%s", r.game.GetOtherPlayerSymbol(), initData))
					} else {
						r.clients[0].send <- []byte(fmt.Sprintf("GAME_START\nYou are player %s\n%d", r.game.GetCurrentPlayerSymbol(), 2))
						r.clients[1].send <- []byte(fmt.Sprintf("GAME_START\nYou are player %s\n%d", r.game.GetOtherPlayerSymbol(), 3))
					}

				} else {
					r.players[client] = games.PLAYER1
					r.clients = append(r.clients, client)
					r.clients[0].send <- []byte("GAME_WAIT\nWaiting for the other player to join")
				}
			}

		case clmsg := <-r.unregister:
			client := clmsg.client
			if _, ok := r.players[client]; ok {
				if len(r.clients) == 2 && unregisterCount < 1 {
					unregisterCount++
				} else {
					return
				}

			}

		case clmsg := <-r.broadcast:
			client := clmsg.client
			msg := clmsg.msg
			moveby, ok := r.players[client]
			if ok && len(r.players) == 2 {
				if strings.Compare(string(msg), "opponent left the game") == 0 {
					i := 0
					for i < len(r.clients) {
						select {
						case r.clients[i].send <- []byte(fmt.Sprintf("GAME_EVENT\nFORFEITED\n%s left the game", currentPlayer)):
							i++
						default:
							i++
						}
					}

				} else if moveby == r.game.GetCurrentPlayer() {

					msgSplit := strings.Split(string(msg), " ")
					// fmt.Println(string(msg))
					msgArgs := make([]int, 0)
					for _, msgElement := range msgSplit {
						if k, err := strconv.Atoi(msgElement); err == nil {
							msgArgs = append(msgArgs, k)
						}
					}

					status, data := r.game.PerformMove(msgArgs...)
					if status >= 0 {
						if status != games.IN_PROGRESS {
							for _, cl := range r.clients {
								if status == games.DRAW {
									cl.send <- []byte(fmt.Sprintf("GAME_EVENT\n%s", "Draw"))
								} else {
									if status == r.players[cl] {
										cl.send <- []byte(fmt.Sprintf("GAME_EVENT\n%s", "you won"))
									} else {
										cl.send <- []byte(fmt.Sprintf("GAME_EVENT\n%s", "you lost"))
									}
								}
							}
						} else {
							prevPlayer := r.game.GetCurrentPlayer()
							r.game.UpdateCurrentPlayer()
							currentPlayer = r.game.GetCurrentPlayerSymbol()

							for _, cl := range r.clients {
								cl.send <- r.buildMoveEventResponse(currentPlayer, prevPlayer, data) // cl.send <- r.buildMoveEventResponse(currentPlayer, r.game.GetCurrentPlayer(), msgArgs) //[]byte(fmt.Sprintf("MOVE_EVENT\n%s\n%d\n%d\n%d\n%d", currentPlayer, r.game.GetCurrentPlayer(), msgArgs[0], msgArgs[1], msgArgs[2]))
							}

						}
					}

					ticker.Reset(readWait)
				}

			}

		}
	}
}

func (r *Room) buildMoveEventResponse(currentPlayerSymbol string, prevPlayer int, data []string) []byte {
	sb := strings.Builder{}
	sb.WriteString(fmt.Sprintf("MOVE_EVENT\n%s\n%d", currentPlayerSymbol, prevPlayer))

	for _, el := range data {
		sb.WriteString(fmt.Sprintf("\n%s", el))
	}
	return []byte(sb.String())
}
