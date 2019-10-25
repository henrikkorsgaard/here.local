package context

import (
	"crypto/rand"
	"crypto/tls"
	"fmt"
	"log"
	"net/rpc"
	"time"

	"github.com/henrikkorsgaard/here.local/device"
	"github.com/henrikkorsgaard/here.local/initialise/configs"
)

var (
	deviceCache   = cache.New(cache.NoExpiration, cache.NoExpiration)
	locationCache = cache.New(10*time.Minute, 30*time.Minute) // prolly should only rarely expire!
)

type Reply struct {
	Message string
	Peers   []string
}

type ContextServer struct {
}

// Possible solution https://gist.github.com/ncw/9253562
func Run(addr string, tls.Config) {
	fmt.Println("Running context server")

	locationCache.onEvicted()

	server := new(ContextServer)
	rpc.Register(server)
	
	config.Rand = rand.Reader
	service := addr
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

func (c *ContextServer) SendProximityData(d interface{}, r *Reply) error {
	fmt.Println(d)
	return nil
}

func (c *ContextServer) ConnectLocation(l interface{}, r *Reply) error {
	fmt.Println(l)
	return nil
}

func locationEvicted(mac string, l interface{}) {
	location := l(*device.Location)
	fmt.Printf("Location evicted %+v\n", location)
}







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
