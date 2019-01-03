package main

import (
	"fmt"
	"log"
	"net/rpc"
)

func main() {
	client, err := rpc.DialHTTP("tcp", "localhost:1339")
	if err != nil {
		log.Fatal("dialing:", err)
	}

	err = client.Call("Context.Hello", "Hi", &reply)

	if err != nil {
		log.Fatal("ContextServer error:", err)
	}
	fmt.Println(reply)

}
