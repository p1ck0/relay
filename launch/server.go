package launch

import (
    "fmt"
    "net"
    "log"
    "strings"
    tcp "github.com/p1ck0/selay/tcpconn"
)


var (
    buff        = 1024
	aconns      = make(map[string]net.Conn)
	tcpconns    = make(chan net.Conn)
	udpconns    = make(chan net.UDPConn)
	serverstcp  = make(map[string]string)
	dconns      = make(chan net.Conn)
	msgs        = make(chan tcp.PackageTCP)
)

//Server - launch server
func Server(addr ,servers string) {
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
    handle(addr, serverstcp)

}

func handle(addr string, serverstcp map[string]string) {
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