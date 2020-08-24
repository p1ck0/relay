package main

import (
	"github.com/p1ck0/selay/cli"
	"github.com/p1ck0/selay/launch"
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
	launch.Server(addr, servers)
}
