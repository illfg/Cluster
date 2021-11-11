package network

import (
	"github.com/golang/glog"
	"net"
	"sync"
)

//connectors 维护所有的connector，负责connector的创建
var (
	connectors map[string]*defaultConnector
	lockConn   = sync.Mutex{}
)

const containerLogFlag = "container"

//initContainerAndListen 初始化并监听端口
func initContainerAndListen(IPPort string) {
	connectors = make(map[string]*defaultConnector, 8)
	go listen(IPPort)
}

//Send 发送事件
func send(IPPort string, event Event) bool {

	connector := getConn(IPPort)
	if connector == nil {
		addConnector(IPPort, newConnectorWithIPPort(IPPort))
		if getConn(IPPort) == nil {
			return false
		}
		connector = getConn(IPPort)
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
	glog.Infof("[%s]: listen on %s\n", containerLogFlag, IPPort)
	for {
		conn, err := listen.Accept()
		if err != nil {
			glog.Errorf("[%s]:Accept() failed, err is %s\n", containerLogFlag, err)
			continue
		}
		glog.Infof("[%s]:accept connect[%s]", containerLogFlag, conn.RemoteAddr().String())
		addConnector(conn.RemoteAddr().String(), newConnectorWithConn(conn))
	}
}

func getConn(IPPort string) *defaultConnector {
	lockConn.Lock()
	defer lockConn.Unlock()
	return connectors[IPPort]
}
func addConnector(IPPort string, connector *defaultConnector) {
	lockConn.Lock()
	defer lockConn.Unlock()
	connectors[IPPort] = connector
}
func deleteFaultyConn(IPPort string) {
	lockConn.Lock()
	defer lockConn.Unlock()
	delete(connectors, IPPort)
}
