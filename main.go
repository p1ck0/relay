package main

import (
	"fmt"
	tcp "github.com/p1ck0/relay/tcpconn"
	"github.com/urfave/cli/v2"
	"github.com/fatih/color"
	"log"
	"net"
	"os"
	"strings"
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
	servers     string
)

func init() {
	color.Cyan(`
#######################################


   _____      __    __  __
  / ___/___  / /___ \ \/ /
  \__ \/ _ \/ / __ '/\  / 
 ___/ /  __/ / /_/ / / /  
/____/\___/_/\__,_/ /_/ 
	  

######################################
	  `)
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
			&cli.StringFlag{
				Name:        "conn",
				Value:       "",
				Usage:       "the port on which the server will run",
				Aliases:     []string{"c"},
				Destination: &servers,
			},
		},
		Action: func(c *cli.Context) error {
			if len(port) > 0 {
				color.Green("*** uses port " + port + " ***")
			} else {
				color.Green("*** uses port " + port + " ***")
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
	if len(servers) > 0 {
		serversArr := strings.Split(servers, " ")
		tcp.ConnectServer(serversArr, addr, aconns, serverstcp)
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
			go tcp.ReciveConn(conn, msgs, dconns, aconns, serverstcp, addr)

		case msg := <-msgs:
			go tcp.RedirectPackages(msg, aconns, serverstcp)

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