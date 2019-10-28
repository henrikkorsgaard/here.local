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
	"github.com/henrikkorsgaard/here.local/models"

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
	initSqliteDB()
	Salt = randSeq(8)

	locationCache.OnEvicted(locationEvicted)

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

func (c *ContextServer) DeviceReading(rd models.Reading, r *Reply) error {

	mac := Salt + rd.DeviceMAC

	//we only do that on insert
	hash, err := bcrypt.GenerateFromPassword([]byte(mac), 10)
	if err != nil {
		fmt.Println("should log, but here is error")
		fmt.Println(err)
	}

	var location models.Location
	var device models.Device

	if l, ok := locationCache.Get(rd.LocationMAC); ok {
		location = l.(models.Location)
	} else {
		fmt.Printf("unable to fetch location %s\n", rd.LocationMAC)
		return nil
	}

	devices := deviceCache.Items()
	for m, d := range devices {

		err := bcrypt.CompareHashAndPassword([]byte(m), []byte(mac))
		if err != nil {
			continue
		}

		device = d.Object.(models.Device)
		break
	}

	location.Devices[string(hash)] = models.Device{ID: string(hash), Signal: rd.Signal, Vendor: device.Vendor, Name: device.Name, IP: device.IP}
	device.Locations[rd.LocationMAC] = models.Location{MAC: rd.LocationMAC, IP: location.IP, Name: location.Name, Signal: rd.Signal}
	fmt.Println("Setting device and location")
	deviceCache.SetDefault(string(hash), device)

	locationCache.SetDefault(rd.LocationMAC, location)

	fmt.Printf("Device reading - Mac: %s, Salt: %s, Hash: %s\n", rd.DeviceMAC, Salt, string(hash))
	insertReading(rd.LocationMAC, device.ID, rd.Signal, rd.Timestamp)

	return nil
}

func (c *ContextServer) DeviceEvent(e models.DeviceEvent, r *Reply) error {

	mac := Salt + e.DeviceMAC

	hash, err := bcrypt.GenerateFromPassword([]byte(mac), 10)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Printf("Device evebt - Mac: %s, Salt: %s, Hash: %s\n", e.DeviceMAC, Salt, string(hash))

	if e.Event == models.DEVICE_JOINED {
		fmt.Println("Device joined")
		d := models.Device{ID: string(hash), Locations: make(map[string]models.Location)}
		deviceCache.Add(string(hash), d, cache.DefaultExpiration)

		//initiate a nmap scan as a go routine --- I need to do a simulated nmap as well.
	} else if e.Event == models.DEVICE_LEFT {
		if d, ok := deviceCache.Get(string(hash)); ok {
			device := d.(models.Device)
			//remove location
			fmt.Println(device)

		}

		if l, ok := locationCache.Get(e.LocationMAC); ok {
			location := l.(models.Location)
			fmt.Println(location)
			//remove device
		}

		//we need to remove the location from the device -> monday
		//we need to remove the device from the location -> monday

	}

	return nil
}

func (c *ContextServer) ConnectLocation(l models.Location, r *Reply) error {
	fmt.Println("connecting location")
	l.Devices = make(map[string]models.Device)
	err := locationCache.Add(l.MAC, l, cache.DefaultExpiration)
	if err != nil {
		fmt.Println("ConnectLocation err")
		fmt.Println(err)
	}
	return nil
}

func locationEvicted(mac string, i interface{}) {
	fmt.Println("location evicted !!!!!!!!!!!!!!!!!!!")
	location := i.(models.Location)
	fmt.Printf("Location evicted %+v\n", location)
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
