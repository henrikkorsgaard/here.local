package context

import (
	"log"
	"net/rpc"
	"testing"

	"../proximity/proximity"
)

var (
	client   *rpc.Client
	testSalt string
)

func init() {
	go Run()
	var err error
	client, err = rpc.DialHTTP("tcp", "localhost:1339")
	if err != nil {
		log.Fatal("Error setting up RPC connection: ", err)
	}
}

func TestEcho(t *testing.T) {
	var result Reply

	msg := "Foo"
	err := client.Call("ContextServer.Echo", msg, &result)

	if err != nil {
		t.Logf("Failed test with error: %s", err)
		t.Fail()
	}

	if result.Message != msg {
		t.Fatalf("Expected 'Foo' but got %s", result.Message)
	}
}

func TestSendProximityData(t *testing.T) {

	var result Reply

	rd := proximity.RawDevice{Mac: "bc:ee:7b:c4:ab:b8", Signal: -100}

	err := client.Call("ContextServer.SendProximityData", rd, &result)
	if err != nil {
		t.Logf("Failed test with error: %s", err)
		t.Fail()
	}
}
