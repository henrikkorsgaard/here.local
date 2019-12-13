package context

import (
	"crypto/rand"
	"crypto/tls"
	"fmt"
	"log"
	mrand "math/rand"
	"net/rpc"
	"strconv"
	"time"

	"github.com/grandcat/zeroconf"
	"github.com/henrikkorsgaard/here.local/configuration"
	"github.com/henrikkorsgaard/here.local/models"

	"github.com/patrickmn/go-cache"
	"golang.org/x/crypto/bcrypt"
)

var (
	deviceCache   = cache.New(cache.NoExpiration, cache.NoExpiration)
	locationCache = cache.New(10*time.Second, 30*time.Second) // prolly should only rarely expire!
	salt          string
)

type Reply struct {
	Message string
	Peers   []string
}

type nmapDevice struct {
	MAC    string
	IP     string
	Vendor string
	Name   string
}

type ContextServer struct{}

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
	salt = randSeq(8)

	locationCache.OnEvicted(locationEvicted)

	StartTLSService()
	StartContextServerRPC()

	//nmapChannel = make(chan nmapDevice)

}

//StartContextServerRPC will start the context server RPC component
func StartContextServerRPC() {
	server := new(ContextServer)
	rpc.Register(server)
	//need to fix this
	config, err := serverTLSConfif
	if err != nil {
		fmt.Println(err)
	}

	config.Rand = rand.Reader
	service := configuration.IP + ":" + strconv.Itoa(configuration.CS_PORT)
	listener, err := tls.Listen("tcp", service, &config)
	if err != nil {
		log.Fatalf("server: listen: %s", err)
	}
	log.Print("server: listening")

	_, err = zeroconf.Register("here.local.context.server", "_http._tcp", "local.", CS_PORT, []string{"txtv=0", "lo=1", "la=2"}, nil)

	//this will block
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

//DeviceReading gets a reading from proximity nodes (testing and training data)
func (c *ContextServer) DeviceReading(rd models.Reading, r *Reply) error {

	mac := salt + rd.DeviceMAC

	//we only do that on insert
	hash, err := bcrypt.GenerateFromPassword([]byte(mac), 10)
	if err != nil {
		fmt.Println("should log, but here is error")
		fmt.Println(err)
	}

	var location models.Location

	if l, ok := locationCache.Get(rd.LocationMAC); ok {
		location = l.(models.Location)
	} else {
		fmt.Printf("unable to fetch location %s\n", rd.LocationMAC)
		return nil
	}

	device, err := getDeviceFromCache(rd.DeviceMAC)
	if err != nil {
		return nil
	}

	location.Devices[string(hash)] = models.Device{ID: string(hash), Signal: rd.Signal, Vendor: device.Vendor, Name: device.Name, IP: device.IP}
	device.Locations[rd.LocationMAC] = models.Location{MAC: rd.LocationMAC, IP: location.IP, Name: location.Name, Signal: rd.Signal}
	deviceCache.SetDefault(string(hash), device)
	locationCache.SetDefault(rd.LocationMAC, location)

	insertReading(rd.LocationMAC, device.ID, rd.Signal, rd.Timestamp)

	return nil
}

//DeviceEvent receives events from proximity nodes
func (c *ContextServer) DeviceEvent(e models.DeviceEvent, r *Reply) error {

	if e.Event == models.DEVICE_JOINED {
		mac := salt + e.DeviceMAC

		hash, err := bcrypt.GenerateFromPassword([]byte(mac), 10)
		if err != nil {
			fmt.Println(err)
		}

		d := models.Device{ID: string(hash), Locations: make(map[string]models.Location)}
		fmt.Printf("Device joined mac: %s, device: %+v\n", e.DeviceMAC, d)
		deviceCache.Add(string(hash), d, cache.DefaultExpiration)

		nmapChannel := make(chan nmapDevice)
		go nmapScan(e.DeviceMAC, nmapChannel)
		go func() {
			for {
				select {
				case d := <-nmapChannel:
					device, err := getDeviceFromCache(d.MAC)
					if err != nil {
						return
					}
					device.IP = d.IP
					device.Name = d.Name
					device.Vendor = d.Vendor
					deviceCache.SetDefault(device.ID, device)
				}
			}
		}()

		//EMMIT API EVENT: Target: location and device, event: device joined

	} else if e.Event == models.DEVICE_LEFT {
		var device models.Device

		device, err := getDeviceFromCache(e.DeviceMAC)
		if err != nil {
			return nil
		}

		delete(device.Locations, e.LocationMAC)
		if len(device.Locations) == 0 {
			deviceCache.Delete(device.ID)
		}

		if l, ok := locationCache.Get(e.LocationMAC); ok {
			location := l.(models.Location)
			delete(location.Devices, device.ID)
		}

		//EMMIT API EVENT: Target: location and device, event: device left
	}

	return nil
}

func getDeviceFromCache(MACaddr string) (device models.Device, err error) {

	smac := salt + MACaddr

	devices := deviceCache.Items()
	for m, d := range devices {

		err := bcrypt.CompareHashAndPassword([]byte(m), []byte(smac))
		if err != nil {
			continue
		}

		device = d.Object.(models.Device)
		break
	}

	return
}

func locationEvicted(mac string, i interface{}) {
	location := i.(models.Location)
	fmt.Printf("Location evicted %+v\n", location)
}

func randSeq(n int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ123456789")
	mrand.Seed(time.Now().Unix())
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[mrand.Intn(len(letters))]
	}
	return string(b)
}

func nmapScan(mac string, ch chan nmapDevice) {
	if configuration.MODE == configuration.DEVELOPER_MODE {
		simulateNmap(mac, ch)
	} else {
		fmt.Println("Run normal nmap scan")
	}
}
