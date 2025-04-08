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
		"poc":     "poc.html",
	}
	gameConstructors = map[string]func() games.Game{
		"tttt":    games.NewTTTTGame,
		"nttt":    games.NewNTTTGame,
		"reversi": games.NewReversiGame,
		"maxit":   games.NewMaxitGame,
		"poc":     games.NewPOCGame,
	}

	templateData *TemplateData

	homePage string
)

type envelope map[string]interface{}

var addr = flag.String("addr", ":8080", "http service address")
var redisAddr = flag.String("redis-addr", os.Getenv("GAMESITE_REDIS_ADDR"), "REDIS ADDR")
var serverID = flag.String("server-id", os.Getenv("GAMESITE_SERVER_ID"), "SERVER ID")

func main() {
	flag.Parse()

	if *serverID == "" {
		log.Fatalln("server id not provided. Server id must be set either as a flag `--server-id` or as an environment variable `GAMESITE_SERVER_ID`.")
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
			{GameImageURL: "poc.png", GameId: "poc", GameName: "pawns only chess"},
		},
	}

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
		data, err := fs.Read("/" + v)
		if err != nil {
			log.Fatal(err)
		}
		gamePages[k] = string(data)
	}

	if *redisAddr == "" {
		*redisAddr = "localhost:6379"
	}

	rdb, err := getRedisClient()
	if err != nil {
		log.Fatalln("Error connecting to redis.")
	}

	stopHub := make(chan struct{})
	mainWg := &sync.WaitGroup{}
	internal.StartHub(rdb, *serverID, mainWg, stopHub)

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
	err = srv.ListenAndServe()
	log.Println("Listening on port", addr)
	if err != nil {
		log.Fatalf("HTTP server ListenAndServe: %v", err)
	}
	<-idleConnsClosed

	//send terminating signal to the hub and
	close(stopHub)

	//wait for the hub to shutdown
	mainWg.Wait()

	rdb.Close()
}
