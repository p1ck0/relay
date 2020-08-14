package tcpconn

import (
	"fmt"
	"bufio"
	"encoding/json"
	"net"
	"sort"
)

//ReciveConn - receive connection and reads the packets
func ReciveConn(conn net.Conn, msgs chan PackageTCP, dconns chan net.Conn, aconns map[string]net.Conn, serverstcp map[string][]string, addr string) {
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
		} else if pack.Server == true {
			serverstcp[pack.TCPport] = pack.Conns
			fmt.Println(serverstcp)
			server := []string{pack.TCPport}
			index := sort.SearchStrings(pack.Servers, addr)
			if index == len(pack.Servers) {
				ConnectServer(server, addr, aconns, serverstcp)
				fmt.Println(serverstcp)
			}
		} else {
			msgs <- pack
		}
	}
	dconns <- conn
}

//RedirectPackages - redirects packets to recipient
func RedirectPackages(msg PackageTCP, aconns map[string]net.Conn, serverstcp map[string][]string) {
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
