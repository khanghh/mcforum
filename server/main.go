package main

import (
	"bbs-go/internal/server"
	_ "bbs-go/internal/service/eventhandler"
)

func main() {
	server.Init()
	server.NewServer()
}
