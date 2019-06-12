package proximity

import (
	"log"
	"net/rpc"
	"os/user"
	"runtime"

	"github.com/henrikkorsgaard/here.local/logging"
)

type RawDevice struct {
	Mac    string
	Signal int
}

type NmapDevice struct {
	Mac      string
	Ip       string
	Hostname string
	Vendor   string
}

func Run() {

	usr, err := user.Current()
	if err != nil {
		logging.Fatal(err)
	}

	if runtime.GOOS != "linux" {
		simulateProximityData()
	} else {
		if usr.Uid != "0" && usr.Gid != "0" {

		} else {
			log.Fatal("YOU NEED TO RUN HERE.LOCAL AS ROOT")
		}
	}
}

func setupRPCClient() {

	//we need to do this in a go routine

	var err error
	client, err = rpc.DialHTTP("tcp", "localhost:1339")
	if err != nil {
		log.Fatal("Error setting up RPC connection: ", err)
	}
}

func simulateProximityData() {

}

func scan() {

}
