package main

import (
	"Cluster/conf"
	"log"
	"os"
)

func main() {
	content,err := os.ReadFile("conf.yaml")
	if err!= nil {
		log.Fatal(err)
	}
	conf.ParseYaml(string(content))

}