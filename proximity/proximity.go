package proximity

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/rpc"
	"time"

	"github.com/henrikkorsgaard/here.local/configuration"
	"github.com/henrikkorsgaard/here.local/logging"
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
	} else if mode != configuration.CONFIG_MODE {
		monitorNetwork()
	}

}

func connectRPC() {
	fmt.Println("Connecting proximity sensor to context server")
	conn, err := tls.Dial("tcp", "here.local:1337", &configTLS)
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
	fmt.Println("Detected device with mac ", MAC, " and signal ", Signal)
	var device Device
	if obj, ok := deviceCache.Get(MAC); ok {
		device = obj.(Device)
		ksig := device.kalman.Filter(float64(Signal), 0)
		device.Signal = int(ksig)
	} else {
		device = Device{MAC: MAC, Signal: Signal, kalman: kalmango.NewKalmanFilter(0.5, 8, 1, 0, 1), Discovered: time.Now()}
		var result context.Reply
		rpcClient.Call("Context.DeviceEvent", models.DeviceEvent{Event: models.DEVICE_JOINED, DeviceMAC: MAC, LocationMAC: configuration.MAC, Timestamp: time.Now()}, &result)
	}

	var result context.Reply
	rpcClient.Call("Context.DeviceReading", models.Reading{MAC, configuration.MAC, Signal, time.Now()}, &result)
	deviceCache.Set(MAC, device, cache.DefaultExpiration)
}

func deviceEvicted(MAC string, i interface{}) {
	device := i.(Device)
	fmt.Printf("Evicting device: %+v", device)
	var result context.Reply
	rpcClient.Call("Context.DeviceEvent", models.DeviceEvent{Event: models.DEVICE_LEFT, DeviceMAC: MAC, LocationMAC: configuration.MAC, Timestamp: time.Now()}, &result)
}

func monitorWifiNetwork() {
	handle, err := pcap.OpenLive("here-monitor", snapLen, true, pcap.BlockForever)
	logging.Fatal(err)

	monitorHandle = handle
	defer monitorHandle.Close()

	// Set filter
	var filter = "not broadcast" //TODO: ADD DESTINATION/SOURCE TO THE FILTER TO AVOID GETTING TOO MANY PACKETS
	err = monitorHandle.SetBPFFilter(filter)
	if err != nil {
		log.Fatal(err)
	}

	packetSource := gopacket.NewPacketSource(monitorHandle, monitorHandle.LinkType())

	parser := gopacket.NewDecodingLayerParser(
		layers.LayerTypeRadioTap,
		&radioLayer,
		&dot11layer,
	)

	foundLayerTypes := []gopacket.LayerType{}

	for packet := range packetSource.Packets() {

		parser.DecodeLayers(packet.Data(), &foundLayerTypes)

		if len(foundLayerTypes) >= 2 && radioLayer.DBMAntennaSignal != 0 {

			station := configuration.STATION
			addr1 := dot11layer.Address1.String()
			addr2 := dot11layer.Address2.String()
			addr3 := dot11layer.Address3.String()

			//see: https://networkengineering.stackexchange.com/questions/25100/four-layer-2-addresses-in-802-11-frame-header
			//The ideal case for capturing signal strength between the node and device
			//is when addr1 and addr3 is equal to the station mac and addr2 is not equal
			//to station mac

			if addr1 == station && addr3 == station && addr2 != station {
				signal := radioLayer.DBMAntennaSignal
				mac := dot11layer.Address2.String()
				sendDevice(mac, int(signal))
			}
		}
	}
}
