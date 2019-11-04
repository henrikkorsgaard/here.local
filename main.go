package main

import (
	"github.com/henrikkorsgaard/here.local/proximity"
	"github.com/henrikkorsgaard/here.local/server/context"
)

func main() {
	//TODO:
	//Decide on configs

	go context.Run() //now we have the rpc server running?
	go proximity.Run()

	for {
	}
	//server.Run()
	// We will fix the config at a later point!
	//context.Run()
	//Run context server
	//RUn api server (do this last)
	//Run simulated scanner

}
