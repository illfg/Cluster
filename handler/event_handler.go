package handler

//
//import (
//	"Cluster/context"
//)
//
//
//
//
//
//const (
//	HEARTBEAT State = iota
//)
//
//var handlers = make(map[string]EventHandler, 8)
//
//func Start() {
//	for _, handler := range handlers {
//		go handler.DoProcess(context.GetContext())
//	}
//}
//
//func RegisterHandler(indicator string, handler EventHandler) {
//	handlers[indicator] = handler
//}
//
//func DoDistribute(event Event) {
//	for _, handler := range handlers {
//		if handler.CheckEvent(event) {
//			handler.PutEvent(event)
//		}
//	}
//}
