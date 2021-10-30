package tcpconn

import (
	"fmt"
	"net"
)

type TCPClient struct {
	conn     net.Conn
	protocol TCPProtocol
	decoder  defaultDecoder
}

func newTCPClientWithConn(conn net.Conn) *TCPClient {
	protocol := TCPProtocol{}
	client := &TCPClient{
		conn:     conn,
		protocol: protocol,
		decoder: defaultDecoder{
			protocol: protocol,
		},
	}
	//监听该连接并解析数据
	go client.ListenAndParsePkg()
	return client
}
func newTCPClientWithIP(IPPort string) *TCPClient {
	conn, err := net.Dial("tcp", IPPort)
	if err != nil {
		fmt.Println("err : ", err)
		return nil
	}
	return newTCPClientWithConn(conn)
}

func (t *TCPClient) Send(event Event) {
	data := t.protocol.Make(event)
	_, err := t.conn.Write([]byte(data)) // 发送数据
	if err != nil {
		return
	}
}

func (t *TCPClient) ListenAndParsePkg() {
	t.decoder.createEventFromConn(t.conn)
}

func (t *TCPClient) Close() {
	t.conn.Close()
}
