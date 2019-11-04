package proximity

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"log"
	"net/rpc"
	"strconv"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	"github.com/grandcat/zeroconf"

	"github.com/henrikkorsgaard/here.local/configuration"
	"github.com/henrikkorsgaard/here.local/logging"
	"github.com/henrikkorsgaard/here.local/models"
	contextserver "github.com/henrikkorsgaard/here.local/server/context"
	"github.com/henrikkorsgaard/kalmango"
	"github.com/patrickmn/go-cache"
)

var (
	rpcClient     *rpc.Client
	configTLS     tls.Config
	deviceCache   = cache.New(10*time.Second, 30*time.Second)
	rpcServerIP   string
	rpcServerPort int
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

	deviceCache.OnEvicted(deviceEvicted)

		fmt.Println("soo far 1")
	var err error
	configTLS, err = configuration.GetTLSClientConfig()
	if err != nil {
		log.Fatalf("configurationGetTLSClientConfig failed: %s", err)
	}
	rpcServerIP, rpcServerPort, err = detectContextServerIP()
	fmt.Println(err)
	fmt.Println("soo far 2")
	if err == nil && (rpcServerIP != "" || rpcServerPort != 0) {
		fmt.Println("soo far 3")
		connectRPC()
	}

	mode := configuration.MODE
	if mode == configuration.DEVELOPER_MODE {
		//simulate()
	} else if mode != configuration.CONFIG_MODE {
		//monitorWifiNetwork()
	}

}

func connectRPC() {

	fmt.Println("Connecting proximity sensor to context server")
	fmt.Println(rpcServerIP, " - ", rpcServerPort)

	conn, err := tls.Dial("tcp", "here.local:"+strconv.Itoa(rpcServerPort), &configTLS)
	if err != nil {
		fmt.Println("Unable to establish RCP connection")
		fmt.Println("Trying again in a second")
		fmt.Println(err)
		return
	}

	//defer conn.Close()
	rpcClient = rpc.NewClient(conn)

	l := models.Location{MAC: configuration.MAC, IP: configuration.IP, Name: configuration.LOCATION()}
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
		var result contextserver.Reply
		rpcClient.Call("Context.DeviceEvent", models.DeviceEvent{Event: models.DEVICE_JOINED, DeviceMAC: MAC, LocationMAC: configuration.MAC, Timestamp: time.Now()}, &result)
	}

	var result contextserver.Reply
	rpcClient.Call("Context.DeviceReading", models.Reading{MAC, configuration.MAC, Signal, time.Now()}, &result)
	deviceCache.Set(MAC, device, cache.DefaultExpiration)
}

func deviceEvicted(MAC string, i interface{}) {
	device := i.(Device)
	fmt.Printf("Evicting device: %+v", device)
	var result contextserver.Reply
	rpcClient.Call("Context.DeviceEvent", models.DeviceEvent{Event: models.DEVICE_LEFT, DeviceMAC: MAC, LocationMAC: configuration.MAC, Timestamp: time.Now()}, &result)
}

func detectContextServerIP() (ip string, port int, err error) {
	fmt.Println("detecting context server ip")
	resolver, err := zeroconf.NewResolver(nil)
	if err != nil {
		return
	}

	entries := make(chan *zeroconf.ServiceEntry)

	go func(results <-chan *zeroconf.ServiceEntry) {
		for entry := range results {
			if entry.Instance == "here.local.context.server" {
				if len(entry.AddrIPv4) != 0 {
					fmt.Printf("%+v\n", entry)
					ip = entry.AddrIPv4[0].String()
					port = entry.Port

				} else {
					err = errors.New("ServiceEntry does not contain a IPv4 address for context server")
				}

				break
			}
		}
	}(entries)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()

	err = resolver.Browse(ctx, "_http._tcp", "local.", entries)
	if err != nil {
		return
	}


	<-ctx.Done()
	fmt.Println("from detection ", ip, " ", port)
	return
}

func monitorWifiNetwork() {
	fmt.Println("starting to monitor WiFi network on here-monitor interface")
	handle, err := pcap.OpenLive("here-monitor", 96, true, pcap.BlockForever)
	logging.Fatal(err)

	defer handle.Close()

	// Set filter
	/*
	var filter = "not broadcast" //TODO: ADD DESTINATION/SOURCE TO THE FILTER TO AVOID GETTING TOO MANY PACKETS
	err = handle.SetBPFFilter(filter)
	if err != nil {
		log.Fatal(err)
	}*/


	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())

	var radioLayer layers.RadioTap
	var dot11layer layers.Dot11

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
			//addr4 := dot11layer.Address4.String()
			//test := "38:f9:d3:20:3f:e9"
						
			//fmt.Println(station, addr1, addr2, addr3)
			//see: https://networkengineering.stackexchange.com/questions/25100/four-layer-2-addresses-in-802-11-frame-header
			//The ideal case for capturing signal strength between the node and device
			//is when addr1 and addr3 is equal to the station mac and addr2 is not equal
			//to station mac

			//fmt.Println("Addr1: ", addr1, " Addr 2: ", addr2, " Addr3 ", addr3, " Addr4 ", addr4, " SIG ", radioLayer.DBMAntennaSignal)

			if addr1 == station || addr3 == station && addr2 != station {
				signal := radioLayer.DBMAntennaSignal
				mac := dot11layer.Address2.String()
				fmt.Println(mac, " - ", signal)
				/*
				if rpcClient != nil {
					fmt.Println(mac)
					sendDevice(mac, int(signal))
				} else {
					connectRPC()
				}*/
			}
		}
	}
}
