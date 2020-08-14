package main

import (
	//"bufio"
	//"encoding/json"
	"fmt"
	tcp "github.com/p1ck0/relay/tcpconn"
	"github.com/urfave/cli/v2"
	"log"
	"net"
	"os"
)

var (
	buff        = 1024
	aconns      = make(map[string]net.Conn)
	tcpconns    = make(chan net.Conn)
	udpconns    = make(chan net.UDPConn)
	serverstcp  = make(map[string][]string)
	dconns      = make(chan net.Conn)
	msgs        = make(chan tcp.PackageTCP)
	command     = make(chan string)
	serversconn []string
	port        string
)

func init() {
	port = "8888"
	app := &cli.App{
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "port",
				Value:       "8888",
				Usage:       "the port on which the server will run",
				Aliases:     []string{"p"},
				Destination: &port,
			},
		},
		Action: func(c *cli.Context) error {
			if len(port) > 0 {
				fmt.Println("use port", port)
			} else {
				fmt.Println("use port", port)
			}
			return nil
		},
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	addr := "127.0.0.1:" + port
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalln(err.Error())
	}
	defer ln.Close()
	go func() {
		for {
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
			go tcp.ReciveConn(conn, msgs, dconns, aconns)

		case msg := <-msgs:
			go tcp.RedirectPackages(msg, aconns)

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
