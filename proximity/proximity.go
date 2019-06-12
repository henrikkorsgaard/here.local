package proximity

import (
	"fmt"
	"log"
	"net/rpc"
	"os/user"
	"runtime"
	"time"

	"github.com/henrikkorsgaard/here.local/configuration"
	"github.com/henrikkorsgaard/here.local/logging"
)

var (
	rpcClient *rpc.Client
)

func Run() {

	usr, err := user.Current()
	if err != nil {
		logging.Fatal(err)
	}

	if runtime.GOOS != "linux" {
		//go setupRPCClient()
		//simulate()
	} else {
		if usr.Uid != "0" && usr.Gid != "0" {

		} else {
			log.Fatal("YOU NEED TO RUN HERE.LOCAL AS ROOT")
		}
	}
}

func setupRPCClient() {
	for {
		c, err := rpc.DialHTTP("tcp", configuration.ContextServerAddress)

		if err != nil || c == nil {
			fmt.Println(err)
			time.Sleep(10 * time.Second)
		} else {
			rpcClient = c
			break
		}
	}
}

func scan() {

}
