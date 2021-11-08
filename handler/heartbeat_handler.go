package handler

//
//import (
//	"Cluster/context"
//	"Cluster/network"
//	"time"
//)
//
//const HEARTBEAT_INTEVAL = time.Second
//
//type heartbeatHandler struct {
//	list          chan Event
//	receiveEnable bool
//}
//
//func (h heartbeatHandler) CheckEvent(event Event) bool {
//	if event.EType == HEARTBEAT && h.receiveEnable {
//		return true
//	}
//	return false
//}
//func (h heartbeatHandler) PutEvent(event Event) {
//	select {
//	case h.list <- event:
//	default:
//	}
//
//}
//func (h heartbeatHandler) DoProcess(ctx *context.Context) {
//	if ctx.Field["isMaster"].(bool) {
//		h.receiveHeartbeat(ctx)
//	} else {
//		h.reportHeartbeat(ctx)
//	}
//}
//
//func (h heartbeatHandler) receiveHeartbeat(ctx *context.Context) {
//	availMap := ctx.Field["available"].(*map[string]time.Time)
//	for {
//		event := <-h.list
//		(*availMap)[event.From] = time.Now()
//	}
//}
//func (h heartbeatHandler) reportHeartbeat(ctx *context.Context) {
//	for {
//		master := ctx.Field["Master"].(network.defaultConnector)
//		event := Event{
//			EType:   HEARTBEAT,
//			UUID:    "",
//			From:    ctx.Field["IPPort"].(string),
//			Content: "",
//		}
//		master.Send(event)
//		time.Sleep(HEARTBEAT_INTEVAL)
//	}
//}
