package context

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"net"
	"net/http"
	"net/rpc"
	"time"

	"github.com/henrikkorsgaard/here.local/logging"
)

var (
	pemBytes      []byte
	serverTLSConf *tls.Config
)

type TLSService struct{}

type TLSReply struct {
	PemBytes []byte
}

func init() {
	var err error
	serverTLSConf, pemBytes, err = generateTLSCertificates()
	logging.Fatal(err)
}

//StartTLSService will start the rpc service for clients to register in
func StartTLSService() {
	tlsService := new(TLSService)
	err := rpc.Register(tlsService)
	logging.Fatal(err)

	rpc.HandleHTTP()
	listener, err := net.Listen("tcp", ":1234")
	logging.Fatal(err)
	go http.Serve(listener, nil)
}

//RegisterNodeAndGetPemBytes will register the location in the database and return pemBytes to generate tls certificate
func (r *TLSService) RegisterNodeAndGetPemBytes(location string, reply *TLSReply) error {
	reply.PemBytes = pemBytes
	return nil
}

/*

TODO : INTEGRATE INTO ABOVE
func (c *ContextServer) ConnectLocation(l models.Location, r *Reply) error {
	fmt.Println("connecting location")
	l.Devices = make(map[string]models.Device)
	err := locationCache.Add(l.MAC, l, cache.DefaultExpiration)
	if err != nil {
		fmt.Println("ConnectLocation err")
		fmt.Println(err)
	}
	return nil
}*/

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
