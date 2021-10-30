package tcpconn

import (
	"Cluster/handler"
	"fmt"
	"net"
)

//负责该链接的接收与发送，具体的操作由client完成
type DefaultConnector struct {
	IPPort string
	Queue  chan handler.Event
	client *TCPClient
}

const (
	CHANNEL_SIZE = 100
)

func NewConnectorWithConn(conn net.Conn) *DefaultConnector {
	return createInstance(newTCPClientWithConn(conn))
}
func NewConnectorWithIPPort(IPPort string) *DefaultConnector {
	return createInstance(newTCPClientWithIP(IPPort))
}
func createInstance(client *TCPClient) *DefaultConnector {
	connector := &DefaultConnector{
		IPPort: client.conn.RemoteAddr().String(),
		Queue:  make(chan handler.Event, CHANNEL_SIZE),
		client: client,
	}
	go connector.process()
	return connector
}

//将待发送的事件插入队列中
func (c *DefaultConnector) Send(event handler.Event) bool {
	select {
	case c.Queue <- event:
		return true
	default:
		return false
	}
}

//处理发送事件
func (c *DefaultConnector) process() {
	for {
		event, err := <-c.Queue
		if !err {
			fmt.Println(c.IPPort + "chanel closed")
			break
		}
		fmt.Println("sending")
		c.client.Send(event)
	}
}

func (c *DefaultConnector) Close() {
	close(c.Queue)
	c.client.Close()
}
