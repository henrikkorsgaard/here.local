package configs

import (
	"crypto/tls"
	"crypto/x509"
	"log"
)

/**
** GetTLSClientConfig() will return a valid tls.Config or error
**/

func GetTLSClientConfig() (config tls.Config, err error) {
	cert, err := tls.LoadX509KeyPair("crts/client.crt", "crts/client.key")
	if err != nil {
		log.Fatalf("client: loadkeys: %s", err)
		return
	}

	certCA, err := tls.LoadX509KeyPair("crts/here.local.crt", "crts/here.local.key")
	if err != nil {
		log.Fatalf("client: loadCA: %s", err)
		return
	}

	ca, err := x509.ParseCertificate(certCA.Certificate[0])
	if err != nil {
		log.Fatal(err)
		return
	}
	certPool := x509.NewCertPool()
	certPool.AddCert(ca)
	config = tls.Config{
		Certificates: []tls.Certificate{cert},
		RootCAs:      certPool,
	}
	return
}

/**
** GetTLSHostConfig() will return a valid tls.Config or error
**/

func GetTLSHostConfig() (config tls.Config, err error) {
	cert, err := tls.LoadX509KeyPair("crts/localhost.crt", "crts/localhost.key")
	if err != nil {
		log.Fatalf("server: loadkeys: %s", err)
		return
	}

	certCA, err := tls.LoadX509KeyPair("crts/here.local.crt", "crts/here.local.key")
	if err != nil {
		log.Fatalf("server: loadCA: %s", err)
		return
	}

	ca, err := x509.ParseCertificate(certCA.Certificate[0])
	if err != nil {
		log.Fatal(err)
		return
	}

	certPool := x509.NewCertPool()
	certPool.AddCert(ca)
	config = tls.Config{
		Certificates: []tls.Certificate{cert},
		ClientAuth:   tls.RequireAndVerifyClientCert,
		ClientCAs:    certPool,
	}

	return
}
