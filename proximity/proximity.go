package proximity

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/rpc"
	"time"

	"github.com/henrikkorsgaard/here.local/device"
	"github.com/henrikkorsgaard/here.local/server/context"

	"github.com/henrikkorsgaard/here.local/initialise/configs"
)

var (
	rpcClient   *rpc.Client
	configTLS   tls.Config
	deviceCache = cache.New(10*time.Second, 30*time.Second)
)

type Device struct {
	MAC   string
	IP    string
	Name  string
	Wired bool

	Signal  *int
	Signals []int

	kalman kalmango.KalmanFilter
}

type ProximityEvent struct {
	Event    string
	Device   Device
	Location Location
}

//NOTE: we want to do some "bulk" sends of data, so to avoid near syncronous transmits that lead to increased db writes on the context server.
//NOTE: we could do some normalisation here with a kalman filter?

func Run() {

	var err error
	configTLS, err = configs.GetTLSClientConfig()
	if err != nil {
		log.Fatalf("configurationGetTLSClientConfig failed: %s", err)
	}

	mode := configs.MODE
	if mode == configs.DEVELOPER_MODE {
		simulate()
	} else {

	}

	deviceCache.on
}

func connectRPC() {
	conn, err := tls.Dial("tcp", configs.CONTEXT_SERVER_ADDR, &configTLS)
	if err != nil {
		fmt.Println("Unable to establish RCP connection")
		fmt.Println("Trying again in a second")
		return
	}
	//defer conn.Close()
	rpcClient = rpc.NewClient(conn)
	l := device.Location{MAC: configs.NODE_MAC_ADDR, IP: configs.NODE_IP_ADDR, Name: configs.NODE_NAME}
	var result context.Reply
	rpcClient.Call("ContextServer.ConnectLocation", l, &result)

}

func sendDevice(d Device) {

}

func deviceEvicted() {
	//send the event please
}
