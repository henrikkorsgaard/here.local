package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path"
	"text/template"

	"github.com/henrikkorsgaard/here.local/logging"
)

type config struct {
	Location          string   `json:"location,omitempty"`
	SSID              string   `json:"ssid,omitempty"`
	Password          string   `json:"password,omitempty"`
	Document          string   `json:"document,omitempty"`
	BasicAuthLogin    string   `json:"ba_login,omitempty"`
	BasicAuthPassword string   `json:"ba_password,omitempty"`
	SSIDs             []string `json:"ssids,omitempty"`
	Reboot            bool     `json:"rebbot,omitempty"`
}

func configHandler(w http.ResponseWriter, r *http.Request) {

	//create the config construct from configuration.Getters!

	ssids := []string{"SSID1", "SSID2", "SSID3", "SSID4"}
	//generate the config from configurations
	configs := config{"Henrik's office", "SSID3", "secret", "webstrate", "admin", "pass", ssids, false}

	if r.Method == "POST" {
		var clientConfigs config
		err := json.NewDecoder(r.Body).Decode(&clientConfigs)
		logging.Fatal(err)
		fmt.Printf("%+v", clientConfigs)

		//we want to return a message if a) we are good on the changes or b) if we reboot
		//setheader to return json
		w.WriteHeader(200)
	} else if r.Method == "GET" {

		t := template.New("")
		t.ParseFiles(path.Join(templatePath, "config.tmpl"))
		if err := t.ExecuteTemplate(w, "config", configs); err != nil {
			logging.Fatal(err)
		}
	}
	return
}
