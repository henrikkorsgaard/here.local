package context

import (
	"fmt"
	"math/rand"
	"time"
)

var (
	ips       = []int{2, 3, 5, 7, 11, 13, 17, 19, 23, 29, 31, 37, 41, 43, 47, 53, 59, 61, 67, 71}
	hostnames = []string{"Bob-smartphone", "Alice-smartphone", "Eve-smartphone", "Bob-ipad", "Alice-ipad", "Eve-ipad", "Bob-laptop", "Alice-laptop", "Eve-laptop", "SmartTV"}
	vendors   = []string{"Apple", "Samsung", "Sony", "Cisco", "ASUSTek", "Apple", "Samsung", "Sony", "Cisco", "ASUSTek"}
)

func simulateNmap(mac string, ch chan nmapDevice) {
	time.Sleep(666 * time.Millisecond)
	i := rand.Intn(10)
	device := nmapDevice{
		MAC:    mac,
		IP:     fmt.Sprintf("10.0.0.%d", ips[i]),
		Vendor: vendors[i],
		Name:   hostnames[i],
	}

	ch <- device
}
