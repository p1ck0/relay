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

type PackageTCP struct {
	From string
	To   []string
	Body string
}

func main() {
	wg.Add(1)
	myaddr := "127.0.0.1:6667"

	tcpAddr, err := net.ResolveTCPAddr("tcp", myaddr)
	if err != nil {
		panic(err)
	}
	servAddr, err := net.ResolveTCPAddr("tcp", "127.0.0.1:8081")
	if err != nil {
		panic(err)
	}
	conn, err := net.DialTCP("tcp", tcpAddr, servAddr)
	if err != nil {
		panic(err)
	}
	conn.Write([]byte("vasya"))
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
		fmt.Println(recPack.From, ":", recPack.Body)
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
			From: myaddr,
			To:   anotheraddr,
			Body: message,
		}
		data, _ := json.Marshal(pack)
		_, err := conn.Write(data)
		if err != nil {
			panic(err)
		}
	}
}
