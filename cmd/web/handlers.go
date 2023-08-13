package main

import (
	"fmt"
	"hunaidsav/gamesite/internal"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func (app *Config) serveHome(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(homePage))
}

func (app *Config) createRoom(w http.ResponseWriter, r *http.Request) {
	gameType := chi.URLParam(r, "gametype")

	if gameConstructor, ok := gameConstructors[gameType]; ok {
		ch := make(chan string)
		msg := &internal.CreateRoomMsg{GameConstructor: gameConstructor, GameType: gameType, ReceiveChan: ch}
		//sending the message to create a room
		internal.HubChan <- msg
		//waiting for the message roomId from the hub
		roomId := <-ch
		err := writeJSON(w, http.StatusOK, envelope{"roomid": roomId}, nil)
		if err != nil {
			app.ErrorLog.Println(err)
			http.Error(w, "The server encountered a problem and could not process your request", http.StatusInternalServerError)
		}
		return
	}
	http.NotFound(w, r)
}

func (app *Config) getRoom(w http.ResponseWriter, r *http.Request) {
	roomId := chi.URLParam(r, "roomId")

	ch := make(chan string)
	msg := &internal.RoomExistsMsg{RoomId: roomId, ReceiveChan: ch}
	//sending the message to create a room
	internal.HubChan <- msg
	//waiting for the message gameType from the hub
	gameType := <-ch

	if gameType != "" {
		gamePage := gamePages[gameType]

		w.Write([]byte(fmt.Sprintf(gamePage, root_url, root_url, roomId)))
		return
	}
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte(fmt.Sprintf("No room found with Id= %s", roomId)))
}

func (app *Config) serveWs(w http.ResponseWriter, r *http.Request) {

	roomId := chi.URLParam(r, "roomId")

	fmt.Println("Main: connecting to hub", roomId)

	conn, err := upgrader.Upgrade(w, r, nil)
	log.Println("connected", conn.RemoteAddr())

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("could not upgrade to websocket connection"))
		return
	}

	ch := make(chan bool)
	msg := &internal.CreateClientMsg{RoomId: roomId, Conn: conn, ReceiveChan: ch}
	//sending the message to create a client
	internal.HubChan <- msg
	//waiting for the confirmation if client is created from the hub
	ok := <-ch

	if !ok {
		log.Println("client created:", ok)
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(fmt.Sprintf("No room found with Id= %s", roomId)))
		return
	}
}
