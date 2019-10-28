package api

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"

	"github.com/henrikkorsgaard/here.local/configuration"

	"github.com/gorilla/mux"
)

type UserDevice struct {
	Name     string       `json:"name,omitempty"`
	IP       string       `json:"ip"`
	Vendor   string       `json:"vendor,omitempty"`
	Signal   int          `json:"signal,omitempty"`
	Location UserLocation `json:"location,omitempty"`
}

type UserLocation struct {
	Name    string       `json:"name,omitempty"`
	IP      string       `json:"ip"`
	Signal  int          `json:"signal,omitempty"`
	MAC     string       `json:"mac,omitempty"`
	Devices []UserDevice `json:"devices,omitempty"`
}

func RootHandler(w http.ResponseWriter, r *http.Request) {

	fmt.Fprintf(w, "API documentation\nTo be implemented")
}

func DeviceHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["key"]
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(ip)
	if len(key) == 0 {
		ul := UserLocation{Name: configuration.NODE_NAME, IP: configuration.NODE_IP_ADDR, MAC: configuration.NODE_MAC_ADDR, Signal: -65}
		ud := UserDevice{Name: "test-device", IP: ip, Vendor: "Apple", Location: ul}
		b, err := json.Marshal(ud)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(b)
	} else {
		fmt.Println("try to decipher the key and fetch the device")
		fmt.Fprintf(w, "Device searched by id\nTo be implemented")
	}

}

func DevicesHandler(w http.ResponseWriter, r *http.Request) {

	fmt.Fprintf(w, "Devices\nTo be implemented")
}

func LocationHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Location\nTo be implemented")
}

func LocationsHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Locations\nTo be implemented")
}
