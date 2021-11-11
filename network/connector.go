package network

import (
	"github.com/golang/glog"
	"net"
)

//defaultConnector 负责该链接的接收与发送，提供发送缓存，具体的操作由client完成
type defaultConnector struct {
	IPPort   string
	Queue    chan Event
	conn     net.Conn
	protocol protocol
	decoder  decoder
}

const (
	CHANNEL_SIZE     = 100
	connectorLogFlag = "connector"
)

//NewConnectorWithConn 通过net.Conn创建Connector
func newConnectorWithConn(conn net.Conn) *defaultConnector {
	if conn == nil {
		return nil
	}
	connector := &defaultConnector{
		IPPort:   conn.RemoteAddr().String(),
		Queue:    make(chan Event, CHANNEL_SIZE),
		conn:     conn,
		protocol: &defaultProtocol{},
		decoder: &defaultDecoder{
			protocol: &defaultProtocol{},
			IPPort:   conn.RemoteAddr().String(),
		},
	}
	//处理发送事件
	go connector.process()
	//监听该连接并解析数据
	go connector.listenAndParsePkg()
	return connector
}

//NewConnectorWithIPPort 通过IPPort创建Connector
func newConnectorWithIPPort(IPPort string) *defaultConnector {
	conn, err := net.Dial("tcp", IPPort)
	if err != nil {
		glog.Errorf("[%s]:established connection failed,err is %s", connectorLogFlag, err.Error())
		return nil
	}
	return newConnectorWithConn(conn)
}

//Send 将待发送的事件插入队列中,由process协程完成。
func (c *defaultConnector) Send(event Event) bool {
	select {
	case c.Queue <- event:
		return true
	default:
		return false
	}
}

//process 处理发送事件
func (c *defaultConnector) process() {
	for {
		event, err := <-c.Queue
		if !err {
			glog.Errorf("[%s]: send channel have been close\n", connectorLogFlag)
			deleteFaultyConn(c.IPPort)
			break
		}
		glog.Infof("[%s]: sending event[%v]\n", connectorLogFlag, event)
		data := c.protocol.make(event)
		_, sendErr := c.conn.Write([]byte(data)) // 发送数据
		if sendErr != nil {
			glog.Errorf("[%s]: detect err when writing data, err is %s\n", connectorLogFlag, sendErr)
		}
	}
}

//Close 关闭链接并释放资源
func (c *defaultConnector) Close() {
	close(c.Queue)
	c.conn.Close()
}

//listenAndParsePkg 监听并解析接收到的数据
func (c *defaultConnector) listenAndParsePkg() {
	c.decoder.createEventFromConn(c.conn)
}
