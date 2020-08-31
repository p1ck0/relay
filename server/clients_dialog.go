package server

import (
	"bufio"
	"encoding/json"
	"fmt"
    "net"
    pack "github.com/p1ck0/selay/packagetcp"
)

//ReciveConn - receive connection and reads the packets
func (serv *Server) ReciveConn(conn *net.Conn) {
	rd := bufio.NewReader(*conn)
	for {
		var (
			buffer  = make([]byte, serv.Buff)
			message string
			pack    *pack.PackageTCP
		)
		length, err := rd.Read(buffer)
		if err != nil {
			break
		}
        message += string(buffer[:length])
		err = json.Unmarshal([]byte(message), &pack)
		if err != nil {
			fmt.Println(err)
        }
        fmt.Println(pack.Head.UserMod.RegUser)
		switch {
		case pack.Head.ServerInfo.Server:
            serv.RegistServer(pack, conn)
        case pack.Head.UserMod.RegUser:
			serv.RegistUser(pack, conn)
		case pack.Head.UserMod.NewUser:
			mut.Lock()
			serv.ServersTCP[pack.Head.UserMod.User] = pack.Head.From
			mut.Unlock()
		case pack.Head.UserMod.DelUser:
			delete(serv.ServersTCP, pack.Head.UserMod.User)
		default:
			serv.Msgs <- pack
		}
	}
	serv.Dconns <- *conn
}

//RedirectPackages - redirects packets to recipient
func (serv *Server) RedirectPackages(msg *pack.PackageTCP) {
	for _, to := range msg.Head.To {
		go func(to string, msg pack.PackageTCP) {
			mut.Lock()
			conn, ok := serv.Aconns[to]
			mut.Unlock()
			if !ok {
				serv.ConnAnotherServer(to, &msg)
				return
			}
			data, err := json.Marshal(msg)
			if err != nil {
				panic(err)
			}
			conn.Write([]byte(data))
		}(to, *msg)
	}
}

//RegistUser - registers a new client
func (serv *Server) RegistUser(pack *pack.PackageTCP, conn *net.Conn) {
	mut.Lock()
	serv.Aconns[pack.Head.UserMod.User] = *conn
	mut.Unlock()
	if len(serv.ServersTCP) > 0 {
		serv.NewUser(pack)
	}
}

//RegistServer - registers a new server
func (serv *Server) RegistServer(pack *pack.PackageTCP, conn *net.Conn) {
	mut.Lock()
	if len(pack.Head.ServerInfo.Conns) > 0 {
		for _, conn := range pack.Head.ServerInfo.Conns {
			serv.ServersTCP[conn] = pack.Head.ServerInfo.TCPport
		}
	} else {
		serv.ServersTCP[""] = pack.Head.ServerInfo.TCPport
	}
	mut.Unlock()
	fmt.Println(serv.ServersTCP)
	server := []string{pack.Head.ServerInfo.TCPport}
	if _, ok := pack.Head.ServerInfo.Servers[serv.Addr]; !ok {
		serv.ConnectServer(server)
		fmt.Println(serv.ServersTCP)
	}
}
