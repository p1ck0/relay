package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"strings"
	"sync"
)

var (
	BUFF = 1024

	wg sync.WaitGroup

	buffer = make([]byte, BUFF)
)

//PackageTCP - tcp package for processing
type PackageTCP struct {
	Head 	Head
	Body    interface{}
}

//Head - struct for PackageTCP
type Head struct {
	From string
	To []string
	UserMod User
	ServerInfo Server
}

//User - struct for PackageTCP
type User struct {
	RegUser bool
	DelUser bool
	NewUser bool
	User string
}

//Server - struct for PackageTCP
type Server struct {
	Servers map[string]bool
	Conns []string
	TCPport string
	Server bool
}

func main() {
	wg.Add(1)
	myaddr := "127.0.0.1:6667"

	tcpAddr, err := net.ResolveTCPAddr("tcp", myaddr)
	if err != nil {
		panic(err)
	}
	servAddr, err := net.ResolveTCPAddr("tcp", "127.0.0.1:8888")
	if err != nil {
		panic(err)
	}
	conn, err := net.DialTCP("tcp", tcpAddr, servAddr)
	if err != nil {
		panic(err)
	}
	var pack = PackageTCP{
		Head: Head{
			UserMod : User{
				RegUser : true,
				User : "vasya",
			},
		},
	}
	data, _ := json.Marshal(pack)
	_, err = conn.Write(data)
	if err != nil {
		panic(err)
	}
	go Read(conn)
	go Write(conn)

	wg.Wait()

}

func Read(conn net.Conn) {
	var recPack PackageTCP
	for {
		len, err := conn.Read(buffer)
		if err != nil {
			fmt.Printf("dissconect")
			wg.Done()
			return
		}
		json.Unmarshal(buffer[:len], &recPack)
		fmt.Println(recPack.Head.From, ":", recPack.Body)
	}
}

func Write(conn net.Conn) {
	myaddr := "vasya"
	anotheraddr := []string{"petya","nik"}
	reader := bufio.NewReader(os.Stdin)
	for {
		message, _ := reader.ReadString('\n')
		message = strings.ReplaceAll(message, "\n", "")
		var pack = PackageTCP{
			Head : Head {
				From: myaddr,
				To:   anotheraddr,
			},
			Body : message,
		}
		data, _ := json.Marshal(pack)
		_, err := conn.Write(data)
		if err != nil {
			panic(err)
		}
	}
}
