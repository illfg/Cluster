package handler

import (
	"Cluster/context"
)

type EventHandler interface {
	CheckEvent(event Event) bool
	PutEvent(event Event)
	DoProcess(ctx *context.Context)
}
type State int

type Event struct {
	EType   State
	UUID    string
	From    string
	Content string
}

const (
	HEARTBEAT State = iota
)

var handlers = make(map[string]EventHandler, 8)

func Start() {
	for _, handler := range handlers {
		go handler.DoProcess(context.GetContext())
	}
}

func RegisterHandler(indicator string, handler EventHandler) {
	handlers[indicator] = handler
}

func DoDistribute(event Event) {
	for _, handler := range handlers {
		if handler.CheckEvent(event) {
			handler.PutEvent(event)
		}
	}
}
