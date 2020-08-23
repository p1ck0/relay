package tcpconn

import "sync"

var buff = 1024

var mut sync.Mutex

//PackageTCP - tcp package for processing
type PackageTCP struct {
	Head *Head
	Body interface{}
}

//Head - struct for PackageTCP
type Head struct {
	From       string
	To         []string
	UserMod    *User
	ServerInfo *Server
}

//User - struct for PackageTCP
type User struct {
	RegUser bool
	DelUser bool
	NewUser bool
	User    string
}

//Server - struct for PackageTCP
type Server struct {
	Servers map[string]bool
	Conns   []string
	TCPport string
	Server  bool
}
