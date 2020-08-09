package tcpconn

//PackageTCP - tcp package for processing
//From - who the package came from
//To - to whom to send the package
//Body - body package
type PackageTCP struct {
    From string
    To   []string
	Body interface{}
}