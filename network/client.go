package network

import "github.com/golang/glog"

//为事件的发送与接收提供回调
var (
	processors    = make(map[string]EventPostProcessor, 8)
	eventChanel   = make(chan Event, CHANNEL_SIZE*10)
	handlers      = make(map[string]EventHandler, 8)
	clientLogFlag = "client"
)

type EventPostProcessor interface {
	Callback(event Event)
}
type EventHandler interface {
	PutEvent(event Event)
	DoProcess()
}

func InitClient(IPPort string) {
	initContainerAndListen(IPPort)
	for name, handler := range handlers {
		glog.Infof("[%s]:init handler[%s]\n", clientLogFlag, name)
		handler.DoProcess()
	}
}
func Send(IPPort string, event Event, processor EventPostProcessor) bool {
	if IPPort == "" {
		return false
	}
	if processor != nil {
		processors[event.UUID] = processor
		glog.Infof("[%s]:event post processor registered[UUID:%s]", clientLogFlag, event.UUID)
	}
	return send(IPPort, event)
}

func RegisterHandler(name string, eventHandler EventHandler) {
	if name != "" && eventHandler != nil {
		handlers[name] = eventHandler
	}
}

func doReceive(event Event) {
	glog.Infof("[%s]:received event[%v]\n", clientLogFlag, event)
	processor := processors[event.UUID]
	if processor != nil {
		glog.Infof("[%s]:invoke callback[UUID:%s]", clientLogFlag, event.UUID)
		processor.Callback(event)
	}
	for _, handler := range handlers {
		handler.PutEvent(event)
	}
}
