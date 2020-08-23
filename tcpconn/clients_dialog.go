package tcpconn

import (
	"fmt"
	"bufio"
	"encoding/json"
	"net"
)

//ReciveConn - receive connection and reads the packets
func ReciveConn(conn net.Conn, msgs chan PackageTCP, dconns chan net.Conn, aconns map[string]net.Conn, serverstcp map[string]string, addr string) {
	rd := bufio.NewReader(conn)
	for {
		var (
			buffer  = make([]byte, buff)
			message string
			pack    PackageTCP
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
		switch {
		case pack.Head.UserMod.RegUser:
			aconns[pack.Head.UserMod.User] = conn
			if len(serverstcp) > 0 {
				NewUser(pack, addr, serverstcp)
			} 
		case pack.Head.ServerInfo.Server:
			if len(pack.Head.ServerInfo.Conns) > 0 {
				for _, conn := range pack.Head.ServerInfo.Conns {
					serverstcp[conn] = pack.Head.ServerInfo.TCPport
				}
			} else {
				serverstcp[""] = pack.Head.ServerInfo.TCPport
			}
			fmt.Println(serverstcp)
			server := []string{pack.Head.ServerInfo.TCPport}
			if _,ok := pack.Head.ServerInfo.Servers[addr];!ok {
				ConnectServer(server, addr, aconns, serverstcp)
				fmt.Println(serverstcp)
			}
		case pack.Head.UserMod.NewUser:
			fmt.Println(pack.Head.UserMod.User)
			serverstcp[pack.Head.UserMod.User] = pack.Head.From
		case pack.Head.UserMod.DelUser:
			delete(serverstcp, pack.Head.UserMod.User)
		default:
			msgs <- pack
		}
	}
	dconns <- conn
}

//RedirectPackages - redirects packets to recipient
func RedirectPackages(msg PackageTCP, aconns map[string]net.Conn, serverstcp map[string]string) {
	for _, to := range msg.Head.To {
		go func(to string, msg PackageTCP) {
			conn, ok := aconns[to]
			if !ok {
				ConnAnotherServer(to, msg, serverstcp)
				return
			}
			data, err := json.Marshal(msg)
			if err != nil {
				panic(err)
			}
			conn.Write([]byte(data))
		}(to, msg)
	}
}
