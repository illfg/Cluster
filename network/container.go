package network

import (
	"github.com/golang/glog"
	"net"
)

//connectors 维护所有的connector，负责connector的创建
var connectors map[string]*defaultConnector

const containerLogFlag = "container"

//initContainerAndListen 初始化并监听端口
func initContainerAndListen(IPPort string) {
	connectors = make(map[string]*defaultConnector, 8)
	go listen(IPPort)
}

//Send 发送事件
func send(IPPort string, event Event) bool {
	connector := connectors[IPPort]
	if connector == nil {
		connectors[IPPort] = newConnectorWithIPPort(IPPort)
		if connectors[IPPort] == nil {
			return false
		}
		connector = connectors[IPPort]
	}
	return connector.Send(event)
}

//listen 监听端口事件
func listen(IPPort string) {
	listen, err := net.Listen("tcp", IPPort)
	if err != nil {
		glog.Fatalf("[%s]:Listen failed, err is %s\n", containerLogFlag, err)
		return
	}
	for {
		conn, err := listen.Accept()
		if err != nil {
			glog.Errorf("[%s]:Accept() failed, err is %s\n", containerLogFlag, err)
			continue
		}
		glog.Infof("[%s]:accept connect[%s]", containerLogFlag, conn.RemoteAddr().String())
		connectors[conn.RemoteAddr().String()] = newConnectorWithConn(conn)
	}
}
