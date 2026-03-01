package internal

import (
	"fmt"
	"hunaidsav/gamesite/games"
	"strconv"
	"strings"
	"time"
)

const moveTimeout = 60 * time.Second

type ClientMessage struct {
	client *Client
	msg    []byte
}

type Room struct {
	ID         string
	game       games.Game
	gameType   string
	hub        *Hub
	register   chan *Client
	unregister chan *ClientMessage
	broadcast  chan *ClientMessage
	stop       chan bool

	clients []*Client
	players map[*Client]int

	currentPlayerSymbol string
}

func NewRoom(id string, game games.Game, gameType string, hub *Hub) *Room {
	return &Room{
		ID:                  id,
		game:                game,
		gameType:            gameType,
		hub:                 hub,
		register:            make(chan *Client),
		unregister:          make(chan *ClientMessage),
		broadcast:           make(chan *ClientMessage),
		stop:                make(chan bool),
		clients:             []*Client{},
		players:             map[*Client]int{},
		currentPlayerSymbol: game.GetCurrentPlayerSymbol(),
	}
}

func (r *Room) Run() {
	ticker := time.NewTicker(moveTimeout)
	defer ticker.Stop()

	for {
		select {
		case client := <-r.register:
			fmt.Println("loop: player registered")
			r.handleRegister(client, ticker)

		case clMsg := <-r.unregister:
			fmt.Println("loop: player unregistered")
			if toContinue := r.handleUnregister(clMsg); !toContinue {
				return
			}

		case clMsg := <-r.broadcast:
			fmt.Println("loop: player made a move")
			if toContinue := r.handleBroadcast(clMsg, ticker); !toContinue {
				return
			}

		case <-ticker.C:
			fmt.Println("loop: ticker time out")
			r.handleTimeout(ticker)

		case <-r.stop:
			fmt.Println("loop: room called to stop")
			return
		}
	}
}

func (r *Room) handleRegister(client *Client, ticker *time.Ticker) {
	if len(r.players) >= 2 {
		client.send <- []byte("ROOM_UNAVAILABLE")
		return
	}
	ticker.Reset(moveTimeout)
	r.clients = append(r.clients, client)
	if len(r.players) == 0 {
		r.players[client] = games.PLAYER1
		client.send <- []byte("GAME_WAIT\nWaiting for the other player to join")
	} else {
		r.players[client] = games.PLAYER2
		initData := r.game.GetInitData()
		if initData != "" {
			r.clients[0].send <- []byte(fmt.Sprintf("GAME_START\nYou are player %s\n%s", r.game.GetCurrentPlayerSymbol(), initData))
			r.clients[1].send <- []byte(fmt.Sprintf("GAME_START\nYou are player %s\n%s", r.game.GetOtherPlayerSymbol(), initData))
		} else {
			r.clients[0].send <- []byte(fmt.Sprintf("GAME_START\nYou are player %s\n%d", r.game.GetCurrentPlayerSymbol(), 2))
			r.clients[1].send <- []byte(fmt.Sprintf("GAME_START\nYou are player %s\n%d", r.game.GetOtherPlayerSymbol(), 3))
		}
	}
}

func (r *Room) handleUnregister(msg *ClientMessage) bool {
	client := msg.client
	if _, ok := r.players[client]; !ok {
		return true
	}

	for _, cl := range r.clients {
		if cl != client {
			cl.send <- []byte("GAME_EVENT\nFORFEITED\nOpponent left the game")
		}
	}
	r.removeClient(client)
	return false
}

func (r *Room) handleBroadcast(msg *ClientMessage, ticker *time.Ticker) bool {
	client := msg.client
	moveBy, ok := r.players[client]
	if !ok || len(r.players) != 2 {
		return true
	}

	if moveBy != r.game.GetCurrentPlayer() {
		// fmt.Println("hub:line: 126 :: moveBy != r.game.GetCurrentPlayer()")
		return true
	}

	ticker.Reset(moveTimeout)

	args := parseMove(msg.msg)
	fmt.Println("hub:line:132 :: ", args)
	status, data := r.game.PerformMove(args...)

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
			return false
		} else {
			prevPlayer := r.game.GetCurrentPlayer()
			r.game.UpdateCurrentPlayer()
			currentPlayer := r.game.GetCurrentPlayerSymbol()

			for _, cl := range r.clients {
				cl.send <- r.buildMoveEventResponse(currentPlayer, prevPlayer, data)
			}
			return true
		}
	}
	return true
}

func (r *Room) handleTimeout(ticker *time.Ticker) {
	if len(r.players) == 1 {
		r.clients[0].send <- []byte("TIME_OUT\nNo one joined the game")
	} else if len(r.players) >= 2 {
		for _, cl := range r.clients {
			if r.players[cl] == r.game.GetCurrentPlayer() {
				cl.send <- []byte("GAME_EVENT\nYou lost\nTook too long to make a move")
			} else {
				cl.send <- []byte(fmt.Sprintf("TIME_OUT\n%s took too long to make a move", r.currentPlayerSymbol))
			}

		}
	} else {
		ticker.Reset(moveTimeout)
	}

}

func (r *Room) removeClient(client *Client) {
	delete(r.players, client)

	for i, cl := range r.clients {
		if cl == client {
			r.clients = append(r.clients[:i], r.clients[i+1:]...)
			break
		}
	}
}

func (r *Room) buildMoveEventResponse(currentSymbol string, prevPlayer int, data []string) []byte {
	sb := strings.Builder{}
	sb.WriteString(fmt.Sprintf("MOVE_EVENT\n%s\n%d", currentSymbol, prevPlayer))
	for _, d := range data {
		sb.WriteString(fmt.Sprintf("\n%s", d))
	}
	return []byte(sb.String())
}

func parseMove(msg []byte) []int {
	parts := strings.Fields(string(msg))
	var args []int
	for _, p := range parts {
		if n, err := strconv.Atoi(p); err == nil {
			args = append(args, n)
		}
	}
	return args
}

func (r *Room) Stop() {
	close(r.stop)
}
