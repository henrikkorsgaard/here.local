package proximity

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/rpc"

	"github.com/henrikkorsgaard/here.local/initialise/configs"
)

var (
	rpcClient *rpc.Client
	configTLS tls.Config
)

//NOTE: we want to do some "bulk" sends of data, so to avoid near syncronous transmits that lead to increased db writes on the context server.
//NOTE: we could do some normalisation here with a kalman filter?

func Run() {

	var err error
	configTLS, err = configs.GetTLSClientConfig()
	if err != nil {
		log.Fatalf("configurationGetTLSClientConfig failed: %s", err)
	}

	mode := configs.MODE
	if mode == configs.DEVELOPER_MODE {
		simulate()
	} else {

	}
}

func connectRPC() {
	conn, err := tls.Dial("tcp", configs.CONTEXT_SERVER_ADDR, &configTLS)
	if err != nil {
		fmt.Println("Unable to establish RCP connection")
		fmt.Println("Trying again in a second")
		return
	}
	//defer conn.Close()
	rpcClient = rpc.NewClient(conn)
}
