package main

import (
	"hunaidsav/gamesite/internal"
	"log"

	"github.com/knadh/stuffbin"
	"github.com/redis/go-redis/v9"
)

type App struct {
	fs       stuffbin.FileSystem
	infoLog  *log.Logger
	errorLog *log.Logger
	hub      *internal.Hub
	rdb      *redis.Client
	serverID string
}
