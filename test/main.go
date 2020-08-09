package main

import (
	"encoding/json"
	"bufio"
	"fmt"
	"log"
	"net"
	"sort"
	"strings"
	"os"
)

var BUFF = 1024

const (
	connect = "conn"
	disconnect = "disconn"
)

type PackageTCP struct{
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
}

var (
	aconns = make(map[string]net.Conn)
    tcpconns  = make(chan net.Conn)
	udpconns = make(chan net.UDPConn)
	servers = make(map[string][]string)
	dconns = make(chan net.Conn)
	msgs   = make(chan PackageTCP)
	command = make(chan string)
)

func comms() {
	com, _ := bufio.NewReader(os.Stdin).ReadString('\n')
	command <- com
}


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
            udpconns <- *lnudp
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
			fmt.Println(conn.RemoteAddr().String())
			
			go func(conn net.Conn) {
				rd := bufio.NewReader(conn)
				fmt.Println(conn.RemoteAddr().String())
				for {
					var (
						buffer = make([]byte, BUFF)
						message string
						pack PackageTCP
					)
                    length, err := rd.Read(buffer)
                    if err != nil { 
						fmt.Println(err)
						break 
					}
					message += string(buffer[:length])
					err  = json.Unmarshal([]byte(message), &pack)
					if err != nil {
						aconns[message] = conn
					}
					msgs <- pack
                }
				dconns <- conn
            }(conn)

			case udpconn := <-udpconns:
                go func(udpconn net.UDPConn) {
                    for {
						var (
							buffer = make([]byte, BUFF)
							message string
							pack PackageUDP
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
						err = json.Unmarshal([]byte(message), &pack)
						if err != nil {
							json.Unmarshal([]byte(message), &pack2)
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
						} else {
							pack.Status = true
							data, _ := json.Marshal(pack)
							go func(udpconn net.UDPConn, remoteaddr *net.UDPAddr, data []byte) {
								_,err := udpconn.WriteToUDP(data, remoteaddr)
								if err != nil {
									fmt.Printf("Couldn't send response %v", err)
								}
							}(udpconn, remoteaddr, data)
						}
                    }
                }(udpconn)

			case msg := <-msgs:
				fmt.Println(msg)
				for _, to := range msg.To {
					go func(to string, msg PackageTCP){
						conn, ok := aconns[to]
						if !ok {
							//connAnotherServer(to, msg)
							return
						}
						data, err := json.Marshal(msg)
						if err != nil {
							panic(err)
						}
						conn.Write([]byte(data))
					}(to, msg)
				}
				
		case com := <-command:
			splited := strings.Split(com, " ")
			switch splited[0] {
			case connect:
				connectServer(splited[1:])
//			case disconnect:
//				disconnectServer(splited[1:])
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

func connectServer(servers []string) {
	for _, conns := range servers {
		go func(conns string){
			udpservadr,_ := net.ResolveUDPAddr("udp", conns)
			var pack PackageServerUDP
			data,_ := json.Marshal(pack)
			udpconn,_ := net.DialUDP("udp", nil, udpservadr)
			_,err := udpconn.WriteToUDP(data, udpservadr)
			if err != nil {
				fmt.Printf("Couldn't send response %v", err)
			}
		}(conns)
	}
}

func connAnotherServer(name string, msg PackageTCP) {
	for ip, user := range servers {
		index := sort.SearchStrings(user, name)
		if index != len(user) {
			conn, _ := net.Dial("tcp", ip)
			defer conn.Close()
			data,_ := json.Marshal(msg)
			conn.Write(data)
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