package tcpconn

import (
    "bufio"
	"encoding/json"
	"fmt"
	"net"
)

//ReciveConn - receive connection and reads the packets 
func ReciveConn(conn net.Conn, msgs chan PackageTCP, dconns chan net.Conn, aconns map[net.Conn]string) {
    rd := bufio.NewReader(conn)
	fmt.Println(conn.RemoteAddr().String())
	for {
		var (
			buffer  = make([]byte, 1024)
			message string
			pack    PackageTCP
		)
		length, err := rd.Read(buffer)
		if err != nil {
			fmt.Println(err)
			break
		}
		message += string(buffer[:length])
		err = json.Unmarshal([]byte(message), &pack)
		if err != nil {
			aconns[conn] = message
		}
		msgs <- pack
				}
	dconns <- conn
}

//RedirectPackages - redirects packets to recipient
func RedirectPackages(msg PackageTCP, aconns map[net.Conn]string) {
    for conn, name := range aconns {
        for _, to := range msg.To{
            if name == string(to) {
                data, err := json.Marshal(msg)
                if err != nil {
                    panic(err)
                }
                conn.Write([]byte(data))
            }
        }
    }
}