package proximity

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/rpc"
	"time"

	"github.com/henrikkorsgaard/here.local/configuration"
	"github.com/henrikkorsgaard/here.local/models"
	"github.com/henrikkorsgaard/here.local/server/context"
	contextserver "github.com/henrikkorsgaard/here.local/server/context"
	"github.com/henrikkorsgaard/kalmango"
	"github.com/patrickmn/go-cache"
)

var (
	rpcClient   *rpc.Client
	configTLS   tls.Config
	deviceCache = cache.New(10*time.Second, 30*time.Second)
)

type Device struct {
	MAC        string
	Signal     int
	Discovered time.Time

	kalman kalmango.KalmanFilter
}

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

	l := models.Location{MAC: configuration.NODE_MAC_ADDR, IP: configuration.NODE_IP_ADDR, Name: configuration.NODE_NAME}
	var result contextserver.Reply
	rpcClient.Call("Context.ConnectLocation", l, &result)
	fmt.Println(result)

}

func sendDevice(MAC string, Signal int) {
	var device Device
	if obj, ok := deviceCache.Get(MAC); ok {
		device = obj.(Device)
		ksig := device.kalman.Filter(float64(Signal), 0)
		device.Signal = int(ksig)
	} else {
		device = Device{MAC: MAC, Signal: Signal, kalman: kalmango.NewKalmanFilter(0.5, 8, 1, 0, 1), Discovered: time.Now()}
		var result context.Reply
		rpcClient.Call("Context.DeviceEvent", models.DeviceEvent{Event: models.DEVICE_JOINED, DeviceMAC: MAC, LocationMAC: configuration.NODE_MAC_ADDR, Timestamp: time.Now()}, &result)
	}

	var result context.Reply
	rpcClient.Call("Context.DeviceReading", models.Reading{MAC, configuration.NODE_MAC_ADDR, Signal, time.Now()}, &result)
	deviceCache.Set(MAC, device, cache.DefaultExpiration)
}

func deviceEvicted(MAC string, i interface{}) {
	device := i.(Device)
	fmt.Printf("Evicting device: %+v", device)
	var result context.Reply
	rpcClient.Call("Context.DeviceEvent", models.DeviceEvent{Event: models.DEVICE_LEFT, DeviceMAC: MAC, LocationMAC: configuration.NODE_MAC_ADDR, Timestamp: time.Now()}, &result)
}
