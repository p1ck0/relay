package tcpconn

import (
    "net"
    "encoding/json"
)

//ConnectServer -  function for connections to other servers
func ConnectServer(servers []string, myaddr string) {
	for _, serv := range servers {
		go func(serv string) {
			conn, err := net.Dial("tcp", serv)
			if err != nil {
				panic(err)
			}
			var pack = PackageTCP{
				Server: true,
				TCPport: myaddr,
			}
			for names := range aconns {
				pack.Conns = append(pack.Conns, names)
			}
			for ip := range serverstcp {
				pack.Servers = append(pack.Servers, ip)
			}
			data, _ := json.Marshal(pack)
			_, err = conn.Write(data)
			if err != nil {
				panic(err)
			}
		}(serv)
	}
}

//ConnAnotherServer - function for messaging between servers
func ConnAnotherServer(to string, msg PackageTCP) {
	for ip, users := range servers {
		for _, user:= range users{
			if user == to {
				conn, _ := net.Dial("tcp", ip)
				msg.To = []string{to}
				data,_ := json.Marshal(msg)
				conn.Write(data)
				conn.Close()
			}
		}
	}
}
