package context

//
//import "Cluster/network"
//
//type Context struct {
//	//操作数据要加锁
//	Field map[string]interface{}
//	//ctx中使用携程检查数据，例如同步，主从
//	//handler中完成数据的更新
//	connectors map[string]network.defaultConnector
//}
//
//var ctxIns *Context
//
//const ListenAddr = "127.0.0.1:9985"
//
//func Init() {
//	ctxIns = &Context{
//		Field:      make(map[string]interface{}, 8),
//		connectors: make(map[string]network.defaultConnector, 8),
//	}
//	//解析配置???选举
//	if ctxIns.Field["isMaster"].(bool) {
//		network.Listen(ListenAddr)
//	}
//
//}
//func GetContext() *Context {
//	if ctxIns == nil {
//		Init()
//	}
//	return ctxIns
//}
