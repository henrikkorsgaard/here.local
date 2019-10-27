package contextserver

import (
	"crypto/rand"
	"crypto/tls"
	"fmt"
	"log"
	mrand "math/rand"
	"net/rpc"
	"time"

	"github.com/henrikkorsgaard/here.local/configuration"

	"github.com/patrickmn/go-cache"
	"golang.org/x/crypto/bcrypt"
)

var (
	deviceCache   = cache.New(cache.NoExpiration, cache.NoExpiration)
	locationCache = cache.New(10*time.Second, 30*time.Second) // prolly should only rarely expire!
	Salt          string
)

type Reply struct {
	Message string
	Peers   []string
}

type ContextServer struct {
}

type RawDevice struct {
	MAC         string
	Signal      int
	LocationMAC string
	Timestamp   time.Time
}

type Device struct {
	Hash        string
	IP          string
	Vendor      string
	Name        string
	Signal      string
	LocationMAC string
}

// Possible solution https://gist.github.com/ncw/9253562
func Run() {
	fmt.Println("Running context server")

	Salt = randSeq(8)

	server := new(ContextServer)
	rpc.Register(server)

	config, err := configuration.GetTLSServerConfig()
	if err != nil {
		fmt.Println(err)
	}

	config.Rand = rand.Reader
	service := configuration.CONTEXT_SERVER_ADDR
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

func (c *ContextServer) DeviceReading(obj interface{}, r *Reply) error {
	fmt.Println("got a reading!")
	rd := obj.(*RawDevice)

	//what happens if it fails

	mac := Salt + rd.MAC
	hash, err := bcrypt.GenerateFromPassword([]byte(mac), 10)
	if err != nil {
		fmt.Println("should log, but here is error")
		fmt.Println(err)
	}

	//we need to find out if the device exists

	insertReading(rd.LocationMAC, string(hash), rd.Signal, rd.Timestamp)
	return nil
}

func (c *ContextServer) DeviceEvent(obj interface{}, r *Reply) error {

	rd := obj.(*RawDevice)

	//what happens if it fails

	mac := Salt + rd.MAC
	hash, err := bcrypt.GenerateFromPassword([]byte(mac), 10)
	if err != nil {
		fmt.Println("should log, but here is error")
		fmt.Println(err)
	}

	//we need to find out if the device exists

	insertReading(rd.LocationMAC, string(hash), rd.Signal, rd.Timestamp)
	return nil
}

func (c *ContextServer) ConnectLocation(obj interface{}, r *Reply) error {
	fmt.Println("Connecting location LKHFALKSFHSAHKFSHAFKSKADFHKKASFHKDSJFHDKSJFH")
	fmt.Println(obj)
	return nil
}

func locationEvicted(mac string, l interface{}) {
	/*
		location := l(*device.Location)
		fmt.Printf("Location evicted %+v\n", location)
	*/
}

//Run starts the context server

func randSeq(n int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ123456789")
	mrand.Seed(time.Now().Unix())
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[mrand.Intn(len(letters))]
	}
	return string(b)
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
