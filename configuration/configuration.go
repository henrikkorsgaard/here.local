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
	ConfigViper *viper.Viper
	devMode     bool
)

func Bootstrap() {
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
	ConfigViper = viper.New()
	var path string

	if devMode {
		//check if file exists
		fmt.Println("dev file")
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

	ConfigViper.SetConfigFile(path)
	err = ConfigViper.ReadInConfig() // Find and read the config file
	logging.Fatal(err)

	location := ConfigViper.GetString("location")

	if location == "" {
		ConfigViper.Set("location", "HERE-"+randSeq(6))
		err = ConfigViper.WriteConfig()
		logging.Fatal(err)
		err = changeHostname(ConfigViper.GetString("location"))
		logging.Fatal(err)
	}

	hostname, err := os.Hostname()
	logging.Fatal(err)

	if location != hostname {
		err = changeHostname(ConfigViper.GetString("location"))
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
