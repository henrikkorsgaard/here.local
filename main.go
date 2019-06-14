package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"log"
	"net/rpc"

	"github.com/henrikkorsgaard/here.local/configuration"
	"github.com/henrikkorsgaard/here.local/server"
)

func main() {
	configuration.Setup()
	go server.Run()

	cert, err := tls.LoadX509KeyPair("crts/client.crt", "crts/client.key")
	if err != nil {
		log.Fatalf("client: loadkeys: %s", err)
	}
	/*
		if len(cert.Certificate) != 2 {
			log.Fatal("client.crt should have 2 concatenated certificates: client + CA")
		}*/
	/*
		ca, err := x509.ParseCertificate(cert.Certificate[0])
		if err != nil {
			log.Fatal(err)
		}*/

	certCA, err := tls.LoadX509KeyPair("crts/here.local.crt", "crts/here.local.key")
	if err != nil {
		log.Fatalf("client: loadCA: %s", err)
	}

	ca, err := x509.ParseCertificate(certCA.Certificate[0])
	if err != nil {
		log.Fatal(err)
	}
	certPool := x509.NewCertPool()
	certPool.AddCert(ca)
	config := tls.Config{
		Certificates: []tls.Certificate{cert},
		RootCAs:      certPool,
	}
	fmt.Println(configuration.ContextServerAddress)
	conn, err := tls.Dial("tcp", configuration.ContextServerAddress, &config)
	if err != nil {
		log.Fatalf("client: dial: %s", err)
	}
	defer conn.Close()
	log.Println("client: connected to: ", conn.RemoteAddr())
	rpcClient := rpc.NewClient(conn)
	res := new(server.Reply)
	if err := rpcClient.Call("ContextServer.Echo", "hep", &res); err != nil {
		log.Fatal("Failed to call RPC", err)
	}
	log.Printf("Returned result is %d", res.Message)
}
