package configuration

import (
	"fmt"
	"log"
	"os/user"
	"runtime"

	"github.com/henrikkorsgaard/here.local/logging"
)

var (
	MODE                string
	CONTEXT_SERVER_ADDR string
	NODE_MAC_ADDR       string
	NODE_IP_ADDR        string
	NODE_NAME           string
)

const (
	CONFIG_MODE    = "CONFIG"
	SLAVE_MODE     = "SLAVE"
	MASTER_MODE    = "MASTER"
	DEVELOPER_MODE = "DEVELOPER"
)

//Init initialises the sensing node network, certificate and user configuration.
func Init() {
	fmt.Println("Initialising node configuration")
	usr, err := user.Current()
	if err != nil {
		logging.Fatal(err)
	}
	//TODO READ THE CONFIG!
	if runtime.GOOS != "linux" {
		initDevMode()
	} else {
		if usr.Uid != "0" && usr.Gid != "0" {
			fmt.Println("Running in deployment mode!")
		} else {
			log.Fatal("YOU NEED TO RUN HERE.LOCAL AS ROOT")
		}
	}
}

func initDevMode() {
	fmt.Println("Running in developer mode")
	MODE = DEVELOPER_MODE
	CONTEXT_SERVER_ADDR = "localhost:1339"
	NODE_MAC_ADDR = "00:00:00:00:00:00"
	NODE_IP_ADDR = "127.0.0.1"
	NODE_NAME = "test-device"
}
