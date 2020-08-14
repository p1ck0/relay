package tcpconn

var buff = 1024

//PackageTCP - tcp package for processing
type PackageTCP struct {
	Server  bool
	Servers []string
	Conns   []string
	TCPport string
	From    string
	To      []string
	Body    string
}
