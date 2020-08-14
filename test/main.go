package main

import (
	"encoding/json"
	"bufio"
	"fmt"
	"log"
	"net"
	"sort"
	"io"
	//"strings"
	//"os"
)

var BUFF = 1024

const (
	connect = "conn"
	disconnect = "disconn"
)

type PackageTCP struct{
	Server bool
	Servers []string
	Conns []string
	TCPport string
    From string
    To []string
    Body string
}

type PackageUDP struct {
	Status bool
}

type PackageServerUDP struct {
	Status bool
	TCPport string
	Conns []string
	MyConns []string
	MyTCP string
}

var (
	aconns = make(map[string]net.Conn)
    tcpconns  = make(chan net.Conn)
	udpconns = make(chan net.UDPConn)
	servers = make(map[string][]string)
	dconns = make(chan net.Conn)
	msgs   = make(chan PackageTCP)
)


func main() {
	ln, err := net.Listen("tcp", "127.0.0.1:8081")
	if err != nil {
		log.Fatalln(err.Error())
    }
    udpaddr := net.UDPAddr{
        Port: 8082,
        IP: net.ParseIP("127.0.0.1"),
    }
    lnudp, err := net.ListenUDP("udp", &udpaddr)
    if err != nil {
        log.Fatalln(err.Error())
    }
    defer func() {
		ln.Close()
		lnudp.Close()
	}()
	go func() {
		for {
			//go comms()
            //udpconns <- *lnudp
			conn, err := ln.Accept()
			if err != nil {
				log.Fatalln(err.Error())
			}
			tcpconns <- conn
		}
	}()

	for {
		select {
		case conn := <-tcpconns:
			go func(conn net.Conn) {
				rd := bufio.NewReader(conn)
				for {
					var (
						buffer = make([]byte, BUFF)
						message string
						pack PackageTCP
					)
                    length, err := rd.Read(buffer)
                    if err != nil { 
						fmt.Println("ошибка на 100 строке",err)
						break 
					}
					message += string(buffer[:length])
					err  = json.Unmarshal([]byte(message), &pack)
					if err != nil && err != io.EOF {
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
            }(conn)
/*
			case udpconn := <-udpconns:
                go func(udpconn net.UDPConn) {
                    for {
						var (
							buffer = make([]byte, BUFF)
							message string
							pack2 PackageServerUDP
						)
						length, err := udpconn.Read(buffer)
						if err != nil { 
							fmt.Println(err)
							break 
						}
						_,remoteaddr,err := udpconn.ReadFromUDP(buffer)
						fmt.Printf("Read a message from %v %s \n", remoteaddr, buffer)
						if err !=  nil {
							fmt.Printf("Some error  %v", err)
							continue
						}
						message += string(buffer[:length])
						fmt.Println(message)
						err = json.Unmarshal([]byte(message), &pack2)
						if err != nil {
							fmt.Println(err)
						}
						if pack2.Status == true {
							servers[pack2.TCPport] = pack2.Conns
							fmt.Println(servers)
						} else {
							for name := range aconns {
								pack2.Conns = append(pack2.Conns, name)
							}
							pack2.TCPport = ln.Addr().String()
							pack2.Status = true
							data, _ := json.Marshal(pack2)
							go func(udpconn net.UDPConn, remoteaddr *net.UDPAddr, data []byte) {
								_,err := udpconn.WriteToUDP(data, remoteaddr)
								if err != nil {
									fmt.Printf("Couldn't send response %v", err)
								}
							}(udpconn, remoteaddr, data)
							servers[pack2.MyTCP] = pack2.MyConns
						}
                    }
                }(udpconn)
*/
			case msg := <-msgs:
				fmt.Println(msg)
				for _, to := range msg.To {
					go func(to string, msg PackageTCP){
						conn, ok := aconns[to]
						if !ok {
							fmt.Println(to)
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
				

				
		case dconn := <-dconns:
			defer dconn.Close()
			for name, conn := range aconns {
				if conn == dconn {
					log.Printf("Client %v was gone\n", name)
					dconn.Close()
					delete(aconns, name)
				}
			}
		}
	}
}

func connectServer(serverss []string) {
	myaddr := "127.0.0.1:8081"
	for _, serv := range serverss {
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
			for ip := range servers {
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

func connAnotherServer(to string, msg PackageTCP) {
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







/*
func connServer(ipServ string) {
	conn, _ := net.Dial("udp", ipServ)
	ls, _ := net.Listen("udp", "127.0.0.1:8082")
	conn, _ = ls.Accept()
	servers<-conn
}
*/