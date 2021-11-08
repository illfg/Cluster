package main

import (
	"Cluster/network"
	"flag"
	"fmt"
	"time"
)

type testCall struct {
}

func (c *testCall) Callback(event network.Event) {
	fmt.Println("yes i receive response")
	fmt.Println(event)
}

type testHandler struct {
}

func (receiver *testHandler) DoProcess() {
	fmt.Println("init init init")
}
func (receiver *testHandler) PutEvent(event network.Event) {
	fmt.Println("processing event")
	fmt.Println(event)
	network.Send(event.From, event, nil)
}
func main() {
	flag.Parse()
	client()
	for {
		time.Sleep(time.Hour)
	}
}
func server() {
	network.InitClient("127.0.0.1:9998")
	handler := testHandler{}
	network.RegisterHandler("hhhhh", &handler)
}
func client() {
	network.InitClient("127.0.0.1:9999")
	call := testCall{}
	network.Send("127.0.0.1:9998", network.Event{
		EType:   network.State(1),
		UUID:    "123456789",
		Content: "1heartbeat heartbeat heartbeat\n",
	}, &call)
}
