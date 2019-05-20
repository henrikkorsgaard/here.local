package configuration

import (
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"os/user"
	"time"

	"github.com/henrikkorsgaard/here.local/logging"
	"github.com/spf13/viper"
)

var (
	configViper *viper.Viper
	envViper    *viper.Viper
	devMode     bool
)

const (
	CONFIG_MODE = "CONFIG"
	SLAVE_MODE  = "SLAVE"
	MASTER_MODE = "MASTER"
)

func Setup() {
	//we want to make sure here.local is run as root

	usr, err := user.Current()
	if err != nil {
		log.Fatal("YOU NEED TO RUN HERE.LOCAL AS ROOT")
	}

	if usr.Uid != "0" && usr.Gid != "0" {
		log.Fatal("YOU NEED TO RUN HERE.LOCAL AS ROOT")
	}

	loadConfiguration()
	configureNetworkDevices()
}

func loadConfiguration() {
	devMode = viper.GetBool("dev") // retreive value from viper
	configViper = viper.New()
	envViper = viper.New()
	var path string

	if devMode {
		path = "./here.local.config.toml"
	} else {
		fmt.Println("in here!")
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

	location := configViper.GetString("location")

	if location == "" {
		configViper.Set("location", "HERE-"+randSeq(6))
		err = configViper.WriteConfig()
		logging.Fatal(err)
		err = changeHostname(configViper.GetString("location"))
		logging.Fatal(err)
	}

	hostname, err := os.Hostname()
	logging.Fatal(err)

	if location != hostname {
		err = changeHostname(configViper.GetString("location"))
		logging.Fatal(err)
	}

	//if the location is not set then we have the whole host issue all over again
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
		_, _, err = runCommand("sudo reboot now")
		//logging.Fatal(err)
	}

	return nil
}

func GetLocation() string {
	return configViper.GetString("location")
}

func GetBasicAuthLogin() string {
	return configViper.GetString("basic-auth.login")
}

func GetBasicAuthPassword() string {
	return configViper.GetString("basic-auth.password")
}

func GetIP() string {
	return envViper.GetString("ip")
}

func GetMode() string {
	return envViper.GetString("mode")
}
