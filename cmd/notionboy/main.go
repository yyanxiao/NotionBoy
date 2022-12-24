package main

import (
	"io"
	"log"
	"notionboy/internal/server"
	"notionboy/internal/telegram"
	"notionboy/internal/wechat"
)

func main() {
	log.SetOutput(io.Discard)
	go wechat.Serve()
	go telegram.Serve()
	server.Serve()
}
