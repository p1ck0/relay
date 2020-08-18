package tcpconn

import (
	"encoding/json"
	"net"
)

//ConnectServer -  function for connections to other servers
func ConnectServer(servers []string, myaddr string, aconns map[string]net.Conn, serverstcp map[string]string) {
	for _, serv := range servers {
		go func(serv string) {
			conn, err := net.Dial("tcp", serv)
			if err != nil {
				panic(err)
			}
			defer conn.Close()
			var pack = PackageTCP{
				Server:  true,
				TCPport: myaddr,
			}
			for names := range aconns {
				pack.Conns = append(pack.Conns, names)
			}
			pack.Servers = make(map[string]bool)
			for _,ip := range serverstcp {
				pack.Servers[ip] = false
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
func ConnAnotherServer(to string, msg PackageTCP, serverstcp map[string]string) {
	if ip,ok := serverstcp[to];ok{
		conn, _ := net.Dial("tcp", ip)
		msg.To = []string{to}
		data, _ := json.Marshal(msg)
		conn.Write(data)
		conn.Close()
	}
}

//NewUser - function for registering users on all servers
func NewUser(user string, addr string, serverstcp map[string]string) {
	var pack = PackageTCP{
		NewUser : true,
		From : addr,
		User : user,
	}
	data, _ := json.Marshal(pack)
	for _,ip := range serverstcp {
		conn, _ := net.Dial("tcp", ip)
		conn.Write(data)
		conn.Close()
	}
}

//DelUser - function for deleting users on all servers
func DelUser(user string, addr string, serverstcp map[string]string) {
	var pack = PackageTCP{
		DelUser : true,
		User : user,
	}
	data, _ := json.Marshal(pack)
	for _,ip := range serverstcp {
		conn, _ := net.Dial("tcp", ip)
		conn.Write(data)
		conn.Close()
	}
}