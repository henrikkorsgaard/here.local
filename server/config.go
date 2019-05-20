package server

import (
	"fmt"
	"net/http"
	"path"
	"text/template"

	"github.com/henrikkorsgaard/here.local/logging"
)

type config struct {
	Location          string
	SSID              string
	Password          string
	Document          string
	BasicAuthLogin    string
	BasicAuthPassword string
	SSIDs             []string
}

func configHandler(w http.ResponseWriter, r *http.Request) {

	//create the config construct from configuration.Getters!

	t := template.New("")
	t.ParseFiles(path.Join(templatePath, "config.tmpl"))
	ssids := []string{"SSID1", "SSID2", "SSID3", "SSID4"}
	//generate the config from configurations
	configs := config{"Henrik's office", "SSID3", "secret", "webstrate", "admin", "pass", ssids}

	if r.Method == "POST" {
		err := r.ParseForm()
		logging.Fatal(err)

		reboot := r.FormValue("reboot")
		if reboot == "true" {
			fmt.Println("reboot node")
		} else {

			newConfigs := config{r.FormValue("location"), r.FormValue("ssid"), r.FormValue("password"), r.FormValue("document"), r.FormValue("balogin"), r.FormValue("bapassword"), nil}

			//compare everything now!

			configs = newConfigs
			configs.SSIDs = ssids

		}

		fmt.Println("hep")
	}

	if err := t.ExecuteTemplate(w, "config", configs); err != nil {
		logging.Fatal(err)
	}

}
