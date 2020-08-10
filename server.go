package main

import (
	//"bufio"
	//"encoding/json"
	"fmt"
	"log"
	"net"
	tcp "./tcpconn"
)

var BUFF = 1024

type PackageTCP struct {
	From string
	To   []string
	Body interface{}
}

var (
	aconns  = make(map[string]net.Conn)
	conns   = make(chan net.Conn)
	dconns  = make(chan net.Conn)
	servers = make(chan net.Conn)
	msgs    = make(chan tcp.PackageTCP)
)

func main() {
	ln, err := net.Listen("tcp", "127.0.0.1:8081")
	if err != nil {
		log.Fatalln(err.Error())
	}
	defer ln.Close()
	go func() {
		for {
			conn, err := ln.Accept()
			if err != nil {
				log.Fatalln(err.Error())
			}
			conns <- conn
		}
	}()

	for {
		select {
		case conn := <-conns:
			fmt.Println(conn.RemoteAddr().String())
			go tcp.ReciveConn(conn, msgs, dconns, aconns)

		case msg := <-msgs:
			go tcp.RedirectPackages(msg, aconns)

		case dconn := <-dconns:
			defer dconn.Close()
			for name, conn := range aconns {
				if conn == dconn {
					log.Printf("Client %v was gone\n", name)
					dconn.Close()
					delete(aconns, name)
				}
			}
		}
	}
}

func connServer(ipServ string) {
	conn, _ := net.Dial("udp", ipServ)
	ls, _ := net.Listen("udp", "127.0.0.1:8082")
	conn, _ = ls.Accept()
	servers <- conn
}
