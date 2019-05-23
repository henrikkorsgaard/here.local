package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path"
	"text/template"

	"github.com/henrikkorsgaard/here.local/logging"
)

type settings struct {
	Location          string   `json:"location,omitempty"`
	SSID              string   `json:"ssid,omitempty"`
	Password          string   `json:"password,omitempty"`
	Document          string   `json:"document,omitempty"`
	BasicAuthLogin    string   `json:"ba_login,omitempty"`
	BasicAuthPassword string   `json:"ba_password,omitempty"`
	SSIDs             []string `json:"ssids,omitempty"`
	Reboot            bool     `json:"rebbot,omitempty"`
}

var (
	configs := configuration.GetConfiguration()
)

func configHandler(w http.ResponseWriter, r *http.Request) {

	nodeSettings := settings{configs.GetLocation(),configs.GetSSID(),configs.GetPassword(), configs.GetDocument(), configs.GetBasicAuthLogin(), configs.GetBasicAuthPassword(), configs.GetSSIDs(), false}

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
		if err := t.ExecuteTemplate(w, "config", nodeSettings); err != nil {
			logging.Fatal(err)
		}
	}
	return
}
