package configs

var (
	MODE                string
	CONTEXT_SERVER_ADDR string
	NODE_MAC_ADDR       string
	NODE_IP_ADDR        string
	NODE_NAME           string
)

const (
	CONFIG_MODE    = "CONFIG"
	SLAVE_MODE     = "SLAVE"
	MASTER_MODE    = "MASTER"
	DEVELOPER_MODE = "DEVELOPER"
)
