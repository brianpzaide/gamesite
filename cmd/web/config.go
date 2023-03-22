package main

import (
	"log"

	"github.com/knadh/stuffbin"
)

type Config struct {
	fs       stuffbin.FileSystem
	InfoLog  *log.Logger
	ErrorLog *log.Logger
}
