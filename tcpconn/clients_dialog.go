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
			aconns[message] = conn
			if len(serverstcp) > 0 {
				NewUser(message, addr, serverstcp)
			}
		}
		switch {
		case pack.Server:
			if len(pack.Conns) > 0 {
				for _, conn := range pack.Conns {
					serverstcp[conn] = pack.TCPport
				}
			} else {
				serverstcp[""] = pack.TCPport
			}
			fmt.Println(serverstcp)
			server := []string{pack.TCPport}
			if _,ok := pack.Servers[addr];!ok {
				ConnectServer(server, addr, aconns, serverstcp)
				fmt.Println(serverstcp)
			}
		case pack.NewUser:
			serverstcp[pack.User] = pack.From
		case pack.DelUser:
			delete(serverstcp, pack.User)
		default:
			msgs <- pack
		}
	}
	dconns <- conn
}

//RedirectPackages - redirects packets to recipient
func RedirectPackages(msg PackageTCP, aconns map[string]net.Conn, serverstcp map[string]string) {
	for _, to := range msg.To {
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
