package server

import (
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"log"
	"net"
	"net/rpc"
)

var (
	peers = []string{"1"}
	crt   = "../crt/here.local.server.crt"
	key   = "../crt/here.local.server.key"
)

type Reply struct {
	Message string
	Peers   []string
}

type ContextServer struct {
}

// Possible solution https://gist.github.com/ncw/9253562
func runContextServer(addr string) {

	cert, err := tls.LoadX509KeyPair("crts/localhost.crt", "crts/localhost.key")
	if err != nil {
		log.Fatalf("server: loadkeys: %s", err)
	}

	certCA, err := tls.LoadX509KeyPair("crts/here.local.crt", "crts/here.local.key")
	if err != nil {
		log.Fatalf("server: loadCA: %s", err)
	}

	ca, err := x509.ParseCertificate(certCA.Certificate[0])
	if err != nil {
		log.Fatal(err)
	}

	certPool := x509.NewCertPool()
	certPool.AddCert(ca)
	config := tls.Config{
		Certificates: []tls.Certificate{cert},
		ClientAuth:   tls.RequireAndVerifyClientCert,
		ClientCAs:    certPool,
	}

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

		go handleClient(conn)
	}
}

func (c *ContextServer) Echo(msg string, reply *Reply) error {
	reply.Message = msg
	reply.Peers = peers
	return nil
}

func (c *ContextServer) Echo(msg string, reply *Reply) error {
	reply.Message = msg
	reply.Peers = peers
	return nil
}

func handleClient(conn net.Conn) {
	defer conn.Close()
	rpc.ServeConn(conn)
	log.Println("server: conn: closed")
}

/*
//Hello returns "World"


func (c *ContextServer) SendProximityData(rd device.Device, r *Reply) error {

	//UpsertRawDevice(rd)
	fmt.Println(rd)
	return nil
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
