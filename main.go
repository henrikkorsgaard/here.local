package main

import (
	"github.com/henrikkorsgaard/here.local/configuration"
	"github.com/henrikkorsgaard/here.local/proximity"
	"github.com/henrikkorsgaard/here.local/server/contextserver"
)

func main() {
	//TODO:
	//Decide on configs
	configuration.Init()

	go contextserver.Run() //now we have the rpc server running?
	proximity.Run()
	// We will fix the config at a later point!
	//context.Run()
	//Run context server
	//RUn api server (do this last)
	//Run simulated scanner

}
