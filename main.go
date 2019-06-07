package main

import (
	"time"

	"github.com/henrikkorsgaard/here.local/server/context"
	"github.com/henrikkorsgaard/here.local/simulation"
)

func main() {
	go simulation.Run()
	time.Sleep(20 * time.Second)
	context.Run()
	//create a
	//server.Run()
}
