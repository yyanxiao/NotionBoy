package main

import (
	"io"
	"log"
	"notionboy/internal/server"
	"notionboy/internal/wechat"
)

func main() {
	log.SetOutput(io.Discard)
	go wechat.Serve()
	server.Serve()
}
