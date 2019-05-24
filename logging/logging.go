package logging

import (
	"flag"
	"log"
	"os"
	"runtime"

	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var logger *logrus.Logger

func init() {

	logger = logrus.New()

	flag.Bool("dev", false, "Run in developer mode")

	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	pflag.Parse()
	viper.BindPFlags(pflag.CommandLine)
	devMode := viper.GetBool("dev")

	if devMode || runtime.GOOS != "linux" {

		//This has implications toward logging!
		logrus.Info("Running in developer mode.")
		logger.Out = os.Stdout
		logger.SetLevel(logrus.DebugLevel)
	} else {
		file, err := os.OpenFile("/boot/here.local.log", os.O_CREATE|os.O_WRONLY, 0666)
		if err == nil {
			logger.Out = file
			logger.SetLevel(logrus.FatalLevel)
		} else {
			log.Fatal("Unable to open log file /boot/here.local.log!")
		}
	}
}

func Fatal(err error) {
	if err != nil {
		logger.Fatal(err)
	}
}

func Info(msg string) {
	logger.Info(msg)
}
