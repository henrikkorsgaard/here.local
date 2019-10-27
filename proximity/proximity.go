package proximity

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/rpc"
	"time"

	"github.com/henrikkorsgaard/here.local/configuration"
	"github.com/henrikkorsgaard/here.local/server/contextserver"
	"github.com/henrikkorsgaard/kalmango"
	"github.com/patrickmn/go-cache"
)

var (
	rpcClient   *rpc.Client
	configTLS   tls.Config
	deviceCache = cache.New(10*time.Second, 30*time.Second)
)

type Location struct {
	MAC  string
	IP   string
	Name string
}

type Device struct {
	MAC        string
	Signal     int
	Discovered time.Time

	kalman kalmango.KalmanFilter
}

type Reading struct {
	MAC       string
	Signal    int
	Timestamp time.Time
}

type DeviceEvent struct {
	Event       string
	DeviceMAC   string
	LocationMAC string
	Timestamp   time.Time
}

const (
	DEVICE_JOINED = "device joined"
	DEVICE_LEFT   = "device left"
)

//NOTE: we want to do some "bulk" sends of data, so to avoid near syncronous transmits that lead to increased db writes on the context server.
//NOTE: we could do some normalisation here with a kalman filter?

func Run() {

	var err error
	configTLS, err = configuration.GetTLSClientConfig()
	if err != nil {
		log.Fatalf("configurationGetTLSClientConfig failed: %s", err)
	}

	connectRPC()
	deviceCache.OnEvicted(deviceEvicted)
	mode := configuration.MODE
	if mode == configuration.DEVELOPER_MODE {
		simulate()
	} else {

	}

}

func connectRPC() {
	fmt.Println("Connecting proximity sensor to context server")
	conn, err := tls.Dial("tcp", configuration.CONTEXT_SERVER_ADDR, &configTLS)
	if err != nil {
		fmt.Println("Unable to establish RCP connection")
		fmt.Println("Trying again in a second")
		return
	}

	//defer conn.Close()
	rpcClient = rpc.NewClient(conn)
	l := Location{MAC: configuration.NODE_MAC_ADDR, IP: configuration.NODE_IP_ADDR, Name: configuration.NODE_NAME}
	var result contextserver.Reply
	rpcClient.Call("ContextServer.ConnectLocation", l, &result)
	fmt.Println(result)

}

func sendDevice(MAC string, Signal int) {
	fmt.Printf("Getting device %s\n", MAC)
	var device Device
	if obj, ok := deviceCache.Get(MAC); ok {
		device = obj.(Device)
		ksig := device.kalman.Filter(float64(Signal), 0)
		device.Signal = int(ksig)
	} else {
		var result contextserver.Reply
		device = Device{MAC: MAC, Signal: Signal, kalman: kalmango.NewKalmanFilter(0.5, 8, 1, 0, 1), Discovered: time.Now()}
		rpcClient.Call("ContextServer.DeviceEvent", DeviceEvent{Event: DEVICE_JOINED, DeviceMAC: MAC, LocationMAC: configuration.NODE_MAC_ADDR, Timestamp: time.Now()}, &result)
	}
	var result contextserver.Reply
	deviceCache.Set(MAC, device, cache.DefaultExpiration)
	rpcClient.Call("ContextServer.DeviceReading", Reading{MAC, Signal, time.Now()}, &result)
}

func deviceEvicted(MAC string, i interface{}) {
	fmt.Println("calling evicted")

	//send the event please
}
