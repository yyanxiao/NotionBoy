package main

import (
	"notionboy/internal/chatgpt"
	"notionboy/internal/server"
)

func main() {
	chatgpt.RefreshSession()
	server.Serve()
}
