package device

import (
	"fmt"
	"log"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type RawDevice struct {
	MAC         string
	IP          string
	Signal      int
	Hostname    string
	Vendor      string
	LocationMAC string
}

type Device struct {
	ID        string    //macHashed with the device specific
	hash      string    //hashed mac address
	Signal    int       //device signal strength
	Updated   time.Time //timestamp for keeping up to data in the database
	Locations map[string]int
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

	device := Device{hash: string(hash), Signal: rd.Signal}
	fmt.Println(device)
	fmt.Println("hep")
	//we ne
}
