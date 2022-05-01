package network

import (
	"Cluster/conf"
	"Cluster/utils"
	"fmt"
	"sync"
	"time"

	"github.com/golang/glog"
)

//为事件的发送与接收提供回调
var (
	processors    = make(map[string]EventPostProcessor, 8)
	eventChanel   = make(chan Event, CHANNEL_SIZE*10)
	handlers      = make(map[string]EventHandler, 8)
	clientLogFlag = "client"
	lockSend      = sync.Mutex{}
	lockMap       = sync.Mutex{}
)

type EventPostProcessor interface {
	Callback(event Event)
}
type EventHandler interface {
	PutEvent(event Event)
	DoProcess()
}
type SyncClient struct {
}
type SyncSendEventPostProcessor struct {
	response chan Event
}

func (receiver *SyncSendEventPostProcessor) Callback(event Event) {
	receiver.response <- event
}

func InitClient(IPPort string) {
	initContainerAndListen(IPPort)
	for name, handler := range handlers {
		glog.Infof("[%s]:init handler[%s]\n", clientLogFlag, name)
		handler.DoProcess()
	}
}
func Send(IPPort string, event Event, processor EventPostProcessor) bool {
	lockSend.Lock()
	defer lockSend.Unlock()
	if IPPort == "" {
		return false
	}
	if event.aUUID == "" {
		event.aUUID = utils.GetRand()
	}
	if processor != nil {
		addCallBack(event.aUUID, processor)
		// glog.Infof("[%s]:event post processor registered[UUID:%s]", clientLogFlag, event.aUUID)
	}
	addFromAddrInContent(&event)
	return send(IPPort, event)
}

func (*SyncClient) SendSync(IPPort string, event Event) Event {
	event.aUUID = utils.GetRand()
	processor := SyncSendEventPostProcessor{response: make(chan Event, 8)}
	Send(IPPort, event, &processor)

	ticker := time.NewTicker(time.Second / 10)
	select {
	case res := <-processor.response:
		return res
	case <-ticker.C:
		delCallBack(event.aUUID)
		return Event{}
	}
}

func RegisterHandler(name string, eventHandler EventHandler) {
	if name != "" && eventHandler != nil {
		handlers[name] = eventHandler
	}
}

func doReceive(event Event) {
	// glog.Infof("[%s]:received event[%v]\n", clientLogFlag, event)
	parseFromAddrInContent(&event)
	processor := getCallback(event.aUUID)
	if processor != nil {
		// glog.Infof("[%s]:invoke callback[UUID:%s]", clientLogFlag, event.aUUID)
		go processor.Callback(event)
		delCallBack(event.aUUID)
	}
	for _, handler := range handlers {
		go handler.PutEvent(event)
	}
}

func addFromAddrInContent(event *Event) {
	ip := (*conf.GetConf())["IPPort"]
	event.Content = event.Content + fmt.Sprintf("LogicFrom: %s\n", ip)
}
func parseFromAddrInContent(event *Event) {
	yaml := conf.ParseYaml(event.Content)
	event.From = yaml["LogicFrom"]
}

func addCallBack(name string, processor EventPostProcessor) {
	lockMap.Lock()
	defer lockMap.Unlock()
	processors[name] = processor
}
func delCallBack(name string) {
	lockMap.Lock()
	defer lockMap.Unlock()
	delete(processors, name)
}
func getCallback(name string) EventPostProcessor {
	return processors[name]
}
