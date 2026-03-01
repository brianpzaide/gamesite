package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/gorilla/websocket"
)

const redisTTL = 180 * time.Second

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func (app *App) serveHome(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(homePage))
}

func (app *App) createRoom(w http.ResponseWriter, r *http.Request) {
	gameType := chi.URLParam(r, "gametype")

	if gameConstructor, ok := gameConstructors[gameType]; ok {
		roomId := app.hub.CreateRoom(gameConstructor, gameType)

		app.rdb.Set(context.Background(), roomId, app.serverID, redisTTL).Err()

		err := writeJSON(w, http.StatusOK, envelope{"roomid": roomId}, nil)
		if err != nil {
			app.errorLog.Println(err)
			http.Error(w, "The server encountered a problem and could not process your request", http.StatusInternalServerError)
		}
		return
	}
	http.NotFound(w, r)
}

func (app *App) getRoom(w http.ResponseWriter, r *http.Request) {
	roomId := chi.URLParam(r, "roomId")

	gameType, roomExists := app.hub.RoomExists(roomId)

	if roomExists {
		gamePage := gamePages[gameType]
		fmt.Fprintf(w, gamePage, roomId)
		return
	}
	w.WriteHeader(http.StatusNotFound)
	fmt.Fprintf(w, "No room found with Id= %s", roomId)
}

func (app *App) serveWs(w http.ResponseWriter, r *http.Request) {

	roomId := chi.URLParam(r, "roomId")

	conn, err := upgrader.Upgrade(w, r, nil)
	log.Println("connected", conn.RemoteAddr())

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("could not upgrade to websocket connection"))
		return
	}

	err = app.hub.AddClient(roomId, conn)

	if err != nil {
		log.Println("client not created")
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "No room found with Id= %s", roomId)
		return
	}
}
