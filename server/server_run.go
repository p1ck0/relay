package server

import (
	"fmt"
	"log"
	"net"
	"sync"
	"strings"
	pack "github.com/p1ck0/selay/packagetcp"
)

var mut sync.Mutex

//Server - server structure
type Server struct {
	Buff       int
	Addr       string
	Servers    string
	Aconns     map[string]net.Conn
	TCPconns   chan net.Conn
	ServersTCP map[string]string
	Dconns     chan net.Conn
	Msgs       chan *pack.PackageTCP
}

//Run - starts the server
func (serv *Server) Run() {
	serv.Aconns = make(map[string]net.Conn)
	serv.TCPconns = make(chan net.Conn)
	serv.ServersTCP = make(map[string]string)
	serv.Dconns = make(chan net.Conn)
	serv.Msgs = make(chan *pack.PackageTCP)
	serv.Listen()
}

//Listen - listens for incoming connections
func (serv *Server) Listen() {
	ln, err := net.Listen("tcp", serv.Addr)
	if err != nil {
		log.Fatalln(err.Error())
	}
	if len(serv.Servers) > 0 {
		serversArr := strings.Split(serv.Servers, " ")
		serv.ConnectServer(serversArr)
	}
	defer ln.Close()
	go func() {
		for {
			conn, err := ln.Accept()
			if err != nil {
				log.Fatalln(err.Error())
			}
			serv.TCPconns <- conn
		}
	}()
	serv.Handle()
}

//Handle - handles incoming connections
func (serv *Server) Handle() {
	for {
		select {
		case conn := <-serv.TCPconns:
			fmt.Println(conn.RemoteAddr().String())
			go serv.ReciveConn(conn)

		case msg := <-serv.Msgs:
			go serv.RedirectPackages(msg)

		case dconn := <-serv.Dconns:
			defer dconn.Close()
			for name, conn := range serv.Aconns {
				if conn == dconn {
					log.Printf("Client %v was gone\n", name)
					dconn.Close()
					delete(serv.Aconns, name)
					serv.DelUser(name)
				}
			}
		}
	}
}
