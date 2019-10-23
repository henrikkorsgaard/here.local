package server

import (
	"net/http"
)

func configHandler(w http.ResponseWriter, r *http.Request) {
	/*
		userSettings := models.GetUserSettings()

		if r.Method == "POST" {
			var clientSettings models.UserSettings
			err := json.NewDecoder(r.Body).Decode(&clientSettings)
			logging.Fatal(err)
			fmt.Printf("%+v", clientSettings)

			reboot := configuration.UpdateUserSettings(clientSettings)

			if reboot {
				fmt.Fprint(w, `{"reboot": true}`)
			} else {
				fmt.Fprint(w, `{"reboot": false}`)
			}
		} else if r.Method == "GET" {
			t := template.New("")
			t.ParseFiles(path.Join(templatePath, "config.tmpl"))
			if err := t.ExecuteTemplate(w, "config", userSettings); err != nil {
				logging.Fatal(err)
			}
		}
		return*/
}

func configHandlerPost(w http.ResponseWriter, r *http.Request) {

}
