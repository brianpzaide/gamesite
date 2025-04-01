package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (app *Config) routes() http.Handler {

	r := chi.NewRouter()

	fileServer := http.StripPrefix("/static/", app.fs.FileServer())
	r.Get("/static/*", func(w http.ResponseWriter, r *http.Request) {
		// fmt.Printf("requested url %s\n", r.URL)
		fileServer.ServeHTTP(w, r)
	})
	r.Get("/", app.serveHome)
	r.Get("/gamesite/create/{gametype}", app.createRoom)
	r.Get("/gamesite/rooms/{roomId}", app.getRoom)
	r.Get("/gamesite/rooms/{roomId}/ws", app.serveWs)

	return r
}
