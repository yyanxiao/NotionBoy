package main

import (
	"io"
	"log"

	"notionboy/internal/scheduler"
	"notionboy/internal/server"
	"notionboy/internal/telegram"
)

func main() {
	log.SetOutput(io.Discard)
	// go wechat.Serve()
	go telegram.Serve()
	go scheduler.Run()
	server.Serve()
}
