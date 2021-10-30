package tcpconn

import (
	"fmt"
	"net"
)

type TCPServer struct{}

func (t *TCPServer) Listen(IPPort string) {
	listen, err := net.Listen("tcp", IPPort)
	if err != nil {
		fmt.Println("Listen() failed, err: ", err)
		return
	}
	for {
		conn, err := listen.Accept()
		if err != nil {
			fmt.Println("Accept() failed, err: ", err)
			continue
		}

		NewConnectorWithConn(conn)
	}
}
