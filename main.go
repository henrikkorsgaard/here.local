package main

import (
	"fmt"

	"github.com/henrikkorsgaard/here.local/server"
)

func main() {
	fmt.Println("here")
	//configuration.Setup()
	server.Run()
}
