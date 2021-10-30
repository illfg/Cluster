package tcpconn

import (
	"fmt"
	"net"
)

//负责该链接的接收与发送，具体的操作由client完成
type defaultConnector struct {
	IPPort string
	Queue  chan Event
	client *TCPClient
}

const (
	CHANNEL_SIZE = 100
)

func NewConnectorWithConn(conn net.Conn) *defaultConnector {
	return createInstance(newTCPClientWithConn(conn))
}
func NewConnectorWithIPPort(IPPort string) *defaultConnector {
	return createInstance(newTCPClientWithIP(IPPort))
}
func createInstance(client *TCPClient) *defaultConnector {
	connector := &defaultConnector{
		IPPort: client.conn.RemoteAddr().String(),
		Queue:  make(chan Event, CHANNEL_SIZE),
		client: client,
	}
	go connector.process()
	return connector
}

//将待发送的事件插入队列中
func (c *defaultConnector) Send(event Event) bool {
	select {
	case c.Queue <- event:
		return true
	default:
		return false
	}
}

//处理发送事件
func (c *defaultConnector) process() {
	for {
		event,err := <-c.Queue
		if !err{
			fmt.Println(c.IPPort+"chanel closed")
			break
		}
		fmt.Println("sending")
		c.client.Send(event)
	}
}

func (c *defaultConnector) Close() {
	close(c.Queue)
	c.client.Close()
}
