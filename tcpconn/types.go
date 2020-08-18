package tcpconn

var buff = 1024

//PackageTCP - tcp package for processing
type PackageTCP struct {
    User    string
	Server  bool
	DelUser bool
	NewUser bool
	Servers map[string]bool
	Conns   []string
	TCPport string
	From    string
	To      []string
	Body    string
}
