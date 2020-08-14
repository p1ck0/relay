package tcpconn

import (
    "bufio"
	"encoding/json"
	"fmt"
	"net"
	"sort"
)

//ReciveConn - receive connection and reads the packets 
func ReciveConn(conn net.Conn, msgs chan PackageTCP, dconns chan net.Conn, aconns map[string]net.Conn) {
	rd := bufio.NewReader(conn)
	for {
		var (
			buffer = make([]byte, buff)
			message string
			pack PackageTCP
		)
		length, err := rd.Read(buffer)
		if err != nil { 
			break 
		}
		message += buffer[:length]
		err  = json.Unmarshal(message, &pack)
		if err != nil {
			aconns[message] = conn
		} else if pack.Server == true {
			servers[pack.TCPport] = pack.Conns
			fmt.Println(servers)
			server := []string{pack.TCPport}
			index := sort.SearchStrings(pack.Servers, "127.0.0.1:8081")
			if index == len(pack.Servers) {
				connectServer(server)
			}
		} else {
			msgs <- pack
		}
	}
	dconns <- conn
}

//RedirectPackages - redirects packets to recipient
func RedirectPackages(msg PackageTCP, aconns map[string]net.Conn) {
	for _, to := range msg.To {
		go func(to string, msg PackageTCP){
			conn, ok := aconns[to]
			if !ok {
				connAnotherServer(to, msg)
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