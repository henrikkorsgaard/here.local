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

//https://golang.org/pkg/net/rpc/
/*
//Context is the context server struct used for rpc communication. This could be the server
type ContextServer struct {
	srv *grpc.Server
}
/
//https://bbengfort.github.io/programmer/2017/03/03/secure-grpc.html
//https://github.com/bbengfort/sping/blob/master/server.go
func (cs *ContextServer) Serve(addr string) error {
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("could not list on %s: %s", addr, err)
	}

	// Create the TLS credentials
	creds, err := credentials.NewServerTLSFromFile(crt, key)
	if err != nil {
		return fmt.Errorf("could not load TLS keys: %s", err)
	}

	// Create the gRPC server with the credentials
	srv := grpc.NewServer(grpc.Creds(creds))

	// Serve and Listen
	if err := srv.Serve(lis); err != nil {
		return fmt.Errorf("grpc serve error: %s", err)
	}

	return nil
}
*/

//https://gist.github.com/fntlnz/cf14feb5a46b2eda428e000157447309
//https://gist.github.com/artyom/6897140
type Reply struct {
	Message string
	Peers   []string
}

type ContextServer struct {
}

func runContextServer(addr string) {

	cert, err := tls.LoadX509KeyPair("crt/server.crt", "crt/server.key")
	if err != nil {
		log.Fatalf("server: loadkeys: %s", err)
	}

	if len(cert.Certificate) != 2 {
		log.Fatal("server.crt should have 2 concatenated certificates: server + CA")
	}
	ca, err := x509.ParseCertificate(cert.Certificate[1])
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
	/*
		server := new(ContextServer)
		rpc.Register(server)
		rpc.HandleHTTP()
		l, e := net.Listen("tcp", addr)
		if e != nil {
			log.Fatal("listen error:", e)
		}
		http.Serve(l, nil)*/
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
