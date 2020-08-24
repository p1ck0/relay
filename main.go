package main

import (
	"github.com/p1ck0/selay/cli"
	"github.com/p1ck0/selay/server"
)
var (
	port        string
	servers     string
)

func init() {
	port = "8888"
	cli.App(&port, &servers)
}

func main() {
	addr := "127.0.0.1:" + port
	server := &server.Server{
		Addr : addr,
		Buff : 1024,
		Servers : servers,
	}
	server.Run()
}
