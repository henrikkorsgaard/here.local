package main

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"log"
	"math/big"
	"net"
	"net/http"
	"net/rpc"
	"time"
)

var (
	//ServerTLSConf is the sever tls config for RPC
	serverTLSConf *tls.Config
	//ClientTLSConf is the client tls for RPC
	PEMBytes []byte
)

type TLSRPCServer struct {
}

type RPCServer struct {
}

type RPCREPLY struct {
	Bytes []byte
}

type TLSREPLY struct {
	MSG string
}

func (r *TLSRPCServer) Hello(msg string, reply *TLSREPLY) error {
	reply.MSG = "hello"

	return nil
}

func (r *RPCServer) GetPEMBytes(msg string, reply *RPCREPLY) error {
	reply.Bytes = PEMBytes

	return nil
}

func main() {
	// get our ca and server certificate
	var err error
	serverTLSConf, PEMBytes, err = generateTLSCertificates()
	if err != nil {
		panic(err)
	}

	rpcserver := new(RPCServer)
	// Publish the receivers methods
	err = rpc.Register(rpcserver)
	if err != nil {
		log.Fatal(err)
	}
	// Register a HTTP handler
	rpc.HandleHTTP()
	// Listen to TPC connections on port 1234
	listener, e := net.Listen("tcp", ":1234")
	if e != nil {
		log.Fatal("Listen error: ", e)
	}
	log.Printf("Serving RPC server on port %d", 1234)
	// Start accept incoming HTTP connections
	go http.Serve(listener, nil)
	go startTLSServer()
	startClientServer()
	for {
	}
}

func startTLSServer() {
	tlsserver := new(TLSRPCServer)
	rpc.Register(tlsserver)

	serverTLSConf.Rand = rand.Reader

	tlslistener, err := tls.Listen("tcp", ":1235", serverTLSConf)
	if err != nil {
		log.Fatalf("server: listen: %s", err)
	}
	log.Print("server: listening")
	for {
		conn, err := tlslistener.Accept()
		if err != nil {
			log.Printf("server: accept: %s", err)
			break
		}
		log.Printf("server: accepted from %s", conn.RemoteAddr())

		go rpc.ServeConn(conn)
	}
}

func startClientServer() {
	time.Sleep(2 * time.Second)
	// Create a TCP connection to localhost on port 1234
	var reply RPCREPLY

	client, err := rpc.DialHTTP("tcp", "localhost:1234")
	if err != nil {
		log.Fatal("Connection error: ", err)
	}
	//var i interface{}

	client.Call("RPCServer.GetPEMBytes", "", &reply)

	certpool := x509.NewCertPool()
	certpool.AppendCertsFromPEM(reply.Bytes)
	clientTLSConf := &tls.Config{
		RootCAs: certpool,
	}

	fmt.Println("Connecting proximity sensor to context server")

	conn, err := tls.Dial("tcp", "192.168.1.134:1235", clientTLSConf)
	if err != nil {
		fmt.Println("Unable to establish RCP connection")
		fmt.Println("Trying again in a second")
		fmt.Println(err)
		return
	}

	//defer conn.Close()
	rpcClient := rpc.NewClient(conn)
	var result TLSREPLY
	rpcClient.Call("TLSRPCServer.Hello", "", &result)
	fmt.Println(result)

}

func generateTLSCertificates() (serverTLSConf *tls.Config, PEMBytes []byte, err error) {

	ca := &x509.Certificate{
		SerialNumber: big.NewInt(2019),
		Subject: pkix.Name{
			Organization:  []string{"Computer Science, Aarhus University"},
			Country:       []string{"DK"},
			Province:      []string{""},
			Locality:      []string{"Aarhus N"},
			StreetAddress: []string{"Aabogade 34"},
			PostalCode:    []string{"8200"},
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(10, 0, 0),
		IsCA:                  true,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		BasicConstraintsValid: true,
	}

	caPrivKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return nil, nil, err
	}

	caBytes, err := x509.CreateCertificate(rand.Reader, ca, ca, &caPrivKey.PublicKey, caPrivKey)
	if err != nil {
		return nil, nil, err
	}

	caPEM := new(bytes.Buffer)
	pem.Encode(caPEM, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: caBytes,
	})

	caPrivKeyPEM := new(bytes.Buffer)
	pem.Encode(caPrivKeyPEM, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(caPrivKey),
	})

	// set up our server certificate
	cert := &x509.Certificate{
		SerialNumber: big.NewInt(2019),
		Subject: pkix.Name{
			Organization:  []string{"Computer Science, Aarhus University"},
			Country:       []string{"DK"},
			Province:      []string{""},
			Locality:      []string{"Aarhus N"},
			StreetAddress: []string{"Aabogade 34"},
			PostalCode:    []string{"8200"},
		},
		DNSNames: []string{"localhost", "here.local"},
		//WE NEED TO ADD THE LAN IP TO THE LIST -- replace "192.168.1.134 with local lan ip.
		IPAddresses:  []net.IP{net.IPv4(127, 0, 0, 1), net.ParseIP("192.168.1.134"), net.IPv6loopback},
		NotBefore:    time.Now(),
		NotAfter:     time.Now().AddDate(10, 0, 0),
		SubjectKeyId: []byte{1, 2, 3, 4, 6},
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:     x509.KeyUsageDigitalSignature,
	}

	certPrivKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return nil, nil, err
	}

	certBytes, err := x509.CreateCertificate(rand.Reader, cert, ca, &certPrivKey.PublicKey, caPrivKey)
	if err != nil {
		return nil, nil, err
	}

	certPEM := new(bytes.Buffer)
	pem.Encode(certPEM, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: certBytes,
	})

	certPrivKeyPEM := new(bytes.Buffer)
	pem.Encode(certPrivKeyPEM, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(certPrivKey),
	})

	serverCert, err := tls.X509KeyPair(certPEM.Bytes(), certPrivKeyPEM.Bytes())
	if err != nil {
		return nil, nil, err
	}

	serverTLSConf = &tls.Config{
		Certificates: []tls.Certificate{serverCert},
	}

	PEMBytes = caPEM.Bytes()

	return
}
