package main

import (
	"github.com/henrikkorsgaard/here.local/configuration"
)

func main() {
	//TODO:
	//Decide on configs
	configuration.Init()

	//go context.Run() //now we have the rpc server running?
	//go proximity.Run()
	//server.Run()
	// We will fix the config at a later point!
	//context.Run()
	//Run context server
	//RUn api server (do this last)
	//Run simulated scanner

}
