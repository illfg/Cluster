package main

import (
	"Cluster/tcpconn"
	"os"
	"time"
)

func main() {
	f, _ := os.OpenFile("C:/Users/Administrator/GolandProjects/Cluster/fmt.log", os.O_WRONLY|os.O_CREATE|os.O_SYNC|os.O_APPEND,0755)
	os.Stdout = f
	os.Stderr = f

	//(&tcpconn.TCPServer{}).Listen("127.0.0.1:3345")
	send()
}
func send() {
	ip := tcpconn.NewConnectorWithIPPort("127.0.0.1:3345")


	ip.Send(tcpconn.Event{Name: "yq1",Content: "123456789"})
	ip.Send(tcpconn.Event{Name: "yq2",Content: "123456789"})
	ip.Send(tcpconn.Event{Name: "yq3",Content: "123456789"})
	ip.Send(tcpconn.Event{Name: "yq4",Content: "123456789"})
	time.Sleep(1000000)
}