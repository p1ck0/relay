package server

import (
	"encoding/json"
    "net"
    pack "github.com/p1ck0/selay/packagetcp"
)

//ConnectServer -  function for connections to other servers
func (serv *Server) ConnectServer(servers []string) {
	for _, server := range servers {
		go func(server string) {
			conn, err := net.Dial("tcp", server)
			if err != nil {
				panic(err)
			}
			defer conn.Close()
			var pack = &pack.PackageTCP{
				Head: pack.Head{
					ServerInfo: pack.Server{
						Server:  true,
						TCPport: serv.Addr,
					},
				},
			}
			for names := range serv.Aconns {
				pack.Head.ServerInfo.Conns = append(pack.Head.ServerInfo.Conns, names)
			}
			pack.Head.ServerInfo.Servers = make(map[string]bool)
			for _, ip := range serv.ServersTCP {
				pack.Head.ServerInfo.Servers[ip] = false
			}
			data, _ := json.Marshal(pack)
			_, err = conn.Write(data)
			if err != nil {
				panic(err)
            }
		}(server)
	}
}

//ConnAnotherServer - function for messaging between servers
func (serv *Server) ConnAnotherServer(to string, msg *pack.PackageTCP) {
	if ip, ok := serv.ServersTCP[to]; ok {
        conn, err := net.Dial("tcp", ip)
        if err != nil {
            delete(serv.ServersTCP, to)
        } else {
            defer conn.Close()
            msg.Head.To = []string{to}
            data, _ := json.Marshal(msg)
            conn.Write(data)
        }
	}
}

//NewUser - function for registering users on all servers
func (serv *Server) NewUser(user *pack.PackageTCP) {
	var pack = &pack.PackageTCP{
		Head: pack.Head{
			From: serv.Addr,
			UserMod: pack.User{
				NewUser: true,
				User:    user.Head.UserMod.User,
			},
		},
	}
	data, _ := json.Marshal(pack)
	for _, ip := range serv.ServersTCP {
		conn, _ := net.Dial("tcp", ip)
		conn.Write(data)
		conn.Close()
	}
}

//DelUser - function for deleting users on all servers
func (serv *Server) DelUser(user string) {
	var pack = &pack.PackageTCP{
		Head: pack.Head{
			UserMod: pack.User{
				DelUser: true,
				User:    user,
			},
		},
	}
	data, _ := json.Marshal(pack)
	for _, ip := range serv.ServersTCP {
		conn, _ := net.Dial("tcp", ip)
		conn.Write(data)
		conn.Close()
	}
}

