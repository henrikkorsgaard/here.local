package main

import (
	"github.com/henrikkorsgaard/here.local/configuration"
	"github.com/henrikkorsgaard/here.local/proximity"
	"github.com/henrikkorsgaard/here.local/server"
)

func main() {
	configuration.Setup()
	server.Run()
	proximity.Run()
}
