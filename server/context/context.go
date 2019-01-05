package context

import (
	"crypto/rand"

	"crypto/rsa"
	"log"
	mrand "math/rand"
	"net"
	"net/http"
	"net/rpc"
	"time"
)

var (
	privateKey *rsa.PrivateKey
	publicKey  rsa.PublicKey
	salt       string
)

const (
	letterBytes   = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

func init() {
	var err error
	privateKey, err = rsa.GenerateKey(rand.Reader, 512)
	if err != nil {
		log.Fatal("Unable to generate private key: ", err)
	}
	publicKey = privateKey.PublicKey
	salt = randString(512)
}

//https://golang.org/pkg/net/rpc/

//Context is the context server struct used for rpc communication
type Context struct {
	name string
}

//Hello returns "World"
func (c *Context) Echo(msg *string, reply *string) error {
	*reply = *msg
	return nil
}

//GetPublicKey returns public key
func (c *Context) GetPublicKey(msg *string, key *rsa.PublicKey) error {
	*key = publicKey
	return nil
}

func (c *Context) GetSalt(msg *string, slt *string) error {
	*slt = salt
	return nil
}

//Run starts the context server
func Run() {
	server := new(Context)
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

func randString(n int) string {
	var src = mrand.NewSource(time.Now().UnixNano())
	b := make([]byte, n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return string(b)
}
