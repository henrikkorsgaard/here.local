package server

import (
	"net/http"

	"github.com/henrikkorsgaard/here.local/configuration"

	"github.com/gorilla/mux"
	"github.com/henrikkorsgaard/here.local/logging"
	"github.com/henrikkorsgaard/here.local/server/api"
	"github.com/rs/cors"
)

func Run() {

	r := mux.NewRouter()

	if configuration.MODE == configuration.DEVELOPER_MODE {
		r.HandleFunc("/api/", api.RootHandler)                   //Documentation
		r.HandleFunc("/api/device", api.DeviceHandler)           //Device making the query
		r.HandleFunc("/api/device/{key}", api.DeviceHandler)     //Device by id, name, ip *IF* public
		r.HandleFunc("/api/devices", api.DevicesHandler)         //All public devices within promity to the location query device is closes too
		r.HandleFunc("/api/location", api.LocationHandler)       //Location proximity to query device device
		r.HandleFunc("/api/location/{key}", api.LocationHandler) //Location based on mac, name or ip
		r.HandleFunc("/api/locations", api.LocationsHandler)
	} else {
		apiRoute := r.Host("api.here.local").Subrouter()
		apiRoute.HandleFunc("/", api.RootHandler)                   //Documentation
		apiRoute.HandleFunc("/device", api.DeviceHandler)           //Device making the query
		apiRoute.HandleFunc("/device/{key}", api.DeviceHandler)     //Device by id, name, ip *IF* public
		apiRoute.HandleFunc("/devices", api.DevicesHandler)         //All public devices within promity to the location query device is closes too
		apiRoute.HandleFunc("/location", api.LocationHandler)       //Location proximity to query device device
		apiRoute.HandleFunc("/location/{key}", api.LocationHandler) //Location based on mac, name or ip
		apiRoute.HandleFunc("/locations", api.LocationsHandler)     //all locations seen by query device
	}

	//jsRoute := r.Host("js.here.local").Subrouter() //SDK files serve on static tho

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://here.local"},
		AllowCredentials: true,
	})

	handler := c.Handler(r)
	logging.Info("Setting up server on port 1337")
	err := http.ListenAndServe(":1337", handler)
	logging.Fatal(err)

	/*
		//runContextServer(models.ContextServerAddress)

			fmt.Println("whaa")
			var exPath string
			ex, err := os.Executable()

			logging.Fatal(err)
			exPath = filepath.Dir(ex)
			publicPath = filepath.Join(exPath, "./html/public")
			templatePath = filepath.Join(exPath, "./html/templates")

			r := mux.NewRouter()
			r.PathPrefix("/public/").Handler(http.StripPrefix("/public/", http.FileServer(http.Dir(publicPath))))

			setupConfigurationServer(r)

			c := cors.New(cors.Options{
				AllowedOrigins:   []string{"http://here.local"},
				AllowCredentials: true,
			})

			handler := c.Handler(r)
			logging.Info("Setting up server on port 1337")
			err = http.ListenAndServe(":1337", handler)
			logging.Fatal(err)
	*/

}

/*
func setupConfigurationServer(r *mux.Router) {
	r.HandleFunc("/config", BasicAuth(configHandler)).Methods("GET", "POST")
}

func setupSlave() {

	ipRouter.HandleFunc("/config", BasicAuth(configHandler)).Methods("GET", "POST")
	hostRouter.HandleFunc("/config", BasicAuth(configHandler)).Methods("GET", "POST")

}*

/*
		ipRouter.HandleFunc("/config", BasicAuth(configHandler)).Methods("GET", "POST")
		hostRouter.HandleFunc("/config", BasicAuth(configHandler)).Methods("GET", "POST")
		globalRouter.HandleFunc("/config", BasicAuth(configHandler)).Methods("GET", "POST")

	//run config server
	//run api server
	//run context server
}

//BasicAuth handles authentication to the config server
//It will fetch the user/pass information from the configuration file on every query
//to ensure that any updates are reflected without a server restart

func BasicAuth(handler http.HandlerFunc) http.HandlerFunc {

	username := configuration.GetAuthenticationLogin()
	password := configuration.GetAuthenticationPassword()

	return func(w http.ResponseWriter, r *http.Request) {

		user, pass, ok := r.BasicAuth()

		if !ok || subtle.ConstantTimeCompare([]byte(user), []byte(username)) != 1 || subtle.ConstantTimeCompare([]byte(pass), []byte(password)) != 1 {
			w.Header().Set("WWW-Authenticate", `Basic realm="Please enter username and password"`)
			w.WriteHeader(401)
			w.Write([]byte("Unauthorised.\n"))
			return
		}

		handler(w, r)
	}

}
*/
/*
func faviconHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, path.Join(publicPath, "./images/favicon.ico"))
}*/
