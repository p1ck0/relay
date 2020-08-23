package main

import (
	"fmt"
	tcp "github.com/p1ck0/selay/tcpconn"
	"github.com/p1ck0/selay/cli"
	"log"
	"net"
	"strings"
)

var (
	buff        = 1024
	aconns      = make(map[string]net.Conn)
	tcpconns    = make(chan net.Conn)
	udpconns    = make(chan net.UDPConn)
	serverstcp  = make(map[string]string)
	dconns      = make(chan net.Conn)
	msgs        = make(chan tcp.PackageTCP)
	command     = make(chan string)
	serversconn []string
	port        string
	servers     string
)

func init() {
	port = "8888"
	cli.App(&port, &servers)
}

func main() {
	addr := "127.0.0.1:" + port
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalln(err.Error())
	}
	if len(servers) > 0 {
		serversArr := strings.Split(servers, " ")
		tcp.ConnectServer(serversArr, addr, aconns, serverstcp)
	}
	defer ln.Close()
	go func() {
		for {
			conn, err := ln.Accept()
			if err != nil {
				log.Fatalln(err.Error())
			}
			tcpconns <- conn
		}
	}()
	for {
		select {
		case conn := <-tcpconns:
			fmt.Println(conn.RemoteAddr().String())
			go tcp.ReciveConn(conn, msgs, dconns, aconns, serverstcp, addr)

		case msg := <-msgs:
			go tcp.RedirectPackages(&msg, aconns, serverstcp)

		case dconn := <-dconns:
			defer dconn.Close()
			for name, conn := range aconns {
				if conn == dconn {
					log.Printf("Client %v was gone\n", name)
					dconn.Close()
					delete(aconns, name)
					tcp.DelUser(name, addr, serverstcp)
				}
			}
		}
	}
}
