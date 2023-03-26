package main

import (
	"context"
	"flag"
	"hunaidsav/gamesite/games"
	"hunaidsav/gamesite/internal"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

//const uiDir = "./ui/html"

type TemplateData struct {
	Games []struct {
		GameImageURL string
		GameId       string
		GameName     string
	}
}

var (
	gamePages = map[string]string{
		"tttt":    "tttt.html",
		"nttt":    "nttt.html",
		"reversi": "reversi.html",
		"maxit":   "maxit.html",
	}
	gameConstructors = map[string]func() games.Game{
		"tttt":    games.NewTTTTGame,
		"nttt":    games.NewNTTTGame,
		"reversi": games.NewReversiGame,
		"maxit":   games.NewMaxitGame,
	}
	templateData = &TemplateData{
		Games: []struct {
			GameImageURL string
			GameId       string
			GameName     string
		}{
			{GameImageURL: "tttt.png", GameId: "tttt", GameName: "tttt"},
			{GameImageURL: "nttt.png", GameId: "nttt", GameName: "nttt"},
			{GameImageURL: "reversi.png", GameId: "reversi", GameName: "reversi"},
			{GameImageURL: "maxit.png", GameId: "maxit", GameName: "maxit"},
		},
	}

	homePage string
)

type envelope map[string]interface{}

var addr = flag.String("addr", ":8080", "http service address")

func main() {

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	fs := initFS()

	app := &Config{
		fs:       fs,
		InfoLog:  infoLog,
		ErrorLog: errorLog,
	}

	setHomePage(fs)

	for k, v := range gamePages {
		data, err := fs.Read("/" + v) //os.ReadFile(filepath.Join(uiDir, v))
		if err != nil {
			log.Fatal(err)
		}
		gamePages[k] = string(data)
	}

	stopHub := make(chan struct{})
	mainWg := &sync.WaitGroup{}
	internal.StartHub(mainWg, stopHub)

	flag.Parse()

	srv := &http.Server{
		Addr:    *addr,
		Handler: app.routes(),
	}

	idleConnsClosed := make(chan struct{})
	go func() {
		//make a chan to listen for term signals and wait for the signals
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		<-quit
		//shut down the server
		if err := srv.Shutdown(context.Background()); err != nil {
			log.Printf("HTTP server Shutdown: %v", err)
		}
		close(idleConnsClosed)
	}()

	//run the server in the main go routine
	err := srv.ListenAndServe()
	log.Println("Listening on port", addr)
	if err != nil {
		log.Fatalf("HTTP server ListenAndServe: %v", err)
	}
	<-idleConnsClosed

	//send terminating signal to the hub and
	close(stopHub)

	//wait for the hub to shutdown
	mainWg.Wait()
}
