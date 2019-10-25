package proximity

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/henrikkorsgaard/here.local/device"
)

//we use these constants to give some weights to the movement algorithm
const (
	COMING  = 1
	GOING   = -1
	STAYING = 0
)

type simDevice struct {
	Device
	Direction int
}

var (
	macs           = []string{"b5:a8:34:bc:8e:e3", "32:f8:ad:53:9d:bf", "4e:7f:19:3a:17:9d", "85:73:47:3d:01:2e", "96:4a:46:2c:0f:37", "34:cb:2a:e5:b9:73", "f6:35:9e:d3:21:b4", "25:ca:6b:d6:81:b8", "d1:13:4b:f1:50:6b", "91:59:93:0f:e4:59"}
	ips            = []int{2, 3, 5, 7, 11, 13, 17, 19, 23, 29, 31, 37, 41, 43, 47, 53, 59, 61, 67, 71}
	hostnames      = []string{"Bob-smartphone", "Alice-smartphone", "Eve-smartphone", "Bob-ipad", "Alice-ipad", "Eve-ipad", "Bob-laptop", "Alice-laptop", "Eve-laptop", "SmartTV"}
	vendors        = []string{"Apple", "Samsung", "Sony", "Cisco", "ASUSTek", "Apple", "Samsung", "Sony", "Cisco", "ASUSTek"}
	rawDevicePool  map[string]Device
	nmapDevicePool map[string]Device
	activeDevices  map[string]Device
)

func simulate() {
	//generate devices
	activeDevices = make(map[string]simDevice)
	rawDevicePool = make(map[string]Device)
	nmapDevicePool = make(map[string]Device)

	for i := 0; i < 10; i++ {
		rd := device.RawDevice{MAC: macs[i], Signal: 0}
		rawDevicePool[macs[i]] = rd
		nd := Device{MAC: macs[i], IP: string(ips[i]), Hostname: hostnames[i], Vendor: vendors[i]}
		nmapDevicePool[macs[i]] = nd
	}

	ticker := time.NewTicker(1000 * time.Millisecond)

	for range ticker.C {
		fmt.Println("Simulating proximity data send every 1000 seconds")
		//should do a roll for adding a device
		rnd := rand.Intn(len(macs) + 10)
		if rnd < len(macs) {
			mac := macs[rnd]
			if _, ok := activeDevices[mac]; !ok {
				activeDevices[mac] = simDevice{Device: Device{MAC: mac, Signal: -70, IP: "", Vendor: "simulation"}, Direction: COMING}
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
			if d.Signal < -70 {
				fmt.Println("Deleting device: ", d.MAC)
				delete(activeDevices, d.MAC)
			} else {
				if rpcClient != nil {
					sendDevice(d.Device)
				} else {
					connectRPC()
				}
			}
		}
	}

}

func SimNmap() {
	//should return a list of devices if they are active
}
