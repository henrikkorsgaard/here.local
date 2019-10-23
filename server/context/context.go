package context

import (
	"crypto/rand"
	"crypto/tls"
	"fmt"
	"log"
	"net/rpc"

	"github.com/henrikkorsgaard/here.local/device"
	"github.com/henrikkorsgaard/here.local/initialise/configs"
)

/*
var (
	peers   = []string{"1"}
	crt     = "../crt/here.local.server.crt"
	key     = "../crt/here.local.server.key"
	Salt    string
	src     rand.Source
	devices map[string]Device
)*/

type Reply struct {
	Message string
	Peers   []string
}

type ContextServer struct {
}

// Possible solution https://gist.github.com/ncw/9253562
func Run() {
	fmt.Println("Running context server")

	server := new(ContextServer)
	rpc.Register(server)
	config, err := configs.GetTLSHostConfig()
	if err != nil {
		log.Fatalf("configuration.GetTLSHostConfig failed: %s", err)
	}
	config.Rand = rand.Reader
	service := configs.CONTEXT_SERVER_ADDR
	listener, err := tls.Listen("tcp", service, &config)
	if err != nil {
		log.Fatalf("server: listen: %s", err)
	}
	log.Print("server: listening")
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("server: accept: %s", err)
			break
		}
		log.Printf("server: accepted from %s", conn.RemoteAddr())

		go rpc.ServeConn(conn)
	}
}

func (c *ContextServer) SendProximityData(rd device.RawDevice, r *Reply) error {
	device.Upsert(rd)
	return nil
}

/*
//Hello returns "World"






//Run starts the context server
func runContextServer() {
	server := new(ContextServer)
	rpc.Register(server)
	rpc.HandleHTTP()
	l, e := net.Listen("tcp", configuration.ContextServer)
	if e != nil {
		log.Fatal("listen error:", e)
	}
	http.Serve(l, nil)
}

/*

Clients need to know the ip and port -- we should use the existing


What do we need:

Security/privacy

func GetPublicKey()

Master/Slave scenarion or how to avoid race conditions:

This means that each node needs to know who is below or above them in terms of numbers?

func AnybodyOutThere()
func GetNodeList()
func ReportNode()



*/
