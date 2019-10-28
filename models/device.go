package models

import (
	"time"
)

type Reading struct {
	DeviceMAC   string
	LocationMAC string
	Signal      int
	Timestamp   time.Time
}

type Device struct {
	ID        string
	Signal    int
	Vendor    string
	IP        string
	Name      string
	Locations map[string]Location
}

/*
type RawDevice struct {
	MAC         string
	IP          string
	Signal      int
	Hostname    string
	Vendor      string
	LocationMAC string
}

//this is insert once - query on API calls
type Device struct {
	ID        string
	Name      string
	Vendor    string
	LastSeen  time.Time
	Public    bool
	Signal    int
	Locations []Location
}

//this is insert once - query on API calls
type Location struct {
	MAC      string
	IP       string
	Name     string
	LastSeen time.Time
	Devices  []Device
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

func (d *Device) Upsert() (exists bool) {
	return false
}

func (l *Location) Upsert() (exists bool) {
	return false
}

//we actually only want to use the mysql database to record raw data for training and analysis.
/*
func UpsertRawDevice(rd RawDevice) {

	mac := Salt + rd.MAC
	hash, err := bcrypt.GenerateFromPassword([]byte(mac), 10)

	if err != nil {
		log.Fatalln("Error: ", err)
	}

	reading := Reading{locationmac: rd.LocationMAC, devicehash: string(hash), vendor: rd.Vendor, signal: rd.Signal, timestamp: time.Now()}
	reading.Insert()

	device := Device{ID: string(hash), Vendor: rd.Vendor, LastSeen: time.Now(), Signal: rd.Signal}
	fmt.Printf("%+v\n", device)
	//we ne
}*/
