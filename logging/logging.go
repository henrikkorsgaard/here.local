package logging

import (
	"fmt"
	"log"
	"os"
	"os/user"
	"runtime"

	"github.com/sirupsen/logrus"
)

var logger *logrus.Logger

func init() {

	logger = logrus.New()

	usr, err := user.Current()
	if err != nil {
		logger.Fatal(err)
	}

	if runtime.GOOS != "linux" {
		logrus.Info("Running in developer mode.")
		logger.Out = os.Stdout
		logger.SetLevel(logrus.DebugLevel)

	} else {
		if usr.Uid != "0" && usr.Gid != "0" {
			log.Fatal("YOU NEED TO RUN HERE.LOCAL AS ROOT")
		} else {
			file, err := os.OpenFile("/boot/here.local.log", os.O_CREATE|os.O_WRONLY, 0666)
			if err == nil {
				logger.Out = file
				logger.SetLevel(logrus.InfoLevel)
			} else {
				log.Fatal("Unable to open log file /boot/here.local.log!")
			}


		}
	}
}

func Fatal(err error) {
	if err != nil {
		logger.Fatal(err)
	}
}

func Info(msg string) {
	fmt.Println(msg)
	logger.Info(msg)
}
