package context

import "Cluster/tcpconn"

type Context struct {
	//操作数据要加锁
	Field map[string]interface{}
	//ctx中使用携程检查数据，例如同步，主从
	//handler中完成数据的更新
	connectors map[string]tcpconn.DefaultConnector
}

var ctxIns *Context

const ListenAddr = "127.0.0.1:9985"

func Init() {
	ctxIns = &Context{
		Field:      make(map[string]interface{}, 8),
		connectors: make(map[string]tcpconn.DefaultConnector, 8),
	}
	//解析配置???选举
	if ctxIns.Field["isMaster"].(bool) {
		(&tcpconn.TCPServer{}).Listen(ListenAddr)
	}

}
func GetContext() *Context {
	if ctxIns == nil {
		Init()
	}
	return ctxIns
}
