package models

import (
	"math/rand"
	"time"
)

var (
	/*
		configViber reflects  configurations in the public configuration file on the node
		envViper reflects environement configurations
		we want to avoid mixing these, as Viper writes all the settings to the file.
	*/

	ContextServerAddress string
	NodeLocationName     string
	NodeLocationIP       string
	NodeLocationMAC      string
	Salt                 string //need to be written to sd
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
	Reboot            bool     `json:"reboot,omitempty"`
}

func init() {
	rand.Seed(time.Now().UnixNano())
	Salt = generateSalt(rand.Intn(128) + 64) //does not need to be long, as it is used to salt mac address to counter the finite address space (see https://en.wikipedia.org/wiki/MAC_address_anonymization)

	//this

	ContextServerAddress = "localhost:1339"
	NodeLocationMAC = "DEVELOPERNODE"
}

func generateSalt(n int) string {
	//see https://stackoverflow.com/a/31832326
	letterBytes := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Int63()%int64(len(letterBytes))]
	}

	return string(b)
}
