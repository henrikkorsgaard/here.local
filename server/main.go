package main

import (
	"fmt"

	"./context"
)

func main() {

	context.Run()
	fmt.Println("Servers needed: ")
	fmt.Println("\tExternal API server")
	fmt.Println("\tClient Server")
	fmt.Println("\tContext Server")
	fmt.Println("\twhat else?")

}
