package configuration

import (
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"os/user"
	"regexp"
	"runtime"
	"time"

	"github.com/spf13/viper"

	"github.com/henrikkorsgaard/here.local/logging"
)

var (
	MODE                string
	CONTEXT_SERVER_ADDR string
	MAC            string
	IP		    string
	STATION		    string
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

	var cfpath string //for viper
	var cfpathfull string //for copying config file if not existing
	if runtime.GOOS != "linux" {
		cfpath = "."
		cfpathfull = "./here.local.config.toml"
		MODE = DEVELOPER_MODE
	} else {
		if usr.Uid != "0" && usr.Gid != "0" {
			log.Fatal("YOU NEED TO RUN HERE.LOCAL AS ROOT")
		} else {
			cfpath = "/boot/"
			cfpathfull = "/boot/here.local.config.toml"
		}
	}

	if _, err := os.Stat(cfpathfull); os.IsNotExist(err) {
		input, err := ioutil.ReadFile("here.local.config.toml.template")
		logging.Fatal(err)

		err = ioutil.WriteFile(cfpathfull, input, 0644)
		logging.Fatal(err)
	}

	viper.SetConfigName("here.local.config")
	viper.AddConfigPath(cfpath)
	err = viper.ReadInConfig()
	logging.Fatal(err)

	if MODE != DEVELOPER_MODE {
		var err error

		hostname, err := os.Hostname()
		logging.Fatal(err)

		if LOCATION() == "" {
			newHostname := "HERE-" + randSeq(4)
			setLocation(newHostname)
			err = changeHostname(newHostname)
			logging.Fatal(err)
			reboot()
		} else if LOCATION() != hostname {
			err = changeHostname(LOCATION())
			logging.Fatal(err)
			reboot()
		}

		configureNetworkInterfaces()
	} else {
		CONTEXT_SERVER_ADDR = "localhost:1339"
		MAC = "00:00:00:00:00:00"
	       IP = "127.0.0.1"
        }
}

//SSID returns ssid from the config
func SSID() string {
	return viper.GetString("network.ssid")
}

//SSID_PASSWORD returns network password from the config
func SSID_PASSWORD() string {
	return viper.GetString("network.password")
}

//LOCATION returns location from the config
func LOCATION() string {
	return viper.GetString("node.location")
}

func setLocation(l string) {
	viper.Set("node.location", l)
	viper.WriteConfig()
}

func generateValidHostname(hostname string) string {
	re := regexp.MustCompile(`[^a-zA-Z0-9--]`)
	str := re.ReplaceAllString(hostname, "")
	if len(str) > 32 {
		str = str[0:32]
	}
	return str
}

func delayedReboot() {
	time.Sleep(2 * time.Second)
	reboot()
}

func reboot() {
	runCommand("sudo reboot now")
}

func randSeq(n int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ123456789")
	rand.Seed(time.Now().Unix())
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func changeHostname(hostname string) (err error) {
	_, stderr, _ := runCommand("sudo hostnamectl set-hostname " + hostname)

	if stderr != "" {
		err = fmt.Errorf(stderr)
		return
	}

	err = ioutil.WriteFile("/etc/hostname", []byte(hostname), 0666)

	if err != nil {
		return
	}

	err = ioutil.WriteFile("/etc/hosts", []byte("127.0.0.1\tlocalhost\n127.0.1.1\t"+hostname+"\n"), 0666)

	if err != nil {
		return
	}

	return nil
}
