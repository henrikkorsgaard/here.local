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

func configHandler(w http.ResponseWriter, r *http.Request) {

	nodeSettings := settings{configs.GetLocation(), configs.GetSSID(), configs.GetPassword(), configs.GetDocument(), configs.GetAuthenticationLogin(), configs.GetAuthenticationPassword(), configs.GetSSIDs(), false}

	if r.Method == "POST" {
		var clientConfigs settings
		err := json.NewDecoder(r.Body).Decode(&clientConfigs)
		logging.Fatal(err)
		fmt.Printf("%+v", clientConfigs)
		//WE NEED TO DETERMINE IF THERE IS A REBOOT INCOMING AND WE NEED TO SEND MESSAGE BACK
		rb := configs.WillNeedReboot(clientConfigs.Location, clientConfigs.SSID, clientConfigs.Password)
		if clientConfigs.Reboot || rb {
			fmt.Fprint(w, `{"reboot": true}`)
		} else {
			fmt.Fprint(w, `{"reboot": false}`)
		}
	} else if r.Method == "GET" {

		t := template.New("")
		t.ParseFiles(path.Join(templatePath, "config.tmpl"))
		if err := t.ExecuteTemplate(w, "config", nodeSettings); err != nil {
			logging.Fatal(err)
		}
	}
	return
}
