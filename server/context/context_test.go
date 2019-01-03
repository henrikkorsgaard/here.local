package context

import (
	"crypto/rsa"
	"log"
	"net/rpc"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	client   *rpc.Client
	testSalt string
	testKey  rsa.PublicKey
)

func init() {
	go Run()
	var err error
	client, err = rpc.DialHTTP("tcp", "localhost:1339")
	if err != nil {
		log.Fatal("Error setting up RPC connection: ", err)
	}
}

func TestRPCEcho(t *testing.T) {
	var actualResult string
	msg := "Foo"
	err := client.Call("Context.Echo", msg, &actualResult)

	if err != nil {
		log.Fatal("ContextServer error:", err)
	}

	if actualResult != msg {
		t.Fatalf("Expected 'World' but got %s", actualResult)
	}
}

func TestRPCPublicKey(t *testing.T) {
	var actualResult rsa.PublicKey
	err := client.Call("Context.GetPublicKey", "", &actualResult)
	if err != nil {
		log.Fatal("ContextServer error:", err)
	}

	assert.NotNil(t, actualResult)
	testKey = actualResult

}

func TestRPCSalt(t *testing.T) {

	var actualResult string
	err := client.Call("Context.GetSalt", testKey, &actualResult)
	if err != nil {
		t.Logf("Failed test with error: %s", err)
		t.Fail()
	} else {
		//assert.NotNil(t, actualResult)
	}

}

func TestRPCEncryption(t *testing.T) {

}
