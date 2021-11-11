package handler

type syncHandler struct {
}

func (receiver syncHandler) sync() {
	//从节点返回ack时，携带checkpoint，根据checkpoint同步数据

	//分布式数据库，主节点修改，从节点同步
}
