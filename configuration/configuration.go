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
	"sync"
	"time"

	"github.com/henrikkorsgaard/here.local/logging"
	"github.com/spf13/viper"
)

var (
	/*
		configViber reflects  configurations in the public configuration file on the node
		envViper reflects environement configurations
		we want to avoid mixing these, as Viper writes all the settings to the file.
	*/
	configViper *viper.Viper
	envViper    *viper.Viper
	devMode     bool
	linux       bool
)

const (
	CONFIG_MODE    = "CONFIG"
	SLAVE_MODE     = "SLAVE"
	MASTER_MODE    = "MASTER"
	DEVELOPER_MODE = "DEVELOPER"
)

type Configuration struct {
}

var instance *Configuration
var once sync.Once

//GetConfiguration follows singleton pattern introduced here: http://marcio.io/2015/07/singleton-pattern-in-go/
func GetConfiguration() *Configuration {
	once.Do(func() {

		usr, err := user.Current()
		if err != nil {
			log.Fatal("YOU NEED TO RUN HERE.LOCAL AS ROOT")
		}

		if usr.Uid != "0" && usr.Gid != "0" {
			log.Fatal("YOU NEED TO RUN HERE.LOCAL AS ROOT")
		}

		devMode = viper.GetBool("dev") // retreive value from viper
		if runtime.GOOS != "linux" {
			loadDeveloperConfiguration()
		} else {
			loadConfiguration()
			configureNetworkDevices()
		}

		instance = &Configuration{}
	})
	return instance
}

//SetUserConfigs sets the configs and potentially reboots based on the delta
func (c *Configuration) SetUserConfigs(location string, ssid string, password string, authLogin string, authPassword string, document string) {
	//I dont know if this will reboot?
	reboot := false

	validLocation := generateValidHostname(location)

	if configViper.GetString("node.location") != location && configViper.GetString("location") != validLocation {
		configViper.Set("node.location", validLocation)
		reboot = true
	}

	if configViper.GetString("network.ssid") != ssid {
		configViper.Set("network.ssid", ssid)
		reboot = true
	}

	if configViper.GetString("network.password") != ssid {
		configViper.Set("network.password", password)
		reboot = true
	}

	configViper.Set("authentication.login", authLogin)
	configViper.Set("authentication.password", authPassword)
	configViper.Set("node.document", document)

	if reboot {
		rebootNode()
	}
}

//WillNeedReboot ...
func (c *Configuration) WillNeedReboot(location string, ssid string, password string) bool {
	reboot := false
	validLocation := generateValidHostname(location)
	if configViper.GetString("node.location") != location && configViper.GetString("location") != validLocation {
		reboot = true
	}

	if configViper.GetString("network.ssid") != ssid {
		reboot = true
	}

	if configViper.GetString("network.password") != ssid {
		reboot = true
	}

	return reboot
}

//GetLocation ...
func (c *Configuration) GetLocation() string {
	return configViper.GetString("node.location")
}

//GetDocument ...
func (c *Configuration) GetDocument() string {
	return configViper.GetString("node.document")
}

//GetBasicAuthLogin ...
func (c *Configuration) GetAuthenticationLogin() string {
	return configViper.GetString("authentication.login")
}

//GetBasicAuthPassword ...
func (c *Configuration) GetAuthenticationPassword() string {
	return configViper.GetString("authentication.password")
}

//GetSSID ...
func (c *Configuration) GetSSID() string {
	return configViper.GetString("network.ssid")
}

//GetPassword ...
func (c *Configuration) GetPassword() string {
	return configViper.GetString("network.password")
}

//GetIP ...
func (c *Configuration) GetIP() string {
	return envViper.GetString("ip")
}

//GetMode ...
func (c *Configuration) GetMode() string {
	return envViper.GetString("mode")
}

//GetSSIDs ...
func (c *Configuration) GetSSIDs() []string {
	return getAvailableNetworkSSIDS()
}

func loadConfiguration() {
	devMode = viper.GetBool("dev") // retreive value from viper
	configViper = viper.New()
	envViper = viper.New()
	var path string

	if devMode {
		path = "./here.local.config.toml"
	} else {
		path = "/boot/here.local.config.toml"
	}

	exist, err := fileOrDirExists(path)
	logging.Fatal(err)

	if !exist { //we copy the config file if not existing
		input, err := ioutil.ReadFile("./file-templates/here.local.config.toml.template")
		logging.Fatal(err)

		err = ioutil.WriteFile(path, input, 0666)
		logging.Fatal(err)
		os.Chmod(path, 0666) //need this to make sure we set the file permissions (WriteFile will not do it alone)

	}

	configViper.SetConfigFile(path)
	err = configViper.ReadInConfig() // Find and read the config file
	logging.Fatal(err)

	location := configViper.GetString("node.location")

	if location == "" {
		configViper.Set("node.location", "HERE-"+randSeq(6))
		err = configViper.WriteConfig()
		logging.Fatal(err)
		err = changeHostname(location)
		logging.Fatal(err)
	}

	hostname, err := os.Hostname()
	logging.Fatal(err)
	validLocation := generateValidHostname(location)
	if validLocation != hostname {
		configViper.Set("node.location", validLocation)
		err = configViper.WriteConfig()
		logging.Fatal(err)
		err = changeHostname(validLocation)
		logging.Fatal(err)
	}

	//if the location is not set then we have the whole host issue all over again
}

func loadDeveloperConfiguration() {
	fmt.Println("Loading Developer Configurations")
	fmt.Println("Warning: Key network configurations are missing")
	fmt.Println("Developer mode only allow testing http server part")
	configViper = viper.New()
	envViper = viper.New()
	envViper.Set("mode", DEVELOPER_MODE)
	path := "./here.local.config.toml"

	exist, err := fileOrDirExists(path)
	logging.Fatal(err)

	if !exist { //we copy the config file if not existing
		input, err := ioutil.ReadFile("./file-templates/here.local.config.toml.template")
		logging.Fatal(err)

		err = ioutil.WriteFile(path, input, 0666)
		logging.Fatal(err)
		os.Chmod(path, 0666) //need this to make sure we set the file permissions (WriteFile will not do it alone)

	}

	configViper.SetConfigFile(path)
	err = configViper.ReadInConfig() // Find and read the config file
	logging.Fatal(err)

	location := configViper.GetString("node.location")

	if location == "" {
		hostname, err := os.Hostname()
		logging.Fatal(err)
		configViper.Set("node.location", hostname)
		err = configViper.WriteConfig()
		logging.Fatal(err)
	}

	configViper.Set("network.ssid", "developer")
	configViper.Set("network.password", "secret")

	envViper.Set("mode", DEVELOPER_MODE)
	envViper.Set("ip", "127.0.0.1")
}

// exists returns whether the given file or directory exists
// https://stackoverflow.com/a/10510783
func fileOrDirExists(path string) (bool, error) {
	_, err := os.Stat(path)

	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
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

func changeHostname(hostname string) error {

	_, stderr, err := runCommand("sudo hostnamectl set-hostname " + hostname)
	logging.Fatal(err)

	err = ioutil.WriteFile("/etc/hostname", []byte(hostname), 0666)
	logging.Fatal(err)
	err = ioutil.WriteFile("/etc/hosts", []byte("127.0.0.1\tlocalhost\n127.0.1.1\t"+hostname+"\n"), 0666)
	logging.Fatal(err)

	if stderr != "" {
		fmt.Println(stderr)
		return fmt.Errorf(stderr)
	}

	if devMode {
		fmt.Println("You need to restart to avoid sudo host not recognised errors")
	} else {
		rebootNode()
	}

	return nil
}

//we need to ensure that the hostname is a) a valid hostname and b) a valid ssid
//meaning: < 32 chars and only 0-9,a-b, A-Z, -
func generateValidHostname(hostname string) string {
	re := regexp.MustCompile(`[^a-zA-Z0-9--]`)
	str := re.ReplaceAllString(hostname, "")
	if len(str) > 32 {
		str = str[0:32]
	}
	return str
}

func writeConfig() {
	err := configViper.WriteConfig()
	logging.Fatal(err)
}

func rebootNode() {
	if mode := envViper.GetString("mode"); mode != DEVELOPER_MODE {
		_, _, err := runCommand("sudo reboot now")
		logging.Fatal(err)
	} else {
		fmt.Println("RUNNING IN DEVELOPER MODE -- RESTART COMMAND IGNORED!")
	}

}
