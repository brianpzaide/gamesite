package main

import (
	"hunaidsav/gamesite/internal"
	"log"

	"github.com/knadh/stuffbin"
)

type App struct {
	fs       stuffbin.FileSystem
	InfoLog  *log.Logger
	ErrorLog *log.Logger
	Hub      *internal.Hub
}
