package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/rpc"
	"time"

	"../context/context"
	"../proximity/proximity"
)

type simDevice struct {
	Mac       string
	Signal    int
	Direction int
}

//we use these constants to give some weights to the movement algorithm
const (
	COMING  = 1
	GOING   = -1
	STAYING = 0
)

var (
	macs           = []string{"b5:a8:34:bc:8e:e3", "32:f8:ad:53:9d:bf", "4e:7f:19:3a:17:9d", "85:73:47:3d:01:2e", "96:4a:46:2c:0f:37", "34:cb:2a:e5:b9:73", "f6:35:9e:d3:21:b4", "25:ca:6b:d6:81:b8", "d1:13:4b:f1:50:6b", "91:59:93:0f:e4:59"}
	ips            = []int{2, 3, 5, 7, 11, 13, 17, 19, 23, 29, 31, 37, 41, 43, 47, 53, 59, 61, 67, 71}
	hostnames      = []string{"Bob-smartphone", "Alice-smartphone", "Eve-smartphone", "Bob-ipad", "Alice-ipad", "Eve-ipad", "Bob-laptop", "Alice-laptop", "Eve-laptop", "SmartTV"}
	vendors        = []string{"Apple", "Samsung", "Sony", "Cisco", "ASUSTek", "Apple", "Samsung", "Sony", "Cisco", "ASUSTek"}
	rawDevicePool  map[string]proximity.RawDevice
	nmapDevicePool map[string]proximity.NmapDevice
	activeDevices  map[string]simDevice
	client         *rpc.Client
)

func RunProximitySimulation() {
	//genrate devices
	activeDevices = make(map[string]simDevice)
	rawDevicePool = make(map[string]proximity.RawDevice)
	nmapDevicePool = make(map[string]proximity.NmapDevice)

	for i := 0; i < 10; i++ {
		rd := proximity.RawDevice{Mac: macs[i], Signal: 0}
		rawDevicePool[macs[i]] = rd
		nd := proximity.NmapDevice{Mac: macs[i], Ip: string(ips[i]), Hostname: hostnames[i], Vendor: vendors[i]}
		nmapDevicePool[macs[i]] = nd
	}

	var err error
	client, err = rpc.DialHTTP("tcp", "localhost:1339")
	if err != nil {

		log.Fatal("Error setting up RPC connection: ", err)
	}

	simProximity()

}

func simProximity() {
	ticker := time.NewTicker(1000 * time.Millisecond)

	for range ticker.C {
		fmt.Println("Simulating proximity data send every 1000 seconds")
		//should do a roll for adding a device
		rnd := rand.Intn(len(macs) + 10)
		if rnd < len(macs) {
			mac := macs[rnd]
			if _, ok := activeDevices[mac]; !ok {
				activeDevices[mac] = simDevice{Mac: mac, Signal: -70, Direction: COMING}
			}
		}

		for _, d := range activeDevices {
			rnd = rand.Intn(10)
			if rnd >= 6 && rnd <= 8 {
				d.Direction = STAYING
			} else {
				d.Direction = d.Direction * -1
			}

			noise := rand.Intn(4) - 2
			d.Signal += d.Direction * noise
			fmt.Printv()
			if d.Signal < -70 {
				fmt.Println("Deleting device: ", d.Mac)
				delete(activeDevices, d.Mac)
			} else {
				var result context.Reply
				err := client.Call("ContextServer.SendProximityData", d, &result)
				if err != nil {
					fmt.Println(err)
				}
				//send via the client and do nothing with the response
			}
		}
	}
}

func SimNmap() {
	//should return a list of devices if they are active
}
