package configuration

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"runtime"

	"github.com/henrikkorsgaard/here.local/logging"
	goconfig "github.com/zpatrick/go-config"
)

var (
	MODE                string
	CONTEXT_SERVER_ADDR string
	NODE_MAC_ADDR       string
	NODE_IP_ADDR        string
	NODE_NAME           string
	SSID string
	SSID_PASSWORD string
	LOCATION string
	config *goconfig.Config
)

const (
	CONFIG_MODE    = "CONFIG"
	SLAVE_MODE     = "SLAVE"
	MASTER_MODE    = "MASTER"
	DEVELOPER_MODE = "DEVELOPER"
)

//Init initialises the sensing node network, certificate and user configuration.
func init() {
	fmt.Println("Initialising node configuration")
	usr, err := user.Current()
	if err != nil {
		logging.Fatal(err)
	}

	var cfname string
	if runtime.GOOS != "linux" {
		cfname = "./here.local.config.toml"
		MODE = DEVELOPER_MODE
	} else {
		if usr.Uid != "0" && usr.Gid != "0" {
			log.Fatal("YOU NEED TO RUN HERE.LOCAL AS ROOT")
		} else {
			cfname = "/boot/here.local.config.toml"
		}
	}

	if _, err := os.Stat(cfname); os.IsNotExist(err) {
		input, err := ioutil.ReadFile("./here.local.config.toml.template")
		logging.Fatal(err)

		err = ioutil.WriteFile("./here.local.config.toml", input, 0644)
		logging.Fatal(err)
	}

	cf := goconfig.NewTOMLFile(cfname)
	cl := goconfig.NewOnceLoader(cf)
	config = goconfig.NewConfig([]goconfig.Provider{cl})

	if MODE != DEVELOPER_MODE {
		var err error
		SSID, err =  config.String("network.ssid")
		SSID_PASSWORD, err = config.String("network.password")
		LOCATION, err = config.String("node.location")
		logging.Fatal(err)


		if LOCATION == "" {
		 //generate
		 //change hostname and reboot
		} LOCATION != "hostname" {
		 //change hostname and reboot
		}
	} else {
		CONTEXT_SERVER_ADDR = "localhost:1339"
		NODE_MAC_ADDR = "00:00:00:00:00:00"
		NODE_IP_ADDR = "127.0.0.1"
		NODE_NAME = "dev-home"
	}
}
