package main

import (
	"Cluster/context"
	"flag"
)

func main() {
	flag.Parse()
	context.Run()
	select {}
}
