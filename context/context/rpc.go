package context

import (
	"log"
	"net"
	"net/http"
	"net/rpc"

	"../../proximity/proximity"
)

var (
	peers = []string{"1"}
)

//https://golang.org/pkg/net/rpc/

//Context is the context server struct used for rpc communication. This could be the server
type ContextServer struct {
	name string
}

type Reply struct {
	Message string
	Peers   []string
}

//Hello returns "World"
func (c *ContextServer) Echo(msg string, reply *Reply) error {

	reply.Message = msg
	reply.Peers = peers

	return nil
}

func (c *ContextServer) SendProximityData(rd proximity.RawDevice, reply *Reply) error {
	UpsertRawDevice(rd)
	reply.Peers = peers
	return nil
}

//Run starts the context server
func Run() {
	server := new(ContextServer)
	rpc.Register(server)
	rpc.HandleHTTP()
	l, e := net.Listen("tcp", ":1339")
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
