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

	"github.com/henrikkorsgaard/here.local/logging"
)

var (
	/*
		configViber reflects  configurations in the public configuration file on the node
		envViper reflects environement configurations
		we want to avoid mixing these, as Viper writes all the settings to the file.
	*/

	ContextServerAddress string
)

const (
	CONFIG_MODE    = "CONFIG"
	SLAVE_MODE     = "SLAVE"
	MASTER_MODE    = "MASTER"
	DEVELOPER_MODE = "DEVELOPER"
)

//NodeSettings are the settings that can be changed through the config server!
type UserSettings struct {
	Location          string   `json:"location,omitempty"`
	SSID              string   `json:"ssid,omitempty"`
	Password          string   `json:"password,omitempty"`
	Document          string   `json:"document,omitempty"`
	BasicAuthLogin    string   `json:"ba_login,omitempty"`
	BasicAuthPassword string   `json:"ba_password,omitempty"`
	SSIDs             []string `json:"ssids,omitempty"`
	Reboot            bool     `json:"rebbot,omitempty"`
}

/*
func init() {
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
}

*/

func Setup() {

	usr, err := user.Current()
	if err != nil {
		logging.Fatal(err)
	}

	if runtime.GOOS != "linux" {
		devConfig()
	} else {
		if usr.Uid != "0" && usr.Gid != "0" {
			fmt.Println("Running in deployment mode!")
		} else {
			log.Fatal("YOU NEED TO RUN HERE.LOCAL AS ROOT")
		}
	}
}

func devConfig() {
	ContextServerAddress = "localhost:1339"
}

func config() {

}

//SetUserConfigs sets the configs and potentially reboots based on the delta
func UpdateUserSettings(settings UserSettings) (reboot bool) {
	/*
		reboot = settings.Reboot

		validLocation := generateValidHostname(settings.Location)
		if settings.Location != "" && configViper.GetString("node.location") != validLocation {
			configViper.Set("node.location", validLocation)
			reboot = true
		}

		if settings.SSID != "" && configViper.GetString("network.ssid") != settings.SSID {
			configViper.Set("network.ssid", settings.SSID)
			reboot = true
		}

		if configViper.GetString("network.password") != settings.Password {
			configViper.Set("network.password", settings.Password)
			reboot = true
		}

		configViper.Set("authentication.login", settings.BasicAuthLogin)
		configViper.Set("authentication.password", settings.BasicAuthPassword)
		configViper.Set("node.document", settings.Document)

		err := configViper.WriteConfig()
		logging.Fatal(err)

		if reboot || settings.Reboot {
			//If we are rebooting then we might as well update the hostname
			err = changeHostname(validLocation, false)
			logging.Fatal(err)
			//This will delay the reboot by 2 seconds
			//giving the server time to reply to the client
			go delayedReboot()
		}

		return
	*/

	return
}

//GetUserSettings will return the current settings to the configuration server
func GetUserSettings() (settings UserSettings) {
	/*
		settings.Location = configViper.GetString("node.location")
		settings.SSID = configViper.GetString("network.ssid")
		settings.Password = configViper.GetString("network.password")
		settings.Document = configViper.GetString("node.document")
		settings.BasicAuthLogin = configViper.GetString("authentication.login")
		settings.BasicAuthPassword = configViper.GetString("authentication.password")
		settings.SSIDs = getSSIDList()
		settings.Reboot = false
	*/
	return
}

//GetLocation ...
func GetLocation() string {
	return ""
}

//GetDocument ...
func GetDocument() string {
	return ""
}

//GetBasicAuthLogin ...
func GetAuthenticationLogin() string {
	return ""
}

//GetBasicAuthPassword ...
func GetAuthenticationPassword() string {
	return ""
}

//GetSSID ...
func GetSSID() string {
	return ""
}

//GetPassword ...
func GetPassword() string {
	return ""
}

//GetIP ...
func GetIP() string {
	return ""
}

//GetMode ...
func GetMode() string {
	return ""
}

func loadConfiguration() {
	/*
		logging.Info("Loading configurations")
		devMode = viper.GetBool("dev")
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
			logging.Info("Changing hostname and rebooting")
			location = "HERE-" + randSeq(6)
			configViper.Set("node.location", location)
			err = configViper.WriteConfig()
			logging.Fatal(err)
			err = changeHostname(location, true)
			logging.Fatal(err)
		}

		hostname, err := os.Hostname()
		logging.Fatal(err)
		validLocation := generateValidHostname(location)
		if validLocation != hostname {
			configViper.Set("node.location", validLocation)
			err = configViper.WriteConfig()
			logging.Fatal(err)
			err = changeHostname(validLocation, true)
			logging.Fatal(err)
		}
	*/
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

func changeHostname(hostname string, rebootNode bool) (err error) {
	fmt.Println(hostname)
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

	if rebootNode {
		//we can't ignore a restart here becuase of the "sudo: unable to resolve host" issue
		reboot()
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
	/*
		err := configViper.WriteConfig()
		logging.Fatal(err)*/
}

func delayedReboot() {
	time.Sleep(2 * time.Second)
	reboot()
}

func reboot() {
	runCommand("sudo reboot now")
}
