package device

import (
	"fmt"
	"log"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type RawDevice struct {
	MAC        string
	IP         string
	Signal     int
	Hostname   string
	Vendor     string
	LocationID string
}

//this is insert once - query on API calls
type Device struct {
	ID       string
	Name     string
	Vendor   string
	LastSeen time.Time
	Public   bool
}

//this is insert once - query on API calls
type Location struct {
	MAC      string
	IP       string
	Name     string
	LastSeen time.Time
}

//this is the heavy write persist //don't know how stuff scales, but resolution should be handled by the node.
type Proximity struct {
	DeviceID   string //unique relation
	LocationID string //unique relation
	Signal     int
	Signals    []int //historical signales for additional analysis
}

//this is what is returned from the api
type UserDevice struct {
	ID        string         `json:"id"`
	Name      string         `json:"name"`
	Vendor    string         `json:"vendor"`
	IP        string         `json:"ip"`
	Signal    int            `json:"signal,omitempty"`
	Locations []UserLocation `json:"locations,omitempty"`
}

//this is what is returned from the api
type UserLocation struct {
	ID      string       `json:"id"`
	Name    string       `json:"name"`
	Vendor  string       `json:"vendor"`
	IP      string       `json:"ip"`
	Signal  int          `json:"signal,omitempty"`
	Devices []UserDevice `json:"locations,omitempty"`
}

/*


MAC        string `json:"mac"`
	IP         string `json:"ip"`
	Name       string `json:"name"`
	DeviceType string `json:"type,omitempty"`
	Zone       string `json:"zone,omitempty"`
	Wired      bool   `json:"wired,omitempty"`

	//Locations is never used internally - only for sending via external api
	Locations []Location `json:"locations,omitempty"`
	Signal    *int       `json:"signal,omitempty"`
	Signals   []int      `json:"raw_signals,omitempty"`
	Proximity *Location  `json:"proximity,omitempty"`

	locations map[string]*locationData
	kalman    kalmango.KalmanFilter


*/

func Upsert(rd RawDevice) {

	mac := Salt + rd.MAC
	hash, err := bcrypt.GenerateFromPassword([]byte(mac), 10)

	if err != nil {
		log.Fatalln("Error: ", err)
	}

	//we need to update the device in the device database
	//we need to update the proximity in the proximity database

	device := Device{hash: string(hash), Signal: rd.Signal}
	fmt.Println(device)
	fmt.Println("hep")
	//we ne
}
