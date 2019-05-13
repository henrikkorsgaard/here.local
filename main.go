package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/henrikkorsgaard/wifi"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
)

func main() {
	devMode := flag.Bool("dev", false, "Run in developer mode")
	flag.Parse()

	var logger = logrus.New()

	if *devMode {
		//This has implications toward logging!
		log.Info("Running in developer mode.")
		logger.Out = os.Stdout
		logger.SetLevel(log.InfoLevel)
	} else {

		file, err := os.OpenFile("/boot/here.local.log", os.O_CREATE|os.O_WRONLY, 0666)
		if err == nil {
			logger.Out = file
			logger.SetLevel(log.FatalLevel)
		} else {
			logger.Fatal("Unable to open log file /boot/here.local.log!")
		}
	}

	configureNetwork()
}

func configureNetwork() {
	_, err := wifi.New()
	if err != nil {
		fmt.Println(err)
		log.Panic("Unable to create nl80211 comunication client.")
	}
	//fmt.Println(c)
	log.Info("Configuring network.")
}
func detectMode() {
	log.Info("Detecting operation mode.")
}
