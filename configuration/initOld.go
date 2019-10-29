
package configuration

//THIS FILE IS REDUNDANT AS SOON AS I HAVE COPIED OVER THE CODE

/*

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"regexp"
	"time"
)

const (
	CONFIG_MODE    = "CONFIG"
	SLAVE_MODE     = "SLAVE"
	MASTER_MODE    = "MASTER"
	DEVELOPER_MODE = "DEVELOPER"
)

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

		err := configViper.WriteConfig()
		logging.Fatal(err)
}

func delayedReboot() {
	time.Sleep(2 * time.Second)
	reboot()
}

func reboot() {
	runCommand("sudo reboot now")
}*/
