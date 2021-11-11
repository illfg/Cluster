package main

import (
	"Cluster/context"
	"Cluster/network"
	"fmt"
	"strconv"
	"time"
)

type testHandler struct {
}

func (receiver testHandler) PutEvent(event network.Event) {
	event.Content = ""
	network.Send(event.From, event, nil)
}
func (receiver testHandler) DoProcess() {

}

func main() {
	//flag.Parse()
	context.Run()
	for {
		time.Sleep(time.Hour)
	}
	client()
	time.Sleep(time.Hour)
}

func client() {
	network.InitClient("127.0.0.1:8888")
	i := 1
	for {
		syncClient := network.SyncClient{}
		sync := syncClient.SendSync("127.0.0.1:7777",
			network.Event{
				EType: network.HEARTBEAT,
			})
		fmt.Println(sync)
		fmt.Println(fmt.Sprintf("Time: %s\n", strconv.Itoa(i)))
		time.Sleep(time.Second)
		i++
	}
}
func serve() {
	network.RegisterHandler("hhhh", testHandler{})
	network.InitClient("127.0.0.1:7777")
}
